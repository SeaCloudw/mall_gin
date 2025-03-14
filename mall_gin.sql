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


-- 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;