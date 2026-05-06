## 1. TUI list height fix

- [x] 1.1 In `internal/ui/list.go` `Choose()`: for `len(items) <= 4` set height to `len(items) + 1` and call `l.SetShowHelp(false)`
- [x] 1.2 For `len(items) > 4` change height formula from `len(items) + 6` to `min(len(items) + 3, 30)`
- [x] 1.3 Verify: run `bmc ec2connect` and confirm no blank rows in "Connection method" picker (2 items) and "Select SSH user" picker (4 items)

## 2. profsel wrapper hint

- [x] 2.1 In `cmd/profsel.go` after `fmt.Printf("export AWS_PROFILE=...")`: check `term.IsTerminal(int(os.Stdout.Fd()))` and if true, print hint to stderr
- [x] 2.2 Verify: `bmc profsel` directly shows hint on stderr; `eval "$(bmc profsel)"` shows no hint and sets AWS_PROFILE

## 3. install-shell-integration permission denied fallback

- [x] 3.1 In `cmd/install_shell.go`: add a `const` block with manual snippets for home-manager zsh, home-manager bash, manual zsh/bash, and Fish
- [x] 3.2 When `os.OpenFile` returns a permission error, print the explanation and snippets instead of returning the error
- [x] 3.3 Verify: run `bmc install-shell-integration` with `~/.zshrc` made read-only; confirm helpful output appears

## 5. TUI list height: terminal-aware sizing

- [x] 5.1 Add `wantHeight` field to `listModel` to store desired height before terminal clamp
- [x] 5.2 In `WindowSizeMsg` handler: also call `SetHeight(min(wantHeight, msg.Height-2))`
- [x] 5.3 In `Choose()`: open `/dev/tty` once, read terminal height via `term.GetSize`, clamp initial height, pass tty to bubbletea
- [x] 5.4 Use `SetShowPagination(false)` for ≤4 items to suppress spurious pagination dots

## 6. Documentation

- [x] 6.1 Add NixOS shell integration section to `README.md` under Setup with home-manager and Fish snippets
- [x] 6.2 Update CHANGELOG under `## NEXT VERSION`
