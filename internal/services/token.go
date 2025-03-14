package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Secret key for signing JWT tokens
var jwtKey = []byte("7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069")

// Claims represents the JWT claims (payload)
type Claims struct {
	Customer_id int `json:"customer_id"`
	jwt.RegisteredClaims
}

// CreateToken generates a JWT token for the given customer_id
func CreateToken(customer_id int) (string, error) {
	expirationTime := time.Now().Add(50 * time.Minute)
	claims := &Claims{
		Customer_id: customer_id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// log.Println("create  time id =", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies the JWT token and extracts the claims
func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	if !token.Valid {
		log.Printf("Token is invalid: %+v", token)
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// AuthMiddleware verifies the JWT token for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// 检查 Authorization 头是否有效
		if len(tokenString) < 8 || tokenString[:7] != "Bearer " {
			// 如果缺少或无效的 Authorization 头，重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 提取 Token
		tokenString = tokenString[7:]

		// 验证 Token
		claims, err := VerifyToken(tokenString)
		if err != nil {
			// 如果 Token 验证失败，重定向到登录页面
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		log.Println("tokn=", claims)
		// 将用户信息存储在上下文中
		c.Set("customer_id", claims.Customer_id)

		// 继续处理请求
		c.Next()
	}
}

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func GetIDbyToken(c *gin.Context) (int, error) {
	// 从上下文中获取 customer_id
	customerIdInterface, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "customer_id not found"})
		return 0, &AppError{Code: http.StatusInternalServerError, Message: "customer_id not found"}
	}

	// 类型断言，将 interface{} 转换为实际的类型，这里假设 customer_id 是 int 类型
	customerId, ok := customerIdInterface.(int)
	if !ok {
		log.Printf("error: type assertion failed for customer_id\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "type assertion failed for customer_id"})
		return 0, &AppError{Code: http.StatusInternalServerError, Message: "type assertion failed for customer_id"}
	}

	return customerId, nil
}
