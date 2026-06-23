## Context

`bmc` uses a single global `totp_script` in `~/.config/bmc/config.json` to acquire TOTP codes during MFA session refresh. The script is invoked by `acquireTOTP()` in `internal/mfa/session.go`. The caller, `EnsureValid()`, already knows the `sourceProfile` name — it just never passes it down.

Users with multiple AWS identities (e.g. `technative` and `wvandrtoorren`) keep separate entries in their password manager. Today they must either share one TOTP entry across all profiles or fall back to manual input for secondary accounts.

## Goals / Non-Goals

**Goals:**
- Allow per-source-profile totp script overrides via config
- Preserve backwards compatibility — existing single-script configs work unchanged
- Resolve the correct script automatically without user-facing flags or prompts

**Non-Goals:**
- Per-profile clipboard commands (`copy_command` / `paste_command`) — global values remain sufficient
- CLI flag override at runtime — config-driven is enough; a flag adds surface area without proportional value
- Support for profile name wildcards or regex matching

## Decisions

### Add `ProfileScripts map[string]string` to `MFAConfig`

**Decision**: Add a new optional field to `MFAConfig`:
```go
ProfileScripts map[string]string `json:"profile_scripts,omitempty"`
```

**Alternatives considered**:
- Separate top-level config section: more indirection, no benefit for a single map
- Array of `{profile, script}` objects: harder to look up, no advantage over a map

**Rationale**: A map is the most direct representation — key is the source profile name, value is the script command string. `omitempty` ensures zero-value marshalling stays clean.

### Lookup order: profile-specific → global fallback → manual

**Decision**: In `acquireTOTP()`, resolve script as:
1. `cfg.MFA.ProfileScripts[sourceProfile]` if the map is non-nil and the key exists
2. `cfg.MFA.TOTPScript` (existing global)
3. Prompt user for manual input

**Rationale**: This is the minimal change needed. Existing behaviour is preserved as the fallback path. Users who only set `totp_script` see no change.

### Pass `sourceProfile` to `acquireTOTP()`

**Decision**: Change signature from `acquireTOTP(cfg, outfd)` to `acquireTOTP(cfg, sourceProfile, outfd)`.

**Rationale**: `acquireTOTP` is a private function called only from `EnsureValid`, which already has `sourceProfile`. The signature change is contained entirely within `internal/mfa/session.go` — no external callers are affected.

## Risks / Trade-offs

- **Profile name sensitivity**: The map key must exactly match the source profile name as resolved by `awsconfig.ResolveSourceProfile()`. A typo in config silently falls back to the global script. Mitigation: `bmc doctor` can be extended later to validate profile names against known profiles.
- **Nil map on zero config**: If `profile_scripts` is absent from JSON, `ProfileScripts` will be `nil`. The lookup must guard with a nil check before key access. Mitigation: standard Go nil-map read (`m[key]` on nil map returns zero value, no panic) handles this safely.

## Migration Plan

No migration required. Existing `config.json` files without `profile_scripts` continue to work. Users who want per-profile scripts add the new key:

```json
{
  "mfa": {
    "totp_script": "rbw code \"Technative AWS (new)\"",
    "profile_scripts": {
      "wvandrtoorren": "rbw code \"wvandrtoorren AWS Personal\""
    }
  }
}
```
