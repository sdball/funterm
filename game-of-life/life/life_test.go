package life

import "testing"

func equalSets(a, b Set) bool {
	if len(a) != len(b) {
		return false
	}

	for c := range a {
		if !b.Contains(c) {
			return false
		}
	}

	return true
}

func TestNeighborsCount(t *testing.T) {
	ns := Neighbors(Cell{2, 2})
	if len(ns) != 8 {
		t.Fatalf("expected 8 neighbors, got %d", len(ns))
	}
}

func TestNeighborsAreUnique(t *testing.T) {
	ns := Neighbors(Cell{2, 2})
	seen := map[Cell]struct{}{}
	for _, n := range ns {
		if _, ok := seen[n]; ok {
			t.Fatalf("duplicate neighbor: %+v", n)
		}
	}
}

func TestNeighborsAreCorrect(t *testing.T) {
	ns := Neighbors(Cell{0, 0})
	want := map[Cell]bool{
		{-1, -1}: true, {0, -1}: true, {1, -1}: true,
		{-1, 0}: true, {1, 0}: true,
		{-1, 1}: true, {0, 1}: true, {1, 1}: true,
	}
	for _, n := range ns {
		if !want[n] {
			t.Fatalf("unexpected neighbor: %+v", n)
		}
	}
}

func TestBlinker(t *testing.T) {
	start := NewSet(Cell{1, 0}, Cell{1, 1}, Cell{1, 2})
	want := NewSet(Cell{0, 1}, Cell{1, 1}, Cell{2, 1})
	got := Step(start)
	if !equalSets(got, want) {
		t.Fatalf("blinker step mismatch\ngot: %+v\nwant: %+v", got, want)
	}
}

func TestBlock(t *testing.T) {
	start := NewSet(Cell{1, 1}, Cell{1, 2}, Cell{2, 1}, Cell{2, 2})
	want := start
	got := Step(start)
	if !equalSets(got, start) {
		t.Fatalf("block step mismatch\ngot: %+v\nwant: %+v", got, want)
	}
}
