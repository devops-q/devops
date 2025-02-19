package json_models

type RegisterUserBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pwd      string `json:"pwd"`
}

type GetFollowsResponse struct {
	Follows []string `json:"follows"`
}
