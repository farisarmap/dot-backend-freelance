package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/farisarmap/dot-backend-freelance/internal/api"
	"github.com/farisarmap/dot-backend-freelance/internal/service"
	"github.com/farisarmap/dot-backend-freelance/pkg"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", users))
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var req api.CreateUser

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	user, err := h.userService.CreateUser(ctx, req.Name, req.Email)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, pkg.ResponseSuccess("User created", user))
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	user, err := h.userService.GetUserByID(ctx, uint(id))
	if err != nil {
		return pkg.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", user))
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	var req api.CreateUser

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	user, err := h.userService.UpdateUser(ctx, uint(id), req.Name, req.Email)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("User updated", user))
}

func (h *UserHandler) PartialUpdateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}
	var req api.PartiallyUpdateUser

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	user, err := h.userService.PartialUpdateUser(ctx, uint(id), req.Name, req.Email)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("User partially updated", user))
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err, http.StatusBadRequest)
	}

	err = h.userService.DeleteUser(ctx, uint(id))
	if err != nil {
		return pkg.HandleError(c, err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("User deleted", nil))
}
