# Implementation Tasks

## 1. Research & Design
- [x] 1.1 Review existing command registration in `bmc` and `_bmclib.sh`
- [x] 1.2 Determine completion generation approach (static vs dynamic)
- [x] 1.3 Research bash completion best practices and conventions
- [x] 1.4 Research zsh completion best practices and conventions

## 2. Bash Completion Implementation
- [x] 2.1 Create bash completion script
- [x] 2.2 Implement subcommand completion
- [x] 2.3 Test bash completion with various shells and scenarios
- [x] 2.4 Add installation instructions for bash completion

## 3. Zsh Completion Implementation
- [x] 3.1 Create zsh completion script
- [x] 3.2 Implement subcommand completion
- [x] 3.3 Test zsh completion with various shells and scenarios
- [x] 3.4 Add installation instructions for zsh completion

## 4. Integration & Documentation
- [x] 4.1 Add `gencompletions` command to `bmc` (e.g., `bmc gencompletions bash|zsh`)
- [x] 4.2 Add brief installation suggestions to completion output
- [x] 4.3 Update README.md with completion setup instructions
- [x] 4.4 Update CHANGELOG.md with new feature
- [x] 4.5 Test end-to-end completion workflow

## 5. Validation
- [x] 5.1 Validate with `openspec validate add-shell-completion --strict --no-interactive`
- [x] 5.2 Test on multiple shell versions (bash 4+, zsh 5.8+)
- [x] 5.3 Verify completion works for all registered commands

## 6. Bug Fix - Zsh Completion Function Call
- [x] 6.1 Fix zsh completion script to use `_bmc` instead of `_bmc "$@"`
- [x] 6.2 Test that zsh completion works correctly with the fix
- [x] 6.3 Verify function is called by completion system, not at load time
