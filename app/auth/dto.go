package auth

import "time"

type finishOnboardingReq struct {
	FullName     string    `json:"full_name"`
	Nickname     string    `json:"nickname"`
	Gender       int       `json:"gender"`
	Struggles    string    `json:"struggles"`
	SubjectNames []string  `json:"subject_names"`
	BirthDate    time.Time `json:"birth_date"`
	// Avatar    string    `json:"avatar"`
	// Status    UserStatus  `json:"status"`
	// UpdatedAt null.Time   `json:"updated_at"`
}

type addSubjectToSelfReq struct {
	// Subjects []Subject `json:"subjects"`
	SubjectNames []string `json:"subject_names"`
}

type addSubjectsReq struct {
	SubjectNames []string `json:"subject_names"`
}
