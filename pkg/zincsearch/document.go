package zincsearch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Document 文档结构
type Document struct {
	ID     string                 `json:"_id,omitempty"`
	Index  string                 `json:"_index,omitempty"`
	Source map[string]interface{} `json:"_source,omitempty"`
}

// IndexDocumentRequest 索引文档请求
type IndexDocumentRequest struct {
	Index   string                   `json:"index"`
	Records []map[string]interface{} `json:"records"`
}

// BulkResponse 批量操作响应
type BulkResponse struct {
	Message    string `json:"message"`
	RecordID   string `json:"record_id,omitempty"`
	Successful int    `json:"successful,omitempty"`
	Failed     int    `json:"failed,omitempty"`
}

// IndexDocument 索引单个文档
func (c *Client) IndexDocument(index string, doc map[string]interface{}) (*BulkResponse, error) {
	return c.IndexDocuments(index, []map[string]interface{}{doc})
}

// IndexDocuments 批量索引文档
func (c *Client) IndexDocuments(index string, docs []map[string]interface{}) (*BulkResponse, error) {
	// 为每个文档添加时间戳
	for i := range docs {
		if docs[i] == nil {
			docs[i] = make(map[string]interface{})
		}
		if _, exists := docs[i]["@timestamp"]; !exists {
			docs[i]["@timestamp"] = time.Now().Format(time.RFC3339)
		}
	}

	req := IndexDocumentRequest{
		Index:   index,
		Records: docs,
	}

	endpoint := fmt.Sprintf("/api/%s/_bulk", index)
	resp, err := c.makeRequest("POST", endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("index documents: %w", err)
	}
	defer resp.Body.Close()

	var result BulkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode bulk response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &result, fmt.Errorf("index documents failed, status: %d, message: %s", resp.StatusCode, result.Message)
	}

	return &result, nil
}

// GetDocument 获取单个文档
func (c *Client) GetDocument(index, docID string) (*Document, error) {
	endpoint := fmt.Sprintf("/api/%s/_doc/%s", index, docID)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("document not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get document failed, status: %d", resp.StatusCode)
	}

	var doc Document
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("decode document response: %w", err)
	}

	return &doc, nil
}

// UpdateDocument 更新文档
func (c *Client) UpdateDocument(index, docID string, doc map[string]interface{}) error {
	endpoint := fmt.Sprintf("/api/%s/_doc/%s", index, docID)
	resp, err := c.makeRequest("PUT", endpoint, doc)
	if err != nil {
		return fmt.Errorf("update document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update document failed, status: %d", resp.StatusCode)
	}

	return nil
}

// DeleteDocument 删除文档
func (c *Client) DeleteDocument(index, docID string) error {
	endpoint := fmt.Sprintf("/api/%s/_doc/%s", index, docID)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete document failed, status: %d", resp.StatusCode)
	}

	return nil
}
