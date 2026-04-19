package response

type ShortUrlResponse struct {
	ShortCode string `json:"short_code"`
	LongUrl   string `json:"long_url,omitempty"`
	Message   string `json:"message,omitempty"`
}
