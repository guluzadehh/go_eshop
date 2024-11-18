package render

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
)

type Responder struct {
	log *slog.Logger
}

func NewResponder(log *slog.Logger) *Responder {
	return &Responder{
		log: log,
	}
}

func (r *Responder) JSON(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		r.log.Error("json encode err", sl.Err(err))
		http.Error(w, "failed to return json", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}
