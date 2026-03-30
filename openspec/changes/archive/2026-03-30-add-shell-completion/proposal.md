# Change: Add Shell Completion for Bash and Zsh

## Why
Users should have tab-completion support when using the `bmc` command to improve usability and discoverability of commands and options. Shell completion helps prevent typos, speeds up command entry, and makes the CLI more professional and user-friendly.

## What Changes
- Add `bmc gencompletions bash|zsh` command to generate completion scripts
- Generate bash completion script with support for `bmc` subcommands
- Generate zsh completion script with support for `bmc` subcommands
- Output completion script to stdout along with brief installation instructions
- Generate completion suggestions dynamically based on registered commands
- Add documentation for enabling completions in user shells

## Impact
- Affected specs: new capability `shell-completion`
- Affected code:
  - New command: `gencompletions` in `bmc` script
  - Completion generation logic for bash and zsh
  - Command registration or help system changes to support completion generation
  - Documentation updates in README.md
