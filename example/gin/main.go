package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"syscall"
	"time"

	manager "github.com/channingdefoe/go-process-manager"
	"github.com/gin-gonic/gin"
)

func main() {
	serverOne := NewServer("8080")
	serverTwo := NewServer("8081")

	addRoutes(serverOne)
	addRoutes(serverTwo)

	errorProcess := &ErrorClass{}

	manager := manager.NewManager([]manager.Process{errorProcess, serverOne, serverTwo})
	manager.Start()
}

type ErrorClass struct{}

func (s *ErrorClass) Start() error {
	time.Sleep(20 * time.Second)
	return fmt.Errorf("error thrown after 10 seconds")
}

func (s *ErrorClass) Stop() error {
	log.Println("ErrorClass stopped")
	return nil
}

type Server struct {
	port   string
	router *gin.Engine
	srv    *http.Server
}

func NewServer(port string) *Server {
	router := gin.Default()
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{
		port:   port,
		router: router,
		srv:    srv,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Printf("Shutting down server on port: %s, going to sleep for 3 seconds", s.port)
	time.Sleep(3 * time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

func addRoutes(server *Server) {
	server.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	server.router.GET("/error", func(c *gin.Context) {
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	})
}
