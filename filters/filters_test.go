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
	const input1 string = `1,2,3,4
alpha,bravo,charlie,delta`
	const input2 string = `1   2   3   4
alpha    bravo    charlie    delta`
	const input3 string = `1,2,3,4,5,6,7,8,9,10,11,12
alpha,beta,gamma,delta,epsilon,zeta,eta,theta,iota,kappa,la,mu`

	t.Run("works on CSV data", func(t *testing.T) {
		want := fmt.Sprintf("2 3\nbravo charlie\n")
		fn := func(c rune) bool {
			return c == ','
		}
		got, err := script.
			Echo(input1).
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
			Echo(input1).
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
			Echo(input1).
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
			Echo(input1).
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
			Echo(input2).
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
			Echo(input2).
			FilterScan(Columns(fn, "\t", 2, 3, 6)).
			String()
		if err != nil {
			t.Errorf("expected nil error but got %v", err)
		}

		if got != want {
			t.Errorf("wanted '%v', but got '%v'", want, got)
		}
	})

	t.Run("allows any order columns", func(t *testing.T) {
		want := fmt.Sprintf("5|2|8|3|9\nepsilon|beta|theta|gamma|iota\n")
		fn := func(c rune) bool {
			return c == ','
		}
		got, err := script.
			Echo(input3).
			FilterScan(Columns(fn, "|", 5, 2, 8, 3, 9)).
			String()
		if err != nil {
			t.Errorf("expected nil error but got %v", err)
		}

		if got != want {
			t.Errorf("wanted '%v', but got '%v'", want, got)
		}
	})

	t.Run("allows any order columns and ignores invalid columns", func(t *testing.T) {
		want := fmt.Sprintf("5|2|8|3|9\nepsilon|beta|theta|gamma|iota\n")
		fn := func(c rune) bool {
			return c == ','
		}
		got, err := script.
			Echo(input3).
			FilterScan(Columns(fn, "|", 5, 2, 17, 8, 3, 9, 21)).
			String()
		if err != nil {
			t.Errorf("expected nil error but got %v", err)
		}

		if got != want {
			t.Errorf("wanted '%v', but got '%v'", want, got)
		}
	})
}
