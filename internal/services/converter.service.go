package services

import (
	"archive/zip"
	"bytes"
	"context"
	"homedy/internal/libs/converter"
	"homedy/internal/libs/replylib"
	"homedy/internal/models/payloads"
	"io"
	"path/filepath"
	"strings"

	"github.com/chesta132/goreply/reply"
	"github.com/gin-gonic/gin"
)

type Converter struct{}

type ContextedConverter struct {
	Converter
	c   *gin.Context
	ctx context.Context
}

func NewConverter() *Converter {
	return &Converter{}
}

func (s *Converter) AttachContext(c *gin.Context) *ContextedConverter {
	return &ContextedConverter{*s, c, c.Request.Context()}
}

func (s *ContextedConverter) ConvertMultiple(payload payloads.RequestConvertMultiple) (map[string][]byte, error) {
	converted := make(map[string][]byte)

	for i, file := range payload.Files {
		convertTo := payload.ConvertTo[i]
		// slice ext is safe cz already validated
		ext := filepath.Ext(file.Filename)[1:]
		conv, err := converter.GetConverter(ext, convertTo)
		if err != nil {
			return nil, &reply.ErrorPayload{Code: replylib.CodeBadRequest, Message: err.Error()}
		}
		f, err := file.Open()
		if err != nil {
			return nil, &reply.ErrorPayload{Code: replylib.CodeBadRequest, Message: err.Error()}
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, f); err != nil {
			f.Close()
			return nil, err
		}
		f.Close()

		fileConverted, err := conv(buf.Bytes(), file.Filename)
		if err != nil {
			return nil, err
		}
		converted[strings.TrimSuffix(file.Filename, ext) + convertTo] = fileConverted
	}
	return converted, nil
}

func (s *Converter) StreamMultipleToZip(w io.Writer, multiple map[string][]byte) error {
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	for name, file := range multiple {
		writer, err := zipWriter.Create(name)
		if err != nil {
			return err
		}
		writer.Write(file)
	}
	return nil
}
