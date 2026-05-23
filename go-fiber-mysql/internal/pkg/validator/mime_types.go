package validator

import (
	"fmt"
	"mime/multipart"
	"net/http"
)

type AllowedMimes map[string]bool

const (
	ImageJpgMimeType  = "image/jpg"
	ImageJpegMimeType = "image/jpeg"
	ImagePngMimeType  = "image/png"
	ImageWebpMimeType = "image/webp"

	XlsxMimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	XlsMimeType  = "application/vnd.ms-excel"

	PdfMimeType = "application/pdf"
)

var ImageMimeTypes = AllowedMimes{
	ImageJpgMimeType:  true,
	ImageJpegMimeType: true,
	ImagePngMimeType:  true,
	ImageWebpMimeType: true,
}

var ExcelMimeTypes = AllowedMimes{
	XlsxMimeType: true,
	XlsMimeType:  true,
}

var PdfMimeTypes = AllowedMimes{
	PdfMimeType: true,
}

func DetectMime(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func ValidateMime(fileHeader *multipart.FileHeader, allowedMimes AllowedMimes) (string, error) {
	mime, err := DetectMime(fileHeader)
	if err != nil {
		return "", err
	}

	if !allowedMimes[mime] {
		return "", fmt.Errorf("File type %s is not allowed", mime)
	}

	return mime, nil
}
