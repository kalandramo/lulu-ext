# 1. 健康检查接口
  curl http://localhost:8080/healthz

  # 2. GET 请求（带请求 ID）
  curl -i http://localhost:8080/healthz

  # 3. POST /greet 接口（JSON 格式）
  curl -X POST http://localhost:8080/greet \
    -H "Content-Type: application/json" \
    -d '{"name": "Alice"}'

  # 4. POST /greet 接口（XML 格式）
  curl -X POST http://localhost:8080/greet \
    -H "Content-Type: application/xml" \
    -d '<greetRequest><name>Bob</name></greetRequest>'

  # 5. 查看响应头（包含 Request-ID）
  curl -i -X POST http://localhost:8080/greet \
    -H "Content-Type: application/json" \
    -d '{"name": "Charlie"}'

  # 6. 带自定义 Header 的请求
  curl -i -X POST http://localhost:8080/greet \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer your-token" \
    -d '{"name": "David"}'

  # 7. 测试跨域（CORS）
  curl -i -X POST http://localhost:8080/greet \
    -H "Content-Type: application/json" \
    -H "Origin: https://app.example.com" \
    -d '{"name": "Eve"}'

  # 8. 完整示例（格式化输出）
  curl -s -X POST http://localhost:8080/greet \
    -H "Content-Type: application/json" \
    -d '{"name": "Lulu"}' | jq .

  预期响应示例：

  {
    "message": "Hello, Alice!",
    "request_id": "550e8400-e29b-41d4-a716-446655440000"
  }

  启动服务器后测试：

  # 终端 1：启动服务器
  cd _examples/app-gin-server
  go run main.go

  # 终端 2：测试 API
  curl -X POST http://localhost:8080/greet -H "Content-Type: application/json" -d '{"name": "World"}'