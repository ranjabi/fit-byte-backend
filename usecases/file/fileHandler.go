package file

import (
	"fit-byte/models"
	"fit-byte/utils"
	"net/http"
	"path/filepath"
	"strings"
)

type FileHandler struct {
	fileService FileService
}

var (
	maxFileSize   = int64(100 * 1024) // 100KB
)

func NewFileHandler(fileService FileService) FileHandler {
	return FileHandler{fileService}
}

func (h *FileHandler) HandleUploadFile(w http.ResponseWriter, r *http.Request) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	if err := r.ParseMultipartForm(maxFileSize); err != nil {
		return models.NewError(http.StatusBadRequest, "File is too large")
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return models.NewError(http.StatusBadRequest, "Invalid file")
	}
	defer file.Close()

	fileName := fileHeader.Filename
	fileExt := strings.ToLower(filepath.Ext(fileName))

	allowedExtensions := map[string]bool{
		".png": true,
		".jpg": true,
		".jpeg": true,
	}

	if !allowedExtensions[fileExt] {
		return models.NewError(http.StatusBadRequest, "Only jpeg, jpg, and png files are allowed") 
	}

	if fileHeader.Size > maxFileSize {
		return models.NewError(http.StatusBadRequest, "File exceeds 100KB")
	}

	s3FileKey, err := h.fileService.UploadToS3(file, fileHeader)
	if err != nil {
		return err
	}

	res := struct {
		Uri string `json:"uri"`
	}{
		Uri: utils.GenerateS3FileURL(s3FileKey),
	}
	utils.SetJsonResponse(w, http.StatusOK, res)

	return nil
}