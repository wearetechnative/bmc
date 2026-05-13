## 1. Add Fish constants to install_shell.go

- [x] 1.1 Add `fishWrapper` const with the Fish function for `~/.config/fish/functions/bmc.fish` (including `# bmc shell integration` marker comment)
- [x] 1.2 Add `fishManualNixOS` const with the NixOS explanation and `programs.fish.functions` home-manager snippet

## 2. Add NixOS detection helper

- [x] 2.1 Add `isNixOS() bool` helper function that returns true if `/etc/nixos/` directory exists

## 3. Add Fish case to the shell switch

- [x] 3.1 Add `strings.HasSuffix(shell, "fish")` case in `runInstallShell`
- [x] 3.2 In the Fish case: if `isNixOS()` is true, print `fishManualNixOS` and return nil
- [x] 3.3 In the Fish case: determine target path `~/.config/fish/functions/bmc.fish`
- [x] 3.4 Check if target file already exists; if so, print "already installed" and return nil
- [x] 3.5 Call `os.MkdirAll` on `~/.config/fish/functions/` before writing
- [x] 3.6 Write `fishWrapper` to the target file and print success message

## 4. Verify build and existing behaviour

- [x] 4.1 Confirm `go build ./...` succeeds
- [x] 4.2 Confirm zsh/bash paths are unchanged (no regression)
