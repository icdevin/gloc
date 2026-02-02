package ui

import "github.com/charmbracelet/lipgloss"

// Styles for the UI
var (
	AppStyle = lipgloss.NewStyle().Padding(1, 2)

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5C5C5C")).
			Padding(0, 1)

	HeaderActiveStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#7DC4E4")).
				Padding(0, 1)

	// Title style
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7DC4E4")).
			Bold(true)

	// Row styles
	SelectedRowStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF"))

	NormalRowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC"))

	// Cursor style
	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7DC4E4")).
			Bold(true)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	HelpKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7DC4E4")).
			Bold(true)

	// Divider
	DividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4C4C4C"))

	// Border style for tables
	BorderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4C4C4C"))

	// Number styles for different columns
	CodeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6E3A1"))

	CommentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#89B4FA"))

	BlankStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9399B2"))

	FilesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9E2AF"))

	TotalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CBA6F7"))
)
