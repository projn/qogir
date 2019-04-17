package console

import (
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"github.com/golang/protobuf/protoc-gen-go/grpc"
	"github.com/phachon/go-logger"
	consulapi "github.com/hashicorp/consul/api"
)

func main() {

	consulapi.AgentServiceRegistration

	logger := go_logger.NewLogger()

	fileConfig := &go_logger.FileConfig{
		Filename : "./test.log",
		LevelFileName : map[int]string{
			logger.LoggerLevel("error"): "./error.log",
			logger.LoggerLevel("info"): "./info.log",
			logger.LoggerLevel("debug"): "./debug.log",
		},
		MaxSize : 1024 * 1024,
		MaxLine : 10000,
		DateSlice : "d",
		JsonFormat: false,
		Format: "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}
	logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
	logger.SetAsync()


	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome</h1>")
	})

	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
		ctx.Params().Get()
	})

	// Method:   GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))


	db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	cl := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})


	cl := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: ":6379",
	})

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	sarama.NewSyncProducer
}