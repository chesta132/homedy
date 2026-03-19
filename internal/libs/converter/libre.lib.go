package converter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// libreofficeConvert runs soffice --headless --convert-to <format>.
// Accepts input as bytes + original filename (for extension hints),
// writes to temp, converts, reads result, then cleans up.
func libreofficeConvert(inputData []byte, inputFilename, targetFormat string) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "libreconv-*")
	if err != nil {
		return nil, fmt.Errorf("libreofficeConvert mktemp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, inputFilename)
	if err := os.WriteFile(inputPath, inputData, 0644); err != nil {
		return nil, fmt.Errorf("libreofficeConvert write input: %w", err)
	}

	cmd := exec.Command(
		"soffice",
		"--headless",
		"--convert-to", targetFormat,
		"--outdir", tmpDir,
		inputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("libreofficeConvert soffice error: %w\nstderr: %s", err, stderr.String())
	}

	// Reconstruct output path
	base := strings.TrimSuffix(inputFilename, filepath.Ext(inputFilename))
	outputPath := filepath.Join(tmpDir, base+"."+targetFormat)

	result, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("libreofficeConvert read output: %w", err)
	}

	return result, nil
}
