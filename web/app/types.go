package app

type imageResponse struct {
	ID        string `json:"id"`
	ImageURL  string `json:"image_url"`
	Size      int    `json:"size"`
	Parts     []int  `json:"parts"`
	LabelName string `json:"label_name"`
	SourceURL string `json:"source_url"`
	PhotoURL  string `json:"photo_url"`
	PostedAt  int64  `json:"posted_at"`
	Meta      string `json:"meta"`
}
