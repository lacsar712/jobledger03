package asrun

import "testing"

func TestSkewAccept(t *testing.T) {
	if err := Enforce("e1", "playlist_ms=1000 aired_ms=1080", []string{"spot"}); err != nil {
		t.Fatal(err)
	}
}

func TestSkewReject(t *testing.T) {
	if err := Enforce("e1", "playlist_ms=1000 aired_ms=1600", []string{"spot"}); err == nil {
		t.Fatal("expected reject")
	}
}

func TestManualOK(t *testing.T) {
	if err := Enforce("e1", "playlist_ms=1000 aired_ms=1600", []string{"manual-ok"}); err != nil {
		t.Fatal(err)
	}
}
