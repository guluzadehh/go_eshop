package authhttp

import (
	"errors"
	"net/http"

	requestmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/request"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/api"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/render"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/services/auth"
	"github.com/guluzadehh/go_eshop/services/user/internal/types"
)

type SignupReq struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ConfPassword string `json:"conf_password"`
}

type SignupRes struct {
	api.Response
	Data *SignupData `json:"data"`
}

type SignupData struct {
	User *types.UserView `json:"user"`
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Signup"

	log := sl.HandlerJob(h.Log, op, requestmdw.GetRequestId(r), h.srvc)

	var req SignupReq
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("can't decode json", sl.Err(err))
		h.JSON(w, http.StatusBadRequest, api.Err("failed to parse request body"))
		return
	}

	// TODO: validate

	user, err := h.srvc.Signup(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrEmailExists) {
			h.JSON(w, http.StatusConflict, api.ErrD(
				"user exists",
				[]api.ErrDetail{
					{
						Field: "email",
						Info:  "email is already being used",
					},
				},
			))
			return
		}

		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	h.JSON(w, http.StatusCreated, SignupRes{
		Response: api.Ok(),
		Data: &SignupData{
			User: types.NewUser(user),
		},
	})
}
