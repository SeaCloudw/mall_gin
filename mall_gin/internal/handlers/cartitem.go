package handlers

import (
	"log"
	"mall_gin/internal/database"
	"mall_gin/internal/models"
	"mall_gin/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

//这是购物车的处理
//handlers处理具体的实现方法,可以进一步封装，比如使用server的函数

type CartItemsRequest struct {
	CustomerID uint `json:"customer_id" form:"customer_id"` // 客户ID
	ProductID  uint `json:"product_id" form:"product_id"`   // 商品ID
	// Quantity   int    `json:"quantity" form:"quantity"`       // 商品数量，默认为1
}

func CreateNewCart(c *gin.Context) {
	cusid, err := services.GetIDbyToken(c)
	if err != nil {
		log.Println("fail to get token id")
		return
	}
	//接收数据、(检查参数) 检索、保存到数据库
	var request CartItemsRequest
	//绑定json
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("请求参数: %+v\n", request)
		log.Printf("请求参数解析失败: %v\n", err) // 输出错误信息
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//(检查参数)省略

	// 创建用户
	newcart := models.CartItem{
		CustomerID: cusid,
		ProductID:  int(request.ProductID),
	}
	// 保存到数据库
	log.Println("request.ProductID=", request.ProductID)
	if err := database.CreateModel[*models.CartItem](&newcart); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加购物车失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "添加购物车成功",
		"data":    newcart,
	})

}
func GetAllCartbyToken(c *gin.Context) {
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
	//接收数据、(检查参数) 根据customerId 检索、查询数据库，返回相关数据

	//查询第一次，向CartItem查询 返回切片ProductID
	queryresult, err := database.FindByEqualUnion[*models.CartItem]("CustomerID", customerId, nil)
	// queryresult, err := database.FindByAttributeLike[*models.Product]("ProductName", query.SearchText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	// idList:=make([]int)
	var idList []int
	for _, item := range queryresult {
		idList = append(idList, item.ProductID)
	}
	//查询第二次，向Product查询，返回Product类型的切片
	dataresult, err := database.FindByEqualListUnionT[*models.Product]("ProductID", idList, database.PreloadProductCategory)
	// queryresult, err := database.FindByAttributeLike[*models.Product]("ProductName", query.SearchText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, dataresult)
}
func GetAllCartbyId(c *gin.Context) {

	//接收数据、(检查参数) 根据customerId 检索、查询数据库，返回相关数据
	var request CartItemsRequest
	//绑定json
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("请求参数: %+v\n", request)
		log.Printf("请求参数解析失败: %v\n", err) // 输出错误信息
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// s := strconv.Itoa(int(request.ProductID))
	//查询第一次，向CartItem查询 返回切片ProductID
	queryresult, err := database.FindByEqualUnion[*models.CartItem]("CustomerID", int(request.CustomerID), nil)
	// queryresult, err := database.FindByAttributeLike[*models.Product]("ProductName", query.SearchText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	// idList:=make([]int)
	var idList []int
	for _, item := range queryresult {
		idList = append(idList, item.ProductID)
	}
	//查询第二次，向Product查询，返回Product类型的切片
	dataresult, err := database.FindByEqualListUnion[*models.Product]("ProductID", idList, database.PreloadProductCategory)
	// queryresult, err := database.FindByAttributeLike[*models.Product]("ProductName", query.SearchText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, dataresult)
}
func DeleteCart(c *gin.Context) {
	//删除购物车的一项，

	//接收数据、(检查参数) 根据customerId,ProductID 检索、查询数据库，删除
	var request CartItemsRequest
	//绑定json
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("请求参数: %+v\n", request)
		log.Printf("请求参数解析失败: %v\n", err) // 输出错误信息
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//查询第一次，向CartItem查询 返回多个，这里只需要一个
	Conditions := map[string][]int{
		"CustomerID": {int(request.CustomerID)},
		"ProductID":  {int(request.ProductID)},
	}
	queryresult, err := database.FindByMultipleConditions[*models.CartItem](Conditions, nil)
	// queryresult, err := database.FindByAttributeLike[*models.Product]("ProductName", query.SearchText)
	// log.Printf("queryresult =%v \n", queryresult)
	// log.Printf("err =%v \n", err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	if len(queryresult) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "没有找到符合要求的结果",
		})
		return
	}
	//删除
	// queryresult.cart_item_id
	if err := database.DeleteByIDModel[*models.CartItem]((uint(queryresult[0].CartItemID))); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "删除成功",
	})

}
