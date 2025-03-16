package database

import (
	"errors"
	"fmt"
	"log"
	"mall_gin/internal/models" // 导入 models 包

	"gorm.io/gorm"
)

// GetCustomerByID 根据ID获取客户信息
func GetCustomerByID(id uint) (*models.Customer, error) {
	db := GetDB()
	var customer models.Customer
	result := db.First(&customer, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &customer, nil
}
func GetAllUCustomers() ([]models.Customer, error) {
	db := GetDB()

	var customers []models.Customer

	if err := db.Select("CustomerID,name,email,address,PhoneNumber").Find(&customers).Error; err != nil {
		return nil, errors.New("查询用户数据失败")
	}

	return customers, nil
}
func CreateANewCustomer(customer *models.Customer) error {
	db := GetDB()

	// 在这里可以添加密码加密逻辑
	// customer.Password = hashPassword(customer.Password)

	if err := db.Create(customer).Error; err != nil {
		return errors.New("创建客户失败: " + err.Error())
	}
	return nil
}

// 删除一个用户
func DeleteCustomer(id uint) error {
	db := GetDB()
	result := db.Delete(&models.Customer{}, id)
	if result.Error != nil {

		return result.Error
	}
	return nil
}
func UpdateCustomer(customer *models.Customer) error {
	db := GetDB()
	result := db.Save(customer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// 通用的方法
// Deleter 是一个通用接口，表示可以被删除的对象
// type Deleter interface {
// 	GetID() uint
// }

// DeleteByID 根据ID删除指定类型的记录
// func DeleteByID(model Deleter, id uint) error {
// 	db := GetDB()

// 	// 将Deleter转换为interface{}以便传递给Delete方法
// 	result := db.Delete(model, id)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

// Model 是一个通用接口，表示可以进行CRUD操作的对象
type Model interface {
	GetID() uint
}

// Create 创建一个新的记录
func CreateModel[T Model](model T) error {
	db := GetDB()
	result := db.Create(&model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Create 创建一个新的记录
func CreateModelAndReturn[T Model](model T) (T, error) {
	db := GetDB()
	result := db.Create(&model)
	if result.Error != nil {
		return model, result.Error
	}
	return model, nil
}

// DeleteByID 根据ID删除指定类型的记录
func DeleteByIDModel[T Model](id uint) error {
	db := GetDB()
	var model T
	result := db.Delete(&model, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Update 更新一个记录
func UpdateModel[T Model](model T) error {
	db := GetDB()
	result := db.Save(&model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FindByID 根据ID查找指定类型的记录
func FindByIDModel[T Model](id uint) (T, error) {
	var model T
	db := GetDB()
	result := db.First(&model, id)
	if result.Error != nil {
		return model, result.Error
	}
	return model, nil
}

// FindByAttributeLike 函数，支持按特定的属性名进行模糊搜索
// 例如：products, err := FindByAttributeLike[Product]("ProductName", "apple")
func FindByAttributeLike[T Model](attributeName, keyword string) ([]T, error) {
	db := GetDB()
	var results []T

	// 构建查询条件
	query := db.Where(fmt.Sprintf("%s LIKE ?", attributeName), fmt.Sprintf("%%%s%%", keyword))

	// 执行查询
	if result := query.Find(&results); result.Error != nil {
		return nil, errors.New("查询数据失败")
	}

	return results, nil
}

// 联合多个表进行模糊搜索
func FindByAttributeLikeUnion[T Model](attributeName, keyword string, preloadFuncs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	db := GetDB()
	var results []T

	query := db.Model(new(T))
	// 应用所有预加载函数
	for _, preloadFunc := range preloadFuncs {
		if preloadFunc != nil {
			query = preloadFunc(query)
		}
	}
	// 构建查询条件
	query = query.Where(fmt.Sprintf("%s LIKE ?", attributeName), fmt.Sprintf("%%%s%%", keyword))

	// 执行查询
	if result := query.Find(&results); result.Error != nil {
		return nil, errors.New("查询数据失败")
	}

	return results, nil
}

// 等值查询，只查询一次
func FindByEqualUnion[T Model](attributeName string, keyword int, preloadFuncs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	db := GetDB()
	var results []T

	query := db.Model(new(T))
	// 应用所有预加载函数
	for _, preloadFunc := range preloadFuncs {
		if preloadFunc != nil {
			query = preloadFunc(query)
		}
	}
	// 构建查询条件
	query = query.Where(fmt.Sprintf("%s = ?", attributeName), keyword)
	// query = query.Where(fmt.Sprintf("%s = ?", attributeName), fmt.Sprintf("%d", keyword))

	// 执行查询
	if result := query.Find(&results); result.Error != nil {
		return nil, errors.New("查询数据失败")
	}

	return results, nil
}

// FindByEqualUnion 根据属性名和关键字进行查询，并支持预加载
func FindByEqualUnionString[T Model](attributeName string, keyword string, preloadFuncs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	db := GetDB()
	var results []T

	query := db.Model(new(T))
	// 应用所有预加载函数
	for _, preloadFunc := range preloadFuncs {
		if preloadFunc != nil {
			query = preloadFunc(query)
		}
	}
	// 构建查询条件
	query = query.Where(fmt.Sprintf("%s = ?", attributeName), keyword)

	// 执行查询
	result := query.Find(&results)
	if result.Error != nil {
		return nil, errors.New("查询数据失败: " + result.Error.Error())
	}

	return results, nil
}

// 等值查询,相同条件，一次性查询多个【实际上是查询了n次】
func FindByEqualListUnion[T Model](attributeName string, keyword []int, preloadFuncs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	db := GetDB()
	var results []T

	query := db.Model(new(T))
	// 应用所有预加载函数
	for _, preloadFunc := range preloadFuncs {
		if preloadFunc != nil {
			query = preloadFunc(query)
		}
	}
	for i := 0; i < len(keyword); i++ {
		// 构建查询条件
		query = query.Where(fmt.Sprintf("%s = ?", attributeName), keyword)
		// query = query.Where(fmt.Sprintf("%s = ?", attributeName), fmt.Sprintf("%d", keyword))

		// 执行查询
		tmp := query.Find(&results)
		if tmp.Error != nil {

			return nil, errors.New("查询数据失败")

		}
	}

	return results, nil
}

// FindByEqualListUnion 根据attributeName和一系列关键字查询记录
func FindByEqualListUnionT[T Model](attributeName string, keywords []int, preloadFuncs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	db := GetDB()
	// db = db.Debug() //启动输出

	var records []T

	query := db.Model(new(T))

	// 应用所有预加载函数
	log.Println("len pro=", len(preloadFuncs))
	for _, preloadFunc := range preloadFuncs {
		if preloadFunc != nil {
			query = preloadFunc(query)
		}
	}

	// 如果没有关键词，返回空结果集
	if len(keywords) == 0 {
		return records, nil
	}

	// 使用IN语句一次性查询所有符合条件的结果
	result := query.Where(fmt.Sprintf("%s IN (?)", attributeName), keywords).Find(&records)
	if result.Error != nil {
		return nil, errors.New("查询数据失败")
	}
	// log.Printf("Query Results: %+v\n", records)
	return records, nil
}

// 等值查询,多个条件，返回切片
func FindByMultipleConditions[T Model](conditions map[string][]int, preloadFuncs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	db := GetDB()
	var results []T

	query := db.Model(new(T))

	// 应用所有预加载函数
	for _, preloadFunc := range preloadFuncs {
		if preloadFunc != nil {
			query = preloadFunc(query)
		}
	}

	// 构建查询条件
	for attributeName, keywords := range conditions {
		if len(keywords) > 0 {
			query = query.Where(fmt.Sprintf("%s IN (?)", attributeName), keywords)
		} else {
			// 如果关键词为空，则忽略该条件
			continue
		}
	}

	// 执行查询
	result := query.Find(&results)
	if result.Error != nil {
		return nil, errors.New("查询数据失败: " + result.Error.Error())
	}

	return results, nil
}

// GetAll 获取所有记录的通用函数,可以指定属性名
// 例如 customers, err := service.GetAllWithDefaultColumns[models.Customer](db)
// "CustomerID", "name", "email", "address", "PhoneNumber")
func GetAllModel[T any](columns ...string) ([]T, error) {
	db := GetDB()
	var records []T

	query := db.Model(new(T))

	if len(columns) > 0 {
		query = query.Select(columns)
	}

	if result := query.Find(&records); result.Error != nil {
		return nil, errors.New("查询数据失败")
	}

	return records, nil
}

// GetAllWithDefaultColumns 获取所有记录，默认获取所有列
func GetAllWithDefaultColumns[T any]() ([]T, error) {
	return GetAllModel[T]()
}

func FindByEqualListUnionT_in_order_join_product[T any](attributeName string, keywords []int, preloadFuncs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	// db := GetDB()

	// // 启用 SQL 日志输出
	// db = db.Debug()

	// var records []T

	// query := db.Model(new(T))

	// // 应用所有预加载函数
	// log.Println("len pro=", len(preloadFuncs))
	// for _, preloadFunc := range preloadFuncs {
	// 	if preloadFunc != nil {
	// 		query = preloadFunc(query)
	// 	}
	// }

	// // 如果没有关键词，返回空结果集
	// if len(keywords) == 0 {
	// 	return records, nil
	// }
	// result := query.Where(fmt.Sprintf("%s IN (?)", attributeName), keywords).Find(&records)
	// if result.Error != nil {
	// 	return nil, errors.New("查询数据失败")
	// }
	// // 使用IN语句一次性查询所有符合条件的结果，并且使用 JOIN 预加载 Product
	// // result := query.
	// // 	Model(&models.OrderDetail{}).
	// // 	Select("OrderDetail.*, Products.*"). // 选择所有需要的字段
	// // 	Joins("JOIN Products ON Products.ProductID = OrderDetail.ProductID").
	// // 	Where(fmt.Sprintf("%s IN (?)", attributeName), keywords).
	// // 	Find(&records)

	// log.Printf("Query Results: %+v\n", records)
	// return records, nil

	// 应用所有预加载函数
	// for _, preloadFunc := range preloadFuncs {
	// 	if preloadFunc != nil {
	// 		query = preloadFunc(query)
	// 	}
	// }

	// if result := query.Find(&records); result.Error != nil {
	// 	return nil, errors.New("查询数据失败")
	// }

	// result := query.Where(fmt.Sprintf("%s IN (?)", attributeName), keywords).Find(&records)
	// result := query.Find(&records)
	db := GetDB()
	var records []T
	db = db.Debug()
	// 使用 Preload 预加载 Product 关联数据
	result := db.
		Preload("Product"). // 预加载 Product 关联数据
		Where(fmt.Sprintf("%s IN (?)", attributeName), keywords).
		Find(&records)

	if result.Error != nil {
		return nil, errors.New("查询数据失败")
	}

	log.Printf("Query Results: %+v\n", records)
	return records, nil

	// query := db.Model(new(T))

	// // 应用所有预加载函数
	// for _, preloadFunc := range preloadFuncs {
	// 	if preloadFunc != nil {
	// 		query = preloadFunc(query)
	// 	}
	// }

	// if result := query.Find(&records, 1); result.Error != nil {
	// 	return nil, errors.New("查询数据失败")
	// }

	// return records, nil

}

// 在多个表中交叉查询，需要预加载
func GetAllModelWithPreload[T any](preloadFuncs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	db := GetDB()
	var records []T
	db = db.Debug()
	query := db.Model(new(T))

	// 应用所有预加载函数,相当于圈定了范围
	for _, preloadFunc := range preloadFuncs {
		if preloadFunc != nil {
			query = preloadFunc(query)
		}
	}

	if result := query.Find(&records); result.Error != nil {
		return nil, errors.New("查询数据失败")
	}

	return records, nil
}
func GetPartModelWithPreload[T any](columns []string, preloadFuncs ...func(*gorm.DB) *gorm.DB) ([]T, error) {
	db := GetDB()
	var records []T

	query := db.Model(new(T))

	// 如果提供了列名，则选择这些列
	if len(columns) > 0 {
		query = query.Select(columns)
	}

	// 应用所有预加载函数
	for _, preloadFunc := range preloadFuncs {
		if preloadFunc != nil {
			query = preloadFunc(query)
		}
	}

	if result := query.Find(&records); result.Error != nil {
		return nil, errors.New("查询数据失败")
	}

	return records, nil
}

// 预加载函数
func PreloadProductCategory(db *gorm.DB) *gorm.DB {
	return db.Preload("Category")
}

func PreloadProductSupplier(db *gorm.DB) *gorm.DB {
	return db.Preload("Supplier")
}

func PreloadOrderCustomer(db *gorm.DB) *gorm.DB {
	return db.Preload("Customer")
}

func PreloadOrderDetails(db *gorm.DB) *gorm.DB {
	return db.Preload("OrderDetails")
}
func PreloadProduct(db *gorm.DB) *gorm.DB {
	return db.Preload("Product")
}
func PreloadOrder(db *gorm.DB) *gorm.DB {
	return db.Preload("Order")
}
