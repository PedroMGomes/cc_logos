package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.search.crypto/controllers"
	"go.search.crypto/db"
	"go.search.crypto/models"
)

// Authenticate request middleware.
func Authenticate(config models.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("key") // map[string][]string
		if key != config.Access {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized Access"})
			return
		}
		// Pass to the next chained middleware.
		c.Next()
	}
}

// LoadEnvConfig loads env variables.
// .env files are not optimal if deploying to Heroku. Instead use the environment editor from Heroku dashboard.
// Viper automatically overrides values that have been read from a config file with values from the host environment with AutomaticEnv, only if they match.
func LoadEnvConfig(path string) (config models.AppConfig, err error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	err = viper.ReadInConfig()
	if err != nil {
		print("Failed to read config.yaml: " + err.Error())
	}
	// os.Getenv()
	viper.AutomaticEnv()
	err = viper.Unmarshal(&config)
	if err != nil {
		panic("Failed to read environment: " + err.Error())
	}
	return
}

func main() {
	config, err := LoadEnvConfig(".")
	if err != nil {
		print("Cannot read config.yaml: " + err.Error())
	}
	if len(config.DatabaseURL) == 0 {
		panic("No valid [DATABASE_URL] value.")
	}
	db.ConnectDatabase(config.DatabaseURL)
	// Assign config.yaml to local context.
	// Config = config
	// Gin default Router.
	gin.SetMode(config.GinMode)
	r := gin.Default()
	// Inject middleware.
	r.Use(Authenticate(config))
	// Register endpoints
	r.GET("/currency", controllers.Get)
	r.POST("/currency", controllers.Post)
	// r.GET("/currency/all", controllers.GetAll)
	// Listen and serves on 0.0.0.0:8080 (for windows localhost:8080), if not config is provided.
	builder := strings.Builder{}
	if len(config.Addr) != 0 {
		builder.WriteString(config.Addr)
	}
	if len(config.Port) != 0 {
		builder.WriteString(":")
		builder.WriteString(config.Port)
	}
	// addr := config.Addr + ":" + config.Port
	if len(builder.String()) > 0 {
		err = r.Run(builder.String())
		if err != nil {
			panic("Failed to start the server.")
		}
	} else {
		// invalid empty string as param.
		err = r.Run()
		if err != nil {
			panic("Failed to start the server.")
		}
	}
}
