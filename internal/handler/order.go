package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/farisarmap/dot-backend-freelance/internal/adapter"
	"github.com/farisarmap/dot-backend-freelance/internal/api"
	"github.com/farisarmap/dot-backend-freelance/internal/entity"
	"github.com/farisarmap/dot-backend-freelance/internal/service"
	"github.com/farisarmap/dot-backend-freelance/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService service.OrderService
	cacheManager adapter.CacheManager
}

func NewOrderHandler(orderService service.OrderService, cacheManager adapter.CacheManager) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		cacheManager: cacheManager,
	}
}

func (h *OrderHandler) GetAllOrders(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cachedData, err := h.cacheManager.Get("orders")
	if err == nil && cachedData != "" {
		var orders []entity.Order
		if err := json.Unmarshal([]byte(cachedData), &orders); err == nil {
			return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", orders))
		}
	}
	orders, err := h.orderService.GetAllOrders(ctx)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	if cacheErr := h.cacheManager.Set("orders", orders); cacheErr != nil {
		log.Printf("cache error: %v", cacheErr)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", orders))
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var req api.CreateOrder

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err)
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err)
	}

	order, err := h.orderService.CreateOrder(ctx, req.OrderName, req.UserID)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	h.cacheManager.Delete("orders")

	return c.JSON(http.StatusCreated, pkg.ResponseSuccess("Order created", order))
}

func (h *OrderHandler) GetOrderByID(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	cacheKey := "order:" + idParam
	cachedData, err := h.cacheManager.Get(cacheKey)
	if err == nil && cachedData != "" {
		return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", cachedData))
	}

	order, oErr := h.orderService.GetOrderByID(ctx, uint(id))
	if oErr != nil {
		return pkg.HandleError(c, oErr)
	}

	h.cacheManager.Set(cacheKey, order)

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", order))
}

func (h *OrderHandler) UpdateOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	var req api.CreateOrder

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err)
	}

	order, uErr := h.orderService.UpdateOrder(ctx, uint(id), req.OrderName, req.UserID)
	if uErr != nil {
		return pkg.HandleError(c, uErr)
	}
	h.cacheManager.Delete("orders")
	cacheKey := "order:" + idParam
	h.cacheManager.Delete(cacheKey)

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Order updated", order))
}

func (h *OrderHandler) PartialUpdateOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	var req api.PartiallyUpdateOrder

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err)
	}

	order, pErr := h.orderService.PartialUpdateOrder(ctx, uint(id), req.OrderName, req.UserID)
	if pErr != nil {
		return pkg.HandleError(c, pErr)
	}

	h.cacheManager.Delete("orders")
	cacheKey := "order:" + idParam
	h.cacheManager.Delete(cacheKey)

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Order partially updated", order))
}

func (h *OrderHandler) DeleteOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	if dErr := h.orderService.DeleteOrder(ctx, uint(id)); dErr != nil {
		return pkg.HandleError(c, dErr)
	}

	h.cacheManager.Delete("orders")
	cacheKey := "order:" + idParam
	h.cacheManager.Delete(cacheKey)

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Order deleted", nil))
}

func (h *OrderHandler) CreateUserAndOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var req api.CreateUserAndOrderRequest

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err)
	}

	if err := h.orderService.CreateUserAndOrder(ctx, req); err != nil {
		return pkg.HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, pkg.ResponseSuccess("User and Order created successfully", nil))
}
