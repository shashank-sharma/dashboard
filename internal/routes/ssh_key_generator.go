package routes

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/shashank-sharma/backend/internal/logger"
	"github.com/shashank-sharma/backend/internal/util"
	"golang.org/x/crypto/ssh"
)

func RegisterKeyGenerationRoute(e *core.ServeEvent) {
	e.Router.POST("/api/security-keys/generate", func(c *core.RequestEvent) error {
		token := c.Request.Header.Get("Authorization")
		userID, err := util.GetUserId(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "Not authenticated",
			})
		}

		var requestData struct {
			KeyType string `json:"key_type"`
		}

		if err := json.NewDecoder(c.Request.Body).Decode(&requestData); err != nil {
			if err := c.Request.ParseForm(); err == nil {
				requestData.KeyType = c.Request.Form.Get("key_type")
			}
		}

		if requestData.KeyType == "" {
			requestData.KeyType = "ed25519"
		}

		if requestData.KeyType != "ed25519" && requestData.KeyType != "rsa" {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "Invalid key type. Must be 'ed25519' or 'rsa'",
			})
		}

		var privateKeyString, publicKeyString string
		var genErr error

		if requestData.KeyType == "ed25519" {
			privateKeyString, publicKeyString, genErr = generateED25519KeyPair()
		} else {
			privateKeyString, publicKeyString, genErr = generateRSAKeyPair()
		}

		if genErr != nil {
			logger.LogError("Failed to generate SSH key", "error", genErr, "userId", userID)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Failed to generate SSH key: " + genErr.Error(),
			})
		}

		logger.LogInfo("Generated SSH key", "keyType", requestData.KeyType, "userId", userID)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"private_key": privateKeyString,
			"public_key":  publicKeyString,
			"key_type":    requestData.KeyType,
		})
	})
}

func generateED25519KeyPair() (string, string, error) {
    pub, priv, err := ed25519.GenerateKey(nil)
    if err != nil {
        panic(err)
    }
    p, err := ssh.MarshalPrivateKey(crypto.PrivateKey(priv), "")
    if err != nil {
        panic(err)
    }
    privateKeyPem := pem.EncodeToMemory(p)
    privateKeyString := string(privateKeyPem)
    publicKey, err := ssh.NewPublicKey(pub)
    if err != nil {
        panic(err)
	}
	publicKeyString := "ssh-ed25519" + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())
	logger.LogInfo("Generated ED25519 key pair", "pubKey", publicKeyString, "privKey", privateKeyString)
	return privateKeyString, publicKeyString, nil
}

func generateRSAKeyPair() (string, string, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate RSA key: %w", err)
	}

	sshPubKey, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert RSA public key to SSH format: %w", err)
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(sshPubKey)
	publicKeyString := strings.TrimSuffix(string(pubKeyBytes), "\n")

	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	privateKeyString := string(pem.EncodeToMemory(pemBlock))

	return privateKeyString, publicKeyString, nil
} 