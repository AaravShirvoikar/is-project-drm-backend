package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/AaravShirvoikar/is-project-drm-backend/internal/services"
	"github.com/go-chi/chi"
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
		http.Error(w, "Invalid content data", http.StatusBadRequest)
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

	similarId, created, similarity, err := h.contentService.Create(&content, file, fileExtension, header.Size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if created {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusConflict)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Created    bool    `json:"created"`
		SimilarID  string  `json:"similar_id"`
		Similarity float64 `json:"similarity"`
	}{
		Created:    created,
		SimilarID:  similarId,
		Similarity: similarity,
	})
}

func (h *ContentHandler) ListContent(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)
	contents, err := h.contentService.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filteredContents := make([]struct {
		Id          uuid.UUID `json:"content_id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Price       float64   `json:"price"`
		Purchased   bool      `json:"purchased"`
	}, len(contents))

	for i, content := range contents {
		filteredContents[i].Id = content.ContentID
		filteredContents[i].Title = content.Title
		filteredContents[i].Description = content.Description
		filteredContents[i].Price = content.Price
		isPurchased := h.licenseService.Verify(id, content.ContentID.String())
		if content.CreatorID.String() == id {
			isPurchased = true
		}
		filteredContents[i].Purchased = isPurchased
	}

	json.NewEncoder(w).Encode(filteredContents)
}

func (h *ContentHandler) PurchaseContent(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)
	contentId := chi.URLParam(r, "id")

	err := h.licenseService.Generate(id, contentId, time.Now().Add(time.Hour*24))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ContentHandler) GetContentData(w http.ResponseWriter, r *http.Request) {
	contentId := chi.URLParam(r, "id")

	content, _, err := h.contentService.Get(contentId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		ContentId   string `json:"content_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}{
		ContentId:   content.ContentID.String(),
		Title:       content.Title,
		Description: content.Description,
	})
}

func (h *ContentHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(string)
	contentId := chi.URLParam(r, "id")

	content, fileContent, err := h.contentService.Get(contentId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isValidLicense := h.licenseService.Verify(id, contentId)

	if content.CreatorID.String() == id {
		isValidLicense = true
	}

	if !isValidLicense {
		http.Error(w, "Invalid license", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Disposition", "inline; filename=\"video.mp4\"")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))

	http.ServeContent(w, r, "video.mp4", content.UpdatedAt, bytes.NewReader(fileContent))
}
