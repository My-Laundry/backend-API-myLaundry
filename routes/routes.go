// package routes

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/raihansyahrin/backend_laundry_app.git/controllers"
// 	"github.com/raihansyahrin/backend_laundry_app.git/middlewares"
// )

// func SetupRoutes(router *gin.Engine) {
// 	authRoutes := router.Group("api/auth")
// 	{
// 		authRoutes.POST("/register", controllers.Register)
// 		authRoutes.POST("/login", controllers.Login)
// 	}

// 	userRoutes := router.Group("api/users")
// 	{
// 		userRoutes.GET("/", middlewares.AuthMiddleware(), controllers.GetUsers)
// 		userRoutes.GET("/:id", middlewares.AuthMiddleware(), controllers.GetUser)
// 		userRoutes.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)
// 		userRoutes.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)
// 	}

// 	orderRoutes := router.Group("api/orders")
// 	{
// 		orderRoutes.POST("/", middlewares.AuthMiddleware(), controllers.CreateOrder)
// 		orderRoutes.GET("/", middlewares.AuthMiddleware(), controllers.GetOrders)
// 		orderRoutes.PUT("/status", middlewares.AuthMiddleware(), controllers.UpdateOrderStatus)
// 		orderRoutes.DELETE("/", middlewares.AuthMiddleware(), controllers.DeleteOrder)
// 		orderRoutes.POST("/accept/:id", middlewares.AuthMiddleware(), controllers.AcceptOrder) // New route for accepting orders

// 	}

// 	serviceRoutes := router.Group("api/services")
// 	{
// 		serviceController := &controllers.ServiceController{}
// 		serviceRoutes.GET("/", serviceController.GetServices)
// 		serviceRoutes.GET("/:id", serviceController.GetServiceByID)
// 		serviceRoutes.POST("/", serviceController.CreateService)
// 		serviceRoutes.PUT("/:id", serviceController.UpdateService)
// 		serviceRoutes.DELETE("/:id", serviceController.DeleteService)
// 		serviceRoutes.GET("/category/:category", serviceController.GetServiceByCategory)
// 	}

//		addressRoutes := router.Group("api/addresses")
//		{
//			addressController := &controllers.AddressController{}
//			addressRoutes.POST("/", middlewares.AuthMiddleware(), addressController.CreateAddress)
//			addressRoutes.GET("/user/:user_id", middlewares.AuthMiddleware(), addressController.GetAddressesByUserID)
//			addressRoutes.PUT("/:id", middlewares.AuthMiddleware(), addressController.UpdateAddress)
//			addressRoutes.DELETE("/:id", middlewares.AuthMiddleware(), addressController.DeleteAddress)
//		}
//	}
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/raihansyahrin/backend_laundry_app.git/controllers"
	"github.com/raihansyahrin/backend_laundry_app.git/middlewares"
)

func SetupRoutes(router *gin.Engine) {
	authRoutes := router.Group("api/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
	}

	userRoutes := router.Group("api/users")
	{
		userRoutes.GET("/", middlewares.AuthMiddleware(), controllers.GetUsers)
		userRoutes.GET("/:id", middlewares.AuthMiddleware(), controllers.GetUser)
		userRoutes.PUT("/:id", middlewares.AuthMiddleware(), controllers.UpdateUser)
		userRoutes.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteUser)
	}

	orderRoutes := router.Group("api/orders")
	{
		orderRoutes.POST("/", middlewares.AuthMiddleware(), controllers.CreateOrder)
		orderRoutes.GET("/", middlewares.AuthMiddleware(), controllers.GetOrders)
		orderRoutes.GET("/:customer_id", middlewares.AuthMiddleware(), controllers.GetOrderDetailForCustomer) // Updated route for getting order details for customer

		orderRoutes.PUT("/status", middlewares.AuthMiddleware(), controllers.UpdateOrderStatus)
		orderRoutes.DELETE("/", middlewares.AuthMiddleware(), controllers.DeleteOrder)
		// orderRoutes.POST("/accept/:id", middlewares.AuthMiddleware(), controllers.AcceptOrder)         // Route for accepting orders
		orderRoutes.POST("/courier-arrived", middlewares.AuthMiddleware(), controllers.CourierArrived) // Route for courier arrived and updating weight
		orderRoutes.POST("/accept/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), controllers.AcceptOrder)
		orderRoutes.POST("/payment", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("customer", "courier"), controllers.ProcessPayment)
	}

	serviceRoutes := router.Group("api/services")
	{
		serviceController := &controllers.ServiceController{}
		serviceRoutes.GET("/", serviceController.GetServices)
		serviceRoutes.GET("/:id", serviceController.GetServiceByID)
		serviceRoutes.POST("/", serviceController.CreateService)
		serviceRoutes.PUT("/:id", serviceController.UpdateService)
		serviceRoutes.DELETE("/:id", serviceController.DeleteService)
		serviceRoutes.GET("/category/:category", serviceController.GetServiceByCategory)
	}

	addressRoutes := router.Group("api/addresses")
	{
		addressController := &controllers.AddressController{}
		addressRoutes.POST("/", middlewares.AuthMiddleware(), addressController.CreateAddress)
		addressRoutes.GET("/user/:user_id", middlewares.AuthMiddleware(), addressController.GetAddressesByUserID)
		addressRoutes.PUT("/:id", middlewares.AuthMiddleware(), addressController.UpdateAddress)
		addressRoutes.DELETE("/:id", middlewares.AuthMiddleware(), addressController.DeleteAddress)
	}
}
