package awsops

import (
	"errors"
	"testing"
	"time"

	"github.com/wearetechnative/bmc/internal/config"
)

func TestSessionDuration(t *testing.T) {
	tests := []struct {
		name    string
		seconds int
		want    time.Duration
	}{
		{"unset falls back to one hour", 0, time.Hour},
		{"negative falls back to one hour", -300, time.Hour},
		{"below AWS minimum is raised", 60, MinSessionDuration},
		{"above AWS maximum is capped", 100000, MaxSessionDuration},
		{"in-range value is kept", 4 * 3600, 4 * time.Hour},
		{"AWS minimum is kept", 900, MinSessionDuration},
		{"AWS maximum is kept", 43200, MaxSessionDuration},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sessionDuration(config.ConsoleConfig{SessionDurationSeconds: tt.seconds})
			if got != tt.want {
				t.Errorf("sessionDuration(%d) = %v, want %v", tt.seconds, got, tt.want)
			}
		})
	}
}

// TestSessionDurationNeverSDKDefault guards the actual bug: the SDK substitutes
// a 15-minute DefaultDuration whenever Duration is left at zero, so
// sessionDuration must never return zero.
func TestSessionDurationNeverZero(t *testing.T) {
	for _, seconds := range []int{0, -1, 1, 900, 43200, 999999} {
		if got := sessionDuration(config.ConsoleConfig{SessionDurationSeconds: seconds}); got == 0 {
			t.Errorf("sessionDuration(%d) returned 0, which makes the SDK use its 15m default", seconds)
		}
	}
}

func TestIsDurationRejected(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "MaxSessionDuration exceeded",
			err:  errors.New("operation error STS: AssumeRole, https response error StatusCode: 400, api error ValidationError: The requested DurationSeconds exceeds the MaxSessionDuration set for this role."),
			want: true,
		},
		{
			name: "role chaining one hour limit",
			err:  errors.New("api error ValidationError: The requested DurationSeconds exceeds the 1 hour session limit for roles assumed by role chaining."),
			want: true,
		},
		{
			name: "unrelated failure",
			err:  errors.New("api error AccessDenied: not authorized to perform sts:AssumeRole"),
			want: false,
		},
		{
			name: "expired MFA session",
			err:  errors.New("api error ExpiredToken: The security token included in the request is expired"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDurationRejected(tt.err); got != tt.want {
				t.Errorf("isDurationRejected(%q) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
