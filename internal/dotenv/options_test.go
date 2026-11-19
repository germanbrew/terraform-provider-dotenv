package dotenv_test

import (
	"errors"
	"testing"

	"github.com/germanbrew/terraform-provider-dotenv/internal/dotenv"
)

func TestParseWithOptions_NoLookupEnv(t *testing.T) {
	t.Setenv("DOTENV_TEST_HOST_VAR", "from-host")

	env, err := dotenv.ParseWithOptions(`
LOCAL=local
A=$LOCAL
B=${DOTENV_TEST_HOST_VAR}
C=${DOTENV_TEST_HOST_VAR:-fallback}
D="prefix-${LOCAL}"
`, dotenv.Options{})
	if err != nil {
		t.Fatal(err)
	}

	envShouldContain(t, env, "A", "local")
	envShouldContain(t, env, "B", "")
	envShouldContain(t, env, "C", "fallback")
	envShouldContain(t, env, "D", "prefix-local")
}

func TestParseWithOptions_CustomLookupEnv(t *testing.T) {
	lookup := func(key string) (string, bool) {
		if key == "KNOWN" {
			return "value", true
		}

		return "", false
	}

	env, err := dotenv.ParseWithOptions("A=$KNOWN\nB=${UNKNOWN:-default}\n", dotenv.Options{LookupEnv: lookup})
	if err != nil {
		t.Fatal(err)
	}

	envShouldContain(t, env, "A", "value")
	envShouldContain(t, env, "B", "default")
}

func TestParse_UsesProcessEnv(t *testing.T) {
	t.Setenv("DOTENV_TEST_HOST_VAR", "from-host")

	env := mustParse(t, "A=$DOTENV_TEST_HOST_VAR\n")
	envShouldContain(t, env, "A", "from-host")
}

func TestParseWithOptions_Base(t *testing.T) {
	base := map[string]string{"FIRST": "1", "OVERRIDE": "old"}

	env, err := dotenv.ParseWithOptions("OVERRIDE=new\nREF=${FIRST}-${OVERRIDE}\n", dotenv.Options{Base: base})
	if err != nil {
		t.Fatal(err)
	}

	envShouldContain(t, env, "FIRST", "1")
	envShouldContain(t, env, "OVERRIDE", "new")
	envShouldContain(t, env, "REF", "1-new")

	if base["OVERRIDE"] != "old" {
		t.Error("base map must not be modified")
	}
}

func TestParseWithOptions_Errors(t *testing.T) {
	if _, err := dotenv.ParseWithOptions("this is invalid", dotenv.Options{}); !errors.Is(err, dotenv.ErrInvalidLine) {
		t.Errorf("expected ErrInvalidLine, got %v", err)
	}

	if _, err := dotenv.ParseWithOptions("A=\"unclosed\nB=1", dotenv.Options{}); !errors.Is(err, dotenv.ErrUnclosedQuote) {
		t.Errorf("expected ErrUnclosedQuote, got %v", err)
	}
}
