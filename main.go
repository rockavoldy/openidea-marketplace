package main

import (
	"MarketplaceAPI/user"
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var router = gin.Default()

func main() {
	// init DB
	urlConnDb := "postgres://akhmad:akhmad@localhost:5432/marketplace_api?sslmode=disable"

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, urlConnDb)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	router.Use(gin.Recovery())
	// Router initialization
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	v1 := router.Group("/v1")
	user.SetPool(pool)
	user.Router(v1)

	router.Run("0.0.0.0:8000")
}
