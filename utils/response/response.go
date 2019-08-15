package response

type PaginationMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type IndexResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}
