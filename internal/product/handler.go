package product

import (
	"encoding/json"
	"go-crud-api/m/pkg/responseutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{Service: svc}
}

func (hdl *Handler) GetProduct(response http.ResponseWriter, request *http.Request) {
	products, err := hdl.Service.GetProduct(request.Context())
	if err != nil {
		responseutil.Error(response, http.StatusInternalServerError, "Failed to get products")
		return
	}
	responseutil.Success(response, "Data found", products)
}

func (hdl *Handler) GetProductByID(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responseutil.Error(response, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := hdl.Service.GetProductByID(request.Context(), id)
	if err != nil {
		responseutil.Error(response, http.StatusBadRequest, "Failed to get product")
		return
	}

	if product == nil {
		responseutil.Error(response, http.StatusNotFound, "Product not found")
		return
	}
	responseutil.Success(response, "Data found", product)
}

func (hdl *Handler) CreateProduct(response http.ResponseWriter, request *http.Request) {
	var req struct {
		ProductName  string `json:"product_name"`
		ProductPrice int64  `json:"product_price"`
		ProductStock int64  `json:"product_stock"`
	}
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		responseutil.Error(response, http.StatusBadRequest, "Invalid request body")
		return
	}

	product := Product{
		ProductName:  req.ProductName,
		ProductPrice: req.ProductPrice,
		ProductStock: req.ProductStock,
	}

	created, err := hdl.Service.AddProduct(request.Context(), product)
	if err != nil {
		responseutil.Error(response, http.StatusBadRequest, err.Error())
		return
	}

	responseutil.Success(response, "Data created", created)
}

func (hdl *Handler) UpdateProduct(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responseutil.Error(response, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req struct {
		ProductName  string `json:"product_name"`
		ProductPrice int64  `json:"product_price"`
		ProductStock int64  `json:"product_stock"`
	}

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		responseutil.Error(response, http.StatusBadRequest, "Invalid request body")
	}

	product := Product{
		ID:           id,
		ProductName:  req.ProductName,
		ProductPrice: req.ProductPrice,
		ProductStock: req.ProductStock,
	}

	updated, err := hdl.Service.UpdateProduct(request.Context(), product)
	if err != nil {
		responseutil.Error(response, http.StatusInternalServerError, err.Error())
		return
	}

	responseutil.Success(response, "Data updated", updated)
}

func (hdl *Handler) DeleteProduct(response http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		responseutil.Error(response, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = hdl.Service.DeleteProduct(request.Context(), id)
	if err != nil {
		responseutil.Error(response, http.StatusInternalServerError, err.Error())
		return
	}

	responseutil.Success(response, "Data deleted", map[string]string{"message": "Product Deleted!"})
}
