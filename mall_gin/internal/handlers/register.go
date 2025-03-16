package handlers

import (
	"log"
	"net/http"

	"mall_gin/internal/database"
	"mall_gin/internal/models"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	// 后面是结构体标签，键值对
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	// Password string `json:"password" binding:"required,min=6"`
	Email string `json:"email" binding:"required,email"`
}

func GetBarPage(c *gin.Context) {
	c.HTML(200, "bar.html", nil)

}
func GetTestPage(c *gin.Context) {
	c.HTML(200, "test.html", nil)

}
func GetLoginPage(c *gin.Context) {
	c.HTML(200, "login.html", nil)

}

func GetRegisterPage(c *gin.Context) {
	c.HTML(200, "register.html", nil)

}

// 处理前端传来的JSON数据
func Register(c *gin.Context) {
	var request RegisterRequest
	// 输出 JSON 数据
	// ShouldBindJSON：将请求体中的 JSON 数据解析到 request 结构体中
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("请求参数: %+v\n", request)
		log.Printf("请求参数解析失败: %v\n", err) // 输出错误信息

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	db := database.GetDB()
	if err := db.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", request.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已存在"})
		return
	}

	// 加密密码
	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
	// 	return
	// }

	// 创建用户
	user := models.User{
		Username: request.Username,
		// Password: string(hashedPassword),
		Password: request.Password,
		Email:    request.Email,
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
