package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log"
	"music-echo/api/domain/dao"
	"music-echo/api/domain/dto"
	"music-echo/api/repository"
	"music-echo/utils"
	"music-echo/utils/token"
	"net/http"
	"time"
)

type UserHandler interface {
	CreateUser(e echo.Context) error
	ActivateUser(e echo.Context) error
}

type UserHandlerImpl struct {
	Validators      *validator.Validate
	UsersRepository repository.UsersRepository
	TokenRepository repository.TokenRepository
	Mailer          utils.Mailer
}

func NewUserHandlerImpl(validators *validator.Validate, usersRepository repository.UsersRepository, tokenRepository repository.TokenRepository, mailer utils.Mailer) UserHandler {
	return UserHandlerImpl{
		Validators:      validators,
		UsersRepository: usersRepository,
		TokenRepository: tokenRepository,
		Mailer:          mailer,
	}
}

func (u UserHandlerImpl) CreateUser(e echo.Context) error {
	// read and bind request body
	userRequest := new(dto.UserPostRequest)
	err := utils.ReadJSON(e, userRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// validate request body
	err = u.Validators.Struct(userRequest)
	if err != nil {
		var validationErrors validator.ValidationErrors

		ok := errors.As(err, &validationErrors)
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		errorMap := make(map[string]string)
		for i := 0; i < len(validationErrors); i++ {
			errorMap[validationErrors[i].Field()] = getValidationMessage(validationErrors[i])
		}

		return echo.NewHTTPError(http.StatusNotAcceptable, errorMap)
	}

	// Insert user and token
	users := &dao.Users{
		Email: userRequest.Email,
		Name:  userRequest.Name,
	}
	err = users.Password.Set(userRequest.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = u.UsersRepository.Insert(e.Request().Context(), users)
	if err != nil {
		switch {
		case err.Error() == "duplicate email":
			return echo.NewHTTPError(http.StatusNotAcceptable, "email already exist")
		default:
			return echo.NewHTTPError(http.StatusConflict, err)
		}
	}

	tokens, plainText, err := token.GenerateToken(users.Id, 3*24*time.Hour, token.ScopeActivation)
	if err != nil {
		return err
	}

	err = u.TokenRepository.Insert(e.Request().Context(), tokens)
	if err != nil {
		return err
	}

	// Send email
	utils.Background(func() {
		data := map[string]any{
			"Id":              users.Id,
			"Name":            users.Name,
			"activationToken": plainText,
		}
		err = u.Mailer.Send(u.Mailer.Sender, "user_welcome.tmpl", data)
		if err != nil {
			log.Println(err)
		}
	})

	// Response
	usersResponse := dto.UsersCreateResponse{
		Id:        users.Id,
		CreatedAt: users.CreatedAt,
		Version:   users.Version,
	}

	response := dto.WebResponse{
		Message: "success register new users",
		Data:    usersResponse,
	}

	return e.JSON(http.StatusOK, response)
}

func (u UserHandlerImpl) ActivateUser(e echo.Context) error {
	// Read and bind json request
	userRequest := new(dto.UserRegisterActivated)
	err := utils.ReadJSON(e, userRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Validate request body
	err = u.Validators.Struct(userRequest)
	if err != nil {
		var validationErrors validator.ValidationErrors

		ok := errors.As(err, &validationErrors)
		if !ok {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		errorMap := make(map[string]string)
		for i := 0; i < len(validationErrors); i++ {
			errorMap[validationErrors[i].Field()] = getValidationMessage(validationErrors[i])
		}

		return echo.NewHTTPError(http.StatusNotAcceptable, errorMap)
	}

	// Activate user
	user, err := u.UsersRepository.GetByToken(e.Request().Context(), userRequest.Token, token.ScopeActivation)
	if err != nil {
		log.Println(err)
	}
	user.Activated = true

	err = u.UsersRepository.Update(e.Request().Context(), user)
	if err != nil {
		log.Println(err)
	}

	// Delete activation token for corresponding user (since it is no longer needed)
	err = u.TokenRepository.Delete(e.Request().Context(), user.Id, token.ScopeActivation)
	if err != nil {
		log.Println(err)
	}

	// Response
	return e.JSON(http.StatusOK, user)
}

func getValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "is required"
	case "email":
		return "is not a valid email address"
	case "min":
		return "must be at least " + fieldError.Param() + " characters long"
	case "max":
		return "must be at most " + fieldError.Param() + " characters long"
	default:
		return "is invalid"
	}
}
