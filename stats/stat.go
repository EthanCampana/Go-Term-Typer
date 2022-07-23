package stats

type Stats struct {
	Wpm      int
	Streak   int
	Mistakes int
}

func New() *Stats {
	return &Stats{
		Wpm:      0,
		Streak:   0,
		Mistakes: 0,
	}
}
