package payloads

import "homedy/internal/models"

type RequestSignUp struct {
	Username string `json:"username" validate:"required,username" example:"chesta_ardiona"`
	Email    string `json:"email" validate:"required,email" example:"chestaardi4@gmail.com"`
	Password string `json:"password" validate:"required,password" example:"YourPassword123"`
}

func (p *RequestSignUp) ToUser() models.User {
	return models.User{
		Username: p.Username,
		Email:    p.Email,
		Password: p.Password,
	}
}

type RequestSignIn struct {
	// email or username
	Identifier string `json:"identifier" validate:"required" example:"chesta_ardiona"`
	Password   string `json:"password" validate:"required,password" example:"YourPassword123"`
	RememberMe bool   `json:"remember_me"`
}

type ApprovalAction string

const (
	ApprovalApprove ApprovalAction = "approve"
	ApprovalDeny    ApprovalAction = "deny"
)

type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalDenied   ApprovalStatus = "denied"
)

type RequestSignUpApproval struct {
	Identifier string         `form:"identifier" validate:"required"`
	Action     ApprovalAction `form:"action" validate:"required"`
}

type RequestSignUpApprovalStatus struct {
	Email string `form:"email" validate:"required,email"`
}

type ResponseSignUpApprovalStatus struct {
	Username string         `json:"username,omitempty"`
	Email    string         `json:"email"`
	Status   ApprovalStatus `json:"status"`
}
