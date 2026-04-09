package handlers

import (
	"homedy/config"
	"homedy/internal/libs/ginlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/models"
	"homedy/internal/models/payloads"
	"homedy/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Note struct {
	noteSvc *services.Note
}

func NewNote(noteSvc *services.Note) *Note {
	return &Note{noteSvc}
}

// @Summary      Create one new note
// @Tags         note
// @Accept       json
// @Produce			 json
// @Param				 payload  body	payloads.RequestCreateNote	true	"new note data"
// @Success      201  		{object}  replylib.Envelope{data=models.Note}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /notes [post]
func (h *Note) CreateOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindJSONAndValidate[payloads.RequestCreateNote](c)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	note, err := h.noteSvc.AttachContext(c).CreateOne(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(note).CreatedJSON()
}

// @Summary      Get all client's existing notes with query
// @Tags         note
// @Produce      json
// @Param				 payload  query	payloads.RequestGetManyNote	true	"get many option"
// @Success      200  		{object}  replylib.Envelope{data=[]models.Note}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError,meta=replylib.Pagination}}
// @Router			 /notes [get]
func (h *Note) GetMany(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestGetManyNote](c.ShouldBindQuery)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	notes, err := h.noteSvc.AttachContext(c).GetNotes(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(notes).PaginateCursor(config.LIMIT_RESOURCE_PER_PAGINATION, payload.Offset).OkJSON()
}

// @Summary      Get client's existing notes by id
// @Tags         note
// @Produce      json
// @Param				 param	  path			payloads.RequestGetOneNote	true	"param of note's identification"
// @Success      200  		{object}  replylib.Envelope{data=models.Note}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /notes/{id} [get]
func (h *Note) GetOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestGetOneNote](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	note, err := h.noteSvc.AttachContext(c).GetOne(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(note).OkJSON()
}

// @Summary      Update client's existing note
// @Tags         note
// @Accept       json
// @Produce      json
// @Param				 payload  body	payloads.RequestGetOneNote	true	"updated note"
// @Param				 param		path	models.BaseID	true "param of note's identification"
// @Success      200  		{object}  replylib.Envelope{data=models.Note} "updated note from server"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /notes/{id} [put]
func (h *Note) UpdateOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestUpdateNote](c.ShouldBindUri, c.ShouldBindJSON)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	note, err := h.noteSvc.AttachContext(c).UpdateOne(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(note).OkJSON()
}

// @Summary      Soft delete one client's existing note
// @Tags         note
// @Produce      json
// @Param				 param		path	payloads.RequestDeleteOneNote	true "param of note's identification"
// @Success      200  		{object}  replylib.Envelope{data=models.BaseID} "deleted id in server"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /notes/{id} [delete]
func (h *Note) DeleteOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestDeleteOneNote](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	err = h.noteSvc.AttachContext(c).DeleteOne(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(models.BaseID{ID: payload.ID}).OkJSON()
}

// @Summary      Soft delete one client's existing notes
// @Tags         note
// @Produce      json
// @Param				 param		body	payloads.RequestDeleteManyNote	true "param of notes` identification"
// @Success      200  		{object}  replylib.Envelope "data is null"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /notes [delete]
func (h *Note) DeleteMany(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestDeleteManyNote](c.ShouldBindJSON)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	err = h.noteSvc.AttachContext(c).DeleteMany(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(nil).OkJSON()
}

// @Summary      Restore one client's existing note
// @Tags         note
// @Produce      json
// @Param				 param		path	payloads.RequestRestoreOneNote	true "param of note's identification"
// @Success      200  		{object}  replylib.Envelope{data=models.Note} "restored note data"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /notes/restore/{id} [patch]
func (h *Note) RestoreOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestRestoreOneNote](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	note, err := h.noteSvc.AttachContext(c).RestoreOne(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(note).OkJSON()
}

// @Summary      Restore existing share by name
// @Tags         note
// @Produce      json
// @Param				 param		body	payloads.RequestRestoreManyNote	true "param of notes` identification"
// @Success      200  		{object}  replylib.Envelope{data=[]models.Note} "restored notes data"
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /notes/restore [patch]
func (h *Note) RestoreMany(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestRestoreManyNote](c.ShouldBindJSON, c.ShouldBindQuery)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	notes, err := h.noteSvc.AttachContext(c).RestoreMany(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(notes).OkJSON()
}
