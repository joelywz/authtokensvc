package authtokensvc

type IssueResponse struct {
	RefreshToken string
	AccessToken  string
}

type RefreshResponse = IssueResponse
