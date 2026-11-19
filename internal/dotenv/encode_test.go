package dotenv_test

import (
	"errors"
	"testing"

	"github.com/germanbrew/terraform-provider-dotenv/internal/dotenv"
)

func TestEncode_Format(t *testing.T) {
	got, err := dotenv.Encode(map[string]string{
		"PLAIN":  "value-1.2_3",
		"EMPTY":  "",
		"SPACES": "hello world",
		"QUOTE":  "it's",
		"NL":     "line1\nline2",
		"DOLLAR": "$HOME",
	})
	if err != nil {
		t.Fatal(err)
	}

	want := `DOLLAR='$HOME'
EMPTY=
NL="line1\nline2"
PLAIN=value-1.2_3
QUOTE="it's"
SPACES='hello world'
`
	if got != want {
		t.Errorf("unexpected output:\n%s\nwant:\n%s", got, want)
	}
}

func TestEncode_RoundTrip(t *testing.T) {
	values := []string{
		"",
		"simple",
		"with spaces",
		"  leading and trailing  ",
		"tab\there",
		"hash # not a comment",
		"#leading-hash",
		`double "quotes"`,
		`single 'quotes'`,
		`both 'single' and "double"`,
		"$VAR ${OTHER} ${X:-y}",
		`it's $HOME`,
		`back\slash`,
		`trailing backslash\`,
		`literal \n sequence`,
		"multi\nline\r\nvalue",
		`{"json": [1, 2, 3]}`,
		"unicode äöü ✓",
		"=leading-equals",
		"a=b=c",
	}

	for _, value := range values {
		encoded, err := dotenv.Encode(map[string]string{"KEY": value})
		if err != nil {
			t.Errorf("Encode(%q): %v", value, err)

			continue
		}

		parsed, err := dotenv.ParseWithOptions(encoded, dotenv.Options{})
		if err != nil {
			t.Errorf("Parse(%q): %v", encoded, err)

			continue
		}

		if parsed["KEY"] != value {
			t.Errorf("round trip of %q via %q returned %q", value, encoded, parsed["KEY"])
		}
	}
}

func TestEncode_InvalidKey(t *testing.T) {
	for _, key := range []string{"", "WITH SPACE", "DASH-KEY", "A=B", "ÄÖÜ"} {
		if _, err := dotenv.Encode(map[string]string{key: "v"}); !errors.Is(err, dotenv.ErrInvalidKey) {
			t.Errorf("key %q: expected ErrInvalidKey, got %v", key, err)
		}
	}
}

func FuzzEncodeRoundTrip(f *testing.F) {
	for _, seed := range []string{"", "a b", "it's", `"`, "\\n", "$X", "a\nb"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, value string) {
		encoded, err := dotenv.Encode(map[string]string{"KEY": value})
		if err != nil {
			if !errors.Is(err, dotenv.ErrUnencodableValue) {
				t.Fatalf("unexpected error: %v", err)
			}

			return
		}

		parsed, err := dotenv.ParseWithOptions(encoded, dotenv.Options{})
		if err != nil || parsed["KEY"] != value {
			t.Fatalf("round trip of %q via %q returned %q (%v)", value, encoded, parsed["KEY"], err)
		}
	})
}
