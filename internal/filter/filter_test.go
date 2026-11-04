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

func TestMatch_multiTokenAND(t *testing.T) {
	hosts := []string{"sk1-ex-etapp", "sk1-ex-tapp", "msk2-ex-web03"}

	// Every token must be present for a host to match.
	got := Match("sk1 tapp", hosts)
	if len(got) != 2 || got[0] != "sk1-ex-etapp" || got[1] != "sk1-ex-tapp" {
		t.Fatalf("Match(sk1 tapp) = %v, want both sk1-ex-etapp and sk1-ex-tapp", got)
	}

	// Extra whitespace should behave the same.
	got = Match("  sk1   tapp  ", hosts)
	if len(got) != 2 {
		t.Fatalf("Match with extra whitespace = %v, want 2 hosts", got)
	}

	// A token only one host contains narrows the list down.
	got = Match("sk1 etapp", hosts)
	if len(got) != 1 || got[0] != "sk1-ex-etapp" {
		t.Fatalf("Match(sk1 etapp) = %v, want [sk1-ex-etapp]", got)
	}
}

