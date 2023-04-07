package common

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	if !strings.Contains(ContextKeyTraceID.String(), contextKeyPrefix) {
		t.Errorf("ContextKeyRequestId expect to contain %v", contextKeyPrefix)
	}
}
