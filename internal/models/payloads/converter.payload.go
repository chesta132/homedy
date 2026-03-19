package payloads

import (
	"fmt"
	"homedy/config"
	"homedy/internal/libs/converter"
	"mime/multipart"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

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
		// convert pairs validation
		convertTo := p.ConvertTo[i]
		ext := filepath.Ext(file.Filename)
		if len(ext) < 1 {
			// arg convertTo to process in translate as {fileExt} -> {convertTo} is not a valid pair
			sl.ReportError("unknown", fmt.Sprintf("files[%d]", i), "Files", "convert_pair", convertTo)
			continue
		}
		ext = ext[1:]
		if !converter.IsValidPair(ext, convertTo) {
			sl.ReportError(ext, fmt.Sprintf("files[%d]", i), "Files", "convert_pair", convertTo)
		}

		// size validation
		if limit, ok := config.ConvertFileLimits[ext]; ok && file.Size > limit {
			sl.ReportError(file.Filename, fmt.Sprintf("files[%d]", i), "Files", "size_limit", fmt.Sprint(limit))
		}
	}
}
