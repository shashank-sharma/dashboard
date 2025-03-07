package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pocketbase/pocketbase/core"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/util"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development
	CheckOrigin: func(r *http.Request) bool {
		logger.LogInfo("WebSocket CheckOrigin", "origin", r.Header.Get("Origin"), "host", r.Host)
		return true
	},
	EnableCompression: true,
}

type SSHOutputListener struct {
	UserID       string
	ConnectionID string
	Conn         *websocket.Conn
	mu           sync.Mutex
	closed       bool
}

func (l *SSHOutputListener) Send(data []byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return fmt.Errorf("connection closed")
	}

	return l.Conn.WriteMessage(websocket.TextMessage, data)
}

func (l *SSHOutputListener) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.closed {
		l.closed = true
		l.Conn.Close()
	}
}

var sshOutputListeners = struct {
	sync.RWMutex
	listeners map[string][]*SSHOutputListener
}{
	listeners: make(map[string][]*SSHOutputListener),
}

func RegisterSSHWebSocketRoutes(e *core.ServeEvent) {
	e.Router.GET("/api/ssh/stream", func(c *core.RequestEvent) error {
		logger.LogInfo("WebSocket connection request received", "host", c.Request.Host, "uri", c.Request.RequestURI)
		
		// Parse the connection ID from query parameters first for better error messages
		connectionID := c.Request.URL.Query().Get("connection_id")
		if connectionID == "" {
			logger.LogError("WebSocket missing connection_id parameter")
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Missing connection_id parameter",
			})
		}
		
		logger.LogInfo("WebSocket connection attempt", "connectionID", connectionID)
		
		token := c.Request.URL.Query().Get("token")
		if token == "" {
			token = c.Request.Header.Get("Authorization")
		}
		
		tokenPrefix := ""
		if len(token) > 10 {
			tokenPrefix = token[:10] + "..."
		} else if len(token) > 0 {
			tokenPrefix = token + "..."
		} else {
			tokenPrefix = "(empty)"
		}
		
		logger.LogInfo("WebSocket auth token received", "tokenExists", token != "", "tokenLength", len(token), "tokenPrefix", tokenPrefix)
		
		userID, err := util.GetUserId(token)
		if err != nil {
			logger.LogError("WebSocket auth failed", "error", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "Not authenticated",
			})
		}
		
		logger.LogInfo("WebSocket auth successful", "userID", userID)
		
		// Verify that the connection exists and belongs to the user
		connectionManager.mu.Lock()
		conn, exists := connectionManager.connections[connectionID]		
		connectionCount := len(connectionManager.connections)
		connectionIDs := make([]string, 0, connectionCount)
		for id := range connectionManager.connections {
			connectionIDs = append(connectionIDs, id)
		}
		logger.LogInfo("Connection manager state", "connectionCount", connectionCount, "connectionIDs", connectionIDs)
		
		connectionManager.mu.Unlock()

		if !exists {
			logger.LogError("WebSocket connection not found", "connectionID", connectionID)
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Connection not found",
			})
		}

		logger.LogInfo("WebSocket connection details", "connectionID", connectionID, "requestUserID", userID, "connectionUserID", conn.UserID)
		if conn.UserID != userID {
			logger.LogError("WebSocket unauthorized access", "connectionID", connectionID, "requestUserID", userID, "connectionUserID", conn.UserID)
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "Not authorized to access this connection",
			})
		}
		
		logger.LogInfo("WebSocket connection authorized, upgrading connection", "connectionID", connectionID, "userID", userID)

		// Perform the WebSocket upgrade
		upgrader.CheckOrigin = func(r *http.Request) bool {
			logger.LogInfo("Checking origin for WebSocket", "origin", r.Header.Get("Origin"), "host", r.Host)
			return true
		}
		
		wsConn, err := upgrader.Upgrade(c.Response, c.Request, nil)
		if err != nil {
			logger.LogError("Failed to upgrade connection to WebSocket", "error", err, "connectionID", connectionID)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to upgrade connection: " + err.Error(),
			})
		}
		
		logger.LogInfo("WebSocket connection upgraded successfully", "connectionID", connectionID)

		// Create a new listener for this connection
		listener := &SSHOutputListener{
			UserID:       userID,
			ConnectionID: connectionID,
			Conn:         wsConn,
			closed:       false,
		}

		sshOutputListeners.Lock()
		if _, ok := sshOutputListeners.listeners[connectionID]; !ok {
			sshOutputListeners.listeners[connectionID] = make([]*SSHOutputListener, 0)
		}
		sshOutputListeners.listeners[connectionID] = append(sshOutputListeners.listeners[connectionID], listener)
		sshOutputListeners.Unlock()

		logger.LogInfo("New WebSocket connection established", "userID", userID, "connectionID", connectionID)

		listener.Send([]byte(fmt.Sprintf("\r\nConnected to terminal. Session ID: %s\r\n", connectionID)))

		go func() {
			time.Sleep(500 * time.Millisecond)
			
			shellSession, err := conn.Client.NewSession()
			if err != nil {
				logger.LogError("Failed to create shell session", "error", err)
				listener.Send([]byte("\r\n\x1B[31mFailed to initialize shell: " + err.Error() + "\x1B[0m\r\n"))
				return
			}
			
			conn.Session = shellSession
			
			if err := shellSession.RequestPty("xterm", 80, 24, ssh.TerminalModes{
				ssh.ECHO:          1,
				ssh.TTY_OP_ISPEED: 14400,
				ssh.TTY_OP_OSPEED: 14400,
			}); err != nil {
				logger.LogError("Failed to request PTY", "error", err)
				listener.Send([]byte("\r\n\x1B[31mFailed to initialize terminal: " + err.Error() + "\x1B[0m\r\n"))
				conn.Session = nil
				shellSession.Close()
				return
			}
			
			stdin, err := shellSession.StdinPipe()
			if err != nil {
				logger.LogError("Failed to get stdin pipe", "error", err)
				listener.Send([]byte("\r\n\x1B[31mFailed to initialize shell input: " + err.Error() + "\x1B[0m\r\n"))
				conn.Session = nil
				shellSession.Close()
				return
			}
			
			stdout, err := shellSession.StdoutPipe()
			if err != nil {
				logger.LogError("Failed to get stdout pipe", "error", err)
				listener.Send([]byte("\r\n\x1B[31mFailed to initialize shell output: " + err.Error() + "\x1B[0m\r\n"))
				conn.Session = nil
				shellSession.Close()
				return
			}
			
			stderr, err := shellSession.StderrPipe()
			if err != nil {
				logger.LogError("Failed to get stderr pipe", "error", err)
				listener.Send([]byte("\r\n\x1B[31mFailed to initialize shell error output: " + err.Error() + "\x1B[0m\r\n"))
				conn.Session = nil
				shellSession.Close()
				return
			}
			
			conn.StdinPipe = stdin
			
			if err := shellSession.Shell(); err != nil {
				logger.LogError("Failed to start shell", "error", err)
				listener.Send([]byte("\r\n\x1B[31mFailed to start shell: " + err.Error() + "\x1B[0m\r\n"))
				conn.Session = nil
				conn.StdinPipe = nil
				shellSession.Close()
				return
			}
			
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := stdout.Read(buf)
					if err != nil {
						if err != io.EOF {
							logger.LogError("Error reading from stdout", "error", err)
						}
						break
					}
					if n > 0 {
						SendOutputToListeners(connectionID, buf[:n])
					}
				}
			}()
			
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := stderr.Read(buf)
					if err != nil {
						if err != io.EOF {
							logger.LogError("Error reading from stderr", "error", err)
						}
						break
					}
					if n > 0 {
						SendOutputToListeners(connectionID, buf[:n])
					}
				}
			}()
			
			err = shellSession.Wait()
			if err != nil {
				logger.LogError("Shell session ended with error", "error", err)
			} else {
				logger.LogInfo("Shell session ended normally")
			}
			
			conn.Session = nil
			conn.StdinPipe = nil
		}()

		go func() {
			defer func() {
				removeListener(connectionID, listener)
				logger.LogInfo("WebSocket connection closed", "connectionID", connectionID)
			}()
			
			for {
				messageType, message, err := wsConn.ReadMessage()
				if err != nil {
					break
				}
				
				if messageType == websocket.TextMessage {
					var messageData map[string]interface{}
					if err := json.Unmarshal(message, &messageData); err == nil {
						msgType, hasType := messageData["type"].(string)
						if hasType {
							switch msgType {
							case "auth":
								break
							case "input":
								// Keyboard input
								if data, ok := messageData["data"].(string); ok {
									connectionManager.mu.Lock()
									conn, exists := connectionManager.connections[connectionID]
									
									if exists && conn.StdinPipe != nil {
										_, err := conn.StdinPipe.Write([]byte(data))
										if err != nil {
											logger.LogError("Failed to write to shell stdin", "error", err)
										}
									}
									connectionManager.mu.Unlock()
								}
							case "ping":
								connectionManager.mu.Lock()
								if conn, exists := connectionManager.connections[connectionID]; exists {
									conn.LastUsed = time.Now()
								}
								connectionManager.mu.Unlock()
								
								pongErr := listener.Send([]byte(fmt.Sprintf(`{"type":"pong","time":%d}`, time.Now().UnixMilli())))
								if pongErr != nil {
									logger.LogError("Failed to send pong response", "error", pongErr, "connectionID", connectionID)
								}
							default:
								logger.LogInfo("Unknown message type received", "type", msgType)
							}
						}
					}
				}
			}
		}()

		return nil
	})
}

// SendOutputToListeners sends output to all WebSocket listeners for a connection
func SendOutputToListeners(connectionID string, data []byte) {
	sshOutputListeners.RLock()
	listeners, ok := sshOutputListeners.listeners[connectionID]
	sshOutputListeners.RUnlock()

	if !ok {
		return
	}

	for _, listener := range listeners {
		err := listener.Send(data)
		if err != nil {
			// If we can't send, remove the listener
			removeListener(connectionID, listener)
		}
	}
}

// removeListener removes a WebSocket listener from the list
func removeListener(connectionID string, listener *SSHOutputListener) {
	sshOutputListeners.Lock()
	defer sshOutputListeners.Unlock()

	listeners, ok := sshOutputListeners.listeners[connectionID]
	if !ok {
		return
	}

	for i, l := range listeners {
		if l == listener {
			l.Close()
			
			sshOutputListeners.listeners[connectionID] = append(
				listeners[:i],
				listeners[i+1:]...,
			)
			break
		}
	}

	if len(sshOutputListeners.listeners[connectionID]) == 0 {
		delete(sshOutputListeners.listeners, connectionID)
	}
} 