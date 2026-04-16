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
// @Success      200  		{object}  replylib.Envelope{data=[]payloads.ResponseGetRepo}
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
