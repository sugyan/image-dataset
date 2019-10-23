package app

type imageResponse struct {
	ID        string `json:"id"`
	ImageURL  string `json:"image_url"`
	Size      int    `json:"size"`
	LabelName string `json:"label_name"`
}
