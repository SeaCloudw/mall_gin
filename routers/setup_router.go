package routers

// routers负责添加路由,匹配路径与handler
import (
	"mall_gin/internal/handlers"

	"mall_gin/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {

	// r.GET("/test", handlers.GetTestPage)
	//用于加载导航栏
	r.GET("/bar", handlers.GetBarPage)

	// 使用中间件
	// r.Use(middleware.AuthMiddleware())

	//登录界面，获取token,这里前端需要保存，同时后面所有请求需要加上-H "Authorization: Bearer token
	r.GET("/login", handlers.GetLoginPage)
	r.GET("/register", handlers.GetRegisterPage)

	r.GET("/index", handlers.GetIndexPage) //加载html,服务端渲染（SSR）简单的方式
	//查看个人界面
	r.GET("/profile", handlers.GetProfileList)
	r.GET("/cart", handlers.GetCartPage)
	r.GET("/orders", handlers.GetOrdersPage)

	r.GET("/products/:id", handlers.GetProductsPage) //加载html,服务端渲染（SSR）简单的方式

	//专门用于验证token，成功返回一个string,失败直接跳转到login
	// validate_get('/protect');
	// 查看商品详情
	//查看购物车、订单、个人中心、
	// 加入购物车、下单
	va := r.Group("/protect", services.AuthMiddleware())
	{
		// va.GET("/index", handlers.GetIndexPage)
		va.GET("", handlers.GetValidate)
	}

	v := r.Group("/api")
	{

		//启动JWT验证,需要token的id的操作
		authorized := v.Group("/token", services.AuthMiddleware())
		{
			// 需要拦截的内容，访问一下内容，前端需要
			/* 		    function fetchWithAuth(url, options = {}) {
			    const token = localStorage.getItem('authToken');
			    options.headers = {
			        ...options.headers,
			        'Authorization': token ? `Bearer ${token}` : ''
			    };
			    return fetch(url, options);
			}
			*/

			//返回个人主页
			// authorized.GET("/profile", handlers.GetProfileByToken)
			//使用redis的测试
			authorized.GET("/profile", handlers.GetProfileByToken_Redis)
			//更新更主页
			authorized.PUT("/profile", handlers.UPgrateProfileByToken)
			//返回个人的购物车
			authorized.GET("/cart", handlers.GetAllCartbyToken)
			//将商品添加到购物车，  注：如果已经有了就不能再添加
			authorized.POST("/cart/add", handlers.CreateNewCart)
			//获取个人的订单
			authorized.GET("/orders", handlers.GetAllorderbyToken)

		}

		// 客户登录
		v.POST("/login", handlers.LoginHandler)
		// 返回某个用户的信息

		//GET返回所有顾客信息
		v.GET("/customers", handlers.GetAllCustomersHandler)
		//创建
		v.POST("/customers", handlers.PostNewCustomerHandler)

		v.DELETE("/customers/:id", handlers.DeleteCustomerHandler)

		v.PUT("/customers/:id", handlers.UpgrateCustomerHandler)

		//获取所有商品信息
		v.GET("/products", handlers.GetAllProductsHandler)
		// 根据商品名 模糊查找by name
		v.POST("/products/select", handlers.FindbynameHandler)
		//根据id返回对应属性
		v.GET("/products/:id", handlers.FindbyidHandler)

		//购物车

		// 查看个人购物车界面，这里用了post传递customer_id,返回的还是货物
		v.POST("/cart", handlers.GetAllCartbyId)
		// 、、删除购物车,需要传customer_id,product_id，注：
		v.POST("/cart/delete", handlers.DeleteCart)

		//订单order  添加 注：需要前端发送order_date 具体值
		v.POST("/orders/add", handlers.CreateOrder)
		//获取所有订单
		v.POST("/orders", handlers.GetAllorderbyid)
		//更新订单状态，status
		v.PUT("/orders/updatestatus", handlers.ChangeOrder)

	}
	// r.GET("/index", handlers.GetIndexPage) //加载html,服务端渲染（SSR）简单的方式
	//加载商品详情

	// r.GET("/home", handlers.GetHome)       //加载html,服务端渲染（SSR）简单的方式
	// r.GET("/admin", handlers.GetAdminPage) //加载html,服务端渲染（SSR）简单的方式
	// 定义路由
	/* 	r.GET("/hello", handlers.GetHello)
	   	r.GET("/home", handlers.GetHome)           //加载html,服务端渲染（SSR）简单的方式
	   	r.GET("/users", handlers.GetUsers)         //加载html,服务端渲染（SSR）简单的方式
	   	r.GET("/list", handlers.GetUsersList)      //加载html,服务端渲染（SSR）简单的方式
	   	r.GET("/usersdata", handlers.GetUsersData) // 返回数据
	   	r.GET("/base", handlers.GetHomepage) // 返回数据
	   	// r.POST("/users", handlers.CreateUser) */

	// 用户注册
	// r.GET("/register", handlers.GetRegisterPage) // 返回数据

	// v1 := r.Group("/api")
	// { //提交注册，将数据存到数据库
	// 	v1.POST("/register", handlers.Register)
	// 	//从数据库  获取数据
	// 	v1.GET("/users", handlers.EasyGetUsers)
	// 	//删除数据
	// v1.DELETE("/delete/users/:id", handlers.DeleteUser)
	// 	//更新数据
	// 	v1.PUT("/update/users/:id", handlers.UpdateUser)
	// }

	return r
}
