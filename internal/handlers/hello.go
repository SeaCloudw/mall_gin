package handlers

//handlers处理具体的实现方法,可以进一步封装，比如使用server的函数
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetIndexPage(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
func GetProductsPage(c *gin.Context) {
	c.HTML(200, "productdetail.html", nil)
}
func GetHello(c *gin.Context) {
	c.String(200, "hello world")

}

// 定义返回的结构体
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func GetValidate(c *gin.Context) {
	// 返回一个JSON响应
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "validate success",
	})
}

func GetHome(c *gin.Context) {

	c.HTML(200, "hello.html", gin.H{
		"Name": "Alan",
	})
}
func GetUsers(c *gin.Context) {
	c.HTML(200, "users.html", nil)
}
func GetUsersList(c *gin.Context) {
	c.HTML(200, "list.html", nil)
}
func GetProfileList(c *gin.Context) {
	c.HTML(200, "profile.html", nil)
}
func GetOrdersPage(c *gin.Context) {
	c.HTML(200, "orders.html", nil)
}
func GetCartPage(c *gin.Context) {
	c.HTML(200, "cart.html", nil)
}

// func GetUsersData(c *gin.Context) {
// 	users := []map[string]interface{}{
// 		{"id": 1, "username": "John Doe", "email": "john@example.com"},
// 		{"id": 2, "username": "Jane Smith", "email": "jane@example.com"},
// 		// 更多用户数据...
// 	}
// 	log.Println("%v\n", users)
// 	c.JSON(http.StatusOK, users)
// }

func GetHomepage(c *gin.Context) {
	c.HTML(200, "base.html", gin.H{
		"template": "home.html",
	})
}
