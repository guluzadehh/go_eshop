package profilehttp

import (
	"errors"
	"net/http"

	authmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/auth"
	requestmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/request"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/api"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/service"
	"github.com/guluzadehh/go_eshop/services/user/internal/types"
)

type GetProfileRes struct {
	api.Response
	Data *GetProfileData `json:"data"`
}

type GetProfileData struct {
	Profile *types.ProfileView `json:"user"`
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.profile.GetUser"

	log := sl.HandlerJob(h.Log, op, requestmdw.GetRequestId(r), h.srvc)

	user := authmdw.User(r)
	if user == nil {
		log.Error("failed to get user from auth context")
		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	profile, err := h.srvc.GetProfile(r.Context(), user.Id)
	if err != nil {
		if errors.Is(err, service.ErrProfileNotFound) {
			h.JSON(w, http.StatusNotFound, api.Err("profile not found"))
			return
		}

		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	h.JSON(w, http.StatusOK, GetProfileRes{
		Response: api.Ok(),
		Data: &GetProfileData{
			Profile: types.NewProfile(profile),
		},
	})
}
