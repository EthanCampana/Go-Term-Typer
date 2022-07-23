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

func (s *Stats) GetLongestStreak() int {
	return s.longestStreak
}

func (s *Stats) UpdateStreak(x *Stats) *Stats {
	if s.Streak > s.longestStreak {
		return &Stats{
			Wpm:           x.Wpm,
			Streak:        x.Streak,
			Mistakes:      x.Mistakes,
			longestStreak: x.Streak,
		}
	}
	return s

}
