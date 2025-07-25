package zincsearch

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// IndexMapping 索引映射配置
type IndexMapping struct {
	Properties map[string]interface{} `json:"properties,omitempty"`
	Settings   map[string]interface{} `json:"settings,omitempty"`
}

// CreateIndexRequest 创建索引请求
type CreateIndexRequest struct {
	Name        string       `json:"name"`
	StorageType string       `json:"storage_type,omitempty"`
	Mappings    IndexMapping `json:"mappings,omitempty"`
}

// IndexInfo 索引信息
type IndexInfo struct {
	Name        string `json:"name"`
	StorageType string `json:"storage_type"`
	DocNum      int64  `json:"doc_num"`
	StorageSize int64  `json:"storage_size"`
	WALSize     int64  `json:"wal_size"`
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(name string, mapping *IndexMapping) error {
	req := CreateIndexRequest{
		Name:        name,
		StorageType: "disk",
	}

	if mapping != nil {
		req.Mappings = *mapping
	}

	resp, err := c.makeRequest("POST", "/api/index", req)
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create index failed, status: %d", resp.StatusCode)
	}

	return nil
}

// DeleteIndex 删除索引
func (c *Client) DeleteIndex(name string) error {
	endpoint := fmt.Sprintf("/api/index/%s", name)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("delete index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete index failed, status: %d", resp.StatusCode)
	}

	return nil
}

// ListIndices 列出所有索引
func (c *Client) ListIndices() ([]IndexInfo, error) {
	resp, err := c.makeRequest("GET", "/api/index", nil)
	if err != nil {
		return nil, fmt.Errorf("list indices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list indices failed, status: %d", resp.StatusCode)
	}

	var indices []IndexInfo
	if err := json.NewDecoder(resp.Body).Decode(&indices); err != nil {
		return nil, fmt.Errorf("decode indices response: %w", err)
	}

	return indices, nil
}

// IndexExists 检查索引是否存在
func (c *Client) IndexExists(name string) (bool, error) {
	endpoint := fmt.Sprintf("/api/index/%s", name)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("check index exists: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("check index exists failed, status: %d", resp.StatusCode)
}

// GetIndexInfo 获取索引信息
func (c *Client) GetIndexInfo(name string) (*IndexInfo, error) {
	endpoint := fmt.Sprintf("/api/index/%s", name)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("get index info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get index info failed, status: %d", resp.StatusCode)
	}

	var indexInfo IndexInfo
	if err := json.NewDecoder(resp.Body).Decode(&indexInfo); err != nil {
		return nil, fmt.Errorf("decode index info response: %w", err)
	}

	return &indexInfo, nil
}
