package rbac

import "testing"

func TestNewRBACServiceFromEnvRequiresConfigEnv(t *testing.T) {
	_, err := NewRBACServiceFromEnv()
	if err == nil {
		t.Fatal("expected error when RBAC_CONFIG_PATH is not configured")
	}
}
