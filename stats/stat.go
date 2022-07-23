package stats

type Stats struct {
	Wpm           int
	Streak        int
	Mistakes      int
	longestStreak int
}

func New() *Stats {
	return &Stats{
		Wpm:           0,
		Streak:        0,
		Mistakes:      0,
		longestStreak: 0,
	}
}

func (s *Stats) UpdateStreak() *Stats {
	if s.Streak > s.longestStreak {
		return &Stats{
			Wpm:           s.Wpm,
			Streak:        s.Streak,
			Mistakes:      s.Mistakes,
			longestStreak: s.Streak,
		}
	}
	return s

}
