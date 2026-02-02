package ui

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/devin/gloc/cloc"
)

// Model is the main application model
type Model struct {
	Result           *cloc.Result
	Mode             ViewMode
	SelectedLang     string
	Cursor           int
	FileCursor       int
	Width            int
	Height           int
	TargetPath       string
	Err              error
	SortCol          SortColumn
	SortAsc          bool
	FileSortCol      SortColumn
	FileSortAsc      bool
	ScrollOffset     int
	FileScrollOffset int
	// Dynamic column widths
	ColLanguage int
	ColFiles    int
	ColBlank    int
	ColComment  int
	ColCode     int
	ColFilePath int
}

// NewModel creates a new model with the given path
func NewModel(path string) Model {
	return Model{
		TargetPath:  path,
		Mode:        LanguageView,
		SortCol:     SortByCode,
		SortAsc:     false, // descending by default
		FileSortCol: SortByCode,
		FileSortAsc: false,
	}
}

// ClocResultMsg is the message returned when cloc finishes
type ClocResultMsg struct {
	Result *cloc.Result
	Err    error
}

// RunCloc executes cloc and returns a command
func RunCloc(path string) tea.Cmd {
	return func() tea.Msg {
		result, err := cloc.Run(path)
		return ClocResultMsg{Result: result, Err: err}
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return RunCloc(m.TargetPath)
}

// ContentWidth returns the usable width for content (accounting for AppStyle padding)
func (m *Model) ContentWidth() int {
	// AppStyle has Padding(1, 2) = 2 chars on each side = 4 total
	return m.Width - 4
}

// CalculateColumnWidths calculates dynamic column widths based on terminal width
func (m *Model) CalculateColumnWidths() {
	availableWidth := m.ContentWidth()

	// Language view: distribute space proportionally
	fixedCols := MinColFiles + MinColBlank + MinColComment + MinColCode
	extraSpace := availableWidth - fixedCols - MinColLanguage - 3 // 3 for dot and spacing

	if extraSpace > 0 {
		// Give extra space mostly to language name, some to numeric columns
		m.ColLanguage = MinColLanguage + (extraSpace * 50 / 100)
		remaining := extraSpace - (extraSpace * 50 / 100)
		m.ColFiles = MinColFiles + (remaining * 20 / 100)
		m.ColBlank = MinColBlank + (remaining * 20 / 100)
		m.ColComment = MinColComment + (remaining * 30 / 100)
		m.ColCode = MinColCode + (remaining * 30 / 100)
	} else {
		m.ColLanguage = MinColLanguage
		m.ColFiles = MinColFiles
		m.ColBlank = MinColBlank
		m.ColComment = MinColComment
		m.ColCode = MinColCode
	}

	// File view: give most space to file path
	// 10 * 4 = 40 for numeric columns (blank, comment, code, total)
	m.ColFilePath = availableWidth - 44
	if m.ColFilePath < 40 {
		m.ColFilePath = 40
	}
}

// SortLanguages sorts the languages based on current sort settings
func (m *Model) SortLanguages() {
	if m.Result == nil {
		return
	}

	sort.Slice(m.Result.Languages, func(i, j int) bool {
		var less bool
		switch m.SortCol {
		case SortByName:
			less = strings.ToLower(m.Result.Languages[i].Name) < strings.ToLower(m.Result.Languages[j].Name)
		case SortByFiles:
			less = m.Result.Languages[i].Files < m.Result.Languages[j].Files
		case SortByBlank:
			less = m.Result.Languages[i].Blank < m.Result.Languages[j].Blank
		case SortByComment:
			less = m.Result.Languages[i].Comment < m.Result.Languages[j].Comment
		case SortByCode:
			less = m.Result.Languages[i].Code < m.Result.Languages[j].Code
		case SortByTotal:
			totalI := m.Result.Languages[i].Code + m.Result.Languages[i].Comment + m.Result.Languages[i].Blank
			totalJ := m.Result.Languages[j].Code + m.Result.Languages[j].Comment + m.Result.Languages[j].Blank
			less = totalI < totalJ
		}
		if m.SortAsc {
			return less
		}
		return !less
	})
}

// SortFiles returns sorted files for the given language
func (m *Model) SortFiles(lang string) []cloc.FileInfo {
	files := make([]cloc.FileInfo, len(m.Result.Files[lang]))
	copy(files, m.Result.Files[lang])

	sort.Slice(files, func(i, j int) bool {
		var less bool
		switch m.FileSortCol {
		case SortByName:
			less = strings.ToLower(files[i].Path) < strings.ToLower(files[j].Path)
		case SortByBlank:
			less = files[i].Blank < files[j].Blank
		case SortByComment:
			less = files[i].Comment < files[j].Comment
		case SortByCode:
			less = files[i].Code < files[j].Code
		case SortByTotal:
			totalI := files[i].Code + files[i].Comment + files[i].Blank
			totalJ := files[j].Code + files[j].Comment + files[j].Blank
			less = totalI < totalJ
		default:
			less = files[i].Code < files[j].Code
		}
		if m.FileSortAsc {
			return less
		}
		return !less
	})
	return files
}

// VisibleRows returns the number of visible rows based on terminal height
func (m *Model) VisibleRows() int {
	rows := m.Height - 12
	if rows < 1 {
		return 10
	}
	return rows
}
