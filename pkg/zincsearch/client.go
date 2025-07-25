package zincsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"justus/pkg/setting"
)

// Client ZincSearch 客户端
type Client struct {
	Host     string
	Username string
	Password string
	Timeout  time.Duration
	client   *http.Client
}

// NewClient 创建新的 ZincSearch 客户端
func NewClient() *Client {
	config := setting.ZincSearchSetting

	return &Client{
		Host:     config.Host,
		Username: config.Username,
		Password: config.Password,
		Timeout:  time.Duration(config.Timeout) * time.Second,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// NewCustomClient 创建自定义配置的 ZincSearch 客户端
func NewCustomClient(host, username, password string, timeout int) *Client {
	return &Client{
		Host:     host,
		Username: username,
		Password: password,
		Timeout:  time.Duration(timeout) * time.Second,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

// makeRequest 执行HTTP请求
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := fmt.Sprintf("%s%s", c.Host, endpoint)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}

	return resp, nil
}

// Ping 检查 ZincSearch 服务是否可用
func (c *Client) Ping() error {
	resp, err := c.makeRequest("GET", "/version", nil)
	if err != nil {
		return fmt.Errorf("ping zincsearch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zincsearch not available, status: %d", resp.StatusCode)
	}

	return nil
}

// GetVersion 获取 ZincSearch 版本信息
func (c *Client) GetVersion() (map[string]interface{}, error) {
	resp, err := c.makeRequest("GET", "/version", nil)
	if err != nil {
		return nil, fmt.Errorf("get version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get version failed, status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode version response: %w", err)
	}

	return result, nil
}
