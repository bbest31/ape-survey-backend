package models

import "time"

type Survey struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Nickname string `json:"nickname"`
	Href     string `json:"href"`
}

type SurveysResponse struct {
	Data    []Survey `json:"data"`
	PerPage int      `json:"per_page"`
	Page    int      `json:"page"`
	Total   int      `json:"total"`
	Links   struct {
		Self string `json:"self"`
		Next string `json:"next"`
		Last string `json:"last"`
	} `json:"links"`
}

type SurveyDetailsResponse struct {
	Title         string
	Nickname      string
	Language      string
	FolderID      string
	Category      string
	QuestionCount int
	PageCount     int
	ResponseCount int
	DateCreated   time.Time
	DateModified  time.Time
	ID            string
	ButtonsText   struct {
		NextButton string
		PrevButton string
		DoneButton string
		ExitButton string
	}
	IsOwner         bool
	Footer          bool
	CustomVariables map[string]interface{}
	Href            string
	AnalyzeURL      string
	EditURL         string
	CollectURL      string
	SummaryURL      string
	Preview         string
	Pages           []struct {
		Title         string
		Description   string
		Position      int
		QuestionCount int
		ID            string
		Href          string
		Questions     []string
	}
}

type SaveTokenRequest struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"token"`
}
