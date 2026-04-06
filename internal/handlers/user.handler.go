package handlers

import (
	"homedy/internal/libs/ginlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/models/payloads"
	"homedy/internal/services"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type User struct {
	userSvc *services.User
}

func NewUser(userSvc *services.User) *User {
	return &User{userSvc}
}

// @Summary      Get existing user by id
// @Tags         user
// @Produce      json
// @Param				 param	  path			payloads.RequestGetUser	true	"param of user's identification"
// @Success      200  		{object}  replylib.Envelope{data=payloads.ResponseGetUser}
// @Response     default  {object}  replylib.Envelope{data=reply.ErrorPayload{code=replylib.CodeError}}
// @Router			 /users/{id} [get]
func (h *User) GetUser(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))

	payload, err := ginlib.BindAndValidate[payloads.RequestGetUser](c.ShouldBindUri)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	user, err := h.userSvc.AttachContext(c).GetUser(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}
	rp.Success(user).OkJSON()
}
