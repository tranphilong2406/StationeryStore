package cmd

import (
	"StoreServer/api"
	auth "StoreServer/api/auth"
	"StoreServer/config"
	"StoreServer/job"
	"StoreServer/middleware"
	"StoreServer/models"
	"StoreServer/utils/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupServer() {
	con := config.GetConfig()
	errs := make(chan error)

	//connecting to db
	fmt.Println("Connecting to database...")
	job.DBConnect()
	defer job.Disconnect()
	// init collection
	println("Init collection...")
	models.InitExampleDB()
	models.InitProductDB()
	models.InitCategoryDB()
	models.InitOrderDB()
	models.InitAuthDB()
	models.InitReceivedOrderDB()
	println("Init collection done!")

	s := SetHandler()

	go func() {
		server := &http.Server{
			Addr:              ":" + con.ServerPort,
			Handler:           s,
			WriteTimeout:      time.Second * 30,
			IdleTimeout:       time.Second * 30,
			ReadHeaderTimeout: time.Minute,
		}
		logger.GetLogger().Info("Server running on port: " + server.Addr)

		errs <- server.ListenAndServe()
	}()

	err := <-errs
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}
}

func SetHandler() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPut,
			http.MethodPatch,
			http.MethodOptions,
		},
		AllowHeaders:           []string{"Origin", "Authorization", "Content-Type", "token"},
		AllowBrowserExtensions: true,
		AllowWebSockets:        true,
		AllowFiles:             true,
	}))
	r.Use(gin.Recovery())
	// Example routes
	r.POST("/api/example", middleware.CheckLogin(), api.CreateExample)
	r.GET("/api/example", middleware.CheckLogin(), api.GetExample)
	r.POST("/api/examples", middleware.CheckLogin(), api.CreateListExample)
	r.PUT("/api/example/", middleware.CheckLogin(), api.UpdateExample)
	r.DELETE("/api/example/:id", middleware.CheckLogin(), api.DeleteExample)
	// Product routes
	r.GET("/api/product/", middleware.CheckLogin(), api.GetProduct)
	r.POST("/api/product", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.CreateProduct)
	r.PUT("/api/product/", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.UpdateProduct)
	r.DELETE("/api/product/:id", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.DeleteProduct)
	// Category routes
	r.GET("/api/category/", middleware.CheckLogin(), api.GetCategory)
	r.POST("/api/category", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.CreateCategory)
	r.PUT("/api/category/", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.UpdateCategory)
	r.DELETE("/api/category/:id", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.DeleteCategory)
	//Order routes
	r.POST("/api/order/", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.CreateOrder)
	r.GET("/api/order/", middleware.CheckLogin(), api.GetOrder)
	r.GET("/api/order/:id", middleware.CheckLogin(), api.GetOrderByID)
	r.DELETE("/api/order/:id", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.DeleteOrder)
	// User routes
	r.POST("/api/user", middleware.CheckLogin(), middleware.CheckRole("admin"), api.CreateUser)
	r.PUT("/api/user/", middleware.CheckLogin(), middleware.CheckRole("admin"), api.UpdateUser)
	// Auth routes
	r.POST("/api/auth/login", auth.Login)
	// Received Order
	r.POST("/api/rec_order/", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.CreateReceivedOrder)
	r.GET("/api/rec_order/", middleware.CheckLogin(), api.GetReceivedOrder)
	r.GET("/api/rec_order/:id", middleware.CheckLogin(), api.GetReceivedOrderByID)
	r.DELETE("/api/rec_order/:id", middleware.CheckLogin(), middleware.CheckRole("admin", "manager"), api.DeleteReceivedOrder)
	return r
}
