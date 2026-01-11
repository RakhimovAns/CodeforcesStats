package model

type ContestPhase string

const (
	ContestPhaseBefore   ContestPhase = "BEFORE"
	ContestPhaseCoding   ContestPhase = "CODING"
	ContestPhaseFinished ContestPhase = "FINISHED"
)

type Contest struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	StartTimeSeconds int64   `json:"startTimeSeconds"`
	Phase            string  `json:"phase"`
	WebsiteURL       *string `json:"websiteUrl,omitempty"`
}
