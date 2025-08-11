package internal

import (
	"encoding/json"
	"net/http"

	"github.com/isaacgr/portfolio/internal/pkg/logging"
)

var log = logging.GetLogger("jsonutil", false)

func Encode[T any](
	w http.ResponseWriter,
	r *http.Request,
	status int,
	v T,
) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error("Unable to encode json.", "Error", err)
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Error("Unable to decode json.", "Error", err)
	}
	return v, nil
}
