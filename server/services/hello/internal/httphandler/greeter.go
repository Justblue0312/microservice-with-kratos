package httphandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justblue/luoye/services/hello/internal/domain"
)

type GreeterHandler struct {
    svc domain.Greeter
}

func NewGreeterHandler(svc domain.Greeter) *GreeterHandler {
    return &GreeterHandler{svc: svc}
}

func (h *GreeterHandler) Routes() func(r chi.Router) {
    return func(r chi.Router) {
        r.Get("/hello", h.sayHello)
    }
}

func (h *GreeterHandler) sayHello(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "world"
    }
    reply, err := h.svc.Greet(r.Context(), &domain.GreetRequest{Name: name})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reply)
}
