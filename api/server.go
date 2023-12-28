package api

import (
	"database/sql"
	"fmt"
	db "github/tdadadavid/fingreat/db/sqlc"
	"github/tdadadavid/fingreat/utils"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/go-playground/validator/v10"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)


type Server struct {
	queries *db.Store
	router *gin.Engine
	config *utils.Config
}

var tokenController *utils.JWTToken


func NewServer(envPath string) (*Server) {
  config, err := utils.LoadConfig(envPath)

	if err != nil {
		panic("Could not load config")
	}

	conn, err := sql.Open(config.DbDriver, config.DBSource + config.DBName + "?sslmode=disable");
	if err != nil {
		panic(fmt.Sprintf("Could not open database connection Error {%v}", err));
	}

	tokenController = utils.NewJWTToken(config)
	q := db.NewStore(conn)

	g := gin.Default()
	
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization", "Access-Control-Request-Headers") // Include the Authorization header
	g.Use(cors.New(corsConfig))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", currencyValidator)
	}

	return &Server{
		queries: q,
		router: g,
		config: config,
	}
}

func (s *Server) Start(port int)  {
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"Message": "Welcome to fingreat"})
	})

	User{}.router(s)
	Auth{}.router(s)
	Account{}.router(s)
	Transfer{}.router(s);

	s.router.Run(fmt.Sprintf(":%v", port))
}