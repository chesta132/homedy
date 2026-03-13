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
	shares, err := samba.AddShare(payload.Name, share)
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

func (h *Samba) GetOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	name := c.Param("name")

	shares, err := samba.ReadShare()
	if err != nil {
		rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
		return
	}

	share, ok := samba.FilterShares(shares)[name]
	if !ok {
		rp.Error(replylib.CodeNotFound, "share not found").FailJSON()
		return
	}

	rp.Success(share).OkJSON()
}

func (h *Samba) UpdateOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	name := c.Param("name")
	var payload samba.Share
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	shares, err := samba.UpdateShare(name, payload)
	if err != nil {
		rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
		return
	}
	rp.Success(samba.FilterShares(shares)).OkJSON()
}

func (h *Samba) DeleteOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	name := c.Param("name")

	shares, err := samba.DeleteShare(name)
	if err != nil {
		rp.Error(replylib.CodeServerError, err.Error()).FailJSON()
		return
	}
	rp.Success(samba.FilterShares(shares)).OkJSON()
}
