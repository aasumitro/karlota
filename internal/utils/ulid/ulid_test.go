package ulid_test

import (
	"strings"
	"testing"

	"karlota.aasumitro.id/internal/utils/ulid"
)

func TestGenerate(t *testing.T) {
	prefix := "test"
	id := ulid.Generate(prefix)
	if !strings.HasPrefix(id, prefix) {
		t.Errorf("Expected id to start with '%s', got '%s'", prefix, id)
	}
	if len(id) != len(prefix)+33 { // 34: length of "_01234567_89AB_CDEF_0123_4567" part
		t.Errorf("Unexpected length of id: got %v, want %v", len(id), len(prefix)+34)
	}
	id2 := ulid.Generate(prefix)
	if id == id2 {
		t.Error("Two generated ids are identical, they should be different")
	}
}
