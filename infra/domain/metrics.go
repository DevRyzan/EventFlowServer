package domain
 
type MetricsRequest struct {
	From      int64  `query:"from"`     
	To        int64  `query:"to"`         
	EventName string `query:"event_name"`  
	GroupBy   string `query:"group_by"`    
}
 
type MetricsResponse struct {
	TotalCount      int64        `json:"total_count"`
	UniqueUserCount int64        `json:"unique_user_count"`
	GroupBy         string       `json:"group_by,omitempty"` 
	Buckets        []Bucket     `json:"buckets,omitempty"`  
} 
type Bucket struct {
	Key   string `json:"key"`   
	Count int64  `json:"count"`
}
