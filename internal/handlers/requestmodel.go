package handlers

// OrderRequest 结构体
type OrderRequest struct {
	CustomerID  uint    `json:"customer_id" validate:"required"`        // 客户ID
	OrderDate   string  `json:"order_date" validate:"required"`         // 订单日期
	TotalAmount float64 `json:"total_amount" validate:"required,min=0"` // 总金额
}

// CategoryRequest 结构体
type CategoryRequest struct {
	CategoryName        string `json:"category_name" validate:"required,max=100"`         // 分类名称
	CategoryDescription string `json:"category_description" validate:"omitempty,max=255"` // 分类描述
}

// SupplierRequest 结构体
type SupplierRequest struct {
	Name        string `json:"name" validate:"required,max=100"` // 供应商名称
	Address     string `json:"address" validate:"max=255"`       // 地址
	PhoneNumber string `json:"phone_number" validate:"numeric"`  // 电话
}

// ProductRequest 结构体
type ProductRequest struct {
	CategoryID    uint    `json:"category_id" validate:"required"`             // 分类ID
	SupplierID    uint    `json:"supplier_id" validate:"required"`             // 供应商ID
	ProductName   string  `json:"product_name" validate:"required,max=100"`    // 产品名称
	ProductDetail string  `json:"product_detail" validate:"omitempty,max=255"` // 产品详情
	UnitPrice     float64 `json:"unit_price" validate:"required,min=0"`        // 单价
	Quantity      int     `json:"quantity" validate:"required,min=0"`          // 数量
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
