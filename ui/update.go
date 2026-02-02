package ui

import tea "github.com/charmbracelet/bubbletea"

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.CalculateColumnWidths()

	case ClocResultMsg:
		if msg.Err != nil {
			m.Err = msg.Err
			return m, nil
		}
		m.Result = msg.Result
		m.SortLanguages()
		m.CalculateColumnWidths()
		return m, nil
	}

	return m, nil
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "q":
		if m.Mode == FileView {
			m.Mode = LanguageView
			return m, nil
		}
		return m, tea.Quit
	case "esc":
		if m.Mode == FileView {
			m.Mode = LanguageView
			return m, nil
		}
	case "enter":
		if m.Mode == LanguageView && m.Result != nil && len(m.Result.Languages) > 0 {
			m.SelectedLang = m.Result.Languages[m.Cursor].Name
			m.Mode = FileView
			m.FileCursor = 0
			m.FileScrollOffset = 0
		}
	case "up", "k":
		m.handleUp()
	case "down", "j":
		m.handleDown()
	case "1":
		m.handleSortByName()
	case "2":
		m.handleSortByFiles()
	case "3":
		m.handleSortByBlank()
	case "4":
		m.handleSortByComment()
	case "5":
		m.handleSortByCode()
	case "6":
		m.handleSortByTotal()
	case "home", "g":
		m.handleHome()
	case "end", "G":
		m.handleEnd()
	}

	return m, nil
}

func (m *Model) handleUp() {
	if m.Mode == LanguageView {
		if m.Cursor > 0 {
			m.Cursor--
			if m.Cursor < m.ScrollOffset {
				m.ScrollOffset = m.Cursor
			}
		}
	} else {
		if m.FileCursor > 0 {
			m.FileCursor--
			if m.FileCursor < m.FileScrollOffset {
				m.FileScrollOffset = m.FileCursor
			}
		}
	}
}

func (m *Model) handleDown() {
	if m.Mode == LanguageView && m.Result != nil {
		if m.Cursor < len(m.Result.Languages)-1 {
			m.Cursor++
			visibleRows := m.VisibleRows()
			if m.Cursor >= m.ScrollOffset+visibleRows {
				m.ScrollOffset = m.Cursor - visibleRows + 1
			}
		}
	} else if m.Result != nil {
		files := m.Result.Files[m.SelectedLang]
		if m.FileCursor < len(files)-1 {
			m.FileCursor++
			visibleRows := m.VisibleRows()
			if m.FileCursor >= m.FileScrollOffset+visibleRows {
				m.FileScrollOffset = m.FileCursor - visibleRows + 1
			}
		}
	}
}

func (m *Model) handleSortByName() {
	if m.Mode == LanguageView {
		if m.SortCol == SortByName {
			m.SortAsc = !m.SortAsc
		} else {
			m.SortCol = SortByName
			m.SortAsc = true
		}
		m.SortLanguages()
		m.Cursor = 0
		m.ScrollOffset = 0
	} else {
		if m.FileSortCol == SortByName {
			m.FileSortAsc = !m.FileSortAsc
		} else {
			m.FileSortCol = SortByName
			m.FileSortAsc = true
		}
		m.FileCursor = 0
		m.FileScrollOffset = 0
	}
}

func (m *Model) handleSortByFiles() {
	if m.Mode == LanguageView {
		if m.SortCol == SortByFiles {
			m.SortAsc = !m.SortAsc
		} else {
			m.SortCol = SortByFiles
			m.SortAsc = false
		}
		m.SortLanguages()
		m.Cursor = 0
		m.ScrollOffset = 0
	}
}

func (m *Model) handleSortByBlank() {
	if m.Mode == LanguageView {
		if m.SortCol == SortByBlank {
			m.SortAsc = !m.SortAsc
		} else {
			m.SortCol = SortByBlank
			m.SortAsc = false
		}
		m.SortLanguages()
		m.Cursor = 0
		m.ScrollOffset = 0
	} else {
		if m.FileSortCol == SortByBlank {
			m.FileSortAsc = !m.FileSortAsc
		} else {
			m.FileSortCol = SortByBlank
			m.FileSortAsc = false
		}
		m.FileCursor = 0
		m.FileScrollOffset = 0
	}
}

func (m *Model) handleSortByComment() {
	if m.Mode == LanguageView {
		if m.SortCol == SortByComment {
			m.SortAsc = !m.SortAsc
		} else {
			m.SortCol = SortByComment
			m.SortAsc = false
		}
		m.SortLanguages()
		m.Cursor = 0
		m.ScrollOffset = 0
	} else {
		if m.FileSortCol == SortByComment {
			m.FileSortAsc = !m.FileSortAsc
		} else {
			m.FileSortCol = SortByComment
			m.FileSortAsc = false
		}
		m.FileCursor = 0
		m.FileScrollOffset = 0
	}
}

func (m *Model) handleSortByCode() {
	if m.Mode == LanguageView {
		if m.SortCol == SortByCode {
			m.SortAsc = !m.SortAsc
		} else {
			m.SortCol = SortByCode
			m.SortAsc = false
		}
		m.SortLanguages()
		m.Cursor = 0
		m.ScrollOffset = 0
	} else {
		if m.FileSortCol == SortByCode {
			m.FileSortAsc = !m.FileSortAsc
		} else {
			m.FileSortCol = SortByCode
			m.FileSortAsc = false
		}
		m.FileCursor = 0
		m.FileScrollOffset = 0
	}
}

func (m *Model) handleSortByTotal() {
	if m.Mode == LanguageView {
		if m.SortCol == SortByTotal {
			m.SortAsc = !m.SortAsc
		} else {
			m.SortCol = SortByTotal
			m.SortAsc = false
		}
		m.SortLanguages()
		m.Cursor = 0
		m.ScrollOffset = 0
	} else {
		if m.FileSortCol == SortByTotal {
			m.FileSortAsc = !m.FileSortAsc
		} else {
			m.FileSortCol = SortByTotal
			m.FileSortAsc = false
		}
		m.FileCursor = 0
		m.FileScrollOffset = 0
	}
}

func (m *Model) handleHome() {
	if m.Mode == LanguageView {
		m.Cursor = 0
		m.ScrollOffset = 0
	} else {
		m.FileCursor = 0
		m.FileScrollOffset = 0
	}
}

func (m *Model) handleEnd() {
	if m.Mode == LanguageView && m.Result != nil {
		m.Cursor = len(m.Result.Languages) - 1
		visibleRows := m.VisibleRows()
		if m.Cursor >= visibleRows {
			m.ScrollOffset = m.Cursor - visibleRows + 1
		}
	} else if m.Result != nil {
		files := m.Result.Files[m.SelectedLang]
		m.FileCursor = len(files) - 1
		visibleRows := m.VisibleRows()
		if m.FileCursor >= visibleRows {
			m.FileScrollOffset = m.FileCursor - visibleRows + 1
		}
	}
}
