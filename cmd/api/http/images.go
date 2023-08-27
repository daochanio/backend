package http

import (
	"net/http"

	"github.com/daochanio/backend/domain/entities"
	"github.com/daochanio/backend/domain/usecases"
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
	FileName  string     `json:"fileName"`
	Original  headerJson `json:"original"`
	Formatted headerJson `json:"formatted"`
}

type headerJson struct {
	URL         string `json:"url"`
	ContentType string `json:"contentType"`
}

func toImageJson(image *entities.Image) *imageJson {
	if image == nil {
		return nil
	}

	return &imageJson{
		FileName: image.FileName(),
		Original: headerJson{
			URL:         image.OriginalURL(),
			ContentType: image.OriginalContentType(),
		},
		Formatted: headerJson{
			URL:         image.FormattedURL(),
			ContentType: image.FormattedContentType(),
		},
	}
}
