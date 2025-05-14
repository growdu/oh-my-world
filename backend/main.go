package main

import (
	_ "api/docs" // 引入生成的swagger文档
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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
	Image       string `json:"image"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Visits      int    `json:"visits" gorm:"default:0"`
}

// Category 模型
type Category struct {
	Name  string `json:"name" gorm:"primarykey"`
	Count int    `json:"count" gorm:"default:0"`
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
		file, err := ioutil.ReadFile(configPath)
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

	// 自动迁移
	err = db.AutoMigrate(&Link{}, &Category{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	log.Println("Database migration completed")

	// 初始化默认分类
	if err := initializeCategories(); err != nil {
		log.Fatal("failed to initialize categories:", err)
	}
	log.Println("Categories initialized successfully")
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
	}

	for _, name := range categories {
		var category Category
		result := db.FirstOrCreate(&category, Category{Name: name})
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
	category := c.Query("category")

	query := db
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("visits desc").Find(&links).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch links"})
		return
	}
	c.JSON(http.StatusOK, links)
}

// @Summary 添加新链接
// @Description 添加新的链接信息
// @Tags Links
// @Accept json
// @Produce json
// @Param link body Link true "链接信息"
// @Success 200 {object} gin.H{"message": "Link added successfully", "link": Link}
// @Failure 400 {object} gin.H{"error": "Invalid input"}
// @Failure 500 {object} gin.H{"error": "Unable to add link"}
// @Router /api/v1/links [post]
func addLink(c *gin.Context) {
	var newLink Link

	if err := c.ShouldBindJSON(&newLink); err != nil {
		log.Println("JSON解析失败:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 验证分类是否存在
	var category Category
	if err := db.First(&category, "name = ?", newLink.Category).Error; err != nil {
		log.Printf("分类验证失败: %v, 分类名: %s", err, newLink.Category)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category"})
		return
	}

	// 存储链接信息到数据库
	if err := db.Create(&newLink).Error; err != nil {
		log.Printf("创建链接失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add link"})
		return
	}

	// 更新分类计数
	if err := db.Model(&category).Update("count", gorm.Expr("count + 1")).Error; err != nil {
		log.Printf("更新分类计数失败: %v", err)
		// 不返回错误，继续执行
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link added successfully", "link": newLink})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch categories"})
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
// @Success 200 {object} gin.H{"message": "Visit count updated"}
// @Failure 404 {object} gin.H{"error": "Link not found"}
// @Failure 500 {object} gin.H{"error": "Update failed"}
// @Router /api/v1/links/{id}/visit [post]
func incrementVisits(c *gin.Context) {
	id := c.Param("id")

	result := db.Model(&Link{}).Where("id = ?", id).Update("visits", gorm.Expr("visits + 1"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visit count updated"})
}

func main() {
	initDB()

	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("/api/v1")
	{
		api.GET("/links", getLinks)
		api.POST("/links", addLink)
		api.GET("/categories", getCategories)
		api.POST("/links/:id/visit", incrementVisits)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := config.Server.Port
	if port == "" {
		port = "8092"
	}

	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
