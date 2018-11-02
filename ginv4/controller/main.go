package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorType int

const (
	Content_Type ErrorType = iota
	Session_Id
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var router *gin.Engine

func main() {
	router = gin.Default()
	initializeRoutes()
	router.Run()
}

func initializeRoutes() {
	router.POST("/api", handleVerification)
	router.OPTIONS("/api", handleVerification)
	router.GET("/api", handleGet)
}

func handleGet(c *gin.Context) {
	message, _ := c.GetQuery("m")
	c.String(http.StatusOK, "Get works! you sent: "+message)
}

func handleVerification(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		// setup headers
		c.Header("Allow", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "origin, content-type, accept")
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusOK)
	} else if c.Request.Method == "POST" {
		var checkDefault = checkHeader(c, []ErrorType{Content_Type, Session_Id})
		if !checkDefault.status {

			c.JSON(checkDefault.httpStatus, gin.H{
				"statuscode": checkDefault.internalStatus,
				"statustext": checkDefault.text,
				"response":   nil,
			})
			return
		}
		var u User
		c.BindJSON(&u)
		c.JSON(checkDefault.httpStatus, gin.H{
			"user": u.Username,
			"pass": u.Password,
		})
	}
}

func checkHeader(c *gin.Context, headerType []ErrorType) ErrorMessageDefault {
	var value = getDefaultErorMessage()
	for _, element := range headerType {
		value = getErrorDefault(element, c)
		if !value.status {
			break
		}
	}

	return value
}

func getErrorDefault(valueType ErrorType, c *gin.Context) ErrorMessageDefault {
	var value = getDefaultErorMessage()

	switch valueType {
	case Content_Type:
		if c.GetHeader("Content-Type") != "application/json" {
			value.status = false
			value.httpStatus = http.StatusUnauthorized
			value.internalStatus = 10001
			value.text = "Not allow Content-Type"
		}
	case Session_Id:
		if c.GetHeader("sessionid") == "" {
			value.status = false
			value.httpStatus = http.StatusUnauthorized
			value.internalStatus = 10002
			value.text = "press in insert token"
		}
	default:
	}
	return value
}
func getDefaultErorMessage() ErrorMessageDefault {
	var message ErrorMessageDefault
	message.httpStatus = http.StatusOK
	message.internalStatus = http.StatusOK
	message.text = "Complete!"
	message.status = true
	return message
}
