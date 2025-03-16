-- 删除已存在的 mall_gin 数据库（如果存在）
DROP DATABASE IF EXISTS mall_gin;

-- 创建新的 mall_gin 数据库
CREATE DATABASE mall_gin;
-- 使用指定的数据库
USE mall_gin;

-- 创建Customers表
CREATE TABLE Customers (
    CustomerID INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    Password VARCHAR(100) NOT NULL,
    Address VARCHAR(255),
    PhoneNumber VARCHAR(20),
    Email VARCHAR(100) UNIQUE
);

-- 创建Orders表
CREATE TABLE Orders (
    OrderID INT AUTO_INCREMENT PRIMARY KEY,
    CustomerID INT,
    OrderDate DATE NOT NULL,
    TotalAmount DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (CustomerID) REFERENCES Customers(CustomerID)
);

-- 创建Categories表
CREATE TABLE Categories (
    CategoryID INT AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(100) NOT NULL,
    CategoryDescription TEXT
);

-- 创建Suppliers表
CREATE TABLE Suppliers (
    SupplierID INT AUTO_INCREMENT PRIMARY KEY,
    Name VARCHAR(100) NOT NULL,
    Address VARCHAR(255),
    PhoneNumber VARCHAR(20)
);

-- 创建Products表
CREATE TABLE Products (
    ProductID INT AUTO_INCREMENT PRIMARY KEY,
    CategoryID INT,
    SupplierID INT,
    ProductName VARCHAR(100) NOT NULL,
    ProductDetail TEXT,
    UnitPrice DECIMAL(10, 2) NOT NULL,
    Quantity INT NOT NULL,
    FOREIGN KEY (CategoryID) REFERENCES Categories(CategoryID),
    FOREIGN KEY (SupplierID) REFERENCES Suppliers(SupplierID)
);

-- 创建OrderDetail表
CREATE TABLE OrderDetail (
    OrderDetailID INT AUTO_INCREMENT PRIMARY KEY,
    OrderID INT,
    ProductID INT,
    UnitPrice DECIMAL(10, 2) NOT NULL,
    Amount INT NOT NULL,
    Status ENUM('Pending', 'Shipped', 'Cancelled', 'Returned') NOT NULL,
    FOREIGN KEY (OrderID) REFERENCES Orders(OrderID),
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID)
);

-- 创建 CartItems 表
CREATE TABLE CartItems (
    cart_item_id INT AUTO_INCREMENT PRIMARY KEY,
    CustomerID INT NOT NULL,
    ProductID INT NOT NULL,
    FOREIGN KEY (CustomerID) REFERENCES Customers(CustomerID),
    FOREIGN KEY (ProductID) REFERENCES Products(ProductID)
);

-- 禁用外键检查
SET FOREIGN_KEY_CHECKS = 0;

INSERT INTO Customers (Name, Address, PhoneNumber, Email, Password) VALUES
('张三', '北京市朝阳区某街道1号', '13800138000', 'zhangsan@example.com', 'password123'),
('李四', '上海市浦东新区某路2号', '13900139000', 'lisi@example.com', 'mypassword!'),
('王五', '广州市天河区某巷3号', '13600136000', 'wangwu@example.com', 'securePass456'),
('赵六', '深圳市南山区某弄4号', '13700137000', 'zhaoliu@example.com', 'strongP@ssw0rd'),
('孙七', '杭州市西湖区某街5号', '13500135000', 'sunqi@example.com', 'anotherPass123'),
('1', '杭州市西湖区某街5号', '1', 'test@example.com', '$2a$10$kXRUV87tDakKDk49IW.CQ.RjDhX6s6lbAgMCUU6REWuV0q6jEzYO6');

INSERT INTO Orders (CustomerID, OrderDate, TotalAmount) VALUES
(1, '2025-01-15', 299.99),
(2, '2025-02-10', 149.99),
(3, '2025-03-01', 499.99),
(4, '2025-02-25', 99.99),
(5, '2025-03-05', 199.99),
(6, '2025-03-15', 899.97), 
(6, '2025-03-20', 449.98);

INSERT INTO OrderDetail (OrderID, ProductID, UnitPrice, Amount, Status) VALUES
(1, 1, 2999.99, 1, 'Shipped'),
(2, 2, 199.99, 2, 'Pending'),
(3, 3, 99.99, 3, 'Cancelled'),
(4, 4, 49.99, 1, 'Returned'),
(5, 5, 29.99, 5, 'Shipped'),
(6, 1, 2999.99, 1, 'Shipped'), 
(6, 3, 99.99, 2, 'Pending'), 
(7, 2, 199.99, 2, 'Shipped'), 
(7, 4, 49.99, 1, 'Pending'); 

INSERT INTO Products (CategoryID, SupplierID, ProductName, ProductDetail, UnitPrice, Quantity) VALUES
(1, 1, '智能手机X1', '一款高性能智能手机', 2999.99, 100),
(2, 2, '时尚连衣裙', '适合春夏穿着的时尚连衣裙', 199.99, 50),
(3, 3, '不锈钢锅', '耐用的不锈钢炒锅', 99.99, 200),
(4, 4, '编程入门书', '学习编程的基础教程', 49.99, 300),
(5, 5, '进口巧克力', '来自比利时的高品质巧克力', 29.99, 150);

INSERT INTO Categories (CategoryName, CategoryDescription) VALUES
('电子产品', '包括手机、电脑等电子设备'),
('服装', '男装、女装、童装等各类服装'),
('家居用品', '家具、厨房用具等家居必需品'),
('图书', '各类文学、科技、教育书籍'),
('食品', '零食、饮料、保健品等');
INSERT INTO Suppliers (Name, Address, PhoneNumber) VALUES
('供应商A', '供应商地址A', '13800138001'),
('供应商B', '供应商地址B', '13800138002'),
('供应商C', '供应商地址C', '13800138003'),
('供应商D', '供应商地址D', '13800138004'),
('供应商E', '供应商地址E', '13800138005');

INSERT INTO CartItems (CustomerID, ProductID) VALUES 
(6, 1), -- 智能手机X1
(6, 2), -- 时尚连衣裙
(6, 5); -- 进口巧克力

-- 插入更多Products表的数据

-- 电子产品
INSERT INTO Products (CategoryID, SupplierID, ProductName, ProductDetail, UnitPrice, Quantity) VALUES
(1, 1, '笔记本电脑N2', '一款适合办公与游戏的笔记本电脑', 5999.99, 80),
(1, 2, '无线耳机', '高保真音质无线蓝牙耳机', 199.99, 120),
(1, 3, '智能手表W3', '健康追踪多功能智能手表', 499.99, 60);

-- 服装
INSERT INTO Products (CategoryID, SupplierID, ProductName, ProductDetail, UnitPrice, Quantity) VALUES
(2, 4, '男士皮鞋', '经典款男士商务皮鞋', 149.99, 70),
(2, 5, '运动裤', '舒适透气运动长裤', 89.99, 100),
(2, 1, '女士风衣', '时尚防风防水女士外套', 199.99, 45);

-- 家居用品
INSERT INTO Products (CategoryID, SupplierID, ProductName, ProductDetail, UnitPrice, Quantity) VALUES
(3, 2, '木质餐桌', '现代简约风格四人餐桌', 699.99, 30),
(3, 3, '电热水壶', '快速加热不锈钢电热水壶', 79.99, 200),
(3, 4, '棉被', '冬季保暖加厚棉被', 149.99, 150);

-- 图书
INSERT INTO Products (CategoryID, SupplierID, ProductName, ProductDetail, UnitPrice, Quantity) VALUES
(4, 5, '算法导论', '计算机科学经典教程', 89.99, 100),
(4, 1, '艺术的故事', '了解世界艺术史的最佳读物', 69.99, 50),
(4, 2, '摄影基础', '从零开始学习摄影技巧', 59.99, 75);

-- 食品
INSERT INTO Products (CategoryID, SupplierID, ProductName, ProductDetail, UnitPrice, Quantity) VALUES
(5, 3, '咖啡豆', '精选阿拉比卡咖啡豆', 79.99, 90),
(5, 4, '绿茶', '清香爽口中国绿茶', 39.99, 200),
(5, 5, '坚果混合包', '健康零食坚果组合', 59.99, 120);

-- 继续插入其他产品
INSERT INTO Products (CategoryID, SupplierID, ProductName, ProductDetail, UnitPrice, Quantity) VALUES
(1, 2, '平板电脑T1', '轻便易携的平板电脑', 2499.99, 50),
(2, 3, '儿童连衣裙', '可爱图案设计的儿童连衣裙', 99.99, 60),
(3, 1, '空气净化器', '高效过滤空气中的污染物', 899.99, 40),
(4, 3, '心理学入门', '探索人类心灵的基础书籍', 59.99, 120),
(5, 2, '意大利面', '传统手工制作意大利面', 29.99, 200),
(1, 4, '高清显示器', '色彩鲜艳的高清电脑显示器', 1499.99, 30),
(2, 5, '夏季短袖', '透气舒适的纯棉短袖上衣', 59.99, 150),
(3, 1, '按摩椅', '全身放松的家用按摩椅', 2999.99, 20),
(4, 2, '投资指南', '理财规划与投资策略', 79.99, 80),
(5, 3, '蜂蜜', '天然野生花蜜', 69.99, 100),
(1, 5, 'VR眼镜', '沉浸式虚拟现实体验设备', 1999.99, 25),
(2, 1, '羊毛围巾', '柔软保暖的手工编织围巾', 119.99, 70),
(3, 2, '吸尘器', '强力清洁各种地面污渍', 599.99, 45),
(4, 3, '自我提升', '个人成长与职业发展的指导书籍', 49.99, 90),
(5, 4, '红酒', '法国进口优质红酒', 199.99, 60),
(1, 4, '机械键盘', '专为游戏玩家设计的机械键盘', 799.99, 55),
(2, 5, '户外夹克', '防水防风高性能户外夹克', 249.99, 65),
(3, 1, '保温杯', '长效保温不锈钢保温杯', 89.99, 180),
(4, 2, '历史百科全书', '全面介绍世界历史知识', 149.99, 35),
(5, 3, '手工巧克力', '精致礼盒装手工巧克力', 49.99, 110);


-- 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;