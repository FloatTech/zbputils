package midi

import (
	"os"
	"testing"
)

var (
	testInput = "CCGGAAGR FFEEDDCR GGFFEEDR GGFFEEDR CCGGAAGR FFEEDDCR"
)

func TestTxt2mid(t *testing.T) {
	Txt2mid(40, "test.mid", testInput)
}

func TestGetNumTracks(t *testing.T) {
	data, err := os.ReadFile("test.mid")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(GetNumTracks(data))
}

func TestMid2txt(t *testing.T) {
	data, err := os.ReadFile("test.mid")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(Mid2txt(data, 0))
}

func TestO(t *testing.T) {
	t.Log("O(0, 5):", O(0, 5))
	t.Log("O(1, 6):", O(1, 6))
}

func TestName(t *testing.T) {
	t.Log("Name(60)", Name(60))
	t.Log("Name(64)", Name(64))
}

func TestPitch(t *testing.T) {
	t.Log("Pitch(\"C5\"):", Pitch("C5"))
	t.Log("Pitch(\"C#6\")", Pitch("C#6"))
}
