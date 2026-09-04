package filter

import "testing"

func TestMatch_substringAndCase(t *testing.T) {
	hosts := []string{"nsk1-ex-tav01", "msk2-ex-web03", "nsk1-ex-db02"}

	got := Match("tav", hosts)
	if len(got) != 1 || got[0] != "nsk1-ex-tav01" {
		t.Fatalf("Match(tav) = %v, want [nsk1-ex-tav01]", got)
	}

	got = Match("NSK1", hosts)
	if len(got) != 2 {
		t.Fatalf("Match(NSK1) len = %d, want 2 (%v)", len(got), got)
	}
}

func TestMatch_emptyQueryShowsNothing(t *testing.T) {
	hosts := []string{"a", "b"}
	got := Match("", hosts)
	if len(got) != 0 {
		t.Fatalf("empty query should return no hosts, got %v", got)
	}
	got = Match("   ", hosts)
	if len(got) != 0 {
		t.Fatalf("whitespace query should return no hosts, got %v", got)
	}
}

func TestMatch_shrinksAsQueryGrows(t *testing.T) {
	hosts := []string{"nsk1-ex-tav01", "nsk1-ex-tav02", "msk2-ex-web03"}

	got := Match("nsk", hosts)
	if len(got) != 2 {
		t.Fatalf("after 'nsk': %v", got)
	}
	got = Match("nsk1-ex-tav01", hosts)
	if len(got) != 1 || got[0] != "nsk1-ex-tav01" {
		t.Fatalf("after full name: %v", got)
	}
}
