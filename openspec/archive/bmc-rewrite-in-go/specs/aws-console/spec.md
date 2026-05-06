## ADDED Requirements

### Requirement: Open AWS console via federation URL
The system SHALL open the AWS Management Console in the default browser by performing STS AssumeRole and constructing a federation sign-in URL — without any external tool (no assumego/Granted).

#### Scenario: Open console for selected profile
- **WHEN** user runs `bmc console` and selects a profile
- **THEN** system calls STS AssumeRole, fetches a federation sign-in token, builds the console URL, and opens it in the default browser

#### Scenario: Open console with -p flag
- **WHEN** user runs `bmc console -p <profile-name>`
- **THEN** system skips interactive selection and opens console for the specified profile

#### Scenario: Open console at specific service
- **WHEN** user runs `bmc console -s <service>` (e.g., `-s ec2`)
- **THEN** system appends the service destination to the console URL

#### Scenario: Use existing AWS_PROFILE
- **WHEN** `AWS_PROFILE` is already set and user runs `bmc console` with no flags
- **THEN** system uses the current `AWS_PROFILE` without prompting for selection

### Requirement: Federation URL construction
The system SHALL implement the AWS federation sign-in flow: POST temporary credentials to the federation endpoint, receive a signin token, and build the final sign-in URL.

#### Scenario: Successful federation token retrieval
- **WHEN** STS AssumeRole succeeds and credentials are passed to the federation endpoint
- **THEN** system receives a `SigninToken` from `https://signin.aws.amazon.com/federation`

#### Scenario: Federation endpoint unreachable
- **WHEN** the federation endpoint returns an error
- **THEN** system displays an actionable error and exits with non-zero code

### Requirement: Browser opening cross-platform
The system SHALL open the URL using `xdg-open` on Linux and `open` on macOS.

#### Scenario: Linux browser open
- **WHEN** running on Linux and console URL is ready
- **THEN** system calls `xdg-open <url>`

#### Scenario: macOS browser open
- **WHEN** running on macOS and console URL is ready
- **THEN** system calls `open <url>`
