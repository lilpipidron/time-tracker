package requests

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
	"io"
	"net/http"
)

func Decode(w http.ResponseWriter, r *http.Request, request *interface{}) {
	err := render.DecodeJSON(r.Body, &request)
	if errors.Is(err, io.EOF) {
		log.Debug("request body is empty")

		render.JSON(w, r, "empty request")

		return
	}

	if err != nil {
		log.Error("failed to decode request body", err)

		render.JSON(w, r, "failed to decode request")

		return
	}
}
