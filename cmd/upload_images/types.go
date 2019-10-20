package main

type data struct {
	Angle float32    `json:"angle"`
	Size  int        `json:"size"`
	Parts [68][2]int `json:"parts"`
	Meta  struct {
		FaceID    string `json:"face_id"`
		PhotoID   string `json:"photo_id"`
		SourceURL string `json:"source_url"`
		PhotoURL  string `json:"photo_url"`
		PostedAt  string `json:"posted_at"`
		LabelID   string `json:"label_id"`
		LabelName string `json:"label_name"`
	} `json:"meta"`
}
