package handlers

import (
	"log"
	"mall_gin/internal/database"
	"mall_gin/internal/models"
	"mall_gin/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 这是订单的处理
// handlers处理具体的实现方法,可以进一步封装，比如使用server的函数
type OrderDetailRequest struct {
	CustomerID int       `json:"customer_id" validate:"required"` // 产品ID
	OrderDate  string    `json:"order_date"`
	ProductID  []int     `json:"product_id" validate:"required"`                                      // 产品ID
	UnitPrice  []float64 `json:"unit_price" validate:"required,min=0"`                                // 单价
	Amount     []int     `json:"amount" validate:"required,min=1"`                                    // 数量
	Status     string    `json:"status" validate:"required,oneof=Pending Shipped Cancelled Returned"` // 状态
}
type ChangeOrderDetailRequest struct {
	OrderDetailID int    `json:"order_detail_id" validate:"required"`                                 // 产品ID
	Status        string `json:"status" validate:"required,oneof=Pending Shipped Cancelled Returned"` // 状态
}

func ChangeOrder(c *gin.Context) {
	var request ChangeOrderDetailRequest
	err := c.ShouldBindBodyWithJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  " bind error",
			"message": err.Error(),
		})
	}

	findresult, err := database.FindByEqualUnion[*models.OrderDetail]("OrderDetailID", request.OrderDetailID, nil)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Dont Find;No this detail",
		})
	}

	//更新状态
	findresult[0].Status = request.Status

	if err := database.UpdateModel[*models.OrderDetail](findresult[0]); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": findresult,
	})
}
func CreateOrder(c *gin.Context) {
	//接收数据、(检查参数)
	// 新建订单[用户下单-》后端接收多个n商品信息-》（检查货物剩余数量是否充足）求和，得到TotalAmount-》新建一个Orders,->保存、记录下Orders_id然后再根据商品数量，建立多个n个OrderDetail，填入OrderID、ProductID、UnitPrice、Amount、status，连接外键-》保存n个OrderDetail]
	var request OrderDetailRequest
	//绑定json
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("请求参数: %+v\n", request)
		log.Printf("请求参数解析失败: %v\n", err) // 输出错误信息
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//一次下单的数量
	n := len(request.ProductID)
	totalAmount := 0.0
	neworderdetaillist := make([]models.OrderDetail, n)

	for i := 0; i < n; i++ {
		//(检查参数)商品是否存在、商品剩余总数
		// check()

		// 创建orderdetail

		newitem := models.OrderDetail{
			ProductID: uint(request.ProductID[i]),
			UnitPrice: request.UnitPrice[i],
			Amount:    request.Amount[i],
			Status:    "Pending",
		}
		neworderdetaillist[i] = newitem

		totalAmount += (newitem.UnitPrice) * float64(newitem.Amount)
	}

	// 保存到数据库,先加入Order，再加入n条detail
	neworder := models.Order{
		// OrderID     :request.Amount[]
		CustomerID:  request.CustomerID,
		OrderDate:   request.OrderDate,
		TotalAmount: totalAmount,
	}
	newmodel, err := database.CreateModelAndReturn[*models.Order](&neworder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加订单失败"})
		return
	}
	orderindex := newmodel.OrderID

	for i := 0; i < n; i++ {

		neworderdetaillist[i].OrderID = uint(orderindex) //赋值orderid
		if err := database.CreateModel[*models.OrderDetail](&neworderdetaillist[i]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "添加订单详情失败"})
		}

	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "添加订单与订单详情成功成功",
		"data":    neworderdetaillist,
	})

}

func GetAllorderbyid(c *gin.Context) {
	//展示订单界面
	var request OrderDetailRequest
	//绑定json
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("请求参数: %+v\n", request)
		log.Printf("请求参数解析失败: %v\n", err) // 输出错误信息
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//查询第一次返回orderdetails_id[]
	results, err := database.FindByEqualUnion[*models.Order]("CustomerID", request.CustomerID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	log.Printf("len result 1: %+v\n", len(results))
	var idList []int
	for _, item := range results {
		idList = append(idList, item.OrderID)
	}
	fin, err := database.FindByEqualListUnionT[*models.OrderDetail]("OrderID", idList, database.PreloadProduct)

	log.Printf("len result 2: %+v\n", len(fin))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   fin,
	})
}
func GetAllorderbyToken(c *gin.Context) {
	customerid, err := services.GetIDbyToken(c)

	//查询第一次返回orderdetails_id[]
	results, err := database.FindByEqualUnion[*models.Order]("CustomerID", customerid, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	log.Printf("len result 1: %+v\n", len(results))
	var idList []int
	for _, item := range results {
		idList = append(idList, item.OrderID)
	}
	// fin, err := database.FindByEqualListUnionT[*models.OrderDetail]("OrderID", idList, database.PreloadOrder, database.PreloadProduct)
	// fin, err := database.FindByEqualListUnionT_in_order_join_product[models.OrderDetail]("OrderDetail.OrderID", idList, database.PreloadProduct)
	fin, err := database.GetAllModelWithPreload[models.OrderDetail](database.PreloadProduct)
	for _, detail := range fin {
		if detail.Product.ProductID == 0 {
			log.Println("No Product data found for OrderDetail ID:", detail.OrderDetailID)
		} else {
			log.Printf("OrderDetail ID: %d, Product Name: %s\n", detail.OrderDetailID, detail.Product.ProductName)
		}
	}

	// log.Printf("data =%v\n", fin[0])
	log.Printf("len result 2: %+v\n", len(fin))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   fin,
	})
}
