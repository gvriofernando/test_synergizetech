package apiService

import (
	"github.com/go-redis/redis/v8"
	"github.com/gvriofernando/test_synergizetech/app/core"
	"github.com/gvriofernando/test_synergizetech/app/handler"
	"github.com/gvriofernando/test_synergizetech/middleware"
	gorm "gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Db *gorm.DB
	Rd *redis.Client
}

// * define gin http router here
func GinHttpRouter(cfg Config, httpS *gin.Engine) {
	pCore := core.NewCore(core.Config{
		Rd: cfg.Rd,
		Db: cfg.Db,
	})

	// Handlers
	playerHandler := handler.Handlers{
		PCore: pCore,
	}

	httpS.Use(middleware.CORS())

	// route
	httpS.GET("/", home)
	httpS.POST("/register", playerHandler.Register())
	httpS.POST("/login", playerHandler.Login())
	httpS.POST("/logout", playerHandler.Logout())
	httpS.POST("/add-bank-acount", playerHandler.AddBankAccount())
	httpS.POST("/topup-wallet", playerHandler.TopupWallet())
	httpS.GET("/get-all-player", playerHandler.GetAllPlayers())
	httpS.GET("/get-detail-player/:playerId", playerHandler.GetDetailPlayer())
}

// Returning a health check response
func home(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "healthy",
	})
}
