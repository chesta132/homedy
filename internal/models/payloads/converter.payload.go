package payloads

import (
	"fmt"
	"homedy/config"
	"homedy/internal/libs/converter"
	"mime/multipart"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

func validateFileAndConvertTarget(file *multipart.FileHeader, convertTo string, sl validator.StructLevel) {
	// convert pairs validation
	ext := filepath.Ext(file.Filename)
	if len(ext) < 1 {
		// arg convertTo to process in translate as {fileExt} -> {convertTo} is not a valid pair
		sl.ReportError("unknown", "file", "File", "convert_pair", convertTo)
		return
	}
	ext = ext[1:]
	if !converter.IsValidPair(ext, convertTo) {
		sl.ReportError(ext, "file", "File", "convert_pair", convertTo)
	}

	// size validation
	if limit, ok := config.ConvertFileLimits[ext]; ok && file.Size > limit {
		sl.ReportError(file.Filename, "file", "File", "size_limit", fmt.Sprint(limit))
	}
}

type RequestConvertMultiple struct {
	Files     []*multipart.FileHeader `form:"files" validate:"required,dive,required"`
	ConvertTo []string                `form:"convert_to" validate:"required,dive,required"`
}

func (p *RequestConvertMultiple) ValidateStruct(sl validator.StructLevel) {
	if len(p.Files) != len(p.ConvertTo) {
		sl.ReportError(p.ConvertTo, "convert_to", "ConvertTo", "len_equals", "files")
		return
	}

	for i, file := range p.Files {
		validateFileAndConvertTarget(file, p.ConvertTo[i], sl)
	}
}

type RequestConvertOne struct {
	File      *multipart.FileHeader `form:"file" validate:"required"`
	ConvertTo string                `form:"convert_to" validate:"required"`
}

func (p *RequestConvertOne) ValidateStruct(sl validator.StructLevel) {
	validateFileAndConvertTarget(p.File, p.ConvertTo, sl)
}
