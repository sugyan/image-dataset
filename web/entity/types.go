package entity

import "time"

// Status Type
type Status int

// Kind name
const (
	KindNameImage = "Image"
)

// Status values
const (
	StatusReady Status = iota
	StatusNG
	StatusPending
	StatusOK
)

// Image type
type Image struct {
	ImageURL    string
	SourceURL   string
	PhotoURL    string
	Size        int
	Size0256    bool
	Size0512    bool
	Size1024    bool
	Parts       []int
	LabelName   string
	Status      Status
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Meta        []byte `datastore:",noindex"`
}
