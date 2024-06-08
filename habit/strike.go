package habit

import "time"

type StrikeCountType int

const (
	StrikeTypeInf StrikeCountType = iota
	StrikeTypeMonthly
	StrikeTypeWeekly
)

type Strike struct {
	Type         StrikeCountType
	Count        int
	Best         int
	LastFinished time.Time
}

func NewStrike(t StrikeCountType) Strike {
	return Strike{Type: t}
}

// Update updates strike data upon completion.
func (s *Strike) Update(completeDate time.Time) {
	if AreSameDates(completeDate, s.LastFinished) {
		return
	}

	if !s.IsContinued(completeDate) {
		s.Count = 1
	} else if AreSameDates(completeDate, s.LastFinished.AddDate(0, 0, 1)) {
		s.Count++
	}

	if s.Count > s.Best {
		s.Best = s.Count
	}

	s.LastFinished = completeDate
}

// Downgrade updates strike data when task is unmarked from completion.
func (s *Strike) Downgrade(prevCompletion time.Time) {
	s.Count--
	s.LastFinished = prevCompletion
}

// IsContinued returns whether strike is broken assuming date is today
// meaning there was over 1 day break from finishing it.
func (s Strike) IsContinued(date time.Time) bool {
	out := AreSameDates(s.LastFinished, date) ||
		AreSameDates(date.AddDate(0, 0, -1), s.LastFinished)

	switch s.Type {
	case StrikeTypeInf:
		break
	case StrikeTypeMonthly:
		out = out && date.Month() == s.LastFinished.Month()
	case StrikeTypeWeekly:
		_, dateWeek := date.ISOWeek()
		_, lastFinishWeek := s.LastFinished.ISOWeek()
		out = out && dateWeek == lastFinishWeek
	}

	return out
}
