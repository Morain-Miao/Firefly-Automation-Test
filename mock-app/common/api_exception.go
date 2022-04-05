package common

import (
	"github.com/gin-gonic/gin"
	"log"
)

type HandlerFunc func(c *gin.Context) error

type ApiException struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

func (e *ApiException) Error() string {
	return e.Message
}

func newApiException(code int, data string, message string) *ApiException {
	return &ApiException{
		Code:    code,
		Data:    data,
		Message: message,
	}
}

func Wrapper(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err error
		)
		err = handler(c)
		if err != nil {
			var apiException *ApiException
			if h, ok := err.(*ApiException); ok {
				apiException = h
			} else if e, ok := err.(error); ok {
				if gin.Mode() == "debug" {
					// 错误
					apiException = UnknownException(e.Error())
				} else {
					// 未知错误
					apiException = UnknownException(e.Error())
				}
			} else {
				apiException = ServerException(err.Error())
			}
			log.Println(c.Request.Method + " " + c.Request.URL.String())
			c.JSON(apiException.Code, apiException)
			return
		}
	}
}

//func ErrHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.Next()
//		if length := len(c.Errors); length > 0 {
//			e := c.Errors[length-1]
//			err := e.Err
//			if err != nil {
//				var Err *ApiException
//				if e, ok := err.(*ApiException); ok {
//					Err = e
//				} else if e, ok := err.(error); ok {
//					Err = UnknownException(e.Error())
//				} else {
//					Err = ServerException(err.Error())
//				}
//				// 记录一个错误的日志
//				c.JSON(Err.Code, Err)
//				return
//			}
//		}
//
//	}
//}

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var Err *ApiException
				if e, ok := err.(*ApiException); ok {
					Err = e
				} else if e, ok := err.(error); ok {
					Err = UnknownException(e.Error())
				} else {
					Err = ServerException(e.Error())
				}
				// 记录一个错误的日志
				c.JSON(Err.Code, Err)
				return
			}
		}()
		c.Next()
	}
}
