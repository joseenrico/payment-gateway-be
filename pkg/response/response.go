package response

import "github.com/gin-gonic/gin"

type Response struct {
	ResponseCode    string      `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	Data            interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Error           string `json:"error,omitempty"`
}

func Success(c *gin.Context, code int, responseCode, message string, data interface{}) {
	c.JSON(code, Response{
		ResponseCode:    responseCode,
		ResponseMessage: message,
		Data:            data,
	})
}

func Error(c *gin.Context, code int, responseCode, message, errorDetail string) {
	c.JSON(code, ErrorResponse{
		ResponseCode:    responseCode,
		ResponseMessage: message,
		Error:           errorDetail,
	})
}

const (
	CodeSuccess            = "2004700" 
	CodePaymentSuccess     = "2005100"
	CodeBadRequest         = "4000000"
	CodeUnauthorized       = "4010000"
	CodeNotFound           = "4040000"
	CodeInternalError      = "5000000"
)