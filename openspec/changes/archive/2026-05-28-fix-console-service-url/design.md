## Context

`bmc console -s <service>` builds a federation sign-in URL for the AWS Console. The current implementation constructs the `Destination` parameter as `https://<service>.console.aws.amazon.com/`, treating the service name as a subdomain. The AWS federation endpoint returns a 400 Bad Request for services that do not have a dedicated subdomain (most do not — RDS, EC2, ECS, Lambda, etc. all use path-based routing).

The region is already loaded from the selected AWS profile inside `OpenConsole` via `awsCfg.Region`, but it is not passed through to `buildConsoleURL`.

## Goals / Non-Goals

**Goals:**
- Fix the 400 Bad Request for all services by switching to region-aware path-based console URLs
- Allow the `-s` flag to accept sub-paths (e.g., `systems-manager/parameters`) for deep linking
- Keep the fix minimal and contained to `buildConsoleURL` and its call site

**Non-Goals:**
- A canonical service name mapping (not needed; path-based URLs work universally)
- Aliases or shorthand service names
- Validation of service path values

## Decisions

### Decision: Region-aware path-based URL as the single URL pattern

**Chosen:** `https://<region>.console.aws.amazon.com/<service-path>/home`

**Rationale:** The AWS federation endpoint accepts any `*.console.aws.amazon.com` destination. The region-prefixed subdomain (`<region>.console.aws.amazon.com`) is the canonical modern AWS console URL — it avoids a redirect and lands the user in the correct region immediately. Using the profile's region (already resolved) makes this zero-config.

**Alternatives considered:**
- `https://console.aws.amazon.com/<service>/` — works but triggers a redirect to the region-prefixed URL; slightly worse UX
- Canonical service map — adds maintenance burden for no benefit; the path-based fallback is universal

### Decision: Treat `-s` value as a console path, not a service name

**Chosen:** Pass the value of `-s` directly as the path segment, including any `/` characters.

**Rationale:** A service name like `rds` and a sub-path like `systems-manager/parameters` are structurally identical — both are console URL paths. Treating them uniformly requires no special casing and gives the user deep-link capability for free.

**URL construction logic:**
- If the path contains a `/` → use as-is (user provided a sub-path, no `/home` suffix added)
- If the path has no `/` → append `/home` (`https://<region>.console.aws.amazon.com/rds/home`)

### Decision: Pass region as a parameter to `buildConsoleURL`

**Chosen:** Add a `region string` parameter to `buildConsoleURL`.

**Rationale:** `awsCfg.Region` is already available in `OpenConsole` and is the correct resolved region for the profile. Passing it explicitly keeps `buildConsoleURL` a pure function with no side effects.

## Risks / Trade-offs

- **`/home` suffix may not exist for all services** → For services where `/home` is not a valid path, AWS redirects to the correct page anyway. Low risk.
- **Sub-path format is user responsibility** → If a user passes an invalid path, they get a 404 in the console rather than a bmc error. Acceptable — bmc does not validate console URLs.

## Migration Plan

No migration needed. This is a bug fix with no breaking changes to the CLI interface. The `-s` flag accepts the same short service names as before; existing usage continues to work and now also supports sub-paths.
