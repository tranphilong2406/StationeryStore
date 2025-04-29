package cmd

import (
	"StoreServer/api"
	"StoreServer/config"
	"StoreServer/job"
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
	models.InitExampleDB()
	models.InitProductDB()
	models.InitCategoryDB()

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
	r.POST("/api/example", api.CreateExample)
	r.GET("/api/example", api.GetExample)
	r.POST("/api/examples", api.CreateListExample)
	r.PUT("/api/example/", api.UpdateExample)
	r.DELETE("/api/example/:id", api.DeleteExample)
	// Product routes
	r.GET("/api/product/", api.GetProduct)
	r.POST("/api/product", api.CreateProduct)
	r.PUT("/api/product/", api.UpdateProduct)
	// Category routes
	r.GET("/api/category/", api.GetCategory)
	r.POST("/api/category", api.CreateCategory)
	r.PUT("/api/category/", api.UpdateCategory)

	return r
}
