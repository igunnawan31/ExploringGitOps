package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/igunnawan31/enterprise-product/internal/dto"
	"github.com/igunnawan31/enterprise-product/internal/service"
	"github.com/igunnawan31/enterprise-product/internal/utils"
)

type ProductHandler struct {
	productService service.ProductServiceInterface
}

func NewProductHandler(productService service.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup) {
	// POST
	rg.POST("/products", h.Create)

	// GET
	rg.GET("/products", h.GetAll)
	rg.GET("/products/:id", h.GetByID)
	rg.GET("/internal/products/validation", h.GetValidationData)

	// PATCH
	rg.PATCH("/products/:id", h.Update)

	// DELETE
	rg.DELETE("/products/:id", h.Delete)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "400", "Invalid Input")
		return
	}

	product, err := h.productService.CreateProduct(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "500", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product succesfully created", product)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	search := c.Query("search")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	products, total, err := h.productService.GetAllProduct(page, limit, search)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "500", err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", gin.H{
		"products": products,
		"pagination": gin.H{
			"current_page": page,
			"limit":        limit,
			"total_data":   total,
		},
	})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetProductById(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "400", "Product not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", product)
}

func (h *ProductHandler) GetValidationData(c *gin.Context) {
	var req struct {
		ProductIDs []string `json:"product_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "400", "Invalid array of product ids")
		return
	}

	validations, err := h.productService.GetValidationData(req.ProductIDs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "500", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Array of ProductIds is Valid", validations)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "400", "Invalid Input")
		return
	}

	var updates map[string]interface{}
	reqBytes, _ := json.Marshal(req)
	json.Unmarshal(reqBytes, &updates)

	if len(updates) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "400", "No fields to update")
		return
	}

	product, err := h.productService.UpdateProduct(id, updates)
	if err != nil {
		if errors.Is(err, errors.New("Product Not Found")) {
			utils.ErrorResponse(c, http.StatusNotFound, "400", "Product not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "500", "Internal Server Error")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.productService.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, errors.New("Product Not Found")) {
			utils.ErrorResponse(c, http.StatusNotFound, "400", "Product not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "500", "Internal Server Error")
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}
