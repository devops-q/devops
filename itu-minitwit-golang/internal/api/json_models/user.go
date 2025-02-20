package json_models

type RegisterUserBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pwd      string `json:"pwd"`
}

type GetFollowsResponse struct {
	Follows []string `json:"follows"`
}

type FollowUnfollowBody struct {
	Follow   *string `json:"follow,omitempty"`
	Unfollow *string `json:"unfollow,omitempty"`
}
