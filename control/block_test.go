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

func TestBlock(t *testing.T) {
	err := doBlock(2)
	if err != nil {
		t.Fatal(err)
	}
	if !isBlocked(2) {
		t.Fatal("set failed")
	}
}
