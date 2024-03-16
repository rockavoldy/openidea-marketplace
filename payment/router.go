package payment

import (
	"MarketplaceAPI/product"
	"MarketplaceAPI/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(rg *gin.RouterGroup) {
	productRouter := rg.Group("/product")

	productRouter.POST("/:product_id/buy", buyProductHandler)

}

func buyProductHandler(c *gin.Context) {
	var productUri product.ProductUri
	var paymentRequest Payment
	ctx := c.Request.Context()
	c.BindJSON(&paymentRequest)
	if err := c.ShouldBindUri(&productUri); err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Response(err.Error(), nil))
		return
	}

	paymentRequest.ProductID = productUri.ProductID
	payment, statusCode, err := buyProduct(ctx, paymentRequest)
	if err != nil {
		resp := utils.Response(err.Error(), payment)
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("Payment recorded successfully", payment)

	c.IndentedJSON(statusCode, resp)
}
