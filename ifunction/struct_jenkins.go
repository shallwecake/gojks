package ifunction

type HTTPStatusError struct {
	URL    string
	Code   int
	Status string
}

// SearchResult 表示整个 JSON 对象
type SearchResult struct {
	Class       string       `json:"_class"`      // 映射到 JSON 中的 "_class" 键
	Suggestions []Suggestion `json:"suggestions"` // 映射到 JSON 中的 "suggestions" 键
}

// Suggestion 表示 suggestions 数组中的每个对象
type Suggestion struct {
	Name string `json:"name"` // 映射到 JSON 中的 "name" 键
}
