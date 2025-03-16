package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mall_gin/internal/database"
	"mall_gin/internal/models"
	"mall_gin/internal/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//这是货物的处理
//handlers处理具体的实现方法,可以进一步封装，比如使用server的函数

// ProductQueryRequest 定义了查询商品的请求结构体
type ProductQueryRequest struct {
	Page      int    `json:"page" form:"page"`           // 当前页码，默认为 1
	PageSize  int    `json:"pageSize" form:"pageSize"`   // 每页显示条数，默认为 10
	SortField string `json:"sortField" form:"sortField"` // 排序字段
	SortOrder string `json:"sortOrder" form:"sortOrder"` // 排序顺序，升序（asc）或降序（desc）
	// Filters          []Filter `json:"filters" form:"filters"`         // 过滤条件
	SearchText          string   `json:"searchText" form:"searchText"`                   // 全文搜索关键词
	Fields              []string `json:"fields" form:"fields"`                           // 需要返回的字段
	IncludeAssociations []string `json:"includeAssociations" form:"includeAssociations"` // 需要预加载的关联数据
}

func GetAllProductsHandler(c *gin.Context) {

	products, err := database.GetAllModelWithPreload[models.Product](database.PreloadProductCategory, database.PreloadProductSupplier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   products,
	})
}
func FindbynameHandler(c *gin.Context) {
	var query ProductQueryRequest
	if err := c.ShouldBindBodyWithJSON(&query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	queryresult, err := database.FindByAttributeLikeUnion[*models.Product]("ProductName", query.SearchText, database.PreloadProductCategory)
	// queryresult, err := database.FindByAttributeLike[*models.Product]("ProductName", query.SearchText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return

	}
	c.JSON(200, queryresult)

}

func FindbynameHandler_Store_Hotproduct(c *gin.Context) {
	var query ProductQueryRequest
	if err := c.ShouldBindBodyWithJSON(&query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	queryresult, err := database.FindByAttributeLikeUnion[*models.Product]("ProductName", query.SearchText, database.PreloadProductCategory)
	// queryresult, err := database.FindByAttributeLike[*models.Product]("ProductName", query.SearchText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return

	}
	redisClient := redis.GetClient()
	zetname := "hot_products"
	for _, v := range queryresult {
		id := fmt.Sprintf("%d", v.ProductID)
		//保存到redis中zset
		err := redisClient.ZSetUpgrade(zetname, string(id), 1*time.Hour)
		if err != nil {
			log.Printf("保存到redis中zset err=%v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "save to redis zset fail",
			})
			return
		}

		cacheid := fmt.Sprintf("product:%d", v.ProductID)
		productJSON, err := json.Marshal(v)
		if err != nil {
			log.Printf("json error =%v\n", err)
		}
		err = redisClient.Set(cacheid, productJSON, 1*time.Hour)
		if err != nil {
			log.Printf("save fail err=%v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "save to redis string fail",
			})
			return
		}
	}
	c.JSON(200, queryresult)

}

func GetHotproduct(c *gin.Context) {
	redisClient := redis.GetClient()
	zetname := "hot_products"
	n := 5
	items, err := redisClient.GetFromZSET(zetname, n)
	if err != nil {
		log.Println("error")
		return
	}
	id_list := []int{}
	for _, member := range items {
		// log.Printf("member=%v", member)
		log.Printf("Member: %v, Score: %f\n", member.Member, member.Score)
		id, err := strconv.Atoi(member.Member)
		log.Printf("err=%v\n", err)
		id_list = append(id_list, (id))
	}

	var result []models.Product

	for _, id := range id_list {
		cacheid := fmt.Sprintf("product:%d", id)
		item, err := redisClient.Get(cacheid)
		if err != nil {
			log.Println(err)
			return
		}
		var product models.Product
		err = json.Unmarshal([]byte(item), &product)
		if err != nil {
			log.Printf("err=%v\n", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "failed to unmarshal customer data from Redis",
			})
			return
		}
		result = append(result, product)

	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})

}
func FindbyidHandler(c *gin.Context) {

	id := c.Param("id")
	uid, _ := strconv.Atoi(id)
	result, err := database.FindByEqualUnion[*models.Product]("ProductID", uid, database.PreloadProductCategory)
	// result, err := database.FindByIDModel[*models.Product](uint(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return

	}

	c.JSON(200, result)

}

// 复杂，待实现
func GetAllProductsPartdetailHandler(c *gin.Context) {

	// products, err := database.GetAllModel[models.Product]("CategoryID", "ProductName", "ProductDetail", "UnitPrice", "Quantity", "Category")
	// products, err := database.GetAllWithDefaultColumns[models.Product]()
	// products, err := database.GetAllModelWithPreload[models.Product](database.PreloadProductCategory, database.PreloadProductSupplier)
	products, err := database.GetPartModelWithPreload[models.Product]([]string{"Category"}, database.PreloadProductCategory, database.PreloadProductSupplier)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	// products[0].Category.CategoryName
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   products,
	})
}

// Validate 函数用于验证请求结构体
// func Validate(v *validator.Validate, request interface{}) error {
//     return v.Struct(request)
// }
/* func GetAllUCustomersHandler(c *gin.Context) {
	customers, err := database.GetAllUCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "查询用户数据成功",
		"data":    customers,
	})
}
func PostNewCustomerHandler(c *gin.Context) {
	//接收数据、(检查参数)创建用户、保存到数据库
	var customer CustomerRequest
	//绑定json
	if err := c.ShouldBindJSON(&customer); err != nil {
		log.Printf("请求参数: %+v\n", customer)
		log.Printf("请求参数解析失败: %v\n", err) // 输出错误信息
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//(检查参数)省略

	// 创建用户
	newcus := models.Customer{
		Name:     customer.Name,
		Password: customer.Password,

		Address:     customer.Address,
		PhoneNumber: customer.PhoneNumber,
	}
	// 保存到数据库
	if err := database.CreateANewCustomer(&newcus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "客户创建成功",
		"data":    customer,
	})

}
func DeleteCustomerHandler(c *gin.Context) {
	//删除用户，接收前端数据，查找，删除

	id := c.Param("id")
	uid, _ := strconv.Atoi(id)
	if err := database.DeleteCustomer(uint(uid)); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "删除成功",
	})
	return

}
func UpgrateCustomerHandler(c *gin.Context) {
	//更新用户信息，先根据id进行查找
	//根据json开始解绑
	//检查，赋值
	//写入数据库
	id := c.Param("id")
	uid, _ := strconv.Atoi(id)
	//这个歌oldcus返回的是指针
	oldCustomer, err := database.GetCustomerByID(uint(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	var newCus CustomerRequest
	if err := c.ShouldBindBodyWithJSON(&newCus); err != nil {

		c.JSON(http.StatusBadRequest, err)
		return
	}
	// 检查并赋值,这里虽然是对指针赋值，但是只是从内存，可能是无法同步到数据库的
	if newCus.Name != "" {
		oldCustomer.Name = newCus.Name
	}
	if newCus.Password != "" {
		oldCustomer.Password = newCus.Password // 注意：实际应用中需要对密码进行加密
	}
	if newCus.Email != "" {
		oldCustomer.Email = newCus.Email
	}
	if newCus.Address != "" {
		oldCustomer.Address = newCus.Address
	}
	if newCus.PhoneNumber != "" {
		oldCustomer.PhoneNumber = newCus.PhoneNumber
	}
	//为了持久化，最后需要保存到数据库
	if err := database.UpdateCustomer(oldCustomer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "更新客户失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "客户更新成功",
		"data":    oldCustomer,
	})
}
*/
