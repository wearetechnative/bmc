## Context

The project requires all OpenSpec and documentation files to be in English. Two files need updating:

1. `openspec/specs/tui-color-rendering/spec.md` — written in Dutch when the fix was implemented; requirements are correct but language is wrong.
2. `docs/aws-profile-select.md` — describes the original bash script (`aws-profile-select.sh`). The tool has since been rewritten in Go as `bmc`. The doc must reflect the current CLI, shell integration pattern (`eval "$(bmc profsel)"`), and the `/dev/tty` color rendering behavior added to fix colors in subshell contexts.

## Goals / Non-Goals

**Goals:**
- Translate `tui-color-rendering` spec to English without changing requirements
- Replace `docs/aws-profile-select.md` with accurate documentation for the current Go CLI

**Non-Goals:**
- Changing any requirements or behavior
- Translating archived change directories (historical artifacts, not actively used)
- Adding new features or fixing bugs

## Decisions

**Rewrite docs/aws-profile-select.md from scratch**: The existing content is too tied to the old bash script to be updated incrementally. A fresh write based on the current codebase is cleaner.

**Keep tui-color-rendering spec requirements intact**: Only translate the language, do not alter scenarios, SHALL statements, or behavior descriptions.

## Risks / Trade-offs

- [Risk] Rewriting docs may omit useful context → Read the existing doc fully before rewriting, and cross-reference with current code behavior.
