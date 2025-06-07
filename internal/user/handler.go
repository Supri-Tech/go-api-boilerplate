package user

import (
	"encoding/json"
	"go-crud-api/m/internal/middleware/authjwt"
	"go-crud-api/m/pkg/responseutil"
	"net/http"
)

type Handler struct {
	Service Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{Service: svc}
}

func (hdl *Handler) Login(response http.ResponseWriter, request *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		responseutil.Error(response, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := hdl.Service.Login(request.Context(), req.Username, req.Password)
	if err != nil {
		responseutil.Error(response, http.StatusUnauthorized, err.Error())
		return
	}

	responseutil.Success(response, "Login successful", map[string]string{"token": token})
}

func (hdl *Handler) Register(response http.ResponseWriter, request *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		responseutil.Error(response, http.StatusBadRequest, "Invalid request")
		return
	}

	user := User{
		Username: req.Username,
		Password: req.Password,
		Role:     Role(req.Role),
	}

	createdUser, err := hdl.Service.Register(request.Context(), user)
	if err != nil {
		responseutil.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	json.NewEncoder(response).Encode(createdUser)
}

func (hdl *Handler) Profile(response http.ResponseWriter, request *http.Request) {
	claims, ok := request.Context().Value("user").(authjwt.UserClaims)
	if !ok {
		responseutil.Error(response, http.StatusBadRequest, "unauthorized")
		return
	}

	user, err := hdl.Service.Profile(request.Context(), claims.Username)
	if err != nil || user == nil {
		responseutil.Error(response, http.StatusNotFound, "User not found")
		return
	}

	responseutil.Success(response, "User data found", user)
}
