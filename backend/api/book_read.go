package api

import (
	"net/http"

	"github.com/countsheep123/library/db"
	"github.com/countsheep123/library/obj"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BookReadResponse struct {
	ID        string   `json:"id"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Title     string   `json:"title"`
	ISBN      string   `json:"isbn"`
	Publisher string   `json:"publisher"`
	Pubdate   string   `json:"pubdate"`
	Authors   []string `json:"authors"`
	Cover     string   `json:"cover"`

	Labels       []*label       `json:"labels"`
	Recommenders []*recommender `json:"recommenders"`

	Stocks []*stock `json:"stocks"`
}

// BookRead : Read Book
func (s *Server) BookRead(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	id := c.Param("book_id")

	opt := &db.Option{
		Filters: map[string][]string{
			"id": []string{id},
		},
	}

	var res *BookReadResponse

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		var books []*db.BookReadOutput
		if err := tx.BookRead(ctx, opt, &books); err != nil {
			return err
		}

		if len(books) == 0 {
			return obj.NotFound{}
		}

		var labels []*db.BookLabelReadOutput
		if err := tx.BookLabelRead(ctx, &db.Option{
			Filters: map[string][]string{
				"book_id": []string{books[0].ID},
			},
		}, &labels); err != nil {
			return err
		}

		var recommenders []*db.BookRecommenderReadOutput
		if err := tx.BookRecommenderRead(ctx, &db.Option{
			Filters: map[string][]string{
				"book_id": []string{books[0].ID},
			},
		}, &recommenders); err != nil {
			return err
		}

		var stocks []*db.StockReadOutput
		if err := tx.StockRead(ctx, &db.Option{
			Filters: map[string][]string{
				"book_id": []string{books[0].ID},
			},
		}, &stocks); err != nil {
			return err
		}

		res = &BookReadResponse{
			ID:           books[0].ID,
			CreatedAt:    books[0].CreatedAt,
			UpdatedAt:    books[0].UpdatedAt,
			Title:        books[0].Title,
			ISBN:         strValue(books[0].ISBN),
			Publisher:    strValue(books[0].Publisher),
			Pubdate:      strValue(books[0].Pubdate),
			Authors:      strArrayValue(books[0].Authors),
			Cover:        strValue(books[0].Cover),
			Labels:       convLabels(labels),
			Recommenders: convRecommenders(recommenders),
			Stocks:       convStocks(stocks),
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
