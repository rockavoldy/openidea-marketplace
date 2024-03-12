package utils

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const JWT_SECRET_KEY = "23hj23hk2asdhaskj"

func Response(message string, data gin.H) gin.H {
	return gin.H{
		"message": message,
		"data":    data,
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

type CustomClaims struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authTokenHeader := ctx.Request.Header.Get("Authorization")
		_, authToken, found := strings.Cut(authTokenHeader, "Bearer ")
		if !found {
			return
		}

		token, err := jwt.ParseWithClaims(authToken, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(JWT_SECRET_KEY), nil
		})

		if err != nil {
			log.Println(err)
		} else if claims, ok := token.Claims.(*CustomClaims); ok {
			log.Println(claims)
			if !token.Valid {
				return
			}
		}

		ctx.Next()
	}
}

func CreateToken(username, name string) (string, error) {
	claims := CustomClaims{
		username,
		name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET_KEY))
}

// TODO: prometheus thingies, need to check again later
// func NewRoute(rg *gin.RouterGroup, path, method string, handler gin.HandlerFunc) {
// 	rg.Handle(method, path, HandlerWithMetrics(path, method, handler))
// }

// var requestHistrogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
// 	Name:    "requests",
// 	Help:    "Histrogram for endpoint request duration",
// 	Buckets: prometheus.LinearBuckets(1, 1, 10),
// }, []string{"path", "method", "status"})

// func HandlerWithMetrics(path, method string, handler gin.HandlerFunc) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		startTime := time.Now()

// 		handler(ctx)

// 		duration := time.Since(startTime).Seconds()
// 		statusCode := fmt.Sprintf("%d", ctx.Writer.Status())

// 		requestHistrogram.WithLabelValues(path, method, statusCode).Observe(duration)
// 	}
// }
