# Proposal: Fix Console Service URL Generation

bean: [bmc-drrn](../../../.beans/bmc-drrn--console-service-error.md)

## Why

`bmc console -s <service>` generates federation URLs with a `Destination` of `https://<service>.console.aws.amazon.com/`, which causes a 400 Bad Request for most services because they do not have a dedicated subdomain. The correct console URL format is path-based (`https://console.aws.amazon.com/<service>/`) or region-prefixed, both of which the AWS federation endpoint accepts.

## What Changes

- The `-s` flag value is treated as a **console path** rather than a bare service name, enabling both short names (`rds`) and sub-paths (`systems-manager/parameters`)
- URL generation for the `Destination` parameter switches from the broken subdomain format to a region-aware path-based format: `https://<region>.console.aws.amazon.com/<service-path>/home` (or the path as-is when it already contains `/`)
- Region is resolved from the selected AWS profile (already available in the codebase via `awsCfg.Region`)
- No canonical service mapping is needed; the path-based fallback works universally

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `aws-console-access`: The `-s` flag now accepts a console path (not just a service name), and URL generation uses region-aware path-based URLs instead of the broken subdomain format

## Impact

- `internal/awsops/console.go`: `buildConsoleURL` receives a `region` parameter; URL construction logic changes
- `openspec/specs/aws-console-access/spec.md`: Update `-s` flag requirements to reflect path-based input and region-aware URL construction
