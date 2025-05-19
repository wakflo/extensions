package shared

type SurveyMonkeySurveysResponse struct {
	Data []SurveyMonkeySurvey `json:"data"`
}

type SurveyMonkeySurvey struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}
