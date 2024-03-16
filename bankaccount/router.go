package bankaccount

import (
	"MarketplaceAPI/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BankAccountUri struct {
	BankAccountID int `uri:"bank_account_id" binding:"required"`
}

func Router(rg *gin.RouterGroup) {
	bankAccountRouter := rg.Group("/bank/account")
	bankAccountRouter.Use(utils.Auth())

	bankAccountRouter.POST("/", createHandler)
	bankAccountRouter.GET("/", listHandler)
	bankAccountRouter.PATCH("/:bank_account_id", patchHandler)
	bankAccountRouter.DELETE("/:bank_account_id", deleteHandler)

}

func createHandler(c *gin.Context) {
	var baccountRequest BankAccount
	ctx := c.Request.Context()
	c.BindJSON(&baccountRequest)

	bankAccount, statusCode, err := createBankAccount(ctx, baccountRequest)
	if err != nil {
		resp := utils.Response(err.Error(), baccountRequest)
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("Bank account added successfully", bankAccount)

	c.IndentedJSON(statusCode, resp)
}

func listHandler(c *gin.Context) {
	ctx := c.Request.Context()

	bankAccounts, statusCode, err := listBankAccounts(ctx)
	if err != nil {
		resp := utils.Response(err.Error(), []BankAccount{})
		c.IndentedJSON(statusCode, resp)

		return
	}
	resp := utils.Response("Bank account updated successfully", bankAccounts)

	c.IndentedJSON(statusCode, resp)

}

func patchHandler(c *gin.Context) {
	var baccountRequest map[string]string
	var baccountUri BankAccountUri
	ctx := c.Request.Context()
	c.BindJSON(&baccountRequest)
	if err := c.ShouldBindUri(&baccountUri); err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Response(err.Error(), nil))
		return
	}

	bankAccount, statusCode, err := getBankAccount(ctx, baccountUri.BankAccountID)
	if err != nil {
		resp := utils.Response(err.Error(), baccountRequest)
		c.IndentedJSON(statusCode, resp)

		return
	}
	bankAccount.patchWith(baccountRequest)

	bankAccount, statusCode, err = patchBankAccount(ctx, bankAccount)
	if err != nil {
		resp := utils.Response(err.Error(), baccountRequest)
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("Bank account updated successfully", bankAccount)

	c.IndentedJSON(statusCode, resp)
}

func deleteHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var baccountUri BankAccountUri
	if err := c.ShouldBindUri(&baccountUri); err != nil {
		c.IndentedJSON(http.StatusBadRequest, utils.Response(err.Error(), nil))
		return
	}

	_, statusCode, err := deleteBankAccount(ctx, baccountUri.BankAccountID)
	if err != nil {
		resp := utils.Response(err.Error(), nil)
		c.IndentedJSON(statusCode, resp)

		return
	}

	resp := utils.Response("Bank account deleted successfully", nil)

	c.IndentedJSON(statusCode, resp)
}
