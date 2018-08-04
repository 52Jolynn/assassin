package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/recover"
	"github.com/kataras/iris/middleware/logger"
	"time"
	"os"
	"github.com/kataras/golog"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"flag"
	"github.com/go-sql-driver/mysql"
)

var RootLogger *golog.Logger
var Database *sql.DB

func main() {
	username := flag.String("u", "root", "-u root")
	passwd := flag.String("pwd", "123456", "-pwd 123456")
	host := flag.String("h", "localhost", "-h localhost")
	port := flag.String("p", "3306", "-p 3306")
	dbname := flag.String("db", "assassin", "-db assassin")
	flag.Parse()

	// Creates an application without any middleware by default.
	app := iris.New()

	// Recover middleware recovers from any panics and writes a 500 if there was one.
	app.Use(recover.New())

	requestLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		Query: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		MessageHeaderKeys: []string{"User-Agent"},
	})
	app.Use(requestLogger)

	//添加日志文件
	logFile := newLogFile()
	defer logFile.Close()
	RootLogger = app.Logger()
	RootLogger.AddOutput(logFile)

	//连接数据库
	dbConfig := mysql.Config{User: *username, Passwd: *passwd, Net: "tcp", Addr: fmt.Sprintf("%s:%s", *host, *port), DBName: *dbname}
	RootLogger.Info(fmt.Sprintf("try to connect database: %s", dbConfig.FormatDSN()))

	var err error
	Database, err = sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		panic(err.Error())
	}
	Database.SetMaxOpenConns(100)
	Database.SetMaxIdleConns(10)
	Database.SetConnMaxLifetime(time.Minute * 25)
	defer Database.Close()

	if err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed)); err != nil {
		RootLogger.Warn("Shutdown with error: " + err.Error())
	}
}

func newLogFile() *os.File {
	filename := time.Now().Format("2006-01-02") + ".log"
	// Open the file, this will append to the today's file if server restarted.
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
}
