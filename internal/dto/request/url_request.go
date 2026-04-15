package request

type CreateShortURLRequest struct {
	LongUrl string `json:"long_url"`
}
