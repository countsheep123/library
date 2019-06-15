package api

import (
	"github.com/countsheep123/library/db"
)

type stock struct {
	ID           string `json:"id"`
	IsAvailable  bool   `json:"is_available"`
	OwnerID      string `json:"owner_id"`
	OwnerName    string `json:"owner_name"`
	MarkID       string `json:"mark_id"`
	MarkName     string `json:"mark_name"`
	MarkURL      string `json:"mark_url"`
	LocationID   string `json:"location_id"`
	LocationName string `json:"location_name"`
}

func convStocks(stocks []*db.StockReadOutput) []*stock {
	ss := []*stock{}
	for _, s := range stocks {
		ss = append(ss, &stock{
			ID:           s.ID,
			IsAvailable:  s.IsAvailable,
			OwnerID:      s.UserID,
			OwnerName:    s.UserName,
			MarkID:       s.MarkID,
			MarkName:     s.MarkName,
			MarkURL:      s.MarkURL,
			LocationID:   s.LocationID,
			LocationName: s.LocationName,
		})
	}
	return ss
}
