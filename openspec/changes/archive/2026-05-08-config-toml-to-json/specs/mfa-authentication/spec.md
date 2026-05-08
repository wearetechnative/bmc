# mfa-authentication Delta Specification

## MODIFIED Requirements

### Requirement: Config Loading Path
The MFA subsystem reads its configuration (`enabled`, `totp_script`, `clipboard_command`) from the JSON config file at `~/.config/bmc/config.json` rather than `config.toml`. All field names and defaults remain unchanged.
