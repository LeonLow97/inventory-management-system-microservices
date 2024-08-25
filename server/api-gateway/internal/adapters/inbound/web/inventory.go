package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/LeonLow97/internal/adapters/inbound/web/dto"
	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services/inventory"
	"github.com/LeonLow97/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/status"
)

type InventoryHandler struct {
	InventoryService inventory.Inventory
}

func NewInventoryHandler(inventoryService inventory.Inventory) *InventoryHandler {
	return &InventoryHandler{
		InventoryService: inventoryService,
	}
}

func (h *InventoryHandler) GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		}

		domainResp, err := h.InventoryService.GetProducts(c, userID)
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 5:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}
		}

		resp := &dto.GetProductsResponse{
			Products: make([]dto.Product, len(*domainResp)),
		}
		for i, product := range *domainResp {
			resp.Products[i] = dto.Product{
				BrandName:    product.BrandName,
				CategoryName: product.CategoryName,
				ProductName:  product.ProductName,
				Description:  product.Description,
				Size:         product.Size,
				Color:        product.Color,
				Quantity:     product.Quantity,
				CreatedAt:    product.CreatedAt,
				UpdatedAt:    product.UpdatedAt,
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "products": resp})
	}
}

func (h *InventoryHandler) GetProductByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		}

		productIDString := c.Param("id")
		productID, err := strconv.Atoi(productIDString)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		domainResp, err := h.InventoryService.GetProductByID(c, userID, productID)
		if err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 3:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				case 5:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}
		}

		resp := &dto.GetProductByIDResponse{
			Product: dto.Product{
				BrandName:    domainResp.BrandName,
				CategoryName: domainResp.CategoryName,
				ProductName:  domainResp.ProductName,
				Description:  domainResp.Description,
				Size:         domainResp.Size,
				Color:        domainResp.Color,
				Quantity:     domainResp.Quantity,
				CreatedAt:    domainResp.CreatedAt,
				UpdatedAt:    domainResp.UpdatedAt,
			},
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "product": resp})
	}
}

func (h *InventoryHandler) CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		}

		var req dto.CreateProductRequest
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		// validate http request json property values
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", err.Error())})
			return
		}

		product := domain.Product{
			BrandName:    req.BrandName,
			CategoryName: req.CategoryName,
			ProductName:  req.ProductName,
			Description:  req.Description,
			Size:         req.Size,
			Color:        req.Color,
			Quantity:     req.Quantity,
		}

		if err := h.InventoryService.CreateProduct(c, product, userID); err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 3:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				case 5:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}
		}
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": fmt.Sprintf("Successfully created %s", req.ProductName)})
	}
}

func (h *InventoryHandler) UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		}

		productIDString := c.Param("id")
		productID, err := strconv.Atoi(productIDString)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		var req dto.UpdateProductRequest
		if err := c.BindJSON(&req); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad Request"})
			return
		}

		product := domain.Product{
			BrandName:    req.BrandName,
			CategoryName: req.CategoryName,
			ProductName:  req.ProductName,
			Description:  req.Description,
			Size:         req.Size,
			Color:        req.Color,
			Quantity:     req.Quantity,
		}

		if err := h.InventoryService.UpdateProduct(c, product, userID, productID); err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 3:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				case 5:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}
		}
		c.JSON(http.StatusNoContent, gin.H{"status": http.StatusNoContent, "message": "Updated Product!"})
	}
}

func (h *InventoryHandler) DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		}

		productIDString := c.Param("id")
		productID, err := strconv.Atoi(productIDString)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		if err := h.InventoryService.DeleteProduct(c, userID, productID); err != nil {
			if status, ok := status.FromError(err); ok {
				errorCode := status.Code()
				switch int32(errorCode) {
				case 5:
					c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": fmt.Sprintf("Bad Request: %s", status.Message())})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				}
				return
			} else {
				log.Println("Unable to retrieve error status", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
				return
			}
		}
		c.JSON(http.StatusNoContent, gin.H{"status": http.StatusNoContent, "message": "Deleted Product!"})
	}
}
