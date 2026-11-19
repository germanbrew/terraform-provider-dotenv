package dotenv

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var (
	// ErrInvalidKey is returned by Encode for keys the parser can not read back.
	ErrInvalidKey = errors.New("invalid key")
	// ErrUnencodableValue is returned by Encode for values that can not be
	// written in a form the parser reads back unchanged.
	ErrUnencodableValue = errors.New("value can not be encoded")
)

//nolint:gochecknoglobals // compiled once, read-only
var (
	keyRe      = regexp.MustCompile(`\A[\w.]+\z`)
	plainValRe = regexp.MustCompile(`\A[A-Za-z0-9_./:@%+,=-]+\z`)
	dqEscaper  = strings.NewReplacer(`\`, `\\`, `"`, `\"`, `$`, `\$`, "\n", `\n`, "\r", `\r`)
)

// Encode renders entries in the .env format, one `KEY=value` line per entry,
// sorted by key. Values are written unquoted when possible, otherwise single
// or double quoted. Every line is verified to parse back to the same value.
func Encode(entries map[string]string) (string, error) {
	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	var b strings.Builder

	for _, key := range keys {
		line, err := encodeLine(key, entries[key])
		if err != nil {
			return "", err
		}

		b.WriteString(line)
		b.WriteByte('\n')
	}

	return b.String(), nil
}

func encodeLine(key, value string) (string, error) {
	if !keyRe.MatchString(key) {
		return "", fmt.Errorf("%w %q: keys may only contain letters, digits, underscores and dots", ErrInvalidKey, key)
	}

	for _, candidate := range encodeCandidates(value) {
		line := key + "=" + candidate

		parsed, err := ParseWithOptions(line, Options{})
		if err == nil && len(parsed) == 1 {
			if got, ok := parsed[key]; ok && got == value {
				return line, nil
			}
		}
	}

	return "", fmt.Errorf("%w for key %q", ErrUnencodableValue, key)
}

// encodeCandidates returns the possible representations of value, simplest first.
func encodeCandidates(value string) []string {
	if value == "" {
		return []string{""}
	}

	candidates := make([]string, 0, 3)

	if plainValRe.MatchString(value) {
		candidates = append(candidates, value)
	}

	if !strings.ContainsAny(value, "'\r\n") {
		candidates = append(candidates, "'"+value+"'")
	}

	return append(candidates, `"`+dqEscaper.Replace(value)+`"`)
}
