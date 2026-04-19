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
// @Router			 /deploy/session [post]
func (h *Deploy) CreateSession(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	session, err := h.deploySvc.AttachContext(c).CreateSession()
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(payloads.TemplateWithSession{Session: session}).OkJSON()
}

// @Summary      Create new session
// @Tags         deployment
// @Produce      json
// @Success      200  		{object}  replylib.Envelope "data is null"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/session/{session}/invalidate [delete]
func (h *Deploy) InvalidateSession(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.TemplateWithSession](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	err = h.deploySvc.AttachContext(c).InvalidateSession(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(nil).OkJSON()
}

// @Summary      Get client's existing repositories
// @Tags         deployment
// @Produce      json
// @Param				 param	  path			payloads.TemplateWithSession	true	"session param"
// @Success      200  		{object}  replylib.Envelope{data=[]models.FilteredGHRepo}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/session/{session}/repos [get]
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
// @Param				 payload  body			payloads.RequestSetSelectedRepo	true	"data to set selected repo"
// @Param				 param	  path			payloads.TemplateWithSession	true	"session param"
// @Success      200  		{object}  replylib.Envelope{data=models.SelectedRepoInSession}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/session/{session}/selected-repo [post]
func (h *Deploy) SetSelectedRepo(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestSetSelectedRepo](c.ShouldBindUri, c.ShouldBindJSON)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	repo, err := h.deploySvc.AttachContext(c).SetSelectedRepo(payload)
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
// @Success      200  		{object}  replylib.Envelope{data=models.SelectedRepoInSession}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/session/{session}/selected-repo [get]
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

// @Summary      Get branch of repository
// @Tags         deployment
// @Produce      json
// @Param				 param	  path			payloads.TemplateWithSession	true	"session param"
// @Success      200  		{object}  replylib.Envelope{data=[]string}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/session/{session}/repos/{id}/branches [get]
func (h *Deploy) GetBranches(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestGetBranches](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	branches, err := h.deploySvc.AttachContext(c).GetBranches(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(branches).OkJSON()
}

// @Summary      Get global env and selected repo env
// @Tags         deployment
// @Produce      json
// @Param				 param	  path			payloads.TemplateWithSession	true	"session param"
// @Success      200  		{object}  replylib.Envelope{data=payloads.ResponseSessionEnv}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/session/{session}/env [get]
func (h *Deploy) GetEnv(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.TemplateWithSession](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	env, err := h.deploySvc.AttachContext(c).GetEnv(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(env).OkJSON()
}

// @Summary      Set global and selected repo env. This action is replace behavior
// @Tags         deployment
// @Produce      json
// @Param				 param	  path			payloads.RequestSetSessionEnv	true	"session param"
// @Success      200  		{object}  replylib.Envelope "data is null"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /deploy/session/{session}/env [post]
func (h *Deploy) SetEnv(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestSetSessionEnv](c.ShouldBindUri, c.ShouldBindJSON)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	err = h.deploySvc.AttachContext(c).SetEnv(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	rp.Success(nil).OkJSON()
}
