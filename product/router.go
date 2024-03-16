package product

import (
	"MarketplaceAPI/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductUri struct {
	ProductID int `uri:"product_id" binding:"required"`
}

func Router(rg *gin.RouterGroup) {
	productRouter := rg.Group("/product")

	productRouter.POST("/", createHandler)
	productRouter.PATCH("/:product_id", patchHandler)
	productRouter.DELETE("/:product_id", deleteHandler)

	productRouter.GET("/:product_id", getProductHandler)
	productRouter.GET("/", fetchProductHandler)
	productRouter.POST("/:product_id/stock", updateStockHandler)
}

func createHandler(c *gin.Context) {
	var productRequest Product
	ctx := c.Request.Context()
	c.BindJSON(&productRequest)

	// TODO: validate product input

	product, statusCode, err := createProduct(ctx, productRequest)
	if err != nil {
		resp := utils.Response(err.Error(), gin.H{
			"name":          productRequest.Name,
			"price":         productRequest.Price,
			"imageUrl":      productRequest.ImageURL,
			"stock":         productRequest.Stock,
			"condition":     productRequest.Condition,
			"isPurchasable": productRequest.IsPurchasable,
		})
		c.IndentedJSON(statusCode, resp)

		return
	}

	// TODO: create product to repo
	resp := utils.Response("Product added successfully", product)

	c.IndentedJSON(statusCode, resp)

}

func patchHandler(c *gin.Context) {
	var productRequest map[string]any
	var productUri ProductUri
	ctx := c.Request.Context()
	c.BindJSON(&productRequest)
	if err := c.ShouldBindUri(&productUri); err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Response(err.Error(), nil))
		return
	}

	product, statusCode, err := GetProduct(ctx, productUri.ProductID)
	if err != nil {
		resp := utils.Response(err.Error(), productRequest)
		c.IndentedJSON(statusCode, resp)

		return
	}
	product.patchWith(productRequest)

	product, statusCode, err = PatchProduct(ctx, product)
	if err != nil {
		resp := utils.Response(err.Error(), productRequest)
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("Product edited successfully", product)

	c.IndentedJSON(statusCode, resp)
}

func deleteHandler(c *gin.Context) {
	var productUri ProductUri
	ctx := c.Request.Context()
	if err := c.ShouldBindUri(&productUri); err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Response(err.Error(), nil))
		return
	}

	_, statusCode, err := deleteProduct(ctx, productUri.ProductID)
	if err != nil {
		resp := utils.Response(err.Error(), nil)
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("Product deleted successfully", nil)

	c.IndentedJSON(statusCode, resp)
}

func getProductHandler(c *gin.Context) {
	var productUri ProductUri
	ctx := c.Request.Context()
	if err := c.ShouldBindUri(&productUri); err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Response(err.Error(), nil))
		return
	}

	product, statusCode, err := GetProduct(ctx, productUri.ProductID)
	if err != nil {
		resp := utils.Response(err.Error(), nil)
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("Get product", product)

	c.IndentedJSON(statusCode, resp)
}

func fetchProductHandler(c *gin.Context) {
	ctx := c.Request.Context()
	products, statusCode, err := listProducts(ctx)
	if err != nil {
		resp := utils.Response(err.Error(), nil)
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("Products list", products)

	c.IndentedJSON(statusCode, resp)
}

func updateStockHandler(c *gin.Context) {
	var productUri ProductUri
	var productRequest map[string]uint
	c.BindJSON(&productRequest)
	ctx := c.Request.Context()
	if err := c.ShouldBindUri(&productUri); err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Response(err.Error(), nil))
		return
	}

	product, statusCode, err := GetProduct(ctx, productUri.ProductID)
	if err != nil {
		resp := utils.Response(err.Error(), nil)
		c.IndentedJSON(statusCode, resp)

		return
	}

	product.AdjustStock(productRequest["stock"])
	if err != nil {
		resp := utils.Response(err.Error(), nil)
		c.IndentedJSON(statusCode, resp)

		return
	}

	product, statusCode, err = PatchProduct(ctx, product)
	if err != nil {
		resp := utils.Response(err.Error(), nil)
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("Product stock updated", product)

	c.IndentedJSON(statusCode, resp)
}
