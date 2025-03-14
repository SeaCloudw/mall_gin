package main

import (
	"mall_gin/internal/database"
	"mall_gin/internal/redis"
	"mall_gin/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	// simpledemo1()
	getstart()
}
func getstart() {
	// 初始化数据库
	database.InitDB()

	// 初始化Redis客户端
	redis.InitRedis()

	// 获取Redis客户端并进行操作
	// redisClient := GetClient()

	/* 	// 获取数据库实例
	   	db := database.GetDB()
	   	// 创建用户
	   	user := models.User{
	   		Username: "testuser",
	   		Password: "123456",
	   		Email:    "test@example.com",
	   	}
	   	result := db.Create(&user)
	   	if result.Error != nil {
	   		log.Fatalf("创建用户失败: %v", result.Error)
	   	}
	   	log.Println("用户创建成功") */

	r := gin.Default()
	// 加载模板文件
	r.LoadHTMLGlob("templates/*")

	// 设置静态资源路径
	r.Static("/assets", "./assets")
	r = routers.SetupRouter(r)

	r.Run(":8080") // 指定监听  端口
}
func simpledemo1() {
	r := gin.Default() // 创建一个默认的路由引擎,*gin.Engine,负责管理所有的 HTTP 路由。
	// 默认中间件包括 Logger（记录日志）和 Recovery（恢复从 panic 中恢复）。

	// r.GET(路由路径，处理器handler)//注册一个处理GET请求的路由
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "Hello World!")
		// 向客户端返回状态码为 200 的响应，并将响应体设置为 "Hello World!"。

		c.JSON(200, gin.H{"message": "Get all users"})
	})
	// gin.Context 封装了 HTTP 请求和响应的所有信息
	/* 	// gin.Context的方法列举
		method := c.Request.Method         // 获取请求方法
		userAgent := c.Request.UserAgent() // 获取 User-Agent
		c.Writer.Write([]byte("Hello"))    // 写入响应内容
		c.JSON()：返回 JSON 格式的响应。
	c.String()：返回纯文本响应。
	c.HTML()：返回 HTML 格式的响应*/

	// r.Run() // 默认监听地址为 0.0.0.0:8080
	r.Run(":8080") // 指定监听 9090 端口
}

// 自定义，返回gin.HandlerFunc类型
//
//	注册一个拦截器，添加到路由中，可以全局作用或者指定组
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization") // 获取请求头中的 "Authorization" 字段
		if token == "" {                      // 如果没有提供 token
			c.JSON(401, gin.H{"error": "Unauthorized"}) // 返回 401 错误
			c.Abort()                                   // 终止后续的路由处理器和中间件执行
			return
		}
		c.Next() // 如果有 token，则继续处理下一个中间件或路由
	}
}
func demo() {
	r := gin.Default()
	r.Use(AuthMiddleware()) // 全局注册中间件

	r.GET("/secure", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This is a secure route"})
	})
}

/* func routergroupdemo() {
	r := gin.Default()
	// 两个路由组，都可以访问，大括号是为了保证规范
	v1 := r.Group("/v1")
	{
		// 通过 localhost:8080/v1/hello访问，以此类推
		v1.GET("/hello", sayHello)
		v1.GET("/world", sayWorld)
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello", sayHello)
		v2.GET("/world", sayWorld)
	}
	r.Run(":8080")

} */
