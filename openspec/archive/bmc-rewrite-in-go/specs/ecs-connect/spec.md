## ADDED Requirements

### Requirement: Interactive ECS cluster/service/task/container selection
The system SHALL present sequential filterable lists to guide the user: cluster → service → task → container, with breadcrumb display after each selection.

#### Scenario: Full interactive selection
- **WHEN** user runs `bmc ecsconnect`
- **THEN** system queries ECS and presents: cluster list → service list → running task list → container list, showing breadcrumbs (e.g., `my-cluster > my-service > abc123 > app`) after each step

#### Scenario: No clusters found
- **WHEN** no ECS clusters exist in the region
- **THEN** system displays an error with the current region and exits

#### Scenario: AWS_PROFILE not set
- **WHEN** `AWS_PROFILE` is not set
- **THEN** system invokes profile selection before proceeding

### Requirement: ECS execute-command shell handoff
The system SHALL connect to the selected container by executing `aws ecs execute-command --cluster <cluster> --interactive --container <container> --command /bin/sh --task <task-arn>` via `syscall.Exec`.

#### Scenario: aws CLI or session-manager-plugin not found
- **WHEN** `aws` binary or `session-manager-plugin` is not found before exec
- **THEN** system displays prerequisite error with install instructions (apt, brew, nix-env, nix profile, NixOS config) and exits

#### Scenario: Successful exec handoff
- **WHEN** all selections are made and prerequisites satisfied
- **THEN** system uses `syscall.Exec` to replace bmc process with the `aws ecs execute-command` invocation
