package http

import (
	"io"
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
	bytes, err := io.ReadAll(file)

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	contentType := http.DetectContentType(bytes)

	image, err := h.uploadImageUseCase.Execute(r.Context(), usecases.UploadImageInput{
		Bytes:       &bytes,
		ContentType: contentType,
	})

	if err != nil {
		h.presentBadRequest(w, r, err)
		return
	}

	h.presentJSON(w, r, http.StatusCreated, imageJson{FileName: image.FileName(), Url: image.Url(), ContentType: image.ContentType()})
}

type imageJson struct {
	FileName    string `json:"fileName"`
	Url         string `json:"url"`
	ContentType string `json:"contentType"`
}

func toImageJson(image entities.Image) imageJson {
	return imageJson{
		FileName:    image.FileName(),
		Url:         image.Url(),
		ContentType: image.ContentType(),
	}
}
