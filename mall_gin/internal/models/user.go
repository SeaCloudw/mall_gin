// models/user.go
package models

import "time"

//一共7个表

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`            // 主键
	Username  string    `gorm:"unique;not null" json:"username"` // 用户名
	Password  string    `gorm:"not null" json:"-"`               // 密码（不返回）
	Email     string    `gorm:"unique;not null" json:"email"`    // 邮箱
	CreatedAt time.Time `json:"created_at"`                      // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                      // 更新时间
}

// 在select与选择表时，gorm处理表名与属性名都是蛇形+下划线+小写来处理的！
// 像这种CustomerID  会并从select customer_id,所以只能强制指定它column:数据库的属性的名字
// MYsql表名会区分大小写，但是属性名不会区分大小写，所以想Name这种就可以正常找到
type Customer struct {
	CustomerID  uint   `gorm:"column:CustomerID;primaryKey" json:"id"` // 主键
	Name        string `gorm:"unique;not null" json:"name"`            // 用户名
	Password    string `gorm:"not null" json:"-"`                      // 密码（不返回）
	Email       string `gorm:"unique;not null" json:"email"`           // 邮箱
	Address     string `json:"address"`                                // 地址
	PhoneNumber string `gorm:"column:PhoneNumber" json:"phonenumber"`  // 电话

	// CreatedAt   time.Time `json:"created_at"`                                            // 创建时间
	// UpdatedAt   time.Time `json:"updated_at"`                                            // 更新时间
}

// 显式指定表名为 Customers
func (Customer) TableName() string {
	return "Customers"
}

// Product 模型对应 Products 表
type Product struct {
	ProductID     uint    `gorm:"column:ProductID;primaryKey;autoIncrement" json:"product_id"`
	CategoryID    uint    `gorm:"column:CategoryID" json:"category_id"`
	SupplierID    uint    `gorm:"column:SupplierID" json:"supplier_id"`
	ProductName   string  `gorm:"column:ProductName;not null" json:"product_name"`
	ProductDetail string  `gorm:"column:ProductDetail" json:"product_detail"`
	UnitPrice     float64 `gorm:"column:UnitPrice;not null" json:"unit_price"` // DECIMAL(10, 2) 对应 Go 中的 float64
	Quantity      int     `gorm:"column:Quantity;not null" json:"quantity"`

	// 关联到 Categories 表
	Category Category `gorm:"foreignKey:CategoryID"`

	// 关联到 Suppliers 表
	Supplier Supplier `gorm:"foreignKey:SupplierID"`
}

// 显式指定表名为 Customers
func (Product) TableName() string {
	return "Products"
}

// CartItem 模型对应 CartItems 表
type CartItem struct {
	CartItemID int `gorm:"column:cart_item_id;primaryKey;autoIncrement"`
	CustomerID int `gorm:"column:CustomerID;not null"`
	ProductID  int `gorm:"column:ProductID;not null"`

	// 关联到 Customers 表
	Customer Customer `gorm:"foreignKey:CustomerID"`
	// 关联到 Products 表
	Product Product `gorm:"foreignKey:ProductID"`
}

// TableName 设置模型对应的表名
func (CartItem) TableName() string {
	return "CartItems"
}

// Order 模型对应 Orders 表
type Order struct {
	OrderID     int     `gorm:"column:OrderID;primaryKey;autoIncrement"`
	CustomerID  int     `gorm:"column:CustomerID"`
	OrderDate   string  `gorm:"column:OrderDate;not null"` // 注意：GORM 默认不支持 time.Time 的直接映射，可以使用 string 或自定义类型
	TotalAmount float64 `gorm:"column:TotalAmount;not null"`

	// 关联到 Customers 表
	Customer Customer `gorm:"foreignKey:CustomerID"`
}

// TableName 设置模型对应的表名
func (Order) TableName() string {
	return "Orders"
}

// Category 模型对应 Categories 表
type Category struct {
	CategoryID          int    `gorm:"column:CategoryID;primaryKey;autoIncrement"`
	CategoryName        string `gorm:"column:CategoryName;not null"`
	CategoryDescription string `gorm:"column:CategoryDescription"`
}

// TableName 设置模型对应的表名
func (Category) TableName() string {
	return "Categories"
}

// Supplier 模型对应 Suppliers 表
type Supplier struct {
	SupplierID  int    `gorm:"column:SupplierID;primaryKey;autoIncrement"`
	Name        string `gorm:"column:Name;not null"`
	Address     string `gorm:"column:Address"`
	PhoneNumber string `gorm:"column:PhoneNumber"`
}

// TableName 设置模型对应的表名
func (Supplier) TableName() string {
	return "Suppliers"
}

// OrderDetail 模型对应 OrderDetail 表
type OrderDetail struct {
	OrderDetailID uint    `gorm:"column:OrderDetailID;primaryKey;autoIncrement" json:"order_detail_id"`
	OrderID       uint    `gorm:"column:OrderID" json:"order_id"`
	ProductID     uint    `gorm:"column:ProductID" json:"product_id"`
	UnitPrice     float64 `gorm:"column:UnitPrice;not null" json:"unit_price"` // DECIMAL(10, 2) 对应 Go 中的 float64
	Amount        int     `gorm:"column:Amount;not null" json:"amount"`
	Status        string  `gorm:"column:Status;type:enum('Pending','Shipped','Cancelled','Returned');not null" json:"status"`

	// 关联到 Orders 表
	Order Order `gorm:"foreignKey:OrderID"`

	// 关联到 Products 表
	Product Product `gorm:"foreignKey:ProductID"`
}

// TableName 设置模型对应的表名
func (OrderDetail) TableName() string {
	return "OrderDetail"
}

// 统一实现接口
// 实现 Model 接口
func (c *Customer) GetID() uint {
	return c.CustomerID
}
func (o *Order) GetID() uint {
	return uint(o.OrderID)
}
func (c *Category) GetID() uint {
	return uint(c.CategoryID)
}
func (s *Supplier) GetID() uint {
	return uint(s.SupplierID)
}
func (p *Product) GetID() uint {
	return p.ProductID
}
func (od *OrderDetail) GetID() uint {
	return od.OrderDetailID
}
func (od *CartItem) GetID() uint {
	return uint(od.CartItemID)
}
