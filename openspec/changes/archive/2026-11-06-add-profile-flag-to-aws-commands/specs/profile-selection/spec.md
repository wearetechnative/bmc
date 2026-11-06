## MODIFIED Requirements

### Requirement: Profile resolution priority in ensureAWSProfile
The `ensureAWSProfile()` function SHALL resolve the active AWS profile using the following priority order, from highest to lowest:
1. The `-p`/`--profile` flag value (`globalProfile` variable) — if set, use directly and skip env var and interactive selection
2. The `AWS_PROFILE` environment variable — if set, validate and use
3. Interactive profile picker — if neither flag nor env var is set

Regardless of which resolution path is taken, the MFA validity check SHALL always run before returning the profile.

#### Scenario: Flag value takes priority over env var
- **GIVEN** `AWS_PROFILE=staging` is set in the environment
- **AND** the user passes `-p production` on the command line
- **WHEN** `ensureAWSProfile()` is called
- **THEN** it SHALL use `production` and skip interactive selection
- **AND** SHALL run the MFA check for `production`

#### Scenario: Env var used when no flag set
- **GIVEN** `AWS_PROFILE=staging` is set in the environment
- **AND** no `-p` flag is provided
- **WHEN** `ensureAWSProfile()` is called
- **THEN** it SHALL use `staging` without showing the interactive picker
- **AND** SHALL run the MFA check for `staging`

#### Scenario: Interactive picker shown when neither flag nor env var is set
- **GIVEN** `AWS_PROFILE` is not set in the environment
- **AND** no `-p` flag is provided
- **WHEN** `ensureAWSProfile()` is called
- **THEN** it SHALL show the interactive profile picker
