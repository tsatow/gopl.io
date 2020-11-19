package ex05

import (
	"encoding/json"
	"testing"
)

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func TestMarshal(t *testing.T) {
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr. Strangelove": "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley": "Peter Sellers",
			"Gen. Buck Turgidson": "George C. Scott",
			"Brig. Gen. Jack D. Ripper": "Sterling Hayden",
			`Maj. T.J. "King" Kong`: "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}
	b, err := Marshal(strangelove)
	if err != nil {
		t.Errorf("Marshal(%v), but got err: %v\n", strangelove, err)
	}
	result := Movie{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Errorf("json.Unmarshal(%s), but got err.: %v\n", string(b), err)
		return
	}
	if !movieEquals(result, strangelove) {
		t.Errorf("Marshal(%v) result:%v, but got:%v", strangelove, strangelove, result)
		return
	}
}

func movieEquals(m1, m2 Movie) bool {
	return m1.Title == m2.Title &&
		m1.Subtitle == m2.Subtitle &&
		m1.Year == m2.Year &&
		m1.Color == m2.Color &&
		stringMapEquals(m1.Actor, m2.Actor) &&
		stringsEquals(m1.Oscars, m2.Oscars) &&
		m1.Sequel == m2.Sequel
}

func stringsEquals(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}

	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func stringMapEquals(x, y map[string]string) bool {
	if len(x) != len(y) {
		return false
	}

	for k, xv := range x {
		if yv, ok := y[k]; !ok || xv != yv {
			return false
		}
	}
	return true
}