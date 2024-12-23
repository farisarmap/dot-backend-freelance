package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/farisarmap/dot-backend-freelance/internal/adapter"
	"github.com/farisarmap/dot-backend-freelance/internal/api"
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

	orders, err := h.orderService.GetAllOrders(ctx)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", orders))
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var req api.CreateOrder

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	order, err := h.orderService.CreateOrder(ctx, req.OrderName, req.UserID)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, pkg.ResponseSuccess("Order created", order))
}

func (h *OrderHandler) GetOrderByID(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	order, oErr := h.orderService.GetOrderByID(ctx, uint(id))
	if oErr != nil {
		return pkg.HandleError(c, oErr, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", order))
}

func (h *OrderHandler) UpdateOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	var req api.CreateOrder

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	order, uErr := h.orderService.UpdateOrder(ctx, uint(id), req.OrderName, req.UserID)
	if uErr != nil {
		return pkg.HandleError(c, uErr, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Order updated", order))
}

func (h *OrderHandler) PartialUpdateOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	var req api.PartiallyUpdateOrder

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	order, pErr := h.orderService.PartialUpdateOrder(ctx, uint(id), req.OrderName, req.UserID)
	if pErr != nil {
		return pkg.HandleError(c, pErr, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Order partially updated", order))
}

func (h *OrderHandler) DeleteOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	if dErr := h.orderService.DeleteOrder(ctx, uint(id)); dErr != nil {
		return pkg.HandleError(c, dErr, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Order deleted", nil))
}

func (h *OrderHandler) CreateUserAndOrder(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var req api.CreateUserAndOrderRequest

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	if err := h.orderService.CreateUserAndOrder(ctx, req); err != nil {
		return pkg.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, pkg.ResponseSuccess("User and Order created successfully", nil))
}
