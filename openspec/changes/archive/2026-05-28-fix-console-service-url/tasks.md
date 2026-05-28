## 1. Fix URL Generation

- [x] 1.1 Add `region string` parameter to `buildConsoleURL` in `internal/awsops/console.go`
- [x] 1.2 Replace subdomain URL pattern (`https://<service>.console.aws.amazon.com/`) with region-aware path-based pattern (`https://<region>.console.aws.amazon.com/<service>/home`)
- [x] 1.3 Add sub-path detection: if the service value contains `/`, use it as-is without appending `/home`
- [x] 1.4 Update the call to `buildConsoleURL` in `OpenConsole` to pass `awsCfg.Region`

## 2. Spec Update

- [x] 2.1 Update `openspec/specs/aws-console-access/spec.md` to reflect the new `-s` path semantics and region-aware URL construction (apply delta from change spec)

## 3. Verification

- [x] 3.1 Test `bmc console -s rds` — verify URL is `https://<region>.console.aws.amazon.com/rds/home`
- [x] 3.2 Test `bmc console -s systems-manager/parameters` — verify URL contains the full path without `/home` suffix
- [x] 3.3 Test `bmc console -s s3` — verify it still works correctly
- [x] 3.4 Test `bmc console` without `-s` — verify homepage URL is unaffected
