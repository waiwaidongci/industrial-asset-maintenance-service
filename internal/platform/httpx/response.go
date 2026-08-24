package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/example/asset-maintenance-service/internal/domain"
	"io"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func Decode(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("request body required")
	}
	d := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	d.DisallowUnknownFields()
	if e := d.Decode(v); e != nil {
		return fmt.Errorf("decode json: %w", e)
	}
	return nil
}
func Error(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err == domain.ErrNotFound {
		status = http.StatusNotFound
	} else if err == domain.ErrConflict || err == domain.ErrTransition {
		status = http.StatusConflict
	} else if err == domain.ErrInvalid {
		status = http.StatusBadRequest
	}
	JSON(w, status, map[string]any{"error": err.Error(), "status": status})
}
