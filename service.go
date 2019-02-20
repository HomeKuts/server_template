package server_template

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Info struct {
	Version string `json:"ver"`
}

var info Info

func init() {
	
	// Встановлюємо рівень налагодження gin: debug, realese, test
	switch *GinMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

//	gin.DisableConsoleColor()
}

// Налаштувати сервер на роботу. Налаштування винесено в окрему функцію для потреб тестування
func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(handlerMiddleware())
	r.Use(gin.Recovery())

	r.GET("/", handlerRoot)
	r.GET("/info", handlerInfo)

	return r
}

// Старт роботи сервера. 
// Корректно припинити роботу сервера можливо пославши йому сигнали: SIGINT або SIGTERM:
// kill -2 pid або нажати на клавіатурі Ctrl+C.  На сигнал SIGQUIT (kill -3 pid) сервер видасть 
// в лог файл поточні показаники роботи
func Start(versionMajor, versionMin string) {
	info = Info{Version: versionMajor + "." + versionMin}
	log.Infof("Server ver.%v start", info.Version)

	router := setupRouter()

	srv := &http.Server{
		Addr: *Addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: *Timeout,
		ReadTimeout:  *Timeout,
		IdleTimeout:  *Timeout * 4,
		Handler:      router,
	}

	go func() {
		log.Infof("Server listen: %v", *Addr)

		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: %s\n", err)
		}
	}()

	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exit_chan := make(chan int)
	go func() {
		for {
			s := <-signal_chan
			switch s {

			// kill -SIGQUIT XXXX
			// ps aux | grep -i cmd | grep -v grep | awk {'print $2'} | xargs kill -3
			case syscall.SIGQUIT:
				printStatus()

			// kill -SIGINT XXXX or Ctrl+c
			// ps aux | grep -i cmd | grep -v grep | awk {'print $2'} | xargs kill -2
			case syscall.SIGINT, syscall.SIGTERM:
				exit_chan <- 1

			default:
				log.Info("Unknown signal.")
				exit_chan <- 1
			}
		}
	}()
	<-exit_chan

	log.Info("Server shutdown, wait 5 seconds")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
	}
	log.Info("Server stop")
}

func handlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		request := c.Request
		path := request.URL.Path
		raw := request.URL.RawQuery
		method := request.Method

		origin := request.Header.Get("Origin")

		// Check for valid Origin
		if origin == *Origin {
			//			c.Header("Access-Control-Allow-Origin", "*")
			//			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			//			c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			c.Next()
		} else {
			log.Errorf("error origin: \"%v\"", origin)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "forbidden"})
		}

		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		latency := time.Since(start)

		log.Infof("%s %s %s %s %d %s %s", clientIP, origin, method, path, statusCode, latency, errorMessage)
	}
}

func printStatus() {
	log.Debugf("Server version: %v", info.Version)
}

func handlerRoot(c *gin.Context) {
	c.String(http.StatusOK, "")
}

func handlerInfo(c *gin.Context) {
	c.JSON(http.StatusOK, info)
}
