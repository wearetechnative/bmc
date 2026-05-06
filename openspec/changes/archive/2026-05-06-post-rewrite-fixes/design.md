## Context

Four issues discovered after the Go rewrite of bmc. Three are small, well-understood fixes in existing code. One (console containerized tab) requires deeper investigation into browser isolation mechanisms.

## Goals / Non-Goals

**Goals:**
- Eliminate blank rows in short TUI lists
- Warn the user when `bmc profsel` output is not being eval'd
- Give actionable guidance when `install-shell-integration` cannot write to the rc file
- Document NixOS shell integration options in README

**Non-Goals:**
- NixOS home-manager module in the flake (user's declarative setup is their responsibility)
- Automatic migration of any existing config
- console containerized tab full implementation (scope TBD — tracked as separate concern)

## Decisions

### D1: TUI list height formula

**Decision**: `height = len(items) + 3` for lists with more than 4 items; `height = len(items) + 1` with `SetShowHelp(false)` for lists with ≤ 4 items.

**Rationale**: `bubbles/list` fills its allocated height with item-slots. Surplus slots render as blank rows. Chrome (title + help bar) costs ~3 lines; removing the help bar for trivially short lists reduces chrome to ~1 line. The 30-row cap remains.

### D2: TTY detection for profsel hint

**Decision**: After printing `export AWS_PROFILE=xxx`, check `term.IsTerminal(os.Stdout.Fd())`. If stdout is a terminal (not piped to eval), print a single-line hint to stderr: `Tip: run 'bmc install-shell-integration' to set AWS_PROFILE automatically`.

**Rationale**: stdout is a pipe when called as `eval "$(bmc profsel)"` — the hint must not appear in eval output. stderr is always safe to write to. The hint should only show when the export is not being captured.

### D3: install-shell-integration permission denied fallback

**Decision**: When `os.OpenFile` returns a `permission denied` error, do not surface a generic error. Instead print a message explaining the file is not writable, followed by manual snippets for:
1. home-manager zsh (`programs.zsh.initContent`)
2. home-manager bash (`programs.bash.initContent`)
3. Manual `~/.zshrc` / `~/.bashrc`
4. Fish shell (`~/.config/fish/config.fish`)

**Rationale**: `permission denied` on a dotfile is the reliable signal that the file is managed externally (home-manager on NixOS, chezmoi, etc.). No OS detection needed — the error condition itself is sufficient. Fish shell snippet is included because it requires different syntax.

### D4: console containerized tab

**Decision**: Defer. The mechanism Granted uses (Firefox containers or Chrome profiles via CLI flags) needs investigation. This is a separate capability with unclear scope. Track in bean `bmc-p4wn` issue 2 for a follow-up change.

## Risks / Trade-offs

- **TTY hint shown in unexpected pipes**: any pipe (not just eval) suppresses the hint. Edge case: `bmc profsel | tee log.txt` would not show the hint. Acceptable.
- **Fish snippet in fallback**: bmc does not otherwise support Fish. Including the snippet is documentation, not a code dependency. Low risk.

## Migration Plan

All changes are backwards-compatible. No config changes, no API changes, no data migrations.
