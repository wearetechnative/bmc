## Why

Project documentation and OpenSpec files must be in English, but some files are still in Dutch or refer to the old bash-script version of bmc. This violates the project standard and makes documentation misleading for new contributors.

## What Changes

- Translate `openspec/specs/tui-color-rendering/spec.md` from Dutch to English
- Rewrite `docs/aws-profile-select.md` to reflect the current Go-based tool, shell integration via `eval "$(bmc profsel)"`, and the `/dev/tty` color rendering behavior

## Capabilities

### New Capabilities
<!-- none -->

### Modified Capabilities
- `tui-color-rendering`: Translate spec from Dutch to English (requirements unchanged)

## Impact

- `openspec/specs/tui-color-rendering/spec.md` — translation only, no behavioral changes
- `docs/aws-profile-select.md` — full rewrite to match current bmc Go CLI

Bean: [bmc-c48y](../../../../.beans/bmc-c48y--translate-files.md)
