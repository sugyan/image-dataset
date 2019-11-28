package app

type imageResponse struct {
	ID          string `json:"id"`
	ImageURL    string `json:"image_url"`
	Size        int    `json:"size"`
	Status      int    `json:"status"`
	Parts       []int  `json:"parts"`
	LabelName   string `json:"label_name"`
	SourceURL   string `json:"source_url"`
	PhotoURL    string `json:"photo_url"`
	PublishedAt int64  `json:"published_at"`
	UpdatedAt   int64  `json:"updated_at"`
	Meta        string `json:"meta"`
}
