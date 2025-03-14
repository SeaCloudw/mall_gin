以下是基于您提供的后端代码和路由信息生成的接口文档，采用README格式。这份文档将帮助前端开发者理解如何与这些API进行交互。

# 商品管理API

## 概述
本API提供了商品的查询功能，包括获取所有商品信息以及根据商品名称模糊查找商品。

## API列表


### 1. 顾客
#### 概述
提供了顾客信息的增删改查功能，包括获取所有顾客信息、创建新顾客、删除顾客以及更新顾客信息。

#### 1.1获取所有顾客信息

- **URL**: `/api/customers`
- **方法**: `GET`
- **描述**: 此接口用于获取所有顾客的信息。

##### 响应数据
```json
{
  "status": "success",
  "data": [
    {
      "ID": 1,
      "Name": "顾客姓名",
      "Email": "邮箱地址",
      "Address": "顾客地址",
      "PhoneNumber": "联系电话"
      // 其他可能存在的字段...
    },
    // 更多顾客信息...
  ]
}
```

#### 1.2创建新顾客

- **URL**: `/api/customers`
- **方法**: `POST`
- **描述**: 此接口允许创建一个新的顾客记录。

##### 请求体 (JSON)
```json
{
  "name": "顾客姓名",         // 必须
  "password": "密码",        // 必须
  "email": "邮箱地址",       // 必须
  "address": "顾客地址",     // 必须，最大长度255字符
  "phonenumber": "联系电话"   // 必须，数字格式
}
```

##### 响应数据
成功响应:
```json
{
  "status": "success",
  "message": "客户创建成功",
  "data": {
    "name": "顾客姓名",
    "password": "密码",
    "email": "邮箱地址",
    "address": "顾客地址",
    "phonenumber": "联系电话"
  }
}
```

错误响应示例:
```json
{
  "status": "error",
  "message": "具体错误消息"
}
```

#### 1.3删除顾客

- **URL**: `/api/customers/:id`
- **方法**: `DELETE`
- **描述**: 此接口允许根据顾客ID删除一个顾客记录。

##### 路径参数
- `id`: 顾客的唯一标识符（整数）

##### 响应数据
成功响应:
```json
{
  "status": "success",
  "message": "删除成功"
}
```

错误响应示例:
```json
{
  "status": "error",
  "message": "具体错误消息"
}
```

#### 1.4更新顾客信息

- **URL**: `/api/customers/:id`
- **方法**: `PUT`
- **描述**: 此接口允许根据顾客ID更新顾客的信息。

##### 请求体 (JSON)
```json
{
  "name": "新的顾客姓名",         // 可选
  "password": "新的密码",        // 可选
  "email": "新的邮箱地址",       // 可选
  "address": "新的顾客地址",     // 可选，最大长度255字符
  "phonenumber": "新的联系电话"   // 可选，数字格式
}
```

##### 路径参数
- `id`: 顾客的唯一标识符（整数）

##### 响应数据
成功响应:
```json
{
  "status": "success",
  "message": "客户更新成功",
  "data": {
    "ID": 1,
    "Name": "新的顾客姓名",
    "Email": "新的邮箱地址",
    "Address": "新的顾客地址",
    "PhoneNumber": "新的联系电话"
  }
}
```

错误响应示例:
```json
{
  "status": "error",
  "message": "具体错误消息"
}
```


### 2. 商品
#### 2.1 获取所有商品信息

- **URL**: `/api/products`
- **方法**: `GET`
- **描述**: 此接口用于获取所有商品的信息。

#### 请求参数
此接口无需请求参数。

#### 响应数据
```json
{
  "status": "success",
  "data": [
    {
      "ID": 1,
      "ProductName": "商品名称",
      "ProductCategory": "类别",
      "ProductSupplier": "供应商"
      // 其他可能存在的字段...
    },
    // 更多商品信息...
  ]
}
```

#### 2.2 根据商品名模糊查找

- **URL**: `/api/products/select`
- **方法**: `POST`
- **描述**: 此接口允许通过商品名称的一部分来查找匹配的商品。

#### 请求体 (JSON)
```json
{
  "SearchText": "关键词",  // 必须
  "Page": 1,               // 可选，默认为1
  "PageSize": 10,          // 可选，默认为10
  "SortField": "字段名",   // 可选
  "SortOrder": "asc/desc", // 可选，升序或降序
  "Fields": ["需要返回的字段名"], // 可选
  "IncludeAssociations": ["需要预加载的关联数据"] // 可选
}
```

#### 响应数据
成功响应:
```json
[
  {
    "ID": 1,
    "ProductName": "匹配到的商品名称",
    "ProductCategory": "类别",
    "ProductSupplier": "供应商"
    // 其他可能存在的字段...
  },
  // 更多匹配到的商品信息...
]
```

错误响应示例:
```json
{
  "status": "error",
  "message": "具体错误消息"
}
```


### 3. 购物车
本API提供了购物车的基本操作功能，包括将商品添加到购物车、查看个人购物车、以及从购物车中删除商品。

#### 3.1 将商品添加到购物车

- **URL**: `/api/cart/add`
- **方法**: `POST`
- **描述**: 此接口用于将指定商品添加到顾客的购物车中。如果该商品已经在购物车中，则不能再次添加。

##### 请求体 (JSON)
```json
{
  "customer_id": 1,   // 必须，顾客ID
  "product_id": 101   // 必须，商品ID
}
```

##### 响应数据
成功响应:
```json
{
  "status": "success",
  "message": "添加购物车成功",
  "data": {
    "CustomerID": 1,
    "ProductID": 101
  }
}
```

错误响应示例:
```json
{
  "status": "error",
  "message": "具体错误消息"
}
```

#### 3.2 查看个人购物车界面

- **URL**: `/api/cart`
- **方法**: `POST`
- **描述**: 此接口允许根据顾客ID查询其购物车中的所有商品。

##### 请求体 (JSON)
```json
{
  "customer_id": 1   // 必须，顾客ID
}
```

##### 响应数据
成功响应:
```json
[
  {
    "ID": 101,
    "ProductName": "商品名称",
    "ProductCategory": "类别",
    "Price": 99.99,
    "Description": "商品描述"
    // 其他可能存在的字段...
  },
  // 更多购物车商品信息...
]
```

错误响应示例:
```json
{
  "status": "error",
  "message": "具体错误消息"
}
```

#### 3.3 删除购物车中的商品

- **URL**: `/api/cart/delete`
- **方法**: `POST`
- **描述**: 此接口允许根据顾客ID和商品ID从购物车中删除指定的商品。

##### 请求体 (JSON)
```json
{
  "customer_id": 1,   // 必须，顾客ID
  "product_id": 101   // 必须，商品ID
}
```

##### 响应数据
成功响应:
```json
{
  "status": "success",
  "message": "删除成功"
}
```

错误响应示例:
```json
{
  "status": "error",
  "message": "没有找到符合要求的结果"
}
```

### 4. 订单
本API提供了订单的基本操作功能，包括创建新订单、查看个人订单以及更新订单状态。

#### 4.1 创建新订单

- **URL**: `/api/orders/add`
- **方法**: `POST`
- **描述**: 此接口用于创建一个新的订单。前端需要提供顾客ID、订单日期、产品ID列表、单价列表、数量列表和初始状态。

##### 请求体 (JSON)
```json
{
  "customer_id": 1,                           // 必须，顾客ID
  "order_date": "2025-03-11T10:26:00Z",      // 必须，订单日期
  "product_id": [101, 102],                   // 必须，产品ID列表
  "unit_price": [99.99, 49.99],               // 必须，单价列表
  "amount": [2, 1],                           // 必须，数量列表
  "status": "Pending"                         // 必须，初始状态（Pending, Shipped, Cancelled, Returned）
}
```

##### 响应数据
成功响应:
```json
{
  "status": "success",
  "message": "添加订单与订单详情成功",
  "data": [
    {
      "ProductID": 101,
      "UnitPrice": 99.99,
      "Amount": 2,
      "Status": "Pending"
    },
    {
      "ProductID": 102,
      "UnitPrice": 49.99,
      "Amount": 1,
      "Status": "Pending"
    }
  ]
}
```

错误响应示例:
```json
{
  "status": "error",
  "message": "具体错误消息"
}
```

#### 4.2 查看个人订单

- **URL**: `/api/orders`
- **方法**: `POST`
- **描述**: 此接口允许根据顾客ID查询其所有订单及其详细信息。

##### 请求体 (JSON)
```json
{
  "customer_id": 1   // 必须，顾客ID
}
```

##### 响应数据
成功响应:
```json
{
  "status": "success",
  "data": [
    {
      "OrderID": 1,
      "CustomerID": 1,
      "OrderDate": "2025-03-11T10:26:00Z",
      "TotalAmount": 249.97,
      "Details": [
        {
          "ProductID": 101,
          "UnitPrice": 99.99,
          "Amount": 2,
          "Status": "Pending"
        },
        {
          "ProductID": 102,
          "UnitPrice": 49.99,
          "Amount": 1,
          "Status": "Pending"
        }
      ]
    },
    // 更多订单信息...
  ]
}
```

错误响应示例:
```json
{
  "status": "error",
  "message": "具体错误消息"
}
```

#### 4.3 更新订单状态

- **URL**: `/api/orders/updatestatus`
- **方法**: `PUT`
- **描述**: 此接口允许根据订单详情ID更新订单的状态。

##### 请求体 (JSON)
```json
{
  "order_detail_id": 1,                       // 必须，订单详情ID
  "status": "Shipped"                         // 必须，新的状态（Pending, Shipped, Cancelled, Returned）
}
```

##### 响应数据
成功响应:
```json
{
  "status": "success",
  "message": "订单状态更新成功"
}
```

错误响应示例:
```json
{
  "status": "error",
  "message": "没有找到符合要求的结果"
}
```

## 注意事项
- 在使用`POST`和`PUT`接口时，请确保发送的是有效的JSON格式的请求体。
- 如果在调用接口时遇到任何问题，请检查输入参数是否符合要求，并查看响应中的错误消息以获取更多信息。
- 对于`/api/orders/add`接口，确保提供的商品ID、单价和数量列表长度一致，并且每个商品的数量至少为1。
- 对于`/api/orders/updatestatus`接口，确保传递正确的`order_detail_id`和有效的新状态值。

- 在使用`POST`接口时，请确保发送的是有效的JSON格式的请求体。
- 如果在调用接口时遇到任何问题，请检查输入参数是否符合要求，并查看响应中的错误消息以获取更多信息。
- 对于`/api/cart/add`接口，如果商品已经存在于购物车中，则不会重复添加，前端需要处理这种情况。
- 对于`/api/cart/delete`接口，确保传递正确的`customer_id`和`product_id`组合，以避免误删其他用户的购物车内容。


- 在使用`/products/select`接口时，请确保发送的是一个有效的JSON格式的请求体。
- 如果在调用接口时遇到任何问题，请检查输入参数是否符合要求，并查看响应中的错误消息以获取更多信息。

希望这份文档能帮助你更好地理解和使用商品管理API。如果有任何疑问或需要进一步的帮助，请随时联系我们。