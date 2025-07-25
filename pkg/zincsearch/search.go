package zincsearch

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SearchRequest 搜索请求
type SearchRequest struct {
	Query        map[string]interface{} `json:"query"`
	From         int                    `json:"from,omitempty"`
	Size         int                    `json:"size,omitempty"`
	Sort         []interface{}          `json:"sort,omitempty"`
	Aggregations map[string]interface{} `json:"aggs,omitempty"`
	Highlight    map[string]interface{} `json:"highlight,omitempty"`
	Source       interface{}            `json:"_source,omitempty"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Hits     struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Index     string                 `json:"_index"`
			ID        string                 `json:"_id"`
			Score     float64                `json:"_score"`
			Source    map[string]interface{} `json:"_source"`
			Highlight map[string]interface{} `json:"highlight,omitempty"`
		} `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
}

// Search 执行搜索查询
func (c *Client) Search(index string, req *SearchRequest) (*SearchResponse, error) {
	endpoint := fmt.Sprintf("/api/%s/_search", index)
	resp, err := c.makeRequest("POST", endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed, status: %d", resp.StatusCode)
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	return &result, nil
}

// SimpleSearch 简单文本搜索
func (c *Client) SimpleSearch(index, query string, from, size int) (*SearchResponse, error) {
	req := &SearchRequest{
		Query: map[string]interface{}{
			"match": map[string]interface{}{
				"_all": query,
			},
		},
		From: from,
		Size: size,
	}

	if size <= 0 {
		req.Size = 10 // 默认返回10条
	}

	return c.Search(index, req)
}

// MatchAllSearch 查询所有文档
func (c *Client) MatchAllSearch(index string, from, size int) (*SearchResponse, error) {
	req := &SearchRequest{
		Query: map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		From: from,
		Size: size,
	}

	if size <= 0 {
		req.Size = 10
	}

	return c.Search(index, req)
}

// TermSearch 精确匹配搜索
func (c *Client) TermSearch(index, field, value string, from, size int) (*SearchResponse, error) {
	req := &SearchRequest{
		Query: map[string]interface{}{
			"term": map[string]interface{}{
				field: value,
			},
		},
		From: from,
		Size: size,
	}

	if size <= 0 {
		req.Size = 10
	}

	return c.Search(index, req)
}

// RangeSearch 范围搜索
func (c *Client) RangeSearch(index, field string, gte, lte interface{}, from, size int) (*SearchResponse, error) {
	rangeQuery := make(map[string]interface{})

	if gte != nil {
		rangeQuery["gte"] = gte
	}
	if lte != nil {
		rangeQuery["lte"] = lte
	}

	req := &SearchRequest{
		Query: map[string]interface{}{
			"range": map[string]interface{}{
				field: rangeQuery,
			},
		},
		From: from,
		Size: size,
	}

	if size <= 0 {
		req.Size = 10
	}

	return c.Search(index, req)
}

// BoolSearch 布尔查询
func (c *Client) BoolSearch(index string, must, should, mustNot []map[string]interface{}, from, size int) (*SearchResponse, error) {
	boolQuery := make(map[string]interface{})

	if len(must) > 0 {
		boolQuery["must"] = must
	}
	if len(should) > 0 {
		boolQuery["should"] = should
	}
	if len(mustNot) > 0 {
		boolQuery["must_not"] = mustNot
	}

	req := &SearchRequest{
		Query: map[string]interface{}{
			"bool": boolQuery,
		},
		From: from,
		Size: size,
	}

	if size <= 0 {
		req.Size = 10
	}

	return c.Search(index, req)
}

// SearchWithHighlight 带高亮的搜索
func (c *Client) SearchWithHighlight(index, query string, highlightFields []string, from, size int) (*SearchResponse, error) {
	highlight := map[string]interface{}{
		"fields": make(map[string]interface{}),
	}

	for _, field := range highlightFields {
		highlight["fields"].(map[string]interface{})[field] = map[string]interface{}{}
	}

	req := &SearchRequest{
		Query: map[string]interface{}{
			"match": map[string]interface{}{
				"_all": query,
			},
		},
		From:      from,
		Size:      size,
		Highlight: highlight,
	}

	if size <= 0 {
		req.Size = 10
	}

	return c.Search(index, req)
}
