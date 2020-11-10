package main

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Field1          map[[2]string]string
	Field2          map[Struct]string
	Sequel          *string
}
type Struct struct {
	Field1 string
	Field2 int
}

func main() {
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
		Field1: map[[2]string]string{
			[2]string{"aaa", "bbb"}: "ccc",
			[2]string{"111", "222"}: "333",
		},
		Field2: map[Struct]string{
			Struct{"aaa", 1}: "AAA",
			Struct{"bbb", 2}: "BBB",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapted Screenplay (Nomin.)",
			"Best Director (Nomin.)",
			"Best Picture (Nomin.)",
		},
	}
	Display("strangelove", strangelove)
}
