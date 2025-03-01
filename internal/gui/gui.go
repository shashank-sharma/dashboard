package gui

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shashank-sharma/backend/internal/util"
)

var (
	// Theme colors
	primaryColor   = lipgloss.Color("#3B82F6") // Blue
	successColor   = lipgloss.Color("#10B981") // Green
	errorColor     = lipgloss.Color("#EF4444") // Red
	warningColor   = lipgloss.Color("#F59E0B") // Amber
	infoColor      = lipgloss.Color("#6366F1") // Indigo
	backgroundColor = lipgloss.Color("#1F2937") // Dark blue-gray
	textColor      = lipgloss.Color("#F3F4F6") // Light gray
	mutedTextColor = lipgloss.Color("#9CA3AF") // Medium gray
	accentColor    = lipgloss.Color("#8B5CF6") // Purple

	// Styles
	titleStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(primaryColor).
		Padding(0, 2).
		Bold(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor)

	appStyle = lipgloss.NewStyle().
		Padding(2).
		Background(backgroundColor)

	statusBarStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(backgroundColor).
		Padding(0, 1)

	statusOkStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(successColor).
		Padding(0, 1).
		Bold(true)

	statusErrorStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(errorColor).
		Padding(0, 1).
		Bold(true)

	infoStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(infoColor).
		Padding(0, 1)

	warningStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(warningColor).
		Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(2).
		MarginRight(1)

	activeTabStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(primaryColor).
		Padding(0, 3).
		Bold(true).
		MarginRight(2)

	inactiveTabStyle = lipgloss.NewStyle().
		Foreground(mutedTextColor).
		Background(backgroundColor).
		Padding(0, 3).
		MarginRight(2)

	subtitleStyle = lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		MarginBottom(1).
		Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
		Foreground(mutedTextColor).
		Italic(true).
		Padding(0, 1)

	filterPromptStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)

	filterInputStyle = lipgloss.NewStyle().
		Foreground(textColor).
		Background(backgroundColor).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(accentColor).
		Padding(0, 1)
		
	metadataLabelStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true).
		PaddingRight(1)
		
	metadataValueStyle = lipgloss.NewStyle().
		Foreground(textColor)
		
	metadataBoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1).
		MarginTop(1)
)

// focusedTab is the currently active tab
type focusedTab int

const (
	homeTab focusedTab = iota
	appLogTab
	prometheusTab
)

// Model represents the state of the GUI
type Model struct {
	// General state
	ready          bool
	lastUpdateTime time.Time
	width          int
	height         int
	activeTab      focusedTab

	// App logs view
	appLogViewport  viewport.Model
	logFilePath     string
	logFilterInput  textinput.Model
	filteredLogContent string
	rawLogContent   string
	autoScroll      bool // Track if we should auto-scroll

	// Prometheus logs view
	prometheusViewport viewport.Model
	prometheusContent  string

	// Server metadata
	guiStatus GUIStatus
	serverMetadata ServerMetadata
	
	// Cache for rendered home tab to prevent flickering
	cachedHomeTab string
	homeTabRendered bool
}

type GUIStatus struct {
	ServerRunning bool
	MetricsEnabled bool
}

type ServerMetadata struct {
	ServerURL      string
	ServerVersion  string
	Environment    string
	EnvVariables   map[string]any
	CronJobs       []CronJob
	StartTime      time.Time
	DataDirectory  string
	APIEndpoints   []string
}

type CronJob struct {
	Name     string
	Schedule string
	LastRun  time.Time
	Status   string
}

// Message types
type tickMsg time.Time
type serverStatusMsg bool
type metricsStatusMsg bool
type logContentMsg string
type prometheusLogMsg string

// StartGUI initializes and runs the Bubble Tea GUI
func StartGUI(logFilePath string, guiStatus GUIStatus, metadata ServerMetadata) error {
	model := initialModel(logFilePath, guiStatus, metadata)
	p := tea.NewProgram(
		&model,
		tea.WithAltScreen(),       // Use the full size of the terminal
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	_, err := p.Run()
	return err
}

// Create a new model with initial values
func initialModel(logFilePath string, guiStatus GUIStatus, metadata ServerMetadata) Model {
	// Initialize app logs viewport
	appVp := viewport.New(10, 10) // Start with non-zero dimensions
	appVp.Style = panelStyle
	appVp.SetContent("Initializing application logs...") // Start with some content

	// Initialize prometheus logs viewport
	promVp := viewport.New(10, 10) // Start with non-zero dimensions
	promVp.Style = panelStyle
	promVp.SetContent("Initializing Prometheus metrics...") // Start with some content

	// Initialize filter input
	ti := textinput.New()
	ti.Placeholder = "Filter logs..."
	ti.CharLimit = 50
	ti.Width = 30
	ti.Prompt = "› "
	ti.PromptStyle = filterPromptStyle
	ti.TextStyle = lipgloss.NewStyle().Foreground(textColor)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(mutedTextColor)

	return Model{
		appLogViewport:     appVp,
		prometheusViewport: promVp,
		logFilePath:        logFilePath,
		logFilterInput:     ti,
		guiStatus:          guiStatus,
		serverMetadata:     metadata,
		lastUpdateTime:     time.Now(),
		autoScroll:         true, // Start with auto-scrolling enabled
		activeTab:          homeTab, // Start with the home tab as default
		// Initialize with non-empty content
		rawLogContent:     "Waiting for logs...",
		filteredLogContent: "Waiting for logs...",
		prometheusContent:  "Waiting for metrics...",
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	// Initial load of all data
	return tea.Batch(
		readLogFile(m.logFilePath),
		readPrometheusLogs(),
		tickEvery(1 * time.Second),
		textinput.Blink,
	)
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// First check if filter input is focused
		if m.logFilterInput.Focused() {
			switch msg.Type {
			case tea.KeyEsc, tea.KeyEnter:
				m.logFilterInput.Blur()
				m.applyLogFilter()
				return m, nil
			default:
				var cmd tea.Cmd
				m.logFilterInput, cmd = m.logFilterInput.Update(msg)
				cmds = append(cmds, cmd)
				
				// Apply filter when text changes
				m.applyLogFilter()
				
				// Update app log viewport with filtered content
				m.appLogViewport.SetContent(m.filteredLogContent)
				
				return m, tea.Batch(cmds...)
			}
		}
		
		// If filter input is not focused, handle other keyboard events
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			// Cycle through tabs: home -> app logs -> prometheus -> home
			switch m.activeTab {
			case homeTab:
				m.activeTab = appLogTab
			case appLogTab:
				m.activeTab = prometheusTab
			case prometheusTab:
				m.activeTab = homeTab
			}
		case "1":
			m.activeTab = homeTab
		case "2":
			m.activeTab = appLogTab
		case "3":
			m.activeTab = prometheusTab
		case "r":
			// Manual refresh
			cmds = append(cmds, readLogFile(m.logFilePath), readPrometheusLogs())
		case "g":
			// Go to top of logs in active viewport
			if m.activeTab == appLogTab {
				SafeGotoTop(&m.appLogViewport)
			} else if m.activeTab == prometheusTab {
				SafeGotoTop(&m.prometheusViewport)
			}
			m.autoScroll = false // Disable auto-scroll when manually navigating
		case "G":
			// Go to bottom of logs in active viewport
			if m.activeTab == appLogTab {
				SafeGotoBottom(&m.appLogViewport)
			} else if m.activeTab == prometheusTab {
				SafeGotoBottom(&m.prometheusViewport)
			}
			m.autoScroll = true // Enable auto-scroll when going to the bottom
		case "s":
			// Toggle auto-scroll
			m.autoScroll = !m.autoScroll
			if m.autoScroll {
				// If enabling auto-scroll, go to the bottom
				if m.activeTab == appLogTab {
					SafeGotoBottom(&m.appLogViewport)
				} else if m.activeTab == prometheusTab {
					SafeGotoBottom(&m.prometheusViewport)
				}
			}
		case "f":
			// Focus on filter input (only when in app logs tab)
			if m.activeTab == appLogTab {
				m.logFilterInput.Focus()
				return m, textinput.Blink
			}
		}

	case tea.MouseMsg:
		// When user scrolls with mouse, disable auto-scroll
		if msg.Type == tea.MouseWheelUp || msg.Type == tea.MouseWheelDown {
			// Check if we're at the very bottom before the scroll
			var atBottom bool
			if m.activeTab == appLogTab {
				atBottom = m.appLogViewport.AtBottom()
				
				// Let the viewport handle the mouse event first
				newViewport, cmd := m.appLogViewport.Update(msg)
				m.appLogViewport = newViewport
				cmds = append(cmds, cmd)
			} else {
				atBottom = m.prometheusViewport.AtBottom()
				
				// Let the viewport handle the mouse event first
				newViewport, cmd := m.prometheusViewport.Update(msg)
				m.prometheusViewport = newViewport
				cmds = append(cmds, cmd)
			}
			
			// If we were at the bottom and scrolled down, keep auto-scroll on
			// Otherwise, disable auto-scroll
			if msg.Type == tea.MouseWheelDown && atBottom {
				m.autoScroll = true
			} else {
				m.autoScroll = false
			}
			
			// Return early since we've already updated the viewport
			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Calculate dimensions for UI components
		headerHeight := 12    // Space for header elements with padding and borders
		footerHeight := 5     // Space for help text with padding and border
		filterHeight := 3     // Space for filter input (when visible)
		contentWidth := m.width - 6  // Content width with margins
		
		// Calculate content heights based on tab
		var contentHeight int
		if m.activeTab == appLogTab {
			contentHeight = m.height - headerHeight - footerHeight - filterHeight
		} else {
			contentHeight = m.height - headerHeight - footerHeight
		}
		
		// Update viewport dimensions
		m.appLogViewport.Width = contentWidth
		m.appLogViewport.Height = contentHeight
		
		m.prometheusViewport.Width = contentWidth
		m.prometheusViewport.Height = contentHeight
		
		m.logFilterInput.Width = contentWidth - 10
		
		// Mark as ready after first resize
		if !m.ready {
			m.ready = true
		} else {
			// Clear home tab cache on resize to ensure proper rendering
			m.cachedHomeTab = ""
			m.homeTabRendered = false
		}

	case tickMsg:
		// Update timestamp
		m.lastUpdateTime = time.Time(msg)
		
		// Selectively refresh only the active tab data
		var cmd tea.Cmd
		switch m.activeTab {
		case appLogTab:
			cmd = readLogFile(m.logFilePath)
		case prometheusTab:
			cmd = readPrometheusLogs()
		}
		
		// Schedule the next tick
		cmds = append(cmds, tickEvery(1*time.Second))
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case serverStatusMsg:
		m.guiStatus.ServerRunning = bool(msg)
		
	case metricsStatusMsg:
		m.guiStatus.MetricsEnabled = bool(msg)

	case logContentMsg:
		// Store the raw content
		m.rawLogContent = string(msg)
		
		// Apply filter if any
		m.applyLogFilter()
		
		// Remember if we were at the bottom before updating content
		atBottom := m.appLogViewport.AtBottom() || m.autoScroll
		
		// Update content with filtered logs
		m.appLogViewport.SetContent(m.filteredLogContent)
		
		// Only scroll to bottom if we were already at the bottom or auto-scroll is enabled
		if atBottom && m.activeTab == appLogTab {
			SafeGotoBottom(&m.appLogViewport)
		}
		
	case prometheusLogMsg:
		// Store the content
		m.prometheusContent = string(msg)
		
		// Remember if we were at the bottom before updating content
		atBottom := m.prometheusViewport.AtBottom() || m.autoScroll
		
		// Update content
		m.prometheusViewport.SetContent(m.prometheusContent)
		
		// Only scroll to bottom if we were already at the bottom or auto-scroll is enabled
		if atBottom && m.activeTab == prometheusTab {
			SafeGotoBottom(&m.prometheusViewport)
		}
	}

	// Handle viewport updates based on active tab
	if m.activeTab == appLogTab {
		if _, ok := msg.(tea.MouseMsg); !ok {
			var cmd tea.Cmd
			m.appLogViewport, cmd = m.appLogViewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	} else if m.activeTab == prometheusTab {
		if _, ok := msg.(tea.MouseMsg); !ok {
			var cmd tea.Cmd
			m.prometheusViewport, cmd = m.prometheusViewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// Apply the current filter to the raw log content
func (m *Model) applyLogFilter() {
	filter := m.logFilterInput.Value()
	if filter == "" {
		// No filter, use raw content
		m.filteredLogContent = m.rawLogContent
		return
	}
	
	// Apply filter: only include lines that contain the filter text
	var filteredLines []string
	for _, line := range strings.Split(m.rawLogContent, "\n") {
		if strings.Contains(strings.ToLower(line), strings.ToLower(filter)) {
			filteredLines = append(filteredLines, line)
		}
	}
	
	if len(filteredLines) == 0 {
		m.filteredLogContent = "No logs match filter: " + filter
	} else {
		m.filteredLogContent = strings.Join(filteredLines, "\n")
	}
}

// View renders the UI
func (m *Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Create header with title
	header := titleStyle.Render("Dashboard Server Monitor")
	header = lipgloss.Place(m.width, 3, lipgloss.Center, lipgloss.Center, header)
	
	// Create status indicators
	var serverStatusIndicator, metricsStatusIndicator string
	
	if m.guiStatus.ServerRunning {
		serverStatusIndicator = statusOkStyle.Render("Server: Running")
	} else {
		serverStatusIndicator = statusErrorStyle.Render("Server: Stopped")
	}
	
	if m.guiStatus.MetricsEnabled {
		metricsStatusIndicator = statusOkStyle.Render("Metrics: Enabled")
	} else {
		metricsStatusIndicator = statusErrorStyle.Render("Metrics: Disabled")
	}
	
	// Create tab-specific indicators
	var tabSpecificInfo string
	if m.activeTab == homeTab {
		tabSpecificInfo = infoStyle.Render(" Static Data ")
	} else {
		tabSpecificInfo = infoStyle.Render(fmt.Sprintf("Last Update: %s", m.lastUpdateTime.Format("15:04:05")))
	}
	
	// Create scroll status
	var scrollStatusIndicator string
	if m.activeTab == appLogTab || m.activeTab == prometheusTab {
		if m.autoScroll {
			scrollStatusIndicator = infoStyle.Render("Auto-Scroll: ON")
		} else {
			scrollStatusIndicator = warningStyle.Render("Auto-Scroll: OFF")
		}
	}
	
	// Join status elements
	statusBarElements := []string{serverStatusIndicator, "  ", metricsStatusIndicator}
	if scrollStatusIndicator != "" {
		statusBarElements = append(statusBarElements, "  ", scrollStatusIndicator)
	}
	statusBarElements = append(statusBarElements, "  ", tabSpecificInfo)
	
	statusBar := lipgloss.JoinHorizontal(lipgloss.Center, statusBarElements...)
	statusBar = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, statusBar)
	
	// Create tab bar
	var homeTabContent, appLogTabContent, prometheusTabContent string
	
	if m.activeTab == homeTab {
		homeTabContent = activeTabStyle.Render("1 Home")
		appLogTabContent = inactiveTabStyle.Render("2 App Logs")
		prometheusTabContent = inactiveTabStyle.Render("3 Prometheus")
	} else if m.activeTab == appLogTab {
		homeTabContent = inactiveTabStyle.Render("1 Home")
		appLogTabContent = activeTabStyle.Render("2 App Logs")
		prometheusTabContent = inactiveTabStyle.Render("3 Prometheus")
	} else {
		homeTabContent = inactiveTabStyle.Render("1 Home")
		appLogTabContent = inactiveTabStyle.Render("2 App Logs")
		prometheusTabContent = activeTabStyle.Render("3 Prometheus")
	}
	
	tabBar := lipgloss.JoinHorizontal(lipgloss.Bottom, homeTabContent, appLogTabContent, prometheusTabContent)
	tabBar = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, tabBar)
	
	// Create visually distinct header container
	headerContainer := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		BorderBottom(true).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		Width(m.width).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"", // Spacing
			statusBar,
			"", // Spacing
			tabBar,
		))
	
	// Render content based on active tab
	var contentArea string
	switch m.activeTab {
	case homeTab:
		contentArea = m.renderHomeTab()
	case appLogTab:
		contentArea = m.renderAppLogsTab()
	case prometheusTab:
		contentArea = m.renderPrometheusTab()
	}
	
	// Add padding for visual separation
	contentArea = lipgloss.NewStyle().PaddingTop(1).Render(contentArea)
	
	// Create help footer
	var helpText string
	switch m.activeTab {
	case homeTab:
		helpText = helpStyle.Render("TAB/1/2/3: switch views • Q: quit • R: refresh")
	case appLogTab:
		helpText = helpStyle.Render("TAB/1/2/3: switch views • Q: quit • R: refresh • G: bottom • g: top • S: toggle auto-scroll • F: focus filter")
	case prometheusTab:
		helpText = helpStyle.Render("TAB/1/2/3: switch views • Q: quit • R: refresh • G: bottom • g: top • S: toggle auto-scroll")
	}
	
	footer := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, helpText)
	
	// Add visual separation to footer
	footer = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		BorderTop(true).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Width(m.width).
		PaddingTop(1).
		Render(footer)
	
	// Combine all sections
	return appStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			headerContainer,
			contentArea,
			footer,
		),
	)
}

// renderHomeTab renders the home tab with server metadata
func (m *Model) renderHomeTab() string {
	// Use cached version if available
	if m.homeTabRendered && m.cachedHomeTab != "" {
		return m.cachedHomeTab
	}
	
	metadata := m.serverMetadata
	
	// Server Info Section
	serverInfoTitle := subtitleStyle.Render("Server Information")
	
	serverInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, 
			metadataLabelStyle.Render("Server URL:"), 
			metadataValueStyle.Render(metadata.ServerURL)),
		lipgloss.JoinHorizontal(lipgloss.Left, 
			metadataLabelStyle.Render("Version:"), 
			metadataValueStyle.Render(metadata.ServerVersion)),
		lipgloss.JoinHorizontal(lipgloss.Left, 
			metadataLabelStyle.Render("Environment:"), 
			metadataValueStyle.Render(metadata.Environment)),
		lipgloss.JoinHorizontal(lipgloss.Left, 
			metadataLabelStyle.Render("Start Time:"), 
			metadataValueStyle.Render(renderMetadataValue(metadata.StartTime))),
		lipgloss.JoinHorizontal(lipgloss.Left, 
			metadataLabelStyle.Render("Data Directory:"), 
			metadataValueStyle.Render(metadata.DataDirectory)),
	)
	
	serverInfoBox := metadataBoxStyle.Render(serverInfo)
	
	// Environment Variables Section
	envVarsTitle := subtitleStyle.Render("Environment Variables")
	
	var envVarsList string
	if len(metadata.EnvVariables) == 0 {
		envVarsList = "No environment variables set"
	} else {
		// Sort keys to ensure consistent display order
		var keys []string
		for key := range metadata.EnvVariables {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		
		var envPairs []string
		for _, key := range keys {
			value := metadata.EnvVariables[key]
			envPairs = append(envPairs, fmt.Sprintf("%s: %s", 
				metadataLabelStyle.Render(key), 
				metadataValueStyle.Render(util.AnyToString(value))))
		}
		envVarsList = lipgloss.JoinVertical(lipgloss.Left, envPairs...)
	}
	
	envVarsBox := metadataBoxStyle.Render(envVarsList)
	
	// Cron Jobs Section
	cronJobsTitle := subtitleStyle.Render("Cron Jobs")
	
	var cronJobsList string
	if len(metadata.CronJobs) == 0 {
		cronJobsList = "No cron jobs configured"
	} else {
		// Sort cron jobs by name for consistent order
		sort.Slice(metadata.CronJobs, func(i, j int) bool {
			return metadata.CronJobs[i].Name < metadata.CronJobs[j].Name
		})
		
		var jobEntries []string
		for _, job := range metadata.CronJobs {
			jobInfo := lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.JoinHorizontal(lipgloss.Left, 
					metadataLabelStyle.Render("Name:"), 
					metadataValueStyle.Render(job.Name)),
				lipgloss.JoinHorizontal(lipgloss.Left, 
					metadataLabelStyle.Render("Schedule:"), 
					metadataValueStyle.Render(job.Schedule)),
				lipgloss.JoinHorizontal(lipgloss.Left, 
					metadataLabelStyle.Render("Last Run:"), 
					metadataValueStyle.Render(renderMetadataValue(job.LastRun))),
				lipgloss.JoinHorizontal(lipgloss.Left, 
					metadataLabelStyle.Render("Status:"), 
					metadataValueStyle.Render(job.Status)),
			)
			jobEntries = append(jobEntries, jobInfo)
		}
		cronJobsList = lipgloss.JoinVertical(lipgloss.Left, jobEntries...)
	}
	
	cronJobsBox := metadataBoxStyle.Render(cronJobsList)
	
	// API Endpoints Section
	apiEndpointsTitle := subtitleStyle.Render("API Endpoints")
	
	var endpointsList string
	if len(metadata.APIEndpoints) == 0 {
		endpointsList = "No API endpoints available"
	} else {
		// Sort endpoints for consistent display
		sort.Strings(metadata.APIEndpoints)
		
		var endpoints []string
		for _, endpoint := range metadata.APIEndpoints {
			endpoints = append(endpoints, metadataValueStyle.Render("• "+endpoint))
		}
		endpointsList = lipgloss.JoinVertical(lipgloss.Left, endpoints...)
	}
	
	apiEndpointsBox := metadataBoxStyle.Render(endpointsList)
	
	// Combine all sections - simple vertical layout
	renderedContent := lipgloss.JoinVertical(
		lipgloss.Left,
		serverInfoTitle,
		serverInfoBox,
		"",
		envVarsTitle,
		envVarsBox,
		"",
		cronJobsTitle,
		cronJobsBox,
		"",
		apiEndpointsTitle,
		apiEndpointsBox,
	)
	
	// Cache the rendered content for future use
	m.cachedHomeTab = renderedContent
	m.homeTabRendered = true
	
	return renderedContent
}

// renderAppLogsTab renders the app logs tab with filter
func (m *Model) renderAppLogsTab() string {
	// Create the filter input area
	var filterStatus string
	if m.logFilterInput.Focused() {
		filterStatus = infoStyle.Render(" Editing Filter ")
	} else {
		filterStatus = helpStyle.Render(" Press F to focus filter ")
	}
	
	// Add clear indication of filter status and current filter value
	if m.logFilterInput.Value() != "" {
		filterDesc := fmt.Sprintf(" Filtering by: \"%s\" ", m.logFilterInput.Value())
		filterStatus = infoStyle.Render(filterDesc)
	}
	
	// Style the input box with a border to make it more visible
	inputBox := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Padding(0, 1).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			filterPromptStyle.Render(m.logFilterInput.Prompt),
			m.logFilterInput.View(),
		))
	
	filterSection := lipgloss.JoinVertical(
		lipgloss.Left,
		filterStatus,
		inputBox,
	)
	
	// Add some padding around the filter section
	filterSection = lipgloss.NewStyle().
		Padding(0, 2).
		Render(filterSection)
	
	// Add the logs content
	contentTitle := subtitleStyle.Render("Application Logs")
	contentView := m.appLogViewport.View()
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		filterSection,
		contentTitle,
		contentView,
	)
}

// renderPrometheusTab renders the prometheus tab
func (m *Model) renderPrometheusTab() string {
	contentTitle := subtitleStyle.Render("Prometheus Metrics")
	contentView := m.prometheusViewport.View()
	
	return lipgloss.JoinVertical(
		lipgloss.Left,
		contentTitle,
		contentView,
	)
}

// ------------ Commands ------------

// tickEvery returns a command that ticks at regular intervals
func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// readLogFile reads the log file and returns its content
func readLogFile(path string) tea.Cmd {
	return func() tea.Msg {
		content, err := os.ReadFile(path)
		if err != nil {
			return logContentMsg(fmt.Sprintf("Error reading log file: %v", err))
		}
		return logContentMsg(content)
	}
}

// readPrometheusLogs reads the Prometheus metrics logs
func readPrometheusLogs() tea.Cmd {
	return func() tea.Msg {
		// Try to fetch metrics data from the metrics endpoint
		// First, check if the curl command is available
		curlCmd := exec.Command("curl", "-s", "http://localhost:9091/metrics")
		metricsOutput, err := curlCmd.Output()
		
		if err != nil {
			// Fall back to a default message if we can't fetch metrics
			return prometheusLogMsg("Metrics server is running, but couldn't fetch metrics data.\n" +
				"Try accessing http://localhost:9091/metrics in your browser.")
		}
		
		// Process the metrics output to make it more readable
		metrics := string(metricsOutput)
		if metrics == "" {
			metrics = "No metrics data available."
		}
		
		// Add a header to the metrics data
		return prometheusLogMsg("Prometheus Metrics\n" +
			"===================\n\n" + 
			metrics)
	}
}

// SafeGotoTop is a safe wrapper for GotoTop that checks content length
func SafeGotoTop(vp *viewport.Model) {
	// Only perform operation if viewport is properly initialized
	if vp != nil && vp.Height > 0 && vp.Width > 0 && len(vp.View()) > 0 {
		vp.GotoTop()
	}
}

// SafeGotoBottom is a safe wrapper for GotoBottom that checks content length
func SafeGotoBottom(vp *viewport.Model) {
	// Only perform operation if viewport is properly initialized
	if vp != nil && vp.Height > 0 && vp.Width > 0 && len(vp.View()) > 0 {
		vp.GotoBottom()
	}
}

// renderMetadataValue formats a metadata value based on its type
func renderMetadataValue(value interface{}) string {
	switch v := value.(type) {
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	case []string:
		if len(v) == 0 {
			return "None"
		}
		return strings.Join(v, ", ")
	case map[string]string:
		if len(v) == 0 {
			return "None"
		}
		var pairs []string
		for key, val := range v {
			pairs = append(pairs, fmt.Sprintf("%s: %s", key, val))
		}
		return strings.Join(pairs, "\n")
	default:
		return fmt.Sprintf("%v", v)
	}
}