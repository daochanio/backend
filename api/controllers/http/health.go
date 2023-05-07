package http

import "net/http"

func (h *httpServer) healthRoute(w http.ResponseWriter, r *http.Request) {
	h.presentText(w, r, http.StatusOK, "OK")
}
