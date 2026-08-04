# console-session-duration Specification

## Purpose
TBD - created by archiving change fix-console-session-duration. Update Purpose after archive.
## Requirements
### Requirement: Request an explicit console session duration

`bmc console` SHALL request an explicit AssumeRole duration when resolving credentials for a console session, defaulting to 3600 seconds. It SHALL NOT rely on the AWS SDK's own default duration, which is 15 minutes.

#### Scenario: Profile without duration_seconds

- **WHEN** a user opens a console for a profile that uses `role_arn` and does not set `duration_seconds`
- **THEN** the AssumeRole call SHALL request 3600 seconds
- **AND** the resulting console session SHALL remain usable for approximately one hour

#### Scenario: Profile with a shorter duration_seconds

- **WHEN** a profile sets `duration_seconds` to a value below the configured console session duration
- **THEN** the AssumeRole call SHALL request the profile's shorter value

#### Scenario: Profile with a longer duration_seconds

- **WHEN** a profile sets `duration_seconds` to a value above the configured console session duration
- **THEN** the AssumeRole call SHALL request the configured console session duration

### Requirement: Configurable session duration

The `bmc` configuration SHALL support a `console.session_duration_seconds` integer option, defaulting to 3600, that sets the AssumeRole duration requested for console sessions.

#### Scenario: Option is absent from the config file

- **WHEN** `~/.config/bmc/config.json` does not set `console.session_duration_seconds`
- **THEN** `bmc console` SHALL request 3600 seconds

#### Scenario: Option requests a longer session

- **WHEN** `console.session_duration_seconds` is set to 14400
- **AND** the role's `MaxSessionDuration` permits it
- **THEN** `bmc console` SHALL request 14400 seconds
- **AND** the console session SHALL remain usable for approximately four hours

#### Scenario: Option is below the AWS minimum

- **WHEN** `console.session_duration_seconds` is set to a value below 900
- **THEN** `bmc console` SHALL request 900 seconds

#### Scenario: Option is above the AWS maximum

- **WHEN** `console.session_duration_seconds` is set to a value above 43200
- **THEN** `bmc console` SHALL request 43200 seconds

#### Scenario: Option is zero or negative

- **WHEN** `console.session_duration_seconds` is set to 0 or a negative value
- **THEN** `bmc console` SHALL request 3600 seconds

### Requirement: Fall back when a role caps the session duration

When STS rejects the requested duration because it exceeds the role's `MaxSessionDuration`, `bmc console` SHALL retry once at 3600 seconds instead of reporting an error. Every IAM role permits at least one hour, so the retry succeeds.

#### Scenario: Requested duration exceeds MaxSessionDuration

- **WHEN** the configured duration is longer than one hour
- **AND** STS returns a `ValidationError` stating that the requested `DurationSeconds` exceeds the role's `MaxSessionDuration`
- **THEN** `bmc console` SHALL retry credential resolution requesting 3600 seconds
- **AND** the console SHALL open with a one-hour session

#### Scenario: Role reached by role chaining

- **WHEN** the configured duration is longer than one hour
- **AND** STS rejects it because role chaining limits the session to one hour
- **THEN** `bmc console` SHALL retry requesting 3600 seconds

#### Scenario: Unrelated credential failure

- **WHEN** credential resolution fails for a reason unrelated to the requested duration, such as an expired token or access denied
- **THEN** `bmc console` SHALL report that error to the user
- **AND** SHALL NOT retry with a different duration

### Requirement: Report the session expiry

`bmc console` SHALL report to the user when the opened console session expires, based on the expiry of the credentials that created it.

#### Scenario: Console opened successfully

- **WHEN** a user runs `bmc console` and the console opens
- **THEN** the command SHALL print the session expiry time to stderr

#### Scenario: Credentials without an expiry

- **WHEN** the resolved credentials are long-lived and carry no expiry
- **THEN** the command SHALL NOT print an expiry time
- **AND** SHALL still open the console

### Requirement: Federation DurationSeconds only for long-term credentials

The AWS federation endpoint honours the `DurationSeconds` parameter only for long-term IAM user credentials; with temporary credentials the console session inherits the credentials' remaining lifetime instead. `bmc` SHALL send `DurationSeconds` only when the credentials carry no session token.

#### Scenario: Temporary credentials

- **WHEN** the resolved credentials include a session token
- **THEN** the `getSigninToken` request SHALL omit `DurationSeconds`

#### Scenario: Long-term credentials

- **WHEN** the resolved credentials carry no session token
- **THEN** the `getSigninToken` request SHALL include `DurationSeconds`

### Requirement: Register the real expiry with the watcher

When `bmc console --watch` registers a session, it SHALL use the credential expiry reported when the console was opened, not an assumed duration. The refresh SHALL be scheduled shortly before that expiry and never later than the session's midpoint, so that sessions shorter than the refresh window are still refreshed while alive.

#### Scenario: One-hour session registered

- **WHEN** a user runs `bmc console --watch` and the credentials expire in one hour
- **THEN** the registered session expiry SHALL be that credential expiry
- **AND** the refresh SHALL be scheduled five minutes before it

#### Scenario: Session shorter than twice the refresh window

- **WHEN** a session is registered whose remaining lifetime is shorter than twice the refresh window
- **THEN** the refresh SHALL be scheduled at the midpoint of the remaining lifetime
- **AND** the refresh time SHALL be in the future

#### Scenario: Watcher refreshes a session

- **WHEN** the watcher daemon refreshes a registered session
- **THEN** it SHALL record the new credential expiry
- **AND** SHALL schedule the next refresh from that expiry using the same rule

#### Scenario: Credentials without an expiry

- **WHEN** a session is registered from credentials that carry no expiry
- **THEN** the watcher SHALL assume a one-hour session so that it keeps checking

