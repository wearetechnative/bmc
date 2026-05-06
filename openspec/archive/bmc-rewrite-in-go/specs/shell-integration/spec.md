## ADDED Requirements

### Requirement: profsel outputs export statement
The system SHALL output `export AWS_PROFILE=<profile-name>` to stdout when `profsel` completes successfully, so it can be consumed via `eval`.

#### Scenario: Successful profile selection
- **WHEN** user completes profile selection in profsel
- **THEN** system prints `export AWS_PROFILE=<selected-profile>` to stdout

#### Scenario: No profile selected
- **WHEN** user cancels selection
- **THEN** system prints nothing and exits with code 0

### Requirement: install-shell-integration command
The system SHALL provide `bmc install-shell-integration` which detects the user's shell and appends the profsel wrapper function to the appropriate rc file.

#### Scenario: Zsh detected
- **WHEN** current shell is zsh
- **THEN** system appends the wrapper function to `~/.zshrc` if not already present

#### Scenario: Bash detected
- **WHEN** current shell is bash
- **THEN** system appends the wrapper function to `~/.bashrc` if not already present

#### Scenario: Wrapper already installed
- **WHEN** the wrapper function is already present in the rc file
- **THEN** system reports "already installed" and makes no changes

#### Scenario: Unsupported shell
- **WHEN** shell is neither bash nor zsh
- **THEN** system displays the wrapper function for manual installation and exits with code 0

### Requirement: Shell wrapper function
The installed wrapper SHALL intercept `bmc profsel` calls and eval the output, while passing all other commands through to the binary unchanged.

#### Scenario: Wrapper intercepts profsel
- **WHEN** user calls `bmc profsel` with the wrapper installed
- **THEN** shell function runs `eval "$(command bmc profsel "$@")"`, setting `AWS_PROFILE` in the current shell

#### Scenario: Wrapper passes through other commands
- **WHEN** user calls `bmc ec2ls` or any other command
- **THEN** shell function runs `command bmc ec2ls` directly without eval
