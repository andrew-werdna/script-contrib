package filters

import (
	"fmt"
	"testing"
	"unicode"

	"github.com/bitfield/script"
)

func TestCutFirstN(t *testing.T) {
	sha := "b122cf37282e2dc93d2cf2d89e825911bc0487dc"
	want := sha[:8]
	got, err := script.
		Echo(sha).
		Filter(CutFirstN(8)).
		String()

	if err != nil {
		t.Errorf("expected nil error but got %v", err)
	}

	if want != got {
		t.Errorf("wanted: '%s', but got: '%s'", want, got)
	}
}

func TestColumns(t *testing.T) {
	const test1 string = `1,2,3,4
alpha,bravo,charlie,delta`
	const test2 string = `1   2   3   4
alpha    bravo    charlie    delta`

	t.Run("works on CSV data", func(t *testing.T) {
		want := fmt.Sprintf("2 3\nbravo charlie\n")
		fn := func(c rune) bool {
			return c == ','
		}
		got, err := script.
			Echo(test1).
			FilterScan(Columns(fn, " ", 2, 3)).
			String()
		if err != nil {
			t.Errorf("expected nil error but got %v", err)
		}

		if got != want {
			t.Errorf("wanted '%v', but got '%v'", want, got)
		}
	})

	t.Run("CSV data with bad columns", func(t *testing.T) {
		want := fmt.Sprintf("2 3\nbravo charlie\n")
		fn := func(c rune) bool {
			return c == ','
		}
		got, err := script.
			Echo(test1).
			FilterScan(Columns(fn, " ", 2, 3, 5)).
			String()
		if err != nil {
			t.Errorf("expected nil error but got %v", err)
		}

		if got != want {
			t.Errorf("wanted '%v', but got '%v'", want, got)
		}
	})

	t.Run("doesn't explode with no desired columns", func(t *testing.T) {
		var want string
		fn := func(c rune) bool {
			return c == ','
		}
		got, err := script.
			Echo(test1).
			FilterScan(Columns(fn, " ")).
			String()
		if err != nil {
			t.Errorf("expected nil error but got %v", err)
		}

		if got != want {
			t.Errorf("wanted '%v', but got '%v'", want, got)
		}

	})

	t.Run("doesn't fail with bad column index", func(t *testing.T) {
		var want string
		fn := func(c rune) bool {
			return c == ','
		}
		got, err := script.
			Echo(test1).
			FilterScan(Columns(fn, " ", 0)).
			String()
		if err != nil {
			t.Errorf("expected nil error but got %v", err)
		}

		if got != want {
			t.Errorf("wanted '%v', but got '%v'", want, got)
		}

	})

	t.Run("works on TSV data", func(t *testing.T) {
		want := fmt.Sprintf("2\t3\nbravo\tcharlie\n")
		fn := func(c rune) bool {
			return unicode.IsSpace(c)
		}
		got, err := script.
			Echo(test2).
			FilterScan(Columns(fn, "\t", 2, 3)).
			String()
		if err != nil {
			t.Errorf("expected nil error but got %v", err)
		}

		if got != want {
			t.Errorf("wanted '%v', but got '%v'", want, got)
		}
	})

	t.Run("TSV data with bad columns", func(t *testing.T) {
		want := fmt.Sprintf("2\t3\nbravo\tcharlie\n")
		fn := func(c rune) bool {
			return unicode.IsSpace(c)
		}
		got, err := script.
			Echo(test2).
			FilterScan(Columns(fn, "\t", 2, 3, 6)).
			String()
		if err != nil {
			t.Errorf("expected nil error but got %v", err)
		}

		if got != want {
			t.Errorf("wanted '%v', but got '%v'", want, got)
		}
	})
}
