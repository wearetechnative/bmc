## Context

`BuildFederationURL` in `internal/awsops/console.go` resolves credentials with `awscfg.LoadDefaultConfig` and passes no AssumeRole options. For a profile that uses `role_arn`, the SDK's shared-config resolver fills `stscreds.AssumeRoleOptions.Duration` only from the profile's `duration_seconds` (`config/resolve_credentials.go`), leaving it at zero when that key is absent. `AssumeRoleProvider.Retrieve` then substitutes `stscreds.DefaultDuration`, which is 15 minutes.

An AWS console session created through the federation endpoint from temporary credentials lives exactly as long as those credentials. The requested AssumeRole duration therefore is the session lifetime. The Bash implementation invoked the AWS CLI, whose AssumeRole default is 3600 seconds, so the regression arrived with the Go rewrite.

Two related defects compound the problem:

- `getFederationToken` sends `DurationSeconds=43200`. AWS honours that parameter only for long-term IAM user credentials and ignores it when a session token is present, so it created a false impression that session length was already handled.
- `registerConsoleSession` in `cmd/console.go` hardcoded `time.Now().Add(time.Hour)` as the expiry, even though `BuildFederationURL` already returned the real one and `OpenConsole` discarded it. The watcher therefore scheduled its refresh at 55 minutes, long after a 15-minute session had died.

Measured against a live profile: 14m59s before the fix, 59m59s after.

## Goals / Non-Goals

**Goals:**

- Restore the pre-rewrite session lifetime of one hour by default.
- Let users request longer sessions where their roles allow it.
- Make the watcher schedule refreshes from the real credential expiry.
- Keep a role whose `MaxSessionDuration` is lower than the request working rather than failing.

**Non-Goals:**

- Changing how MFA session tokens are obtained. `mfa.EnsureValid` already requests 43200 seconds via `GetSessionToken` and is unaffected.
- Raising role `MaxSessionDuration` values in AWS, which is an IAM concern outside `bmc`.
- Reworking the watcher's CDP or tab-based refresh mechanism.

## Decisions

### Set the duration through `WithAssumeRoleCredentialOptions`

The duration is applied with `awscfg.WithAssumeRoleCredentialOptions`. The SDK appends this callback after its own shared-config callback, so the option observes whatever `duration_seconds` resolved to and can adjust it.

Alternatives considered:

- Constructing `stscreds.NewAssumeRoleProvider` directly. Rejected: it would mean re-implementing profile resolution, source-profile chaining, and SSO handling that `LoadDefaultConfig` already does.
- Writing `duration_seconds` into `~/.aws/config`. Rejected: it mutates user configuration that other tools share, for a `bmc`-specific concern.

### Treat the configured value as a ceiling

The callback applies `if o.Duration == 0 || o.Duration > limit`. A profile that sets a shorter `duration_seconds` keeps it; one that sets a longer value is capped. A profile that sets nothing gets the configured value instead of the SDK's 15 minutes.

This makes `console.session_duration_seconds` a per-tool ceiling rather than an override, which also lets the fallback path force every profile down to one hour by lowering the limit.

### Retry at one hour on a duration rejection

`isDurationRejected` matches STS `ValidationError` messages by looking for `durationseconds` case-insensitively, covering both "exceeds the MaxSessionDuration set for this role" and the role-chaining variant. Only the retry path depends on this match, and only when the request was longer than an hour; a false negative surfaces the original error, and a false positive costs one extra AssumeRole call.

Alternatives considered:

- Reading each role's `MaxSessionDuration` up front via `iam:GetRole`. Rejected: it adds an IAM permission most console users do not have, plus a round trip on every open.
- Matching on `smithy.APIError` codes. The code is `ValidationError` for several unrelated failures, so the message still has to be inspected; string matching alone keeps `smithy-go` out of the direct dependency list.

### Clamp to the AWS-accepted range

`sessionDuration` clamps to 900–43200 seconds and treats zero or negative as unset. This keeps a mistyped config value from producing an opaque STS `ValidationError`, and guarantees the function never returns zero — the value that caused the original bug.

### Share the refresh-time calculation

`watcher.RefreshTime(expiry)` replaces two copies of `credExpiry.Add(-refreshWindow)` in `daemon.go` and the hardcoded arithmetic in `cmd/console.go`. It caps the window at half the remaining lifetime, so a session shorter than twice the refresh window still gets a refresh while it is alive rather than one scheduled in the past.

### Change the `OpenConsole` signature

`OpenConsole` returns `(time.Time, error)` so callers receive the expiry that `BuildFederationURL` already computed. `cmd/console.go` is the only caller.

## Risks / Trade-offs

- **A longer configured duration silently degrades to one hour on capped roles** → The fallback is deliberate: failing to open the console would be worse than a shorter session. Users who need longer sessions must raise the role's `MaxSessionDuration` in IAM.
- **`isDurationRejected` depends on STS message text** → Scoped so a mismatch only skips the retry and surfaces the real STS error, which names the problem clearly. Unit tests cover both known message variants.
- **The configured value overrides a profile's longer `duration_seconds`** → Documented as a ceiling. Profiles here set no `duration_seconds` at all, so no current profile is affected.
- **Longer-lived credentials sit in memory longer** → The credentials are not written to disk; only the MFA session token is cached, and that behaviour is unchanged. One hour matches what the AWS CLI and the pre-rewrite Bash version already did.
- **Requesting more than 15 minutes needs the role to allow it** → Every IAM role permits at least 3600 seconds, so the default is safe everywhere.

## Migration Plan

No migration or data changes. The new config key is optional and defaults to the previous Bash-era behaviour. Users already running with a `config.json` see one-hour sessions on upgrade without editing anything. Rollback is reverting the commit; nothing persists outside the process.

## Open Questions

None.
