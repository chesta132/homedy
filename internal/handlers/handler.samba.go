package handlers

import (
	"homedy/internal/libs/replylib"
	"homedy/internal/libs/validatorlib"
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

	if errPayload := validatorlib.ValidateStructToReply(payload); errPayload != nil {
		rp.Error(replylib.ErrorPayloadToErrorArg(*errPayload)).FailJSON()
		return
	}

	share := payload.ToShare()
	shares, err := samba.AddShare(payload.Name, share)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(samba.FilterShares(shares)).CreatedJSON()
}

func (h *Samba) GetAll(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	shares, err := samba.ReadShare()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(samba.FilterShares(shares)).OkJSON()
}

func (h *Samba) GetOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	name := c.Param("name")

	if errPayload := validatorlib.ValidateStructToReply(payloads.TemplateShareName{Name: name}); errPayload != nil {
		rp.Error(replylib.ErrorPayloadToErrorArg(*errPayload)).FailJSON()
		return
	}

	shares, err := samba.ReadShare()
	if err != nil {
		replylib.HandleError(err, rp)
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

	if errPayload := validatorlib.ValidateStructToReply(payloads.RequestUpdateShare{
		TemplateShareName: payloads.TemplateShareName{Name: name},
		Share:             payload,
	}); errPayload != nil {
		rp.Error(replylib.ErrorPayloadToErrorArg(*errPayload)).FailJSON()
		return
	}

	shares, err := samba.UpdateShare(name, payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(samba.FilterShares(shares)).OkJSON()
}

func (h *Samba) DeleteOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	name := c.Param("name")

	if errPayload := validatorlib.ValidateStructToReply(payloads.TemplateShareName{Name: name}); errPayload != nil {
		rp.Error(replylib.ErrorPayloadToErrorArg(*errPayload)).FailJSON()
		return
	}

	shares, err := samba.DeleteShare(name)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(samba.FilterShares(shares)).OkJSON()
}

// configuration

func (h *Samba) GetConfiguration(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	config, err := samba.GetConfiguration()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(config).OkJSON()
}

func (h *Samba) UpdateConfig(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	var payload samba.ShareMap
	if err := c.ShouldBindJSON(&payload); err != nil {
		rp.Error(replylib.CodeBadRequest, err.Error()).FailJSON()
		return
	}

	err := samba.UpdateConfig(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(payload).OkJSON()
}
