package main

import (
	"github.com/gin-gonic/gin"
    "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
	"log"
	"net/http"
	"github.com/gin-contrib/cors"
    _ "api/docs" // 引入生成的swagger文档
	"fmt"
	"os"
)

// @title Link Management API
// @version 1.0
// @description 使用 Gin + GORM + PostgreSQL 的链接管理 API。
// @host localhost:8092
// @BasePath /api/v1

// 数据库连接
var db *gorm.DB

// Link 模型
type Link struct {
    ID          uint   `json:"id"`
    Image       string `json:"image"`  // 使用 BYTEA 类型来存储图片数据
    Description string `json:"description"`
    URL         string `json:"url"`
}


// 初始化数据库
func initDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbname, port,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	db.AutoMigrate(&Link{})
}

// @Summary 获取所有链接
// @Description 获取数据库中存储的所有链接信息
// @Tags Links
// @Accept  json
// @Produce  json
// @Success 200 {array} Link
// @Failure 500 {object} map[string]string
// @Router /api/v1/links [get]
func getLinks(c *gin.Context) {
	var links []Link
	if err := db.Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch links"})
		return
	}
	c.JSON(http.StatusOK, links)
}

// @Summary Add a new link
// @Description Upload a new link with an image
// @Tags links
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Image File"
// @Param description formData string true "Link Description"
// @Param url formData string true "Link URL"
// @Success 200 {object} gin.H{"message": "Link added successfully", "link": Link}
// @Failure 400 {object} gin.H{"error": "Invalid input"}
// @Failure 500 {object} gin.H{"error": "Unable to add link"}
// @Router /api/v1/links [post]
func addLink(c *gin.Context) {
	var newLink Link

	// 解析 JSON
	if err := c.ShouldBindJSON(&newLink); err != nil {
		log.Println("JSON解析失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
    // 存储链接信息到数据库
	if err := db.Create(&newLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Link added successfully", "link": newLink})
}


func main() {
	initDB()

	port := os.Getenv("PORT")

	r := gin.Default()

    // CORS 配置
	r.Use(cors.Default())

	// 注册 API 路由
	api := r.Group("/api/v1")
	{
		api.GET("/links", getLinks)
		api.POST("/links", addLink)
	}

	// 集成 Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if port == "" {
		port = "8092"
	}
	// 启动服务器
	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
