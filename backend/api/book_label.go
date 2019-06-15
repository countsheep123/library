package api

import (
	"github.com/countsheep123/library/db"
)

type label struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

func convLabels(labels []*db.BookLabelReadOutput) []*label {
	ls := []*label{}
	for _, l := range labels {
		ls = append(ls, &label{
			ID:    l.ID,
			Label: l.Label,
		})
	}
	return ls
}
