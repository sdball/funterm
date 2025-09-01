package life

type Cell struct{ X, Y int }

type Set map[Cell]struct{}

func NewSet(cells ...Cell) Set {
	s := make(Set, len(cells))
	for _, c := range cells {
		s[c] = struct{}{}
	}
	return s
}

func (s Set) Contains(c Cell) bool {
	_, ok := s[c]
	return ok
}

func (s Set) Add(c Cell) { s[c] = struct{}{} }

func (s Set) Remove(c Cell) { delete(s, c) }

func Neighbors(c Cell) []Cell {
	return []Cell{
		// top row
		{c.X - 1, c.Y - 1},
		{c.X, c.Y - 1},
		{c.X + 1, c.Y - 1},

		// middle row
		{c.X - 1, c.Y},
		/* {c.X, c.Y}, <-- this is us, not a neighbor */
		{c.X + 1, c.Y},

		// bottom row
		{c.X - 1, c.Y + 1},
		{c.X, c.Y + 1},
		{c.X + 1, c.Y + 1},
	}
}

func Step(cells Set) Set {
	countOfNeighbors := make(map[Cell]int, len(cells)*8)
	for c := range cells {
		for _, n := range Neighbors(c) {
			countOfNeighbors[n]++
		}
	}

	// guess at the capacity of the next Set
	next := make(Set, len(countOfNeighbors))

	for cell, n := range countOfNeighbors {
		if n == 3 || (n == 2 && cells.Contains(cell)) {
			next.Add(cell)
		}
	}
	return next
}
