package profilehttp

import (
	"net/http"

	authmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/auth"
	requestmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/request"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/api"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
)

func (h *ProfileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.profile.DeleteProfile"

	log := sl.HandlerJob(h.Log, op, requestmdw.GetRequestId(r), h.srvc)

	user := authmdw.User(r)
	if user == nil {
		log.Error("failed to get user from auth context")
		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	if err := h.srvc.DeleteProfile(r.Context(), user.Id); err != nil {
		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	h.JSON(w, http.StatusNoContent, api.Ok())
}
