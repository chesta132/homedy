package converter

import "fmt"

type ConvertFunc func(input []byte, filename string) ([]byte, error)

type ConvertEntry struct {
	To        string
	Converter ConvertFunc
}

// ConvertPairs maps format asal ke slice of ConvertEntry.
var ConvertPairs = map[string][]ConvertEntry{
	"html": {
		{To: "md", Converter: func(input []byte, _ string) ([]byte, error) { return HTMLToMarkdown(input) }},
	},
	"md": {
		{To: "html", Converter: func(input []byte, _ string) ([]byte, error) { return MarkdownToHTML(input) }},
	},
	"pdf": {
		{To: "xlsx", Converter: func(input []byte, filename string) ([]byte, error) { return PDFToXLSX(input, filename) }},
		{To: "docx", Converter: func(input []byte, filename string) ([]byte, error) { return PDFToDocx(input, filename) }},
		{To: "pptx", Converter: func(input []byte, filename string) ([]byte, error) { return PDFToPptx(input, filename) }},
	},
	"xlsx": {
		{To: "pdf", Converter: func(input []byte, filename string) ([]byte, error) { return XLSXToPDF(input, filename) }},
		{To: "csv", Converter: func(input []byte, _ string) ([]byte, error) { return XLSXToCSV(input) }},
	},
	"docx": {
		{To: "pdf", Converter: func(input []byte, filename string) ([]byte, error) { return DocxToPDF(input, filename) }},
	},
	"pptx": {
		{To: "pdf", Converter: func(input []byte, filename string) ([]byte, error) { return PptxToPDF(input, filename) }},
	},
	"csv": {
		{To: "xlsx", Converter: func(input []byte, _ string) ([]byte, error) { return CSVToXLSX(input) }},
	},
}

// ext = mime
var MimePairs = map[string]string{
	"pdf":  "application/pdf",
	"docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"html": "text/html",
	"md":   "text/markdown",
	"csv":  "text/csv",
}

func IsValidPair(from, to string) bool {
	entries, ok := ConvertPairs[from]
	if !ok {
		return false
	}
	for _, e := range entries {
		if e.To == to {
			return true
		}
	}
	return false
}

func GetConverter(from, to string) (ConvertFunc, error) {
	entries, ok := ConvertPairs[from]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %q", from)
	}
	for _, e := range entries {
		if e.To == to {
			return e.Converter, nil
		}
	}
	return nil, fmt.Errorf("unsupported conversion: %q -> %q", from, to)
}
