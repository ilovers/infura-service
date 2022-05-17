package rpc

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/okex/infura-service/nacos"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gin-gonic/gin"
)

type Service struct {
	config *Config
	router *gin.Engine
	ethRPC *rpc.Server
}

func New(config *Config) (*Service, error) {
	if err := validateConfig(config); err != nil {
		return nil, err
	}
	// gin api
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	// eth rpc server
	ethRPC := rpc.NewServer()
	apis := getAPIs(config)
	for _, api := range apis {
		if err := ethRPC.RegisterName(api.Namespace, api.Service); err != nil {
			panic(err)
		}
	}
	return &Service{
		config: config,
		router: router,
		ethRPC: ethRPC,
	}, nil
}

func (s *Service) Start() {
	// register rpc service to nacos
	if s.config.NacosUrl != "" {
		nacos.Register(s.config.NacosUrl, s.config.NacosNamespaceId, s.config.NacosServiceName, s.config.NacosServiceAddr)
	}

	// register http router
	s.registerRoutes()

	// http server
	srv := &http.Server{
		Addr:    s.config.Address,
		Handler: s.router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func (s *Service) registerRoutes() {
	s.router.POST("/", func(c *gin.Context) {
		s.ethRPC.ServeHTTP(c.Writer, c.Request)
	})
	s.router.OPTIONS("/", func(c *gin.Context) {
		s.ethRPC.ServeHTTP(c.Writer, c.Request)
	})
}
