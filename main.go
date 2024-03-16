package main

import (
	"MarketplaceAPI/bankaccount"
	"MarketplaceAPI/payment"
	"MarketplaceAPI/product"
	"MarketplaceAPI/user"
	"MarketplaceAPI/utils"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var router = gin.Default()

func main() {
	// init DB
	urlConnDb := DbConnStr()
	if urlConnDb == "" {
		log.Fatalln("DB Connection is not right, check the credentials first in the env variables")
	}
	LoadS3FromEnv()

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
	product.SetPool(pool)
	bankaccount.SetPool(pool)
	payment.SetPool(pool)
	user.Router(v1)
	product.Router(v1)
	bankaccount.Router(v1)
	payment.Router(v1)

	// upload image route
	authorized := v1.Group("/", utils.Auth())
	authorized.POST("/image", uploadImage)

	router.Run("0.0.0.0:8000")
}

func uploadImage(c *gin.Context) {
	file, _ := c.FormFile("file")

	// Upload the file to specific dst.
	f, err := file.Open()
	if err != nil {
		log.Println(err)
		c.IndentedJSON(400, err)
		return
	}

	if file.Size < 10*1000 || file.Size > 2000*1000 {
		// validate fileSize below 10kB or 2MB
		c.IndentedJSON(400, "file size more than 2MB or below 10kB")
		return
	}

	filenameSplit := strings.Split(file.Filename, ".")
	fileExt := filenameSplit[len(filenameSplit)-1]
	uuidFilename := fmt.Sprintf("%s.%s", uuid.NewString(), fileExt)
	if strings.ToLower(fileExt) == "jpg" || strings.ToLower(fileExt) == "jpeg" {
		publicUrl, err := uploadToS3(uuidFilename, f)

		if err != nil {
			c.IndentedJSON(500, err)
			return
		}

		c.IndentedJSON(200, gin.H{
			"imageUrl": publicUrl,
		})
		return
	}

	c.IndentedJSON(400, "image must be in .jpg or .jpeg format")
	return
}

func uploadToS3(filename string, fileObj io.Reader) (string, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(S3_REGION),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(S3_ID, S3_SECRET_KEY, ""),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		// _, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(S3_BUCKET_NAME),
		Key:    aws.String(filename),
		ACL:    "public-read",
		Body:   fileObj,
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	return result.Location, nil
}
