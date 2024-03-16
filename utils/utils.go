package utils

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func jwtSecretKey() string {
	jwtSeretKey := os.Getenv("JWT_SECRET")
	if jwtSeretKey == "" {
		return "justrandomstringfordefaultpurposes"
	}

	return jwtSeretKey
}

var JWT_SECRET_KEY = jwtSecretKey()

func Response(message string, data any) gin.H {
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
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authTokenHeader := ctx.Request.Header.Get("Authorization")
		_, authToken, found := strings.Cut(authTokenHeader, "Bearer ")
		if !found {
			ctx.AbortWithStatus(401)
			return
		}

		token, err := jwt.ParseWithClaims(authToken, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(JWT_SECRET_KEY), nil
		})

		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(401)
			return
		} else if claims, ok := token.Claims.(*CustomClaims); ok {
			if !token.Valid {
				ctx.AbortWithStatus(401)
				return
			}
			// when token still valid, pass userId via context
			ctx.Request = ctx.Request.Clone(context.WithValue(ctx.Request.Context(), "userId", claims.UserID))
		}

		ctx.Next()
	}
}

func CreateToken(user_id int, username, name string) (string, error) {
	claims := CustomClaims{
		user_id,
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
