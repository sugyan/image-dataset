package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	kindImage = "Image"

	nameLabelName = "LabelName"
	nameStatus    = "Status"

	nameSize0256 = "Size0256"
	nameSize0512 = "Size0512"
	nameSize1024 = "Size1024"

	nameKey         = "__key__"
	nameUpdatedAt   = "UpdatedAt"
	namePublishedAt = "PublishedAt"

	directionDesc = "desc"
)

type indexes struct {
	Indexes []*index
}

type index struct {
	Kind       string
	Properties []*property
}

type property struct {
	Name      string
	Direction string `yaml:",omitempty"`
}

func main() {
	m := indexes{
		Indexes: []*index{},
	}
	for _, labelName := range []string{"", nameLabelName} {
		properties := []*property{}
		if labelName != "" {
			properties = append(properties, &property{
				Name: labelName,
			})
		}
		for _, status := range []string{"", nameStatus} {
			properties := properties
			if status != "" {
				properties = append(properties, &property{
					Name: status,
				})
			}
			for _, size := range []string{"", nameSize0256, nameSize0512, nameSize1024} {
				properties := properties
				if size != "" {
					properties = append(properties, &property{
						Name: size,
					})
				}
				for _, order := range []string{nameKey, nameUpdatedAt, namePublishedAt} {
					m.Indexes = append(m.Indexes,
						&index{
							Kind: kindImage,
							Properties: append(properties, &property{
								Name: order,
							}),
						},
						&index{
							Kind: kindImage,
							Properties: append(properties, &property{
								Name:      order,
								Direction: directionDesc,
							}),
						},
					)
				}
			}
		}
	}
	out, err := yaml.Marshal(&m)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(os.Stdout, string(out))

	log.Printf("%d indexes created", len(m.Indexes))
}
