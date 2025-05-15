package main

import (
	_ "api/docs" // 引入生成的swagger文档
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Link Management API
// @version 1.0
// @description 使用 Gin + GORM + PostgreSQL 的链接管理 API。
// @host localhost:8092
// @BasePath /api/v1

// Config 配置结构体
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		SSLMode  string `json:"sslmode"`
		TimeZone string `json:"timezone"`
	} `json:"database"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

// 全局变量
var (
	db     *gorm.DB
	config Config
)

// Link 模型
type Link struct {
	ID          uint   `json:"id" gorm:"primarykey"`
	Name        string `json:"name" gorm:"not null"`
	Image       string `json:"image" gorm:"not null"`
	URL         string `json:"url" gorm:"not null"`
	Description string `json:"description" gorm:"type:text"`
	Category    string `json:"category" gorm:"not null"`
	CategoryID  uint   `json:"category_id" gorm:"not null"`
	VisitCount  int    `json:"visits" gorm:"default:0"`
}

// Category 模型
type Category struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Name string `json:"name" gorm:"not null;unique"`
}

// 加载配置
func loadConfig() Config {
	// 首先尝试从环境变量读取
	dbConfig := Config{}
	dbConfig.Database.Host = os.Getenv("DB_HOST")
	dbConfig.Database.Port = os.Getenv("DB_PORT")
	dbConfig.Database.User = os.Getenv("DB_USER")
	dbConfig.Database.Password = os.Getenv("DB_PASSWORD")
	dbConfig.Database.DBName = os.Getenv("DB_NAME")
	dbConfig.Server.Port = os.Getenv("PORT")

	// 如果环境变量不完整，尝试从配置文件读取
	if dbConfig.Database.Host == "" || dbConfig.Database.Port == "" ||
		dbConfig.Database.User == "" || dbConfig.Database.Password == "" ||
		dbConfig.Database.DBName == "" {

		// 获取可执行文件所在目录
		ex, err := os.Executable()
		if err != nil {
			log.Fatal("无法获取可执行文件路径:", err)
		}
		configPath := filepath.Join(filepath.Dir(ex), "config.json")

		// 如果在可执行文件目录找不到配置文件，尝试在当前目录查找
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = "config.json"
		}

		// 读取配置文件
		file, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatal("无法读取配置文件:", err)
		}

		// 解析配置文件
		var fileConfig Config
		if err := json.Unmarshal(file, &fileConfig); err != nil {
			log.Fatal("配置文件格式错误:", err)
		}

		// 使用配置文件中的值作为默认值
		if dbConfig.Database.Host == "" {
			dbConfig.Database.Host = fileConfig.Database.Host
		}
		if dbConfig.Database.Port == "" {
			dbConfig.Database.Port = fileConfig.Database.Port
		}
		if dbConfig.Database.User == "" {
			dbConfig.Database.User = fileConfig.Database.User
		}
		if dbConfig.Database.Password == "" {
			dbConfig.Database.Password = fileConfig.Database.Password
		}
		if dbConfig.Database.DBName == "" {
			dbConfig.Database.DBName = fileConfig.Database.DBName
		}
		if dbConfig.Server.Port == "" {
			dbConfig.Server.Port = fileConfig.Server.Port
		}

		// 设置默认值
		if dbConfig.Database.SSLMode == "" {
			dbConfig.Database.SSLMode = fileConfig.Database.SSLMode
		}
		if dbConfig.Database.TimeZone == "" {
			dbConfig.Database.TimeZone = fileConfig.Database.TimeZone
		}
	}

	// 验证必要的配置是否存在
	if dbConfig.Database.Host == "" || dbConfig.Database.Port == "" ||
		dbConfig.Database.User == "" || dbConfig.Database.Password == "" ||
		dbConfig.Database.DBName == "" {
		log.Fatal("数据库配置不完整")
	}

	return dbConfig
}

// 初始化数据库
func initDB() {
	config = loadConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.Port,
		config.Database.SSLMode,
		config.Database.TimeZone,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	log.Println("Connected to database successfully")

	// 执行数据库迁移
	err = db.AutoMigrate(&Category{})
	if err != nil {
		log.Fatal("failed to migrate categories:", err)
	}

	// 初始化默认分类
	if err := initializeCategories(); err != nil {
		log.Fatal("failed to initialize categories:", err)
	}
	log.Println("Categories initialized successfully")

	// 迁移 Link 表，但不包括 NOT NULL 约束
	if err := db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			image VARCHAR(255),
			url VARCHAR(255),
			description TEXT,
			category VARCHAR(255),
			visit_count INTEGER DEFAULT 0
		)
	`).Error; err != nil {
		log.Fatal("failed to create links table:", err)
	}

	// 更新所有空分类为"其他"
	if err := db.Exec(`
		UPDATE links 
		SET category = '其他' 
		WHERE category IS NULL OR category = ''
	`).Error; err != nil {
		log.Fatal("failed to update null categories:", err)
	}

	// 添加 NOT NULL 约束
	if err := db.Exec(`
		ALTER TABLE links 
		ALTER COLUMN category SET NOT NULL,
		ALTER COLUMN name SET NOT NULL,
		ALTER COLUMN image SET NOT NULL,
		ALTER COLUMN url SET NOT NULL
	`).Error; err != nil {
		log.Fatal("failed to set not null constraints:", err)
	}

	log.Println("Database migration completed")
}

// 初始化默认分类
func initializeCategories() error {
	categories := []string{
		"编程语言",
		"开发工具",
		"框架",
		"数据库",
		"运维",
		"前端",
		"后端",
		"人工智能",
		"其他",
	}

	for _, name := range categories {
		var category Category
		result := db.Where(Category{Name: name}).FirstOrCreate(&category, Category{Name: name})
		if result.Error != nil {
			return fmt.Errorf("failed to create category %s: %v", name, result.Error)
		}
	}
	return nil
}

// @Summary 获取所有链接
// @Description 获取数据库中存储的所有链接信息
// @Tags Links
// @Accept json
// @Produce json
// @Param category query string false "按分类筛选"
// @Success 200 {array} Link
// @Failure 500 {object} map[string]string
// @Router /api/v1/links [get]
func getLinks(c *gin.Context) {
	var links []Link
	categoryID := c.Query("categoryId")

	query := db
	if categoryID != "" {
		if id, err := strconv.ParseUint(categoryID, 10, 32); err == nil {
			query = query.Where("category_id = ?", uint(id))
		}
	}

	if err := query.Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, links)
}

// @Summary 获取所有分类
// @Description 获取所有可用的分类及其链接数量
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {array} Category
// @Failure 500 {object} map[string]string
// @Router /api/v1/categories [get]
func getCategories(c *gin.Context) {
	var categories []Category
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

// @Summary 更新链接访问次数
// @Description 增加链接的访问计数
// @Tags Links
// @Accept json
// @Produce json
// @Param id path int true "链接ID"
// @Success 200 {object} gin.H{"success": bool}
// @Failure 404 {object} gin.H{"error": "Link not found"}
// @Failure 500 {object} gin.H{"error": "Update failed"}
// @Router /api/v1/links/{id}/visit [post]
func incrementVisits(c *gin.Context) {
	id := c.Param("id")
	var link Link
	if err := db.First(&link, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	if err := db.Model(&link).Update("visit_count", link.VisitCount+1).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// @Summary 创建新链接
// @Description 创建一个新的链接
// @Tags Links
// @Accept json
// @Produce json
// @Param link body Link true "链接信息"
// @Success 200 {object} Link
// @Failure 400 {object} map[string]string
// @Router /api/v1/links [post]
func createLink(c *gin.Context) {
	var link Link
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if link.Name == "" || link.Image == "" || link.URL == "" || link.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, image, url and category are required"})
		return
	}

	var category Category
	err := db.Where("name = ?", link.Category).First(&category).Error
	if err != nil {
		// 如果找不到，可能需要创建，或者返回错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	link.CategoryID = category.ID

	if err := db.Create(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, link)
}

// @Summary 删除链接
// @Description 删除指定的链接
// @Tags Links
// @Accept json
// @Produce json
// @Param id path int true "链接ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/links/{id} [delete]
func deleteLink(c *gin.Context) {
	id := c.Param("id")
	var link Link
	if err := db.First(&link, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	if err := db.Delete(&link).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link deleted successfully"})
}

func main() {
	initDB()

	r := gin.Default()

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type"},
	}))

	// API 路由
	v1 := r.Group("/api/v1")
	{
		v1.GET("/links", getLinks)
		v1.POST("/links", createLink)
		v1.DELETE("/links/:id", deleteLink)
		v1.POST("/links/:id/visit", incrementVisits)
		v1.GET("/categories", getCategories)
	}

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := config.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
