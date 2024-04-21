package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/repository"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	jwt_secret = "MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAp0r3VXHEtvmDJ47UQFci+iCu7svucsbrprQMcAUa91RO9gwmnB9QQmAZfVn2wVg16J+fKcFRmcWKJcZNqI882WTzm208n1AEuq+dMIzgNeOjSYooFnCaELAJaPVLPhXFDUzzJgno/E2omIj6SfTt9Lnrbqhgp70R1ygUSO1e+bjNSJ/QtbAqQHkVrGKe9XwpuxApiZq+t/UCoqoTAVF9wqttUyNRlRB0Zb3l8bs/673y5rU7OLsvUIBoNh87EW/Dtf+QnAseRBEou41yXqsoRBKs1n/vELDpZgAOn/jVuTIPf3taUQVc01/T2DXcVn9vZKuXv5IunIYXb9WlM2wsFaASEDHDhZWSBZe/n44hUsqeP3CS56ZpV9EeMWjdb1tJugKL4kMcMMXTw7Xz5UxVED0L3zNDAKCIySfKy379Vhr2y84G+13Oz/51T6AZm8NCVLxxYQFB/a5klKsT4KO1b+nggHV/lOTa/m4KgpYYvS64eATP31meSt7vH308IwVYGpnhc1PEeZEIQ+utKJu9HLgVqSYGmpwJakNLHmENTCvG71seSfCyXgtigooEgz8prrpDpFVRPgfSmgbpfC7KCM3CPSv8k5VHhZwGaVgeecEzrUSBxzmg7EI0b9Q7Vi8CZqISVfoXOpxIz7wXxCw+zqpjvJzZRo1S1hEeW+jiPFMCAwEAAQ=="
)

// Endpoint to register new user to database
// Must validate input params first
// (POST /register)
func (s *Server) Register(ctx echo.Context, params repository.RegistrationParam) error {
	errors := make([]string, 0)
	isNumber := false
	isCapital := false
	isSpecial := false
	errResponse := repository.ErrorResponse{}

	err := json.NewDecoder(ctx.Request().Body).Decode(&params)
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResponse)
	}

	// validate params
	if len(params.PhoneNumber) < 10 || len(params.PhoneNumber) > 13 {
		errors = append(errors, fmt.Sprintf("phone_number: %s, error: phone numbers must be at minimum 10 characters and maximum 13 characters", params.PhoneNumber))
	}

	if !strings.HasPrefix(params.PhoneNumber, "+62") {
		errors = append(errors, fmt.Sprintf("phone_number: %s, error: phone numbers must start with the Indonesia country code (+62)", params.PhoneNumber))
	}

	if len(params.FullName) < 3 || len(params.FullName) > 60 {
		errors = append(errors, fmt.Sprintf("full_name: %s, error: full name must be at minimum 3 characters and maximum 60 characters", params.FullName))
	}

	if len(params.Password) < 3 || len(params.Password) > 64 {
		errors = append(errors, fmt.Sprintf("password: %s, error: passwords must be minimum 6 characters and maximum 64 characters", params.Password))
	}

	for _, letter := range params.Password {
		if unicode.IsDigit(letter) {
			isNumber = true
			continue
		}

		if unicode.IsSymbol(letter) {
			isSpecial = true
			continue
		}

		if unicode.IsUpper(letter) {
			isCapital = true
		}
	}

	if !isNumber || !isCapital || !isSpecial {
		errors = append(errors, fmt.Sprintf("password: %s, error: passwords must contain at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters", params.Password))
	}

	if len(errors) != 0 {
		errResponse.Message = strings.Join(errors, ";")
		return ctx.JSON(http.StatusBadRequest, errResponse)
	}

	// Insert to database
	id, err := s.Repository.Register(ctx.Request().Context(), params)
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResponse)
	}

	return ctx.JSON(http.StatusOK, id)
}

// Endpoint to login
// (PUT /login)
func (s *Server) Login(ctx echo.Context, params repository.LoginParam) error {
	errResponse := repository.ErrorResponse{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&params)
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResponse)
	}

	// Check user data in database
	id, err := s.Repository.Login(ctx.Request().Context(), params)
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResponse)
	}

	err = s.Repository.IncreaseLoginCount(ctx.Request().Context(), id)
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResponse)
	}

	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(jwt_secret))
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResponse)
	}

	return ctx.JSON(http.StatusOK, tokenString)
}

// Endpoint to get my profile
// (GET /profile)
func (s *Server) GetMyProfile(ctx echo.Context) error {
	errResponse := repository.ErrorResponse{}
	reqToken := ctx.Request().Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	hmacSecret := []byte(jwt_secret)
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusForbidden, errResponse)
	}

	var id int64
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id = claims["id"].(int64)
	}

	user, err := s.Repository.GetMyProfile(ctx.Request().Context(), id)
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusForbidden, errResponse)
	}

	return ctx.JSON(http.StatusOK, &user)
}

// Endpoint to update my profile
// (PUT /profile)
func (s *Server) UpdateMyProfile(ctx echo.Context, params repository.UpdateProfileParam) error {
	errResponse := repository.ErrorResponse{}
	reqToken := ctx.Request().Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]

	hmacSecret := []byte(jwt_secret)
	token, err := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusForbidden, errResponse)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		params.ID = claims["id"].(int64)
	}

	err = json.NewDecoder(ctx.Request().Body).Decode(&params)
	if err != nil {
		errResponse.Message = err.Error()
		return ctx.JSON(http.StatusBadRequest, errResponse)
	}

	errors := make([]string, 0)
	// validate params
	if len(params.PhoneNumber) < 10 || len(params.PhoneNumber) > 13 {
		errors = append(errors, fmt.Sprintf("phone_number: %s, error: phone numbers must be at minimum 10 characters and maximum 13 characters", params.PhoneNumber))
	}

	if !strings.HasPrefix(params.PhoneNumber, "+62") {
		errors = append(errors, fmt.Sprintf("phone_number: %s, error: phone numbers must start with the Indonesia country code (+62)", params.PhoneNumber))
	}

	if len(params.FullName) < 3 || len(params.FullName) > 60 {
		errors = append(errors, fmt.Sprintf("full_name: %s, error: full name must be at minimum 3 characters and maximum 60 characters", params.FullName))
	}

	if len(errors) != 0 {
		errResponse.Message = strings.Join(errors, ";")
		return ctx.JSON(http.StatusBadRequest, errResponse)
	}

	err = s.Repository.UpdateMyProfile(ctx.Request().Context(), params)
	if err != nil {
		errResponse.Message = err.Error()
		if strings.Contains(strings.ToLower(err.Error()), "confilct") {
			ctx.JSON(http.StatusConflict, errResponse)
		} else {
			ctx.JSON(http.StatusForbidden, errResponse)
		}
	}

	return ctx.JSON(http.StatusOK, nil)
}
