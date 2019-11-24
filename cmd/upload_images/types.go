package main

type data struct {
	Angle float32    `json:"angle"`
	Size  int        `json:"size"`
	Parts [68][2]int `json:"parts"`
	Meta  struct {
		PhotoID     string `json:"photo_id"`
		PhotoURL    string `json:"photo_url"`
		SourceURL   string `json:"source_url"`
		PublishedAt string `json:"published_at"`
		LabelID     string `json:"label_id"`
		LabelName   string `json:"label_name"`
	} `json:"meta"`
}
