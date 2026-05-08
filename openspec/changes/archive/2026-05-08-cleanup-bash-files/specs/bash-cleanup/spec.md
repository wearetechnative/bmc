## ADDED Requirements

### Requirement: Bash-era files are removed from the repository
The repository SHALL NOT contain `_bmclib.sh`, `_get_var_file.sh`, or `tgselect.sh` after this change is applied.

#### Scenario: Bash library is gone
- **WHEN** a developer clones the repository
- **THEN** `_bmclib.sh` is not present in the working tree

#### Scenario: Terraform helper is gone
- **WHEN** a developer clones the repository
- **THEN** `_get_var_file.sh` is not present in the working tree

#### Scenario: Toggl script is gone
- **WHEN** a developer clones the repository
- **THEN** `tgselect.sh` is not present in the working tree

### Requirement: bmc-go binary is not tracked by git
The repository SHALL have `/bmc-go` listed in `.gitignore` so the Go build artifact is never accidentally committed.

#### Scenario: bmc-go is gitignored
- **WHEN** a developer builds the project and runs `git status`
- **THEN** `bmc-go` does not appear as an untracked or modified file
