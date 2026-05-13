# mfa-clipboard-paste Specification

## Purpose
Defines the automatic paste-after-copy capability for MFA TOTP codes, the execution of `totp_script` via `sh -c` for quoted-argument and interactive-TUI support, and the renaming of `clipboard_command` to `copy_command`.

## Requirements

### Requirement: Automatic paste after clipboard copy
When `mfa.paste_command` is configured and clipboard copy succeeds, the system SHALL wait 300ms then execute `paste_command` to simulate a paste keystroke in the focused window.

#### Scenario: Paste after successful copy
- **WHEN** `mfa.copy_command` is configured and copy succeeds
- **AND** `mfa.paste_command` is configured
- **THEN** the system SHALL wait 300ms after the copy
- **AND** the system SHALL execute `paste_command` without stdin
- **AND** the system SHALL display "-- Pasted to active window" on success

#### Scenario: Paste skipped when copy fails
- **WHEN** `mfa.copy_command` is configured but copy fails
- **AND** `mfa.paste_command` is configured
- **THEN** the system SHALL NOT execute `paste_command`
- **AND** the system SHALL display the copy failure message only

#### Scenario: Paste skipped when paste_command not configured
- **WHEN** `mfa.copy_command` is configured and copy succeeds
- **AND** `mfa.paste_command` is not configured
- **THEN** the system SHALL NOT attempt any paste
- **AND** behaviour SHALL be identical to before this change

#### Scenario: paste_command failure is non-fatal
- **WHEN** `paste_command` is configured and executed
- **AND** `paste_command` exits with a non-zero status
- **THEN** the system SHALL display "-- Note: Paste failed (error details)"
- **AND** the system SHALL continue normally without aborting MFA flow

### Requirement: totp_script supports quoted arguments and interactive TUI tools
The system SHALL execute `totp_script` via `sh -c`, passing the parent process's stdin and stderr to the child, so that scripts with quoted arguments and interactive TUI selection menus work correctly.

#### Scenario: totp_script with spaces in arguments
- **WHEN** `mfa.totp_script` contains a command with quoted arguments (e.g. `rbw code "My Entry (new)"`)
- **THEN** the system SHALL execute the command with correct argument boundaries
- **AND** the quoted string SHALL be passed as a single argument

#### Scenario: totp_script with interactive TUI
- **WHEN** `mfa.totp_script` invokes an interactive selection tool (e.g. gum filter)
- **THEN** the tool SHALL be able to render its UI on the terminal
- **AND** the tool SHALL be able to receive keyboard input
- **AND** the selected value SHALL be captured as the TOTP code

### Requirement: copy_command replaces clipboard_command
The MFA configuration field for the clipboard copy command SHALL be named `copy_command`. The old `clipboard_command` field is removed.

#### Scenario: copy_command used for clipboard copy
- **WHEN** `mfa.copy_command` is set in `~/.config/bmc/config.json`
- **THEN** the system SHALL pipe the TOTP code via stdin to that command
- **AND** behaviour SHALL be identical to the former `clipboard_command`

#### Scenario: Old clipboard_command not recognised
- **WHEN** `mfa.clipboard_command` is present in `~/.config/bmc/config.json`
- **THEN** the system SHALL silently ignore it (JSON unmarshal zero-value)
- **AND** no clipboard copy SHALL occur until the user renames the field to `copy_command`
