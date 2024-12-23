package handler

import (
	"context"
	"encoding/json"
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

type UserHandler struct {
	userService  service.UserService
	cacheManager adapter.CacheManager
}

func NewUserHandler(userService service.UserService, cacheManager adapter.CacheManager) *UserHandler {
	return &UserHandler{userService, cacheManager}
}

func (h *UserHandler) GetAllUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cachedData, err := h.cacheManager.Get("users")
	if err == nil && cachedData != "" {
		var users []entity.User
		if err := json.Unmarshal([]byte(cachedData), &users); err == nil {
			return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", users))
		}
	}

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		return pkg.HandleError(c, err)
	}
	h.cacheManager.Set("users", users)

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", users))
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var req api.CreateUser

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err)
	}

	user, err := h.userService.CreateUser(ctx, req.Name, req.Email)
	if err != nil {
		return pkg.HandleError(c, err)
	}
	h.cacheManager.Delete("users")

	return c.JSON(http.StatusCreated, pkg.ResponseSuccess("User created", user))
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	user, err := h.userService.GetUserByID(ctx, uint(id))
	if err != nil {
		return pkg.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("Success", user))
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	var req api.CreateUser

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err)
	}

	user, err := h.userService.UpdateUser(ctx, uint(id), req.Name, req.Email)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	h.cacheManager.Delete("users")

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("User updated", user))
}

func (h *UserHandler) PartialUpdateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err)
	}
	var req api.PartiallyUpdateUser

	if err := c.Bind(&req); err != nil {
		return pkg.HandleError(c, err)
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return pkg.HandleError(c, err)
	}

	user, err := h.userService.PartialUpdateUser(ctx, uint(id), req.Name, req.Email)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	h.cacheManager.Delete("users")

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("User partially updated", user))
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return pkg.HandleError(c, err)
	}

	err = h.userService.DeleteUser(ctx, uint(id))
	if err != nil {
		return pkg.HandleError(c, err)
	}

	h.cacheManager.Delete("users")

	return c.JSON(http.StatusOK, pkg.ResponseSuccess("User deleted", nil))
}
