package paths

import (
	"strings"
	"os"
	"testing"

)

func TestDataPath(t *testing.T) {
	expected := "data"
	path := DataPath()
	path = strings.ReplaceAll(path, "\\", "/")
	dirnodes := strings.Split(path, "/")
	dirname := dirnodes[len(dirnodes)-1]
	os.Stat(path)
	if dirname != expected {
		t.Errorf("Expected %s, got %s", expected, dirname)
	}
}
