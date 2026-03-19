package converter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// pdfConvertMap defines infilter + convert-to format per target format when source is PDF.
var pdfConvertMap = map[string]struct {
	infilter  string
	convertTo string
}{
	"docx": {infilter: "writer_pdf_import", convertTo: "docx"},
	"pptx": {infilter: "impress_pdf_import", convertTo: "pptx"},
}

// libreofficeConvert runs soffice --headless --convert-to <format>.
// Accepts input as bytes + original filename (for extension hints),
// writes to temp, converts, reads result, then cleans up.
func libreofficeConvert(inputData []byte, inputFilename, targetFormat string) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "libreconv-*")
	if err != nil {
		return nil, fmt.Errorf("libreofficeConvert mktemp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputFilename = filepath.Base(inputFilename)
	inputPath := filepath.Join(tmpDir, inputFilename)
	if err := os.WriteFile(inputPath, inputData, 0644); err != nil {
		return nil, fmt.Errorf("libreofficeConvert write input: %w", err)
	}

	srcExt := strings.ToLower(strings.TrimPrefix(filepath.Ext(inputFilename), "."))

	// Build args
	args := []string{"--headless"}

	if srcExt == "pdf" {
		cfg, ok := pdfConvertMap[targetFormat]
		if !ok {
			return nil, fmt.Errorf("libreofficeConvert: unsupported PDF -> %s conversion", targetFormat)
		}
		args = append(args, "--infilter="+cfg.infilter, "--convert-to", cfg.convertTo)
	} else {
		args = append(args, "--convert-to", targetFormat)
	}

	args = append(args, "--outdir", tmpDir, inputPath)

	cmd := exec.Command("soffice", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("libreofficeConvert soffice error: %w\nstderr: %s", err, stderr.String())
	}

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("libreofficeConvert read dir: %w", err)
	}

	ext := "." + targetFormat
	for _, entry := range entries {
		if !entry.IsDir() && strings.EqualFold(filepath.Ext(entry.Name()), ext) {
			outputPath := filepath.Join(tmpDir, entry.Name())
			result, err := os.ReadFile(outputPath)
			if err != nil {
				return nil, fmt.Errorf("libreofficeConvert read output: %w", err)
			}
			return result, nil
		}
	}

	return nil, fmt.Errorf("libreofficeConvert: output .%s file not found in temp dir", targetFormat)
}