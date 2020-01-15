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

type countResponse struct {
	Size      string `json:"size"`
	Ready     int    `json:"status_ready"`
	NG        int    `json:"status_ng"`
	Pending   int    `json:"status_pending"`
	OK        int    `json:"status_ok"`
	Predicted int    `json:"status_predicted"`
}
