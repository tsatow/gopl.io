package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCharCount(t *testing.T) {
	var tests = []struct{
		input string
		want string
	}{
		// やはりcharCountをちゃんとリファクタしないとつらい...
		{
			"あ",
			"rune	count\n" +
			"'あ'	1\n" +
			"\nlen	count\n" +
			"1	0\n" +
			"2	0\n" +
			"3	1\n",
		},
		{
			"a",
			"rune	count\n" +
			"'a'	1\n" +
			"\nlen	count\n" +
			"1	1\n" +
			"2	0\n" +
			"3	0\n",
		},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)
		w := new(bytes.Buffer)

		charCount(r, w)
		out := w.String()

		if test.want != out {
			t.Errorf("charCount(\"%s\")\nexpected:\n%s\n, but got:\n%s\n", test.input, test.want, out)
		}
	}
}