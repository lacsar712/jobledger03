package asrun

import (
	"fmt"
	"strconv"
	"strings"
)

type Rec struct {
	Title, Body string
	Tags        []string
}

const MaxSkewMS = 150

func Sample() Rec {
	return Rec{Title: "spot-auto-1820", Body: "playlist_ms=1820000 aired_ms=1820080", Tags: []string{"spot"}}
}

func Seed() []Rec {
	return []Rec{
		Sample(),
		{Title: "ident-top", Body: "playlist_ms=0 aired_ms=40", Tags: []string{"ident"}},
	}
}

func AfterWrite(getMin func() (string, error), setMin func(string) error, body string) error {
	return nil
}

func Steps() []string { return []string{"ingest", "match", "discrepancy"} }

func Enforce(title, body string, tags []string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("event title required")
	}
	pl, air, err := parseMS(body)
	if err != nil {
		return err
	}
	skew := pl - air
	if skew < 0 {
		skew = -skew
	}
	manual := false
	for _, t := range tags {
		if t == "manual-ok" {
			manual = true
		}
	}
	if skew > MaxSkewMS && !manual {
		return fmt.Errorf("as-run skew %dms exceeds %dms", skew, MaxSkewMS)
	}
	return nil
}

func parseMS(body string) (playlist, aired int, err error) {
	gotP, gotA := false, false
	for _, part := range strings.Fields(body) {
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		n, conv := strconv.Atoi(v)
		if conv != nil {
			return 0, 0, fmt.Errorf("bad %s", part)
		}
		switch k {
		case "playlist_ms":
			playlist, gotP = n, true
		case "aired_ms":
			aired, gotA = n, true
		}
	}
	if !gotP || !gotA {
		return 0, 0, fmt.Errorf("body must contain playlist_ms and aired_ms")
	}
	return playlist, aired, nil
}
