package user

import (
	"MarketplaceAPI/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(rg *gin.RouterGroup) {
	userRouter := rg.Group("/user")

	// TODO: somehow this one is not working, the metrics can't be found inside grafana, need further check
	// utils.NewRoute(userRouter, "/register", "POST", registerHandler)
	// utils.NewRoute(userRouter, "/login", "POST", loginHandler)

	userRouter.POST("/register", registerHandler)
	userRouter.POST("/login", loginHandler)
}

// Handler for endpoint /register POST
func registerHandler(c *gin.Context) {
	var userRequest User
	ctx := c.Request.Context()
	c.BindJSON(&userRequest)
	// Validate the input, when found the error, return
	if err := ValidateUser(userRequest); err != nil {
		resp := utils.Response(err.Error(), gin.H{
			"username": userRequest.Username,
			"name":     userRequest.Name,
			"password": userRequest.Password,
		})
		c.IndentedJSON(http.StatusBadRequest, resp)

		return
	}

	user, statusCode, err := registerUser(ctx, userRequest.Username, userRequest.Name, userRequest.Password)
	if err != nil {
		resp := utils.Response(err.Error(), gin.H{
			"username": userRequest.Username,
			"name":     userRequest.Name,
			"password": userRequest.Password,
		})
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("User registered successfully", gin.H{
		"username":    user.Username,
		"name":        user.Name,
		"accessToken": user.generateToken(),
	})

	c.IndentedJSON(statusCode, resp)
}

// Handler for endpoint /login POST
func loginHandler(c *gin.Context) {
	var userRequest User
	ctx := c.Request.Context()
	c.BindJSON(&userRequest)

	if err := ValidateUserLogin(userRequest); err != nil {
		resp := utils.Response(err.Error(), gin.H{
			"username": userRequest.Username,
			"password": userRequest.Password,
		})
		c.IndentedJSON(http.StatusBadRequest, resp)

		return
	}

	user, statusCode, err := loginUser(ctx, userRequest.Username, userRequest.Password)
	if err != nil {
		resp := utils.Response(err.Error(), gin.H{
			"username": userRequest.Username,
			"password": userRequest.Password,
		})
		c.IndentedJSON(statusCode, resp)

		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "User logged successfully",
		"data": gin.H{
			"username":    user.Username,
			"name":        user.Name,
			"accessToken": user.generateToken(),
		},
	})
}
