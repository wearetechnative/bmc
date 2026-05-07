## ADDED Requirements

### Requirement: MFA session validated when profile is set via environment
When `AWS_PROFILE` is set in the shell environment before running bmc, the system SHALL still validate and refresh the MFA session before executing any AWS operation.

#### Scenario: Valid MFA session with pre-set profile
- **WHEN** `AWS_PROFILE` is set in the environment
- **AND** the MFA session for that profile is still valid
- **THEN** bmc SHALL proceed without prompting for MFA
- **AND** a message SHALL be shown indicating the session is valid and its expiry time

#### Scenario: Expired MFA session with pre-set profile
- **WHEN** `AWS_PROFILE` is set in the environment
- **AND** the MFA session for that profile has expired
- **THEN** bmc SHALL prompt for a TOTP code (or run the configured TOTP script)
- **AND** refresh the session before executing the AWS operation
- **AND** NOT return `InvalidClientTokenId` from AWS

#### Scenario: MFA disabled with pre-set profile
- **WHEN** `AWS_PROFILE` is set in the environment
- **AND** MFA is disabled in the bmc config
- **THEN** bmc SHALL proceed immediately without any MFA check

#### Scenario: Profile without MFA device with pre-set profile
- **WHEN** `AWS_PROFILE` is set in the environment
- **AND** the profile has no MFA device configured
- **THEN** bmc SHALL proceed without prompting for MFA
- **AND** display "!! AWS MFA Device not found. Can't renew session"
