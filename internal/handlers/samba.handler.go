package handlers

import (
	"homedy/internal/libs/ginlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/libs/sambalib"
	"homedy/internal/models"
	"homedy/internal/models/payloads"
	"homedy/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Samba struct {
	sambaSvc *services.Samba
}

func NewSamba(sambaSvc *services.Samba) *Samba {
	return &Samba{sambaSvc}
}

func (h *Samba) CreateShare(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindJSONAndValidate[payloads.RequestCreateShare](c)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	shares, err := h.sambaSvc.CreateShare(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(sambalib.FilterShares(shares)).CreatedJSON()
}

func (h *Samba) GetShares(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	shares, err := h.sambaSvc.GetShares()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(sambalib.FilterShares(shares)).OkJSON()
}

func (h *Samba) GetShare(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.TemplateShareName](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	shares, err := h.sambaSvc.GetShares()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	share, ok := sambalib.FilterShares(shares)[payload.Name]
	if !ok {
		rp.Error(replylib.CodeNotFound, "share not found").FailJSON()
		return
	}

	rp.Success(share).OkJSON()
}

func (h *Samba) UpdateShare(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestUpdateShare](c.ShouldBindJSON, c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	shares, err := h.sambaSvc.UpdateShare(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(sambalib.FilterShares(shares)).OkJSON()
}

func (h *Samba) DeleteShare(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestDeleteShare](c.ShouldBindJSON, c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	shares, err := h.sambaSvc.DeleteShare(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(sambalib.FilterShares(shares)).OkJSON()
}

// configuration

func (h *Samba) GetConfig(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	config, err := h.sambaSvc.GetConfig()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(config).OkJSON()
}

func (h *Samba) UpdateConfig(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindJSONAndValidate[models.ShareMap](c)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	err = h.sambaSvc.UpdateConfig(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(payload).OkJSON()
}
