package payloads

type RequestGithubOAuthCallback struct {
	State string `form:"state" validate:"required"`
	Code string `form:"code" validate:"required"`
}
