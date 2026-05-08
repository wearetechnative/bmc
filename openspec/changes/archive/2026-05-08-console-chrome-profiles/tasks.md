## 1. Config

- [x] 1.1 Add `ChromeProfiles bool` and `ChromeBinary string` fields to `ConsoleConfig` in `internal/config/config.go`
- [x] 1.2 Set default value for `ChromeBinary` to `"google-chrome"` when empty

## 2. Chrome Profile Logic

- [x] 2.1 Add `sanitizeProfileName(name string) string` helper — replace `/`, `:` and other invalid chars with `-`
- [x] 2.2 Add `chromeProfileDir(profileName string) string` helper — returns `~/.config/bmc/chrome/profiles/<sanitized-name>`
- [x] 2.3 Add `defaultChromeProfilePath() string` helper — returns OS-specific default Chrome profile path (Linux: `~/.config/google-chrome/Default/`, macOS: `~/Library/Application Support/Google/Chrome/Default/`); returns empty string if not found
- [x] 2.4 Add `seedChromeProfile(destDir, srcDir string) error` — copies `Extensions/`, `Local Extension Settings/`, `Preferences` from srcDir to destDir; skips missing files/dirs silently
- [x] 2.5 Add `openChromeProfile(url, profileName string, cfg config.ConsoleConfig) error` — orchestrates: resolve dir, seed if new, launch binary with `--user-data-dir`, `--no-first-run`, `--no-default-browser-check`, and URL

## 3. Browser Open Integration

- [x] 3.1 Extend `openBrowser()` in `internal/awsops/console.go` to call `openChromeProfile()` when `bmcCfg.Console.ChromeProfiles` is true

## 4. Documentation

- [x] 4.1 Add `[console]` Chrome options to README under an "Experimental" section, with example config and usage notes for `chrome_profiles` and `chrome_binary`

## 5. Commit

- [x] 5.1 Commit all changes with a descriptive message
