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

	share := payload.ToShare()
	// TODO: not only add conf but add directory with permission
	shares, err := samba.AddShare(payload.Name, share)
	// TODO: reset smb
	if err != nil {
		rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
		return
	}
	rp.Success(samba.FilterShares(shares)).CreatedJSON()
}

func (h *Samba) GetAll(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	shares, err := samba.ReadShare()
	if err != nil {
		rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
		return
	}
	rp.Success(samba.FilterShares(shares)).OkJSON()
}
