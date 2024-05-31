package api

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/uguremirmustafa/inventory/db"
)

type UploadService struct {
	q  *db.Queries
	db *sql.DB
}

func NewUploadService(q *db.Queries, db *sql.DB) *UploadService {
	// Create the uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads/tmp", os.ModePerm); err != nil {
		fmt.Printf("Unable to create uploads directory: %v\n", err)
	}
	return &UploadService{
		q:  q,
		db: db,
	}
}

func (s *UploadService) HandleUploadImages(w http.ResponseWriter, r *http.Request) error {

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		slog.Error("something went wrong while getting images from form", slog.String("err", err.Error()))
		return err
	}

	var uploadedURLs []string

	files := r.MultipartForm.File["images"]
	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			slog.Error("unable to open file")
			return fmt.Errorf("unable to open file: %v", err.Error())
		}
		defer file.Close()

		uuidFileName := uuid.New().String() + filepath.Ext(fh.Filename)
		filePath := filepath.Join("uploads", "tmp", uuidFileName)
		dst, err := os.Create(filePath)
		if err != nil {
			slog.Error("unable to create file on server")
			return fmt.Errorf("unable to create file on server")
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			slog.Error("unable to save file")
			return fmt.Errorf("unable to save file")
		}
		uploadedURLs = append(uploadedURLs, filePath)
	}

	writeJson(w, http.StatusOK, uploadedURLs)
	return nil
}
