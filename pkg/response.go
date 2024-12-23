package pkg

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseSuccess(message string, data interface{}) Response {
	return Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func ResponseError(message string, data interface{}) Response {
	return Response{
		Status:  "error",
		Message: message,
		Data:    data,
	}
}

func HandleError(c echo.Context, err error) error {
	if ve, ok := err.(validator.ValidationErrors); ok {
		var fieldErrors []string
		for _, fe := range ve {
			fieldErrors = append(fieldErrors, fe.Field()+" failed on the '"+fe.Tag()+"' tag")
		}
		resp := ResponseError("Validation error", fieldErrors)
		return c.JSON(http.StatusBadRequest, resp)
	}
	resp := ResponseError(err.Error(), nil)
	return c.JSON(http.StatusBadRequest, resp)
}
