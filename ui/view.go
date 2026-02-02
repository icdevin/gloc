package ui

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/devin/gloc/colors"
)

// View implements tea.Model
func (m Model) View() string {
	if m.Err != nil {
		return AppStyle.Render(fmt.Sprintf("Error: %v\n\nPress q to quit.", m.Err))
	}

	if m.Result == nil {
		return AppStyle.Render("Loading...")
	}

	var b strings.Builder

	if m.Mode == LanguageView {
		m.renderLanguageView(&b)
	} else {
		m.renderFileView(&b)
	}

	// Status bar
	m.renderStatusBar(&b)

	// Help
	m.renderHelp(&b)

	return AppStyle.Width(m.Width).Render(b.String())
}

func (m Model) renderLanguageView(b *strings.Builder) {
	// Title
	title := TitleStyle.Render(fmt.Sprintf(" 📊 gloc - %s ", m.TargetPath))
	b.WriteString(title)
	b.WriteString("\n\n")

	// Build table data
	visibleRows := m.VisibleRows()
	endIdx := m.ScrollOffset + visibleRows
	if endIdx > len(m.Result.Languages) {
		endIdx = len(m.Result.Languages)
	}

	var rows [][]string
	for i := m.ScrollOffset; i < endIdx; i++ {
		lang := m.Result.Languages[i]
		total := lang.Code + lang.Comment + lang.Blank

		cursor := "  "
		if i == m.Cursor {
			cursor = CursorStyle.Render("▶ ")
		}

		// Color dot
		color := colors.GetColor(lang.Name)
		colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		dot := colorStyle.Render("●")

		rows = append(rows, []string{
			cursor + dot + " " + lang.Name,
			strconv.Itoa(lang.Files),
			strconv.Itoa(lang.Blank),
			strconv.Itoa(lang.Comment),
			strconv.Itoa(lang.Code),
			strconv.Itoa(total),
		})
	}

	// Create table
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers(m.languageHeaders()...).
		Rows(rows...).
		Width(m.ContentWidth()).
		StyleFunc(func(row, col int) lipgloss.Style {
			// Header row
			if row == table.HeaderRow {
				if col == m.sortColIndex(m.SortCol) {
					return HeaderActiveStyle.Align(lipgloss.Center)
				}
				return HeaderStyle.Align(lipgloss.Center)
			}

			// Data rows - first column left-aligned, rest right-aligned
			if col == 0 {
				return lipgloss.NewStyle()
			}

			// Right-align and color numeric columns
			switch col {
			case 1:
				return FilesStyle.Align(lipgloss.Right)
			case 2:
				return BlankStyle.Align(lipgloss.Right)
			case 3:
				return CommentStyle.Align(lipgloss.Right)
			case 4:
				return CodeStyle.Align(lipgloss.Right)
			case 5:
				return TotalStyle.Align(lipgloss.Right)
			default:
				return lipgloss.NewStyle().Align(lipgloss.Right)
			}
		})

	b.WriteString(t.Render())
	b.WriteString("\n")

	// Pad with empty lines if needed
	for i := endIdx - m.ScrollOffset; i < visibleRows; i++ {
		b.WriteString("\n")
	}
}

func (m Model) renderFileView(b *strings.Builder) {
	// Title with language color
	langColor := colors.GetColor(m.SelectedLang)
	titleBg := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(langColor)).
		Padding(0, 1).
		Bold(true)

	title := titleBg.Render(fmt.Sprintf(" 📁 %s Files ", m.SelectedLang))
	b.WriteString(title)
	b.WriteString("\n\n")

	// Build table data
	files := m.SortFiles(m.SelectedLang)
	visibleRows := m.VisibleRows()
	endIdx := m.FileScrollOffset + visibleRows
	if endIdx > len(files) {
		endIdx = len(files)
	}

	var rows [][]string
	for i := m.FileScrollOffset; i < endIdx; i++ {
		file := files[i]
		total := file.Code + file.Comment + file.Blank

		cursor := "  "
		if i == m.FileCursor {
			cursor = CursorStyle.Render("▶ ")
		}

		// Get relative path
		displayPath := file.Path
		if rel, err := filepath.Rel(m.TargetPath, file.Path); err == nil {
			displayPath = rel
		}

		// Truncate if too long
		maxPathLen := 60
		if len(displayPath) > maxPathLen {
			displayPath = "…" + displayPath[len(displayPath)-maxPathLen+1:]
		}

		rows = append(rows, []string{
			cursor + displayPath,
			strconv.Itoa(file.Blank),
			strconv.Itoa(file.Comment),
			strconv.Itoa(file.Code),
			strconv.Itoa(total),
		})
	}

	// Create table
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers(m.fileHeaders()...).
		Rows(rows...).
		Width(m.ContentWidth()).
		StyleFunc(func(row, col int) lipgloss.Style {
			// Header row
			if row == table.HeaderRow {
				if col == m.fileSortColIndex(m.FileSortCol) {
					return HeaderActiveStyle.Align(lipgloss.Center)
				}
				return HeaderStyle.Align(lipgloss.Center)
			}

			// Data rows - first column left-aligned, rest right-aligned
			if col == 0 {
				return lipgloss.NewStyle()
			}

			// Right-align and color numeric columns
			switch col {
			case 1:
				return BlankStyle.Align(lipgloss.Right)
			case 2:
				return CommentStyle.Align(lipgloss.Right)
			case 3:
				return CodeStyle.Align(lipgloss.Right)
			case 4:
				return TotalStyle.Align(lipgloss.Right)
			default:
				return lipgloss.NewStyle().Align(lipgloss.Right)
			}
		})

	b.WriteString(t.Render())
	b.WriteString("\n")

	// Pad with empty lines if needed
	for i := endIdx - m.FileScrollOffset; i < visibleRows; i++ {
		b.WriteString("\n")
	}
}

func (m Model) renderStatusBar(b *strings.Builder) {
	total := m.Result.Total
	totalLines := total.Code + total.Comment + total.Blank

	statusContent := fmt.Sprintf(
		"Total: %s files │ %s blank │ %s comment │ %s code │ %s lines",
		FilesStyle.Render(strconv.Itoa(total.Files)),
		BlankStyle.Render(strconv.Itoa(total.Blank)),
		CommentStyle.Render(strconv.Itoa(total.Comment)),
		CodeStyle.Render(strconv.Itoa(total.Code)),
		TotalStyle.Render(strconv.Itoa(totalLines)),
	)
	b.WriteString("\n")
	b.WriteString(StatusBarStyle.Render(statusContent))
	b.WriteString("\n")
}

func (m Model) renderHelp(b *strings.Builder) {
	var help string
	if m.Mode == LanguageView {
		help = fmt.Sprintf(
			"%s navigate • %s view files • %s sort • %s quit",
			HelpKeyStyle.Render("↑/↓"),
			HelpKeyStyle.Render("enter"),
			HelpKeyStyle.Render("1-6"),
			HelpKeyStyle.Render("q"),
		)
	} else {
		help = fmt.Sprintf(
			"%s navigate • %s sort • %s back • %s quit",
			HelpKeyStyle.Render("↑/↓"),
			HelpKeyStyle.Render("1,3-6"),
			HelpKeyStyle.Render("esc/q"),
			HelpKeyStyle.Render("ctrl+c"),
		)
	}
	b.WriteString(HelpStyle.Render(help))
}

func (m Model) languageHeaders() []string {
	return []string{
		m.sortHeader("[1] Language", SortByName, m.SortCol, m.SortAsc),
		m.sortHeader("[2] Files", SortByFiles, m.SortCol, m.SortAsc),
		m.sortHeader("[3] Blank", SortByBlank, m.SortCol, m.SortAsc),
		m.sortHeader("[4] Comment", SortByComment, m.SortCol, m.SortAsc),
		m.sortHeader("[5] Code", SortByCode, m.SortCol, m.SortAsc),
		m.sortHeader("[6] Total", SortByTotal, m.SortCol, m.SortAsc),
	}
}

func (m Model) fileHeaders() []string {
	return []string{
		m.sortHeader("[1] File", SortByName, m.FileSortCol, m.FileSortAsc),
		m.sortHeader("[3] Blank", SortByBlank, m.FileSortCol, m.FileSortAsc),
		m.sortHeader("[4] Comment", SortByComment, m.FileSortCol, m.FileSortAsc),
		m.sortHeader("[5] Code", SortByCode, m.FileSortCol, m.FileSortAsc),
		m.sortHeader("[6] Total", SortByTotal, m.FileSortCol, m.FileSortAsc),
	}
}

func (m Model) sortHeader(label string, col SortColumn, currentSort SortColumn, asc bool) string {
	if col == currentSort {
		if asc {
			return label + " ▲"
		}
		return label + " ▼"
	}
	return label
}

// sortColIndex returns the column index for highlighting the active sort column
func (m Model) sortColIndex(col SortColumn) int {
	switch col {
	case SortByName:
		return 0
	case SortByFiles:
		return 1
	case SortByBlank:
		return 2
	case SortByComment:
		return 3
	case SortByCode:
		return 4
	case SortByTotal:
		return 5
	default:
		return -1
	}
}

// fileSortColIndex returns the column index for file view
func (m Model) fileSortColIndex(col SortColumn) int {
	switch col {
	case SortByName:
		return 0
	case SortByBlank:
		return 1
	case SortByComment:
		return 2
	case SortByCode:
		return 3
	case SortByTotal:
		return 4
	default:
		return -1
	}
}
