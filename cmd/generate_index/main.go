package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

const (
	collectionImage      = "Image"
	queryScopeCollection = "COLLECTION"

	nameLabelName = "LabelName"
	nameStatus    = "Status"

	nameSize0256 = "Size0256"
	nameSize0512 = "Size0512"
	nameSize1024 = "Size1024"

	nameID          = "ID"
	nameUpdatedAt   = "UpdatedAt"
	namePublishedAt = "PublishedAt"

	orderAsc  = "ASCENDING"
	orderDesc = "DESCENDING"
)

type indexesData struct {
	Indexes        []*index         `json:"indexes"`
	FieldOverrides []*fieldOverride `json:"fieldOverrides"`
}

type index struct {
	CollectionGroup string   `json:"collectionGroup"`
	QueryScope      string   `json:"queryScope"`
	Fields          []*field `json:"fields"`
}

type fieldOverride struct {
	CollectionGroup string        `json:"collectionGroup"`
	FieldPath       string        `json:"fieldPath"`
	Indexes         []interface{} `json:"indexes"`
}

type field struct {
	FieldPath string `json:"fieldPath"`
	Order     string `json:"order"`
}

func main() {
	indexes := []*index{}
	for _, labelName := range []string{"", nameLabelName} {
		fields := []*field{}
		if labelName != "" {
			fields = append(fields, &field{
				FieldPath: labelName,
				Order:     orderAsc,
			})
		}
		for _, status := range []string{"", nameStatus} {
			fields := fields
			if status != "" {
				fields = append(fields, &field{
					FieldPath: status,
					Order:     orderAsc,
				})
			}
			for _, size := range []string{"", nameSize0256, nameSize0512, nameSize1024} {
				fields := fields
				if size != "" {
					fields = append(fields, &field{
						FieldPath: size,
						Order:     orderAsc,
					})
				}
				for _, order := range []string{nameID, nameUpdatedAt, namePublishedAt} {
					if len(fields) == 0 {
						continue
					}
					indexes = append(indexes,
						&index{
							CollectionGroup: collectionImage,
							QueryScope:      queryScopeCollection,
							Fields: append(fields, &field{
								FieldPath: order,
								Order:     orderAsc,
							}),
						},
						&index{
							CollectionGroup: collectionImage,
							QueryScope:      queryScopeCollection,
							Fields: append(fields, &field{
								FieldPath: order,
								Order:     orderDesc,
							}),
						},
					)
				}
			}
		}
	}
	fieldOverrides := []*fieldOverride{
		&fieldOverride{
			CollectionGroup: collectionImage,
			FieldPath:       "Meta",
			Indexes:         []interface{}{},
		},
		&fieldOverride{
			CollectionGroup: collectionImage,
			FieldPath:       "Parts",
			Indexes:         []interface{}{},
		},
	}
	data := &indexesData{
		Indexes:        indexes,
		FieldOverrides: fieldOverrides,
	}
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stdout, string(out))

	log.Printf("%d indexes, %d fieldOverrides", len(indexes), len(fieldOverrides))
}
