package profilehttp

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	authmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/auth"
	requestmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/request"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/api"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/render"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/types"
)

type MakeProfileRequest struct {
	FirstName string `json:"first_name" validate:"max=125"`
	LastName  string `json:"last_name" validate:"max=125"`
	Phone     string `json:"phone" validate:"max=20"` // Add phone validator
}

type MakeProfileResponse struct {
	api.Response
	Data *MakeProfileData `json:"data"`
}

type MakeProfileData struct {
	*types.ProfileView `json:"profile"`
}

func (h *ProfileHandler) MakeProfile(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.profile.MakeProfile"

	log := sl.HandlerJob(h.Log, op, requestmdw.GetRequestId(r), h.srvc)

	var req MakeProfileRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("can't decode json", sl.Err(err))
		h.JSON(w, http.StatusBadRequest, api.Err("failed to parse request body"))
		return
	}

	v := validator.New()
	if err := v.StructCtx(r.Context(), &req); err != nil {
		validateErr, ok := err.(validator.ValidationErrors)
		if !ok {
			log.Error("error happened while validating request", sl.Err(err))
			h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
			return
		}

		log.Info("invalid request")
		h.JSON(w, http.StatusBadRequest, api.ValidationError(validateErr))
		return
	}

	user := authmdw.User(r)
	if user == nil {
		log.Error("failed to get user from auth context")
		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	profile, err := h.srvc.MakeProfile(r.Context(), user.Id, req.FirstName, req.LastName, req.Phone)
	if err != nil {
		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	h.JSON(w, http.StatusCreated, MakeProfileResponse{
		Response: api.Ok(),
		Data: &MakeProfileData{
			ProfileView: types.NewProfile(profile),
		},
	})
}
