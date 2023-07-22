package subject

type addSubjectToSelfReq struct {
	// Subjects []Subject `json:"subjects"`
	SubjectNames []string `json:"subject_names"`
}

type addSubjectsReq struct {
	SubjectNames []string `json:"subject_names"`
}
