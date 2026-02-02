package ui

// Minimum column widths
const (
	MinColLanguage = 20
	MinColFiles    = 8
	MinColBlank    = 8
	MinColComment  = 10
	MinColCode     = 10
)

// SortColumn represents which column to sort by
type SortColumn int

const (
	SortByCode SortColumn = iota
	SortByFiles
	SortByBlank
	SortByComment
	SortByName
	SortByTotal
)

// ViewMode represents the current view
type ViewMode int

const (
	LanguageView ViewMode = iota
	FileView
)
