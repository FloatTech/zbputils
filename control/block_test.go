package control

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	err := os.MkdirAll("data/control", 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = initBlock()
	if err != nil {
		t.Fatal(err)
	}
}
