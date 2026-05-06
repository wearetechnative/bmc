## ADDED Requirements

### Requirement: Lazy per-command prerequisite validation
The system SHALL check only the prerequisites required by the current command, at the point of use (not at startup).

#### Scenario: ec2connect SSH prerequisite check
- **WHEN** user selects SSH connection method in ec2connect
- **THEN** system checks for `ssh` binary before executing; if absent, displays error with install instructions

#### Scenario: ec2connect SSM prerequisite check
- **WHEN** user selects SSM connection method in ec2connect
- **THEN** system checks for `aws` CLI v2 and `session-manager-plugin` before executing; if either absent, displays error

#### Scenario: ecsconnect prerequisite check
- **WHEN** user runs `bmc ecsconnect`
- **THEN** system checks for `aws` CLI v2 and `session-manager-plugin` before the exec handoff

### Requirement: Actionable error messages with per-platform install instructions
The system SHALL display a structured error message when a prerequisite is missing, including the binary name, which command requires it, and install instructions for all supported platforms.

#### Scenario: Missing session-manager-plugin error format
- **WHEN** `session-manager-plugin` is not found
- **THEN** system displays:
  ```
  ✗ session-manager-plugin not found
    Required for: ec2connect (SSM), ecsconnect
    Install:
      apt:              sudo apt install amazon-ssm-agent
      brew:             brew install session-manager-plugin
      nix-env:          nix-env -iA nixpkgs.session-manager-plugin
      nix profile:      nix profile add nixpkgs#session-manager-plugin
      NixOS config:     environment.systemPackages = [ pkgs.session-manager-plugin ];
  ```

#### Scenario: Missing ssh error format
- **WHEN** `ssh` binary is not found
- **THEN** system displays error with install instructions for apt, brew, and all three Nix methods

### Requirement: bmc doctor full system check
The system SHALL provide a `bmc doctor` command that checks all prerequisites and configuration and reports status with actionable remediation for each failure.

#### Scenario: All checks pass
- **WHEN** all prerequisites and config are present and valid
- **THEN** system displays all checks with `✓` prefix and exits with code 0

#### Scenario: Some checks fail
- **WHEN** one or more checks fail
- **THEN** system displays failing checks with `✗` prefix, install/fix instructions, and exits with non-zero code

#### Scenario: Doctor covers all categories
- **WHEN** user runs `bmc doctor`
- **THEN** system checks and reports on:
  - Core: `~/.aws/config` (with profile count), `~/.aws/credentials`, `~/.config/bmc/config.toml`
  - Optional: `ssh`, `aws` CLI v2 (with version), `session-manager-plugin`
  - MFA: enabled status, `totp_script` configured, `clipboard_command` found
  - Shell integration: profsel wrapper installed in shell rc file
  - Legacy: `config.env` present warning

#### Scenario: aws CLI version check
- **WHEN** `aws` is found
- **THEN** system verifies it is v2 (not v1) by checking `aws --version` output; warns if v1 detected
