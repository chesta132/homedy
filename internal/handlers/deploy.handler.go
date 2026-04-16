package handlers

import (
	"homedy/internal/libs/ginlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/models/payloads"
	"homedy/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Deploy struct {
	deploySvc *services.Deploy
}

func NewDeploy(deploySvc *services.Deploy) *Deploy {
	return &Deploy{deploySvc}
}

// @Summary      Create new session
// @Tags         deployment
// @Produce      json
// @Success      200  		{object}  replylib.Envelope{data=payloads.TemplateWithSession}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/new [post]
func (h *Deploy) CreateSession(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	session, err := h.deploySvc.AttachContext(c).CreateSession()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(payloads.TemplateWithSession{Session: session}).OkJSON()
}

// @Summary      Get client's existing repositories
// @Tags         deployment
// @Produce      json
// @Param				 param	  path			payloads.TemplateWithSession	true	"session param"
// @Success      200  		{object}  replylib.Envelope{data=[]models.FilteredGHRepo}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/{session}/repos [get]
func (h *Deploy) GetRepos(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.TemplateWithSession](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	repos, err := h.deploySvc.AttachContext(c).GetRepos(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(repos).OkJSON()
}

// @Summary      Set selected repository
// @Tags         deployment
// @Accept			 json
// @Produce      json
// @Param				 payload  body			payloads.RequestSelectRepo	true	"data to set selected repo"
// @Param				 param	  path			payloads.TemplateWithSession	true	"session param"
// @Success      200  		{object}  replylib.Envelope{data=models.FilteredGHRepo}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/{session}/selected-repo [post]
func (h *Deploy) SelectRepo(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestSelectRepo](c.ShouldBindUri, c.ShouldBindJSON)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	repo, err := h.deploySvc.AttachContext(c).SelectRepo(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(repo).OkJSON()
}

// @Summary      Get selected repository (response not found error if not selected yet)
// @Tags         deployment
// @Produce      json
// @Param				 param	  path			payloads.TemplateWithSession	true	"session param"
// @Success      200  		{object}  replylib.Envelope{data=models.FilteredGHRepo}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/{session}/selected-repo [get]
func (h *Deploy) GetSelectedRepo(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.TemplateWithSession](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	repo, err := h.deploySvc.AttachContext(c).GetSelectedRepo(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(repo).OkJSON()
}
