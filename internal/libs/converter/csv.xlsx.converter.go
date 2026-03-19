package converter

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// XLSXToCSV converts the first sheet of an XLSX file to CSV bytes.
func XLSXToCSV(inputData []byte) ([]byte, error) {
	f, err := excelize.OpenReader(bytes.NewReader(inputData))
	if err != nil {
		return nil, fmt.Errorf("XLSXToCSV open: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("XLSXToCSV: no sheets found")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("XLSXToCSV get rows: %w", err)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.WriteAll(rows); err != nil {
		return nil, fmt.Errorf("XLSXToCSV write csv: %w", err)
	}
	w.Flush()

	return buf.Bytes(), nil
}

// CSVToXLSX converts CSV bytes to XLSX bytes.
func CSVToXLSX(csvContent []byte) ([]byte, error) {
	r := csv.NewReader(bytes.NewReader(csvContent))
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSVToXLSX read csv: %w", err)
	}

	f := excelize.NewFile()
	sheet := "Sheet1"

	for rowIdx, record := range records {
		for colIdx, cell := range record {
			cellName, err := excelize.CoordinatesToCellName(colIdx+1, rowIdx+1)
			if err != nil {
				return nil, fmt.Errorf("CSVToXLSX coord: %w", err)
			}
			if err := f.SetCellValue(sheet, cellName, cell); err != nil {
				return nil, fmt.Errorf("CSVToXLSX set cell: %w", err)
			}
		}
	}

	var buf bytes.Buffer
	if _, err := f.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("CSVToXLSX write: %w", err)
	}

	return buf.Bytes(), nil
}
