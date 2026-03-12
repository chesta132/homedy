package handlers

import (
	"homedy/internal/libs/replylib"
	"homedy/internal/payloads"
	"homedy/internal/services/samba"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Samba struct {
}

func NewSamba() *Samba {
	return &Samba{}
}

func (h *Samba) AddShare(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload payloads.RequestAddShare
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	share := samba.Share{
		Path:       payload.Path,
		ReadOnly:   payload.ReadOnly,
		Browsable:  payload.Browsable,
		GuestUsers: payload.GuestUsers,
		AdminUsers: payload.AdminUsers,
	}
	// TODO: not only add conf but add directory with permission
	shares, err := samba.AddConf(payload.Name, share)
	if err != nil {
		rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
		return
	}
	rp.Success(samba.FilterShares(shares)).CreatedJSON()
}

func (h *Samba) GetAll(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	shares, err := samba.ReadConf()
	if err != nil {
		rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
		return
	}
	rp.Success(samba.FilterShares(shares)).OkJSON()
}
