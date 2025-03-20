package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"

	//"gorm.io/datatypes"

	_ "api/docs" // 引入生成的swagger文档
)

// @title Link Management API
// @version 1.0
// @description 使用 Gin + GORM + PostgreSQL 的链接管理 API。
// @host localhost:8092
// @BasePath /api

// 数据库连接
var db *gorm.DB

// Link 模型
type Link struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	//Image       datatypes.JSON `json:"image"`
	Description string         `json:"description"`
	URL         string         `json:"url"`
}

// 初始化数据库
func initDB() {
	dsn := "host=localhost user=postgres password=double+2=4 dbname=world port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}
	db.AutoMigrate(&Link{}) // 自动迁移
}

// @Summary 获取所有链接
// @Description 获取数据库中存储的所有链接信息
// @Tags Links
// @Accept  json
// @Produce  json
// @Success 200 {array} Link
// @Failure 500 {object} map[string]string
// @Router /links [get]
func getLinks(c *gin.Context) {
	var links []Link
	if err := db.Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch links"})
		return
	}
	c.JSON(http.StatusOK, links)
}

// @Summary 添加新链接
// @Description 向数据库添加一个新的链接
// @Tags Links
// @Accept  json
// @Produce  json
// @Param link body Link true "链接信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /links [post]
func addLink(c *gin.Context) {
	var newLink Link
	if err := c.ShouldBindJSON(&newLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := db.Create(&newLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Link added successfully", "link": newLink})
}

func main() {
	initDB()

	r := gin.Default()

	// 注册 API 路由
	api := r.Group("/api")
	{
		api.GET("/links", getLinks)
		api.POST("/links", addLink)
	}

	// 集成 Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	r.Run(":8092")
}
