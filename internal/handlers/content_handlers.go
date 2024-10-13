package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/services"
	"github.com/gofrs/uuid"
)

type ContentHandler struct {
	contentService services.ContentService
}

func NewContentHandler(contentService services.ContentService) *ContentHandler {
	return &ContentHandler{contentService: contentService}
}

func (h *ContentHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
	var content models.Content

	if err := json.NewDecoder(r.Body).Decode(&content); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	content.CreatorID = id
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()

	err = h.contentService.Create(&content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
