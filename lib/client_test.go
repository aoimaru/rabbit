package lib

import (
	"testing"
)

func TestClient(t *testing.T) {
	client := CreateClient()

	if client.WorkPath != "/home/aoimaru/document/go_project/CLI/rabbit/lib" {
		t.Errorf("Unexpected WorkPath value: %s", client.WorkPath)
	}

}
