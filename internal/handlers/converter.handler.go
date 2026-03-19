package handlers

import (
	"fmt"
	"homedy/internal/libs/converter"
	"homedy/internal/libs/ginlib"
	"homedy/internal/libs/replylib"
	"homedy/internal/models/payloads"
	"homedy/internal/services"
	"net/http"
	"path/filepath"
	"time"

	adapter "github.com/chesta132/goreply/adapter/gin"
	"github.com/gin-gonic/gin"
)

type Converter struct {
	convSvc *services.Converter
}

func NewConverter(convSvc *services.Converter) *Converter {
	return &Converter{convSvc}
}

func (h *Converter) ConvertMultiple(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	svc := h.convSvc.AttachContext(c)

	payload, err := ginlib.BindAndValidate[payloads.RequestConvertMultiple](c.ShouldBind)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	converted, err := svc.ConvertMultiple(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%d_converted_%d.zip"`, len(converted), time.Now().Unix()))

	if err := svc.StreamMultipleToZip(c.Writer, converted); err != nil {
		c.Header("Content-Disposition", "")
		replylib.HandleError(err, rp)
		return
	}
}

func (h *Converter) ConvertOne(c *gin.Context) {
	rp := replylib.Client.Use(adapter.AdaptGin(c))
	svc := h.convSvc.AttachContext(c)

	payload, err := ginlib.BindAndValidate[payloads.RequestConvertOne](c.ShouldBind)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	fileName, fileBytes, err := svc.ConvertOne(payload)
	if err != nil {
		replylib.HandleError(err, rp)
		return
	}

	ext := filepath.Ext(fileName)[1:]
	contentType, ok := converter.MimePairs[ext]
	if !ok {
		contentType = http.DetectContentType(fileBytes)
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	// not using rp due to contentType supports
	c.Data(http.StatusOK, contentType, fileBytes)
}
