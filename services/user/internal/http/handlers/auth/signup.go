package authhttp

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	requestmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/request"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/api"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/render"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/validators"
	"github.com/guluzadehh/go_eshop/services/user/internal/services/auth"
	"github.com/guluzadehh/go_eshop/services/user/internal/types"
)

type SignupReq struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=5,passwordpattern"`
	ConfPassword string `json:"conf_password" validate:"required,eqfield=Password"`
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

	v := validator.New()
	v.RegisterValidation("passwordpattern", validators.PasswordPatternValidator)

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
