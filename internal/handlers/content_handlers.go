package handlers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"time"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/services"
	"github.com/gofrs/uuid"
)

type ContentHandler struct {
	contentService    services.ContentService
	licenseService    services.LicenseService
	sessionKeyService services.SessionKeyService
}

func NewContentHandler(contentService services.ContentService, licenseService services.LicenseService,
	sessionKeyService services.SessionKeyService) *ContentHandler {
	return &ContentHandler{contentService: contentService, licenseService: licenseService, sessionKeyService: sessionKeyService}
}

func (h *ContentHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	data := r.FormValue("data")
	if data == "" {
		http.Error(w, "Missing content data", http.StatusBadRequest)
		return
	}

	var content models.Content
	err = json.Unmarshal([]byte(data), &content)
	if err != nil {
		http.Error(w, "Invalid content metadata JSON", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("content")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	content.CreatorID = id
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()

	fileExtension := filepath.Ext(header.Filename)

	err = h.contentService.Create(&content, file, fileExtension, header.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
