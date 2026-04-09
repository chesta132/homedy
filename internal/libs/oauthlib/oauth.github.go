package oauthlib

import (
	"homedy/config"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var GithubOAuthConfig = &oauth2.Config{
	ClientID:     config.GITHUB_OAUTH_CLIENT_ID,
	ClientSecret: config.GITHUB_OAUTH_CLIENT_SECRET,
	Scopes:       []string{"repo", "read:user", "admin:repo_hook"},
	Endpoint:     github.Endpoint,
	RedirectURL:  (func() string { u, _ := url.JoinPath(config.APP_URL, "/oauth/github/callback"); return u })(),
}
