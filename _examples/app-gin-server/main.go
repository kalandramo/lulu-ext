package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// 注册编解码器（副作用导入）
	_ "github.com/kalandramo/lulu-ext/encoding/json"
	_ "github.com/kalandramo/lulu-ext/encoding/xml"

	"github.com/kalandramo/lulu-ext/ratelimit/tokenbucket"
	httpServer "github.com/kalandramo/lulu-ext/transport/http"
	ginDriver "github.com/kalandramo/lulu-ext/transport/http/driver/gin"
	"github.com/kalandramo/lulu-ext/transport/http/middleware/codec"
	"github.com/kalandramo/lulu-ext/transport/http/middleware/cors"
	"github.com/kalandramo/lulu-ext/transport/http/middleware/logging"
	"github.com/kalandramo/lulu-ext/transport/http/middleware/ratelimit"
	"github.com/kalandramo/lulu-ext/transport/http/middleware/recovery"
	"github.com/kalandramo/lulu-ext/transport/http/middleware/requestid"
	"github.com/kalandramo/lulu-ext/transport/http/middleware/timeout"
)

// 请求/响应结构体
type greetRequest struct {
	Name string `json:"name" xml:"name"`
}

type greetResponse struct {
	Message   string `json:"message" xml:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func main() {
	// 创建限流器：100 QPS，突发 200
	limiter, err := tokenbucket.New(100, 200)
	if err != nil {
		panic(err)
	}

	httpSrv := httpServer.NewServer(
		":8080",
		httpServer.WithDriver(ginDriver.NewDriver()),
	)

	// 中间件链（顺序很重要！）
	httpSrv.Use(
		recovery.Middleware(),  // 1. 最外层：捕获 panic
		requestid.Middleware(), // 2. 生成请求 ID
		logging.Middleware( // 3. 访问日志
			logging.WithSkipPaths("/healthz"),
		),
		cors.Middleware( // 4. 跨域
			cors.WithAllowedOrigins("https://app.example.com"),
			cors.WithAllowCredentials(true),
		),
		timeout.Middleware(30*time.Second), // 5. 超时控制
		ratelimit.Middleware(limiter),      // 6. 限流
		codec.Middleware(),                 // 7. 内容协商（最内层）
	)

	httpSrv.GET("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Healthy")
	})

	httpSrv.POST("/greet", func(w http.ResponseWriter, r *http.Request) {
		var req greetRequest
		if err := codec.ReadBody(r, &req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		codec.Respond(w, r, http.StatusOK, &greetResponse{
			Message:   "Hello, " + req.Name + "!",
			RequestID: requestid.FromContext(r.Context()),
		})
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fmt.Printf("HTTP server listening on %s\n", httpSrv.Endpoint())
	if err := httpSrv.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("server stopped gracefully")
}
