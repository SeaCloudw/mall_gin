package handlers

//这是顾客的处理
//handlers处理具体的实现方法,可以进一步封装，比如使用server的函数
import (
	"encoding/json"
	"fmt"
	"log"
	"mall_gin/internal/database"
	"mall_gin/internal/models"
	"mall_gin/internal/redis"
	"mall_gin/internal/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomerRequest 用于接收和验证客户信息 绑定json的时候使用
type CustomerRequest struct {
	Name        string `json:"name" validate:"required"`     // 用户名
	Password    string `json:"password" validate:"required"` // 密码
	Email       string `json:"email" `                       // 邮箱
	Address     string `json:"address" `                     // 地址
	PhoneNumber string `json:"phonenumber" `                 // 电话
}

/*
	 func GetAllUCustomersHandler(c *gin.Context) {
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
*/
func LoginHandler(c *gin.Context) {
	//登录，若成功，返回to
	// ken
	var request CustomerRequest
	// 输出 JSON 数据
	// ShouldBindJSON：将请求体中的 JSON 数据解析到 request 结构体中
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("请求参数: %+v\n", request)
		log.Printf("请求参数解析失败: %v\n", err) // 输出错误信息

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	log.Printf("接收到的参数: %v\n", request) // 输出错误信息
	//查询数据库返回customer_id
	item, err := database.FindByEqualUnionString[*models.Customer]("name", request.Name, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if len(item) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "查询不到用户",
		})
		return
	}
	// 需要检查用户名与密码是否匹配
	// 示例调用
	log.Printf("请求参数: %v\n", item[0]) // 输出错误信息
	if err := services.CheckPassword(request.Password, item[0].Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "密码错误",
		})
	}
	// customer_id := 1 //这里测试返回一个固定值

	token, err := services.CreateToken(int(item[0].CustomerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  token,
	})
}
func GetProfileByToken(c *gin.Context) {
	// 从上下文中获取 customer_id
	customerIdInterface, exists := c.Get("customer_id")
	if !exists {
		// 如果 customer_id 不存在，则处理错误或返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": "customer_id not found"})
		return
	}

	// 类型断言，将 interface{} 转换为实际的类型，这里假设 customer_id 是 int 类型
	customerId, ok := customerIdInterface.(int)
	if !ok {
		log.Printf("error: assert\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "type assertion failed for customer_id"})
		return
	}
	// log.Println("cus_ID=", customerId)
	customer, err := database.GetCustomerByID(uint(customerId))
	if err != nil {
		log.Printf("error: select\n")

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    customer,
	})
}
func GetProfileByToken_Redis(c *gin.Context) {
	// 从上下文中获取 customer_id
	customerIdInterface, exists := c.Get("customer_id")
	if !exists {
		// 如果 customer_id 不存在，则处理错误或返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": "customer_id not found"})
		return
	}

	// 类型断言，将 interface{} 转换为实际的类型，这里假设 customer_id 是 int 类型
	customerId, ok := customerIdInterface.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "type assertion failed for customer_id"})
		return
	}

	redisClient := redis.GetClient()

	// 构造Redis键
	cacheKey := fmt.Sprintf("customer:%d", customerId)

	// 尝试从Redis获取数据
	cachedData, err := redisClient.Get(cacheKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
	} else if cachedData == "" {
		log.Println("dont find in redis")
		// 缓存未命中，从数据库获取数据
		customer, err := database.GetCustomerByID(uint(customerId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		// 将数据序列化为JSON字符串
		customerJSON, err := json.Marshal(customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "failed to marshal customer data",
			})
			return
		}

		// 将数据存储到Redis中，并设置过期时间为1小时
		log.Println("save data in redis", customerJSON)
		err = redisClient.Set(cacheKey, string(customerJSON), 1*time.Hour)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "failed to set data in Redis",
			})
			return
		}

		// 返回数据给客户端
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    customer,
		})
	} else {
		// 缓存命中，直接返回数据
		log.Println("find in redis cachedata=###", cachedData, "###")
		var customer interface{}
		err = json.Unmarshal([]byte(cachedData), &customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "failed to unmarshal customer data from Redis",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    customer,
		})
	}
}
func UPgrateProfileByToken(c *gin.Context) {
	// 从上下文中获取 customer_id
	customerIdInterface, exists := c.Get("customer_id")
	if !exists {
		// 如果 customer_id 不存在，则处理错误或返回错误信息
		c.JSON(http.StatusInternalServerError, gin.H{"error": "customer_id not found"})
		return
	}

	// 类型断言，将 interface{} 转换为实际的类型，这里假设 customer_id 是 int 类型
	customerId, ok := customerIdInterface.(int)
	if !ok {
		log.Printf("error: assert\n")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "type assertion failed for customer_id"})
		return
	}
	// log.Println("cus_ID=", customerId)
	customer, err := database.GetCustomerByID(uint(customerId))

	var request CustomerRequest
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return

	}

	customer.Name = request.Name
	customer.Address = request.Address
	customer.Email = request.Email
	customer.PhoneNumber = request.PhoneNumber

	err = database.UpdateCustomer(customer)

	if err != nil {
		log.Printf("error: 更新数据失败\n")

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    customer,
	})
}
func GetAllCustomersHandler(c *gin.Context) {

	customers, err := database.GetAllWithDefaultColumns[models.Customer]()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   customers,
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
	log.Printf("传入参数: %+v\n", customer)
	//(检查参数)省略

	// 加密密码
	passwordHash, err := services.HashPassword(customer.Password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加密失败"})
		return
	}
	// 创建用户
	newcus := models.Customer{
		Name:     customer.Name,
		Password: passwordHash,

		Address:     customer.Address,
		PhoneNumber: customer.PhoneNumber,
	}
	// 保存到数据库

	if err := database.CreateModel[*models.Customer](&newcus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "客户创建成功",
		"data":    customer,
	})

}
func DeleteCustomerHandler(c *gin.Context) {
	//删除用户，接收前端数据，查找，删除

	id := c.Param("id")
	uid, _ := strconv.Atoi(id)
	//使用通用的方法，传入一个空的，它满足一个接口
	// customer := &models.Customer{}
	if err := database.DeleteByIDModel[*models.Customer](uint(uid)); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	//原来的方法
	// if err := database.DeleteCustomer(uint(uid)); err != nil {
	// 	c.JSON(http.StatusBadRequest, err)
	// 	return
	// }

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "删除成功",
	})

}
func UpgrateCustomerHandler(c *gin.Context) {
	//更新用户信息，先根据id进行查找
	//根据json开始解绑
	//检查，赋值
	//写入数据库
	id := c.Param("id")
	uid, _ := strconv.Atoi(id)
	//这个oldcus返回的是指针
	// oldCustomer, err := database.GetCustomerByID(uint(uid))
	oldCustomer, err := database.FindByIDModel[*models.Customer](uint(uid))
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
	if err := database.UpdateModel[*models.Customer](oldCustomer); err != nil {
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
