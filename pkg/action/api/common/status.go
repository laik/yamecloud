package common

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	message = "message"
	data    = "data"
	errors  = "errors"
)

func RequestParametersError(g *gin.Context, err error) {
	g.JSON(http.StatusBadRequest,
		gin.H{data: err.Error(), message: err.Error(), errors: err.Error()},
	)
	g.Abort()
}

func InternalServerError(g *gin.Context, _data interface{}, err error) {
	g.JSON(http.StatusInternalServerError,
		gin.H{data: _data, message: err.Error(), errors: err.Error()},
	)
	g.Abort()
}

func Unauthorized(g *gin.Context, _data interface{}) {
	g.JSON(http.StatusUnauthorized,
		gin.H{data: _data, message: "user unauthorized"},
	)
	g.Abort()
}

func Forbidden(g *gin.Context, _data interface{}) {
	g.JSON(http.StatusForbidden,
		gin.H{data: _data, message: "not allow to access"},
	)
	g.Abort()
}

func ResponseJSON(responseWriter http.ResponseWriter, status int, _data interface{}, _message string) {
	responseWriter.WriteHeader(status)
	responseData := make(map[string]interface{}, 6)
	responseData[data] = _data
	responseData[message] = _message
	bytes, err := json.Marshal(responseData)
	if err != nil {
		fmt.Printf("marshal data error (%s)", err)
	}
	responseWriter.Write(bytes)
}

func ResponseJSONFromError(responseWriter http.ResponseWriter, status int, _data interface{}, _error error) {
	responseWriter.WriteHeader(status)
	responseData := make(map[string]interface{}, 6)
	responseData[data] = _data
	responseData[message] = _error.Error()
	responseData[errors] = _error.Error()
	bytes, err := json.Marshal(responseData)
	if err != nil {
		fmt.Printf("marshal data error (%s)", err)
	}
	responseWriter.Write(bytes)
}
