package input

import "time"

type CreateFeedingsInput struct {
	Name      string
	Interval  int
	StartDate time.Time
	StartIx   int
}
