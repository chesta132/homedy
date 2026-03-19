package converter

// PDFToXLSX converts a PDF file to XLSX using LibreOffice.
func PDFToXLSX(inputData []byte, filename string) ([]byte, error) {
	return libreofficeConvert(inputData, filename, "xlsx")
}

// XLSXToPDF converts an XLSX file to PDF using LibreOffice.
func XLSXToPDF(inputData []byte, filename string) ([]byte, error) {
	return libreofficeConvert(inputData, filename, "pdf")
}

// PDFToDocx converts a PDF file to DOCX using LibreOffice.
func PDFToDocx(inputData []byte, filename string) ([]byte, error) {
	return libreofficeConvert(inputData, filename, "docx")
}

// DocxToPDF converts a DOCX file to PDF using LibreOffice.
func DocxToPDF(inputData []byte, filename string) ([]byte, error) {
	return libreofficeConvert(inputData, filename, "pdf")
}

// PDFToPptx converts a PDF file to PPTX using LibreOffice.
func PDFToPptx(inputData []byte, filename string) ([]byte, error) {
	return libreofficeConvert(inputData, filename, "pptx")
}

// PptxToPDF converts a PPTX file to PDF using LibreOffice.
func PptxToPDF(inputData []byte, filename string) ([]byte, error) {
	return libreofficeConvert(inputData, filename, "pdf")
}
