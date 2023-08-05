package http

import (
	"net/http"

	"github.com/daochanio/backend/api/entities"
	"github.com/daochanio/backend/api/usecases"
)

func (h *httpServer) uploadImageRoute(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	defer file.Close()

	image, err := h.uploadImage.Execute(r.Context(), usecases.UploadImageInput{
		Reader: file,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusCreated, toImageJson(image), nil)
}

type imageJson struct {
	FileName     string `json:"fileName"`
	OriginalURL  string `json:"originalUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

func toImageJson(image *entities.Image) *imageJson {
	if image == nil {
		return nil
	}

	return &imageJson{
		FileName:     image.FileName(),
		OriginalURL:  image.OriginalURL(),
		ThumbnailURL: image.ThumbnailURL(),
	}
}
