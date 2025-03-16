package handlers

//handlers处理具体的实现方法,可以进一步封装，比如使用server的函数
import (
	"log"
	"mall_gin/internal/database"
	"mall_gin/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	// 后面是结构体标签，键值对
	Username string `json:"username" binding:"required"`
	Password string `json:"password" ` //binding:"required"
	// Password string `json:"password" binding:"required,min=6"`
	Email string `json:"email" binding:"required,email"`
}

func GetAdminPage(c *gin.Context) {
	c.HTML(200, "admin.html", nil)
}

// 查询
func GetUserByID(id uint) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func GetUsersData(c *gin.Context) {
	db := database.GetDB()
	var count int64
	db.Model(&models.User{}).Count(&count)
	result := make([]models.User, count)
	for i := 0; i < int(count); i++ {
		if user, err := GetUserByID(uint(i)); err != nil {
			log.Println("err= %v\n", err)
		} else {
			result[i] = *user

		}
	}
	// c.JSON()
	log.Println("%v\n", result)
	c.JSON(http.StatusOK, result)
}
func EasyGetUsers(c *gin.Context) {
	db := database.GetDB()

	var users []models.User
	if err := db.Select("id, username,email, created_at").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户数据失败"})
		return
	}
	// log.Println("%v\n", users)

	c.JSON(http.StatusOK, users)
}

func DeleteUser(c *gin.Context) {
	db := database.GetDB()
	uid := c.Param("id")
	result := db.Delete(&models.User{}, uid)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, result.Error)
		return
	}
	c.Status(http.StatusNoContent)

}
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	user, err := GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return

	}
	var ndate UpdateRequest
	if err := c.ShouldBindJSON(&ndate); err != nil {
		log.Printf("请求参数: %+v\n", ndate)
		log.Printf("请求参数解析失败: %v\n", err) // 输出错误信息

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// nname := c.Param("username")
	// nemail := c.Param("email")
	// npass := c.Param("password")
	// 检查
	// 检查用户名是否已存在
	var existingUser models.User
	db := database.GetDB()
	// 如果找到匹配的记录，err将为nil
	if err := db.Where("username = ?", ndate.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if err := db.Where("email = ?", ndate.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已存在"})
		return
	}
	if ndate.Password != "" {
		user.Password = ndate.Password
	}
	user.Username = ndate.Username
	user.Email = ndate.Email

	//保存到数据库
	// db := database.GetDB()
	result := db.Save(user)
	if result.Error != nil {
		c.JSON(400, result.Error)
		return
	} else {
		c.JSON(http.StatusOK, user)

	}
	return

}
