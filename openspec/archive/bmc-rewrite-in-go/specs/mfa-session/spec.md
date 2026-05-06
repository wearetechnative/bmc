## ADDED Requirements

### Requirement: MFA session validity check
The system SHALL read the `expiration` field from `~/.aws/credentials` under `[<sourceProfile>]` and compare it to the current time to determine if the MFA session is still valid.

#### Scenario: Valid session
- **WHEN** `expiration` in `~/.aws/credentials` is in the future
- **THEN** system logs "Current MFA Session Valid, until: <datetime>" and does not refresh

#### Scenario: Expired or missing session
- **WHEN** `expiration` is in the past or the field is absent
- **THEN** system proceeds to MFA refresh

#### Scenario: MFA disabled in config
- **WHEN** `mfa.enabled = false` in `~/.config/bmc/config.toml`
- **THEN** system skips the MFA check entirely

### Requirement: MFA device lookup
The system SHALL read `aws_mfa_device` from the `[<sourceProfile>-long-term]` section in `~/.aws/credentials`.

#### Scenario: MFA device found
- **WHEN** `[<sourceProfile>-long-term]` section exists with `aws_mfa_device`
- **THEN** system uses the device ARN for `GetSessionToken`

#### Scenario: MFA device not found
- **WHEN** no `aws_mfa_device` entry is found
- **THEN** system logs a warning and skips MFA refresh

### Requirement: TOTP code acquisition
The system SHALL obtain the TOTP code either by executing the configured `totp_script` or by prompting the user interactively.

#### Scenario: totp_script configured and succeeds
- **WHEN** `mfa.totp_script` is set in config and the script exits 0
- **THEN** system uses the script output as the TOTP code

#### Scenario: totp_script not configured
- **WHEN** `mfa.totp_script` is empty
- **THEN** system prompts user to enter the MFA code manually via bubbles/textinput

#### Scenario: Clipboard copy configured
- **WHEN** `mfa.clipboard_command` is set and TOTP code is obtained
- **THEN** system pipes the TOTP code to the clipboard command; logs a warning if the command fails but does not abort

### Requirement: STS GetSessionToken and credential write
The system SHALL call `sts.GetSessionToken` with the MFA token and write the resulting credentials to `[<sourceProfile>]` in `~/.aws/credentials` in the same format as aws-mfa (broamski).

#### Scenario: Successful token refresh
- **WHEN** `GetSessionToken` succeeds
- **THEN** system writes `aws_access_key_id`, `aws_secret_access_key`, `aws_session_token`, `expiration` under `[<sourceProfile>]` in `~/.aws/credentials`

#### Scenario: Credential write format compatibility
- **WHEN** credentials are written
- **THEN** the `expiration` value is formatted as `YYYY-MM-DD HH:MM:SS` (UTC), matching the aws-mfa tool format so other AWS tooling remains unaffected

#### Scenario: File lock during write
- **WHEN** writing to `~/.aws/credentials`
- **THEN** system acquires an exclusive file lock before writing and releases it after, to prevent concurrent corruption

#### Scenario: Wrong TOTP code
- **WHEN** `GetSessionToken` returns an authentication error
- **THEN** system displays "Wrong TOTP code?" error and exits with non-zero code
