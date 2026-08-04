## ADDED Requirements

### Requirement: Console sessions last a requested duration

The `bmc console` command SHALL open sessions whose lifetime is determined by an explicitly requested AssumeRole duration, and SHALL report that session's expiry to the user. The duration, its configuration, and its bounds are defined by the `console-session-duration` capability.

#### Scenario: Console opened for a role-based profile

- **WHEN** a user opens the console for a profile that assumes a role
- **THEN** the session SHALL last the requested duration, one hour by default
- **AND** the command SHALL print the expiry time to stderr

#### Scenario: Console opened with --watch

- **WHEN** a user runs `bmc console --watch`
- **THEN** the session registered with the watcher SHALL carry the real credential expiry
- **AND** the watcher SHALL refresh the session before that expiry

## MODIFIED Requirements

### Requirement: Interactive profile selection shows recent profiles
When `AWS_PROFILE` is not set and no `-p` flag is given, the `bmc console` command SHALL present the interactive profile selector with recently used profiles shown at the top of the list.

#### Scenario: Recent profiles exist

- **WHEN** a user runs `bmc console` without `AWS_PROFILE` set and without `-p`
- **AND** the console history contains previously used profiles
- **THEN** the interactive profile selector SHALL list those profiles at the top of the list

#### Scenario: No recent profiles

- **WHEN** a user runs `bmc console` without `AWS_PROFILE` set and without `-p`
- **AND** the console history is empty
- **THEN** the interactive profile selector SHALL list profiles without a recent section
