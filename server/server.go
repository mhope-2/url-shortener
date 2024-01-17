package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
)

type Config struct {
	Port string
}

type Server struct {
	*gin.Engine
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func New() Server {
	server := gin.Default()
	server.GET("/", healthCheck)
	return Server{server}
}

func Start(e *Server, cfg *Config) {
	s := &http.Server{
		Addr:    cfg.Port,
		Handler: e.Engine,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		if err := s.Close(); err != nil {
			log.Println("Failed to Shutdown Server")
		}
		log.Println("Server Shut Down")
	}()

	if err := s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server Closed After Interruption")
		} else {
			log.Println("Unexpected Server Shutdown. err: ", err)
		}
	}
}
