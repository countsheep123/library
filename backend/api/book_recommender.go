package api

import (
	"github.com/countsheep123/library/db"
)

type recommender struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func convRecommenders(recommenders []*db.BookRecommenderReadOutput) []*recommender {
	rs := []*recommender{}
	for _, r := range recommenders {
		rs = append(rs, &recommender{
			ID:   r.UserID,
			Name: r.UserName,
		})
	}
	return rs
}
