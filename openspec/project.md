# Project Context

## Purpose
BMC (Bill McCloud) is a collection of AWS/Terraform DevOps tools designed to simplify working with aws-cli and the AWS Console for cloud engineers and DevOps teams. The project provides:
- AWS profile switching and management with MFA support
- EC2 instance management (list, connect, start/stop/hibernate)
- ECS container connection capabilities
- Interactive TUI for easier AWS resource selection
- Browser extension configuration export for AWS profile switching

## Tech Stack
- Bash (primary scripting language, v4+ required)
- ZSH (v5.8+ supported)
- Nix (packaging and distribution via flakes)
- jq (JSON processing)
- gum (interactive TUI elements)
- awk (text processing and CSV manipulation)
- AWS CLI (core AWS operations)
- Additional tools: jsonify-aws-dotfiles, aws-mfa, assumego

## Project Conventions

### Code Style
- Shell scripts must be POSIX-compatible where possible, with bash/zsh specific features clearly documented
- Function naming: camelCase (e.g., `selectAWSProfile`, `ec2ListInstances`, `setMFA`)
- Variable naming: camelCase for local variables, SCREAMING_SNAKE_CASE for global constants
- Always use double quotes for variable expansion to handle spaces in paths
- Sourcing detection: Scripts that set environment variables must detect if they're sourced vs executed
- Error handling: Use `exit 1` for errors, check command return codes with `$?`
- Comments: Use `#` for inline comments, explain complex logic blocks
- Dependencies: Check for required commands at script start using `checkdeps` function
- Array handling: Use portable awk/sed solutions instead of bash-specific array syntax for zsh compatibility

### Architecture Patterns
- Modular design: Core functionality in `_bmclib.sh`, specific commands in individual scripts
- Central command dispatcher: `bmc` script acts as main entry point with subcommands
- Command registration pattern: `make_command` function registers commands with descriptions for help system
- Shared library pattern: Common functions (AWS profile selection, MFA, EC2 operations) in `_bmclib.sh`
- Interactive selection: Use `gum` for user-friendly menus and tables instead of raw read/select
- Configuration: User config stored in `~/.config/bmc/config.env`
- AWS credentials: Follow standard AWS config file structure in `~/.aws/config` and `~/.aws/credentials`
- Version management: Single VERSION-bmc file for version tracking

### Testing Strategy
- Manual testing workflow (CI testing planned but not yet implemented - see README TODO)
- Test on both bash (v4+) and zsh (v5.8+) for compatibility
- Verify AWS profile switching in both sourced and executed modes
- Test MFA flows with time-based expiration
- Validate OS-specific behavior (Linux vs macOS for date commands)

### Git Workflow
- Main branch: `main` (used for PRs and releases)
- Commit style: lowercase prefix with description (e.g., "fix: replace bash-specific array syntax", "feature: add ec2find option")
- Release cycle: Versions are released periodically, not for every change
- Changelog workflow:
  - New changes are added under "## NEXT VERSION" heading in CHANGELOG.md
  - When ready to release, replace "NEXT VERSION" with actual version number and date
  - Update VERSION-bmc file at the same time
- Changelog format: Categorized changes (Feature, Fix, Breaking, Enhancement) with concise descriptions
- Pull requests: Merge to main branch

## Domain Context
- AWS IAM roles and profiles: Scripts work with AWS named profiles, role assumption, and MFA devices
- AWS profile groups: Custom grouping mechanism for organizing multiple AWS accounts
- MFA session management: Track expiration times and auto-renew sessions
- EC2 hibernation: Special handling for instances with hibernation enabled vs standard stop
- AWS SSM vs SSH: Support both Systems Manager Session Manager and traditional SSH for EC2 connection
- Browser extensions: Export config format for AWS Extend Switch Roles Firefox/Chrome extension

## Important Constraints
- Shell compatibility: Must work on both bash v4+ and zsh v5.8+
- POSIX compatibility: Prefer portable solutions over bash-specific features where possible
- macOS vs Linux: Handle OS differences for date command formatting
- AWS CLI required: Scripts assume aws-cli is configured and available
- Sourcing requirement: Profile switching commands (profsel) must be sourced to export environment variables
- No root required: All operations run as regular user with appropriate AWS credentials
- Terminal requirement: Interactive features require TTY for gum menus

## External Dependencies
- **Required CLI tools**: aws-cli, jq, awk, gum, jsonify-aws-dotfiles, assumego
- **Optional**: aws-mfa (for MFA-enabled profiles), custom TOTP script (configurable)
- **AWS services**: IAM (roles, MFA devices), EC2, ECS, Systems Manager (SSM)
- **Configuration files**: ~/.aws/config, ~/.aws/credentials, ~/.config/bmc/config.env
- **Nix**: For reproducible builds and distribution (flake.nix)
- **Firefox/Chrome extensions**: AWS Extend Switch Roles (for browser integration)
