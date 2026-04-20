package rest

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FilesHandler struct{}

func NewFilesHandler() *FilesHandler {
	return &FilesHandler{}
}

type FileResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Size      int64  `json:"size"`
	Type      string `json:"type"`
	MimeType  string `json:"mime_type"`
	CreatedAt string `json:"created_at"`
}

func (h *FilesHandler) GetFiles(c *gin.Context) {
	workDir, _ := os.Getwd()
	uploadsDir := filepath.Join(workDir, "uploads")

	var allFiles []FileResponse

	err := filepath.Walk(uploadsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(workDir, path)
		url := "/" + strings.ReplaceAll(relPath, "\\", "/")

		fileType := "document"
		mimeType := "application/octet-stream"

		if strings.HasPrefix(url, "/uploads/images/") {
			fileType = "image"
			ext := strings.ToLower(filepath.Ext(info.Name()))
			mimeType = "image/" + strings.TrimPrefix(ext, ".")
		}

		allFiles = append(allFiles, FileResponse{
			ID:        uuid.New().String(),
			Name:      info.Name(),
			URL:       url,
			Size:      info.Size(),
			Type:      fileType,
			MimeType:  mimeType,
			CreatedAt: info.ModTime().Format(time.RFC3339),
		})
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allFiles)
}

func (h *FilesHandler) DeleteFile(c *gin.Context) {
	filepathParam := c.Param("filepath")

	if filepathParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file path required"})
		return
	}

	workDir, _ := os.Getwd()
	fullPath := filepath.Join(workDir, "uploads", filepathParam)

	if !strings.Contains(fullPath, "uploads") {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot delete files outside uploads"})
		return
	}

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
