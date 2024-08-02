package main

import (
	"fmt"

	"github.com/ezep02/microservicios/internal/auth/handler"
	"github.com/ezep02/microservicios/internal/auth/repository"
	"github.com/ezep02/microservicios/internal/auth/service"
	"github.com/ezep02/microservicios/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

	// Iniciar configuraciones de viper
	initConfig()

	// Inicializar la base de datos
	connection := initializeDatabase()

	// Configurar el enrutador
	router := setupRouter(connection)

	// Iniciar el servidor
	router.Run(viper.GetString("PORT"))

}

// Inicializar base de datos
func initializeDatabase() *db.Database {
	connection, err := db.NewDatabase()
	if err != nil {
		panic("[ERROR] Algo salió mal conectando a la base de datos")
	}
	return connection
}

var minSecretKey = 32

// Configurar el enrutador
func setupRouter(connection *db.Database) *gin.Engine {
	r := gin.Default()

	tokenKey := viper.GetString("TOKEN_KEY")

	if len(tokenKey) < minSecretKey {
		fmt.Printf("secret_key must be at least %d characters", minSecretKey)
	}

	// autenticación
	authRepo := repository.NewAuthRepository(connection.GetDB())
	authServices := service.NewAuthService(authRepo)
	authController := handler.NewAuthHandler(authServices, tokenKey)
	handler.AuthRouter(r, authController)

	return r
}

func initConfig() {
	viper.SetConfigName(".env")   // nombre del archivo de configuración (sin extensión)
	viper.SetConfigType("env")    // tipo de archivo de configuración
	viper.AddConfigPath("../../") // ruta donde buscar el archivo de configuración

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("[ERROR] leyendo el archivo de configuración: %v", err)
	}
}
