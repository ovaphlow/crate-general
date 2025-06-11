package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RequestOptions 请求配置选项
type RequestOptions struct {
	Method  string            // HTTP 方法
	URL     string            // 请求 URL
	Headers map[string]string // 请求头
	Body    interface{}       // 请求体
	Timeout time.Duration     // 超时时间
}

// Response 响应结构体
type Response struct {
	StatusCode int               // 状态码
	Headers    map[string]string // 响应头
	Body       []byte            // 响应体
	Text       string            // 响应文本
}

// 基准 request 函数
func request(options RequestOptions) (*Response, error) {
	// 设置默认超时时间
	if options.Timeout == 0 {
		options.Timeout = 30 * time.Second
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: options.Timeout,
	}

	var body io.Reader

	// 处理请求体
	if options.Body != nil {
		switch v := options.Body.(type) {
		case string:
			body = bytes.NewBufferString(v)
		case []byte:
			body = bytes.NewBuffer(v)
		default:
			// JSON 序列化
			jsonData, err := json.Marshal(v)
			if err != nil {
				return nil, fmt.Errorf("序列化请求体失败: %v", err)
			}
			body = bytes.NewBuffer(jsonData)

			// 自动设置 Content-Type 为 JSON
			if options.Headers == nil {
				options.Headers = make(map[string]string)
			}
			if _, exists := options.Headers["Content-Type"]; !exists {
				options.Headers["Content-Type"] = "application/json"
			}
		}
	}

	// 创建请求
	req, err := http.NewRequest(options.Method, options.URL, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 处理响应头
	headers := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       responseBody,
		Text:       string(responseBody),
	}, nil
}

// GET 请求
func Get(url string, headers ...map[string]string) (*Response, error) {
	options := RequestOptions{
		Method: "GET",
		URL:    url,
	}

	if len(headers) > 0 {
		options.Headers = headers[0]
	}

	return request(options)
}

// POST 请求
func Post(url string, body interface{}, headers ...map[string]string) (*Response, error) {
	options := RequestOptions{
		Method: "POST",
		URL:    url,
		Body:   body,
	}

	if len(headers) > 0 {
		options.Headers = headers[0]
	}

	return request(options)
}

// PUT 请求
func Put(url string, body interface{}, headers ...map[string]string) (*Response, error) {
	options := RequestOptions{
		Method: "PUT",
		URL:    url,
		Body:   body,
	}

	if len(headers) > 0 {
		options.Headers = headers[0]
	}

	return request(options)
}

// DELETE 请求
func Delete(url string, headers ...map[string]string) (*Response, error) {
	options := RequestOptions{
		Method: "DELETE",
		URL:    url,
	}

	if len(headers) > 0 {
		options.Headers = headers[0]
	}

	return request(options)
}

// PATCH 请求
func Patch(url string, body interface{}, headers ...map[string]string) (*Response, error) {
	options := RequestOptions{
		Method: "PATCH",
		URL:    url,
		Body:   body,
	}

	if len(headers) > 0 {
		options.Headers = headers[0]
	}

	return request(options)
}

// HEAD 请求
func Head(url string, headers ...map[string]string) (*Response, error) {
	options := RequestOptions{
		Method: "HEAD",
		URL:    url,
	}

	if len(headers) > 0 {
		options.Headers = headers[0]
	}

	return request(options)
}

// 示例函数 - 演示如何使用 HTTP 请求库
func ExampleUsage() {
	// GET 请求示例
	fmt.Println("=== GET 请求示例 ===")
	resp, err := Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Printf("GET 请求失败: %v\n", err)
	} else {
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应体: %s\n", resp.Text)
	}

	// POST 请求示例
	fmt.Println("\n=== POST 请求示例 ===")
	postData := map[string]interface{}{
		"title":  "测试标题",
		"body":   "测试内容",
		"userId": 1,
	}
	resp, err = Post("https://jsonplaceholder.typicode.com/posts", postData)
	if err != nil {
		fmt.Printf("POST 请求失败: %v\n", err)
	} else {
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应体: %s\n", resp.Text)
	}

	// PUT 请求示例
	fmt.Println("\n=== PUT 请求示例 ===")
	putData := map[string]interface{}{
		"id":     1,
		"title":  "更新的标题",
		"body":   "更新的内容",
		"userId": 1,
	}
	resp, err = Put("https://jsonplaceholder.typicode.com/posts/1", putData)
	if err != nil {
		fmt.Printf("PUT 请求失败: %v\n", err)
	} else {
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应体: %s\n", resp.Text)
	}

	// DELETE 请求示例
	fmt.Println("\n=== DELETE 请求示例 ===")
	resp, err = Delete("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Printf("DELETE 请求失败: %v\n", err)
	} else {
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应体: %s\n", resp.Text)
	}

	// 带自定义请求头的示例
	fmt.Println("\n=== 带自定义请求头的示例 ===")
	headers := map[string]string{
		"User-Agent":    "Custom-Client/1.0",
		"Authorization": "Bearer your-token-here",
	}
	resp, err = Get("https://httpbin.org/headers", headers)
	if err != nil {
		fmt.Printf("带请求头的 GET 请求失败: %v\n", err)
	} else {
		fmt.Printf("状态码: %d\n", resp.StatusCode)
		fmt.Printf("响应体: %s\n", resp.Text)
	}
}
