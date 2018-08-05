package main

import (
	"time"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"flag"
	"github.com/go-sql-driver/mysql"
	"52jolynn.com/router"
	"52jolynn.com/core"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

func main() {
	username := flag.String("u", "root", "-u root")
	passwd := flag.String("pwd", "123456", "-pwd 123456")
	host := flag.String("h", "localhost", "-h localhost")
	port := flag.String("p", "3306", "-p 3306")
	dbname := flag.String("db", "assassin", "-db assassin")
	flag.Parse()

	// Creates a router without any middleware by default
	r := gin.New()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	//日志
	r.Use(gin.Logger())
	logFile := newLogFile()
	defer logFile.Close()
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	//连接数据库
	dbConfig := mysql.Config{User: *username, Passwd: *passwd, Net: "tcp", Addr: fmt.Sprintf("%s:%s", *host, *port), DBName: *dbname}
	log.Printf(fmt.Sprintf("try to connect database: %s", dbConfig.FormatDSN()))

	datasource, err := sql.Open("mysql", dbConfig.FormatDSN())
	if err != nil {
		panic(err.Error())
	}
	datasource.SetMaxOpenConns(100)
	datasource.SetMaxIdleConns(10)
	datasource.SetConnMaxLifetime(time.Minute * 25)
	defer datasource.Close()

	//注册路由
	router.RegisterRoutes(core.NewContext(datasource), r)

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
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
