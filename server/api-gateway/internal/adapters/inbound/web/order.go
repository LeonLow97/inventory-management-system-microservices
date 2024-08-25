package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/LeonLow97/internal/adapters/inbound/web/dto"
	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/services/order"
	"github.com/LeonLow97/internal/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	OrderService order.Order
}

func NewOrderHandler(orderService order.Order) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
	}
}

func (h *OrderHandler) GetOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		}

		domainResp, err := h.OrderService.GetOrders(c, userID)
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

		resp := &dto.GetOrdersResponse{
			Orders: make([]dto.Order, len(*domainResp)),
		}
		for i, order := range *domainResp {
			resp.Orders[i] = dto.Order{
				OrderID:      order.OrderID,
				ProductID:    order.ProductID,
				ProductName:  order.ProductName,
				CustomerName: order.CustomerName,
				BrandName:    order.BrandName,
				CategoryName: order.CategoryName,
				Color:        order.Color,
				Size:         order.Size,
				Quantity:     order.Quantity,
				Description:  order.Description,
				Revenue:      order.Revenue,
				Cost:         order.Cost,
				Profit:       order.Profit,
				HasReviewed:  order.HasReviewed,
				Status:       order.Status,
				StatusReason: order.StatusReason,
				CreatedAt:    order.CreatedAt,
			}
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "orders": resp})
	}
}

func (h *OrderHandler) GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		}

		orderIDString := c.Param("id")
		orderID, err := strconv.Atoi(orderIDString)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Internal Server Error"})
			return
		}

		domainResp, err := h.OrderService.GetOrder(c, userID, orderID)
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

		resp := &dto.GetOrderResponse{
			Order: dto.Order{
				OrderID:      domainResp.OrderID,
				ProductID:    domainResp.ProductID,
				ProductName:  domainResp.ProductName,
				CustomerName: domainResp.CustomerName,
				BrandName:    domainResp.BrandName,
				CategoryName: domainResp.CategoryName,
				Color:        domainResp.Color,
				Size:         domainResp.Size,
				Quantity:     domainResp.Quantity,
				Description:  domainResp.Description,
				Revenue:      domainResp.Revenue,
				Cost:         domainResp.Cost,
				Profit:       domainResp.Profit,
				HasReviewed:  domainResp.HasReviewed,
				Status:       domainResp.Status,
				StatusReason: domainResp.StatusReason,
				CreatedAt:    domainResp.CreatedAt,
			},
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "order": resp})
	}
}

func (h *OrderHandler) CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Unauthorized"})
			return
		}

		var req dto.CreateOrderRequest
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

		order := domain.Order{
			CustomerName: req.CustomerName,
			ProductName:  req.ProductName,
			BrandName:    req.BrandName,
			CategoryName: req.CategoryName,
			Color:        req.Color,
			Size:         req.Size,
			Quantity:     req.Quantity,
			Description:  req.Description,
			Revenue:      float32(req.Revenue),
			Cost:         float32(req.Cost),
			Profit:       float32(req.Profit),
			HasReviewed:  req.HasReviewed,
		}

		if err := h.OrderService.CreateOrder(c, order, userID); err != nil {
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
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Successfully created order!"})
	}
}
