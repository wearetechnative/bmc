package history

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const maxHistory = 10

func dataDir() string {
	if d := os.Getenv("XDG_DATA_HOME"); d != "" {
		return filepath.Join(d, "bmc")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "bmc")
}

func historyPath(name string) string {
	return filepath.Join(dataDir(), name+"-history.json")
}

// Load returns up to maxHistory profile names for the given history key,
// most recent first. Returns an empty slice on any error.
func Load(name string) []string {
	data, err := os.ReadFile(historyPath(name))
	if err != nil {
		return nil
	}
	var entries []string
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil
	}
	return entries
}

// Save prepends entry to the history for name, deduplicates, caps at
// maxHistory, and writes atomically. Parent directories are created as needed.
func Save(name, entry string) error {
	existing := Load(name)

	// Prepend and deduplicate
	seen := make(map[string]bool)
	result := make([]string, 0, maxHistory)
	for _, e := range append([]string{entry}, existing...) {
		if !seen[e] {
			seen[e] = true
			result = append(result, e)
		}
		if len(result) == maxHistory {
			break
		}
	}

	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	dir := dataDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	tmp := historyPath(name) + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, historyPath(name))
}
