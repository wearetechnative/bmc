# Fix short-lived AWS console sessions

## Why

Since the migration from Bash to Go, sessions opened with `bmc console` expire after roughly 15 minutes instead of the hour users had before. `BuildFederationURL` loads credentials without specifying an AssumeRole duration, so the AWS SDK for Go substitutes its own `stscreds.DefaultDuration` of 15 minutes. The Bash implementation shelled out to the AWS CLI, which defaults to 3600 seconds — hence the regression. A console session created from temporary credentials lives exactly as long as those credentials, so the requested AssumeRole duration is what determines when the user is logged out.

`bmc console --watch` could not compensate, because it registered a hardcoded one-hour expiry with the watcher and therefore scheduled its refresh roughly 40 minutes after the session had already died.

## What Changes

- `bmc console` requests an explicit AssumeRole duration, defaulting to one hour, instead of letting the SDK fall back to 15 minutes.
- New `console.session_duration_seconds` config option lets users request longer sessions, clamped to the 900–43200 second range AWS accepts.
- A profile's own `duration_seconds` is honoured when it is shorter than the configured value.
- When the requested duration exceeds a role's `MaxSessionDuration`, credential resolution retries at one hour rather than failing outright.
- `bmc console` prints when the session expires, and registers that real expiry with the watcher instead of assuming one hour.
- The watcher schedules its refresh from the real credential expiry, never later than the session's midpoint, so short sessions still refresh while alive.
- `DurationSeconds` is no longer sent to the federation endpoint for temporary credentials, where AWS ignores it.
- Incidental spec repair: the existing `aws-console-access` requirement "Interactive profile selection shows recent profiles" gains the scenarios it was missing, which strict validation requires before the spec can be rebuilt.

## Capabilities

### New Capabilities
- `console-session-duration`: How long a `bmc console` session lasts — the AssumeRole duration requested, its configuration and bounds, the fallback when a role caps it, and how the resulting expiry is reported and handed to the watcher.

### Modified Capabilities
- `aws-console-access`: Opening a console SHALL report the session expiry to the user, and SHALL request a session duration rather than accepting the SDK default.

## Impact

- `internal/awsops/console.go`: `BuildFederationURL` sets the AssumeRole duration and retries on a `MaxSessionDuration` rejection; `OpenConsole` returns the credential expiry; `getFederationToken` omits `DurationSeconds` for temporary credentials.
- `internal/config/config.go`: `ConsoleConfig` gains `SessionDurationSeconds`, default 3600.
- `cmd/console.go`: consumes the returned expiry, prints it, and passes it to the watcher registration.
- `internal/watcher/state.go`, `internal/watcher/daemon.go`: shared `RefreshTime` helper replaces two copies of the refresh-time calculation.
- `docs/content/commands/console.md`: documents session lifetime and the new option.
- Dependency: uses `github.com/aws/aws-sdk-go-v2/credentials/stscreds`, already an existing direct dependency.
- No breaking changes. Users on the default config get one-hour sessions instead of 15-minute ones.
