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

// @Summary      Create new folder share
// @Tags         folder-sharing
// @Accept       json
// @Produce			 json
// @Param				 payload  body	payloads.RequestCreateShare	true	"new folder share config"
// @Success      201  		{object}  replylib.Envelope{data=models.Shares} "all existing shares"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /samba [post]
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

// @Summary      Get all existing shares
// @Tags         folder-sharing
// @Produce      json
// @Success      200  		{object}  replylib.Envelope{data=models.Shares}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /samba [get]
func (h *Samba) GetShares(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	shares, err := h.sambaSvc.GetShares()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(sambalib.FilterShares(shares)).OkJSON()
}

// @Summary      Get existing share by name
// @Tags         folder-sharing
// @Produce      json
// @Param				 param		path	payloads.TemplateShareName	true "param payload"
// @Success      200  		{object}  replylib.Envelope{data=models.Share}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /samba/{name} [get]
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

// @Summary      Update existing share by name
// @Tags         folder-sharing
// @Accept       json
// @Produce      json
// @Param				 payload  body	models.Share	true	"updated folder share config"
// @Param				 param		path	payloads.TemplateShareName	true "param payload"
// @Success      200  		{object}  replylib.Envelope{data=models.Shares} "all existing shares"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /samba/{name} [put]
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

// @Summary      Delete existing share by name
// @Tags         folder-sharing
// @Produce      json
// @Param				 param		path	payloads.TemplateShareName	true "param payload"
// @Success      200  		{object}  replylib.Envelope{data=models.Shares} "all existing shares"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /samba/{name} [delete]
func (h *Samba) DeleteShare(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestDeleteShare](c.ShouldBindUri)
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

// @Summary      Get samba configuration
// @Tags         folder-sharing
// @Produce      json
// @Param				 X-APP-SECRET header string true "app secret authentication for access"
// @Success      200  		{object}  replylib.Envelope{data=models.ShareMap}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /config [get]
func (h *Samba) GetConfig(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	config, err := h.sambaSvc.GetConfig()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(config).OkJSON()
}

// @Summary      Update samba configuration
// @Tags         folder-sharing
// @Accept			 json
// @Produce      json
// @Param				 X-APP-SECRET header string true "app secret authentication for access"
// @Param				 payload  body	models.ShareMap	true	"updated samba config"
// @Success      200  		{object}  replylib.Envelope{data=models.ShareMap}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /config [put]
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

// backup

// @Summary      Backup samba configuration and shares and replace previous backup
// @Tags         folder-sharing
// @Produce      json
// @Param				 X-APP-SECRET header string true "app secret authentication for access"
// @Success      200  		{object}  replylib.Envelope
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /backup [post]
func (h *Samba) Backup(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	err := h.sambaSvc.Backup()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(nil).OkJSON()
}

// @Summary      Restore samba configuration and shares from latest backup
// @Tags         folder-sharing
// @Produce      json
// @Param				 X-APP-SECRET header string true "app secret authentication for access"
// @Success      200  		{object}  replylib.Envelope{data=models.Shares} "all existing shares"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /backup [post]
func (h *Samba) Restore(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	shares, err := h.sambaSvc.Restore()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(sambalib.FilterShares(shares)).OkJSON()
}
