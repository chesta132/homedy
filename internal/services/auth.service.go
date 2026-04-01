package services

import (
	"context"
	"errors"
	"fmt"
	"homedy/config"
	"homedy/internal/libs/authlib"
	"homedy/internal/libs/cryptolib"
	"homedy/internal/libs/dblib"
	"homedy/internal/libs/mail"
	"homedy/internal/libs/replylib"
	"homedy/internal/models"
	"homedy/internal/models/payloads"
	"homedy/internal/repos"
	"net/http"
	"net/url"
	"time"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Auth struct {
	userRepo   *repos.User
	revokeRepo *repos.Revoke
	mailer     *mail.Mailer
}

type ContextedAuth struct {
	Auth
	c   *gin.Context
	ctx context.Context
}

func NewAuth(userRepo *repos.User, revokeRepo *repos.Revoke, mailer *mail.Mailer) *Auth {
	return &Auth{userRepo, revokeRepo, mailer}
}

func (s *Auth) AttachContext(c *gin.Context) *ContextedAuth {
	return &ContextedAuth{*s, c, c.Request.Context()}
}

func (s *Auth) sendSignUpApproval(user models.User) error {
	identifier, err := cryptolib.EncryptGCM([]byte(user.ID), []byte(config.SIGNUP_IDENTIFIER_KEY))
	if err != nil {
		return err
	}
	identifier = url.QueryEscape(identifier)

	html := mail.ApprovalEmail(mail.ApprovalData{
		Title:    "Sign Up Request Approval",
		Subtitle: fmt.Sprintf("%s (%s) request to create account on %s. Please sign your approval", user.Username, user.Email, config.APP_NAME),
		Fields: []mail.InfoField{
			{Label: "Requested by", Value: user.Username},
			{Label: "Requested at", Value: time.Now().Format(time.RFC822)},
		},
		ApproveHref: fmt.Sprintf("%s/signup/review-approval?identifier=%s&action=%s", config.FRONTEND_URL, identifier, url.QueryEscape(string(payloads.ApprovalApprove))),
		DenyHref:    fmt.Sprintf("%s/signup/review-approval?identifier=%s&action=%s", config.FRONTEND_URL, identifier, url.QueryEscape(string(payloads.ApprovalDeny))),
	})
	return s.mailer.Send(config.MAIL_OWNER, fmt.Sprintf("%s Sign Up Request", config.APP_NAME), html)
}

func (s *Auth) sendSignUpApprovalReviewed(user models.User, payload payloads.RequestSignUpApproval) error {
	title := fmt.Sprintf("%s Account Approved", config.APP_NAME)
	subtitle := fmt.Sprintf("Your %s account has been approved. You can now sign in and start using %s.", config.APP_NAME, config.APP_NAME)
	if payload.Action == payloads.ApprovalDeny {
		title = fmt.Sprintf("%s Account Request Denied", config.APP_NAME)
		subtitle = fmt.Sprintf("Your %s account request has been denied. If you think this was a mistake, please contact our support team.", config.APP_NAME)
	}

	html := mail.GeneralInfoEmail(mail.GeneralInfoData{
		Title:    title,
		Subtitle: subtitle,
		Fields: []mail.InfoField{
			{Label: "Requested by", Value: user.Username},
			{Label: "Requested at", Value: user.CreatedAt.Format(time.RFC822)},
		},
	})
	return s.mailer.Send(user.Email, fmt.Sprintf("%s Sign Up Request Has Been Reviewed", config.APP_NAME), html)
}

func (s *ContextedAuth) SignUp(payload payloads.RequestSignUp) error {
	// validate email and username
	email, username, err := s.userRepo.GetEmailOrUsername(s.ctx, payload.Email, payload.Password)
	isErrNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !isErrNotFound {
		return err
	}
	if !isErrNotFound {
		fe := make(reply.FieldsError)
		if email == payload.Email {
			fe["email"] = "email already registered"
		}
		if username == payload.Username {
			fe["username"] = "username already registered"
		}
		return &reply.ErrorPayload{
			Code:    replylib.CodeConflict,
			Message: "email or username already registered",
			Fields:  fe,
		}
	}

	// create pending user (hash in before create)
	newUser := payload.ToUser()
	if err := s.userRepo.Create(s.ctx, &newUser); err != nil {
		return dblib.GormErrorToReplyError(err, &newUser)
	}

	return s.sendSignUpApproval(newUser)
}

func (s *ContextedAuth) SignUpApproval(payload payloads.RequestSignUpApproval) (*models.User, error) {
	id, err := cryptolib.DecryptGCM([]byte(payload.Identifier), []byte(config.SIGNUP_IDENTIFIER_KEY))
	if err != nil {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeBadRequest,
			Message: "invalid identifier format",
		}
	}

	user, err := s.userRepo.GetFirst(s.ctx, "id = ? AND status = ?", id, models.UserPending)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &reply.ErrorPayload{
			Code:    replylib.CodeUnprocessableEntity,
			Message: "token already used",
		}
	}
	if err != nil {
		return nil, err
	}

	if payload.Action == payloads.ApprovalApprove {
		if err = s.userRepo.UpdateByID(s.ctx, id, models.User{Status: models.UserActive}); err != nil {
			return nil, err
		}
	} else {
		if _, err = s.userRepo.DeleteByID(s.ctx, id); err != nil {
			return nil, err
		}
	}

	return &user, s.sendSignUpApprovalReviewed(user, payload)
}

func (s *ContextedAuth) SignUpApprovalStatus(payload payloads.RequestSignUpApprovalStatus) (*payloads.ResponseSignUpApprovalStatus, error) {
	user, err := s.userRepo.GetFirst(s.ctx, "email = ?", payload.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &payloads.ResponseSignUpApprovalStatus{
			Email:  payload.Email,
			Status: payloads.ApprovalDenied,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	res := payloads.ResponseSignUpApprovalStatus{
		Username: user.Username,
		Email:    user.Email,
		Status:   payloads.ApprovalPending,
	}
	if user.Status == models.UserActive {
		res.Status = payloads.ApprovalApproved
	}

	return &res, nil
}

func (s *ContextedAuth) SignIn(payload payloads.RequestSignIn) (*models.User, []http.Cookie, error) {
	user, err := s.userRepo.GetFirst(s.ctx, "email = ? OR username = ?", payload.Identifier, payload.Identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, &reply.ErrorPayload{
				Code:    replylib.CodeNotFound,
				Message: "user not found",
				Fields: reply.FieldsError{
					"identifier": "email or username not found",
				},
			}
		}
		return nil, nil, err
	}

	if user.Status != models.UserActive {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeForbidden,
			Message: "user inactive",
			Fields: reply.FieldsError{
				"identifier": "user is not actived yet",
			},
		}
	}

	if !authlib.ComparePassword(payload.Password, user.Password) {
		return nil, nil, &reply.ErrorPayload{
			Code:    replylib.CodeUnauthorized,
			Message: "password is incorrect",
			Fields: reply.FieldsError{
				"password": "password is incorrect",
			},
		}
	}

	return &user, authlib.CreateTokenCookie(user.ID, payload.RememberMe), nil
}

func (s *ContextedAuth) SignOut() []http.Cookie {
	refresh, _ := s.c.Request.Cookie(config.REFRESH_TOKEN_KEY)
	_ = s.revokeRepo.RevokeToken(s.ctx, refresh.Value, "user already sign out")
	return authlib.InvalidateTokenCookie()
}

func (s *ContextedAuth) Me() (models.User, error) {
	userID, _ := s.c.Get("userID")
	return s.userRepo.GetByID(s.ctx, userID.(string))
}
