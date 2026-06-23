## ADDED Requirements

### Requirement: Per-source-profile TOTP script configuration
The system SHALL support a `profile_scripts` map in the `mfa` config section that maps source profile names to totp script command strings. When a source profile has an entry in this map, that script SHALL be used instead of the global `totp_script`.

#### Scenario: Profile-specific script used when configured
- **WHEN** `mfa.profile_scripts` contains an entry matching the current source profile name
- **THEN** the system SHALL execute the profile-specific script to obtain the TOTP code
- **AND** the system SHALL NOT execute the global `totp_script`

#### Scenario: Global script used as fallback when no profile entry exists
- **WHEN** `mfa.profile_scripts` does not contain an entry for the current source profile
- **AND** `mfa.totp_script` is configured
- **THEN** the system SHALL execute the global `totp_script` to obtain the TOTP code

#### Scenario: Manual entry when neither profile script nor global script is configured
- **WHEN** `mfa.profile_scripts` does not contain an entry for the current source profile
- **AND** `mfa.totp_script` is empty or absent
- **THEN** the system SHALL prompt the user to enter the MFA code manually

#### Scenario: Backwards compatibility when profile_scripts is absent
- **WHEN** `~/.config/bmc/config.json` does not contain a `profile_scripts` key
- **THEN** the system SHALL behave identically to before this change
- **AND** the global `totp_script` SHALL be used for all source profiles

#### Scenario: Empty profile_scripts map uses global fallback
- **WHEN** `mfa.profile_scripts` is present but empty (`{}`)
- **THEN** the system SHALL use the global `totp_script` for all source profiles
