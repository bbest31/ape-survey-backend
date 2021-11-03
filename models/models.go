package models

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

type SaveTokenRequest struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"token"`
}
