package cloc

import (
	"encoding/json"
	"os/exec"
	"sort"
)

// FileInfo contains line count information for a single file
type FileInfo struct {
	Path     string `json:"-"`
	Blank    int    `json:"blank"`
	Comment  int    `json:"comment"`
	Code     int    `json:"code"`
	Language string `json:"language"`
}

// LanguageStats contains aggregate statistics for a language
type LanguageStats struct {
	Name    string
	Files   int `json:"nFiles"`
	Blank   int `json:"blank"`
	Comment int `json:"comment"`
	Code    int `json:"code"`
}

// Result contains the complete cloc analysis result
type Result struct {
	Languages []LanguageStats
	Files     map[string][]FileInfo // Files grouped by language
	Total     LanguageStats
}

// Run executes cloc on the given path and returns parsed results
func Run(path string) (*Result, error) {
	// First, get summary by language
	summaryResult, err := runClocSummary(path)
	if err != nil {
		return nil, err
	}

	// Then, get file-level details
	fileResult, err := runClocByFile(path)
	if err != nil {
		return nil, err
	}

	// Merge results
	summaryResult.Files = fileResult.Files

	return summaryResult, nil
}

func runClocSummary(path string) (*Result, error) {
	cmd := exec.Command("cloc", "--json", path)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var rawResult map[string]json.RawMessage
	if err := json.Unmarshal(output, &rawResult); err != nil {
		return nil, err
	}

	result := &Result{
		Languages: []LanguageStats{},
		Files:     make(map[string][]FileInfo),
	}

	for key, value := range rawResult {
		if key == "header" {
			continue
		}

		var stats LanguageStats
		if err := json.Unmarshal(value, &stats); err != nil {
			continue
		}
		stats.Name = key

		if key == "SUM" {
			result.Total = stats
		} else {
			result.Languages = append(result.Languages, stats)
		}
	}

	// Sort languages by code lines (descending)
	sort.Slice(result.Languages, func(i, j int) bool {
		return result.Languages[i].Code > result.Languages[j].Code
	})

	return result, nil
}

func runClocByFile(path string) (*Result, error) {
	cmd := exec.Command("cloc", "--json", "--by-file", path)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var rawResult map[string]json.RawMessage
	if err := json.Unmarshal(output, &rawResult); err != nil {
		return nil, err
	}

	result := &Result{
		Files: make(map[string][]FileInfo),
	}

	for key, value := range rawResult {
		if key == "header" || key == "SUM" {
			continue
		}

		var fileInfo FileInfo
		if err := json.Unmarshal(value, &fileInfo); err != nil {
			continue
		}
		fileInfo.Path = key

		// Group files by language
		result.Files[fileInfo.Language] = append(result.Files[fileInfo.Language], fileInfo)
	}

	// Sort files within each language by code lines (descending)
	for lang := range result.Files {
		sort.Slice(result.Files[lang], func(i, j int) bool {
			return result.Files[lang][i].Code > result.Files[lang][j].Code
		})
	}

	return result, nil
}
