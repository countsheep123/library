package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/countsheep123/library/db"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type BookListResponse struct {
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

// BookList : List Book
func (s *Server) BookList(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	opt := bookListOpt(c, user)

	res := []*BookListResponse{}
	var total uint64

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		var books []*db.BookReadOutput
		if err := tx.BookRead(ctx, opt, &books); err != nil {
			return err
		}

		count, err := tx.BookCount(ctx, opt)
		if err != nil {
			return err
		}
		total = count

		for _, b := range books {
			var labels []*db.BookLabelReadOutput
			if err := tx.BookLabelRead(ctx, &db.Option{
				Filters: map[string][]string{
					"book_id": []string{b.ID},
				},
			}, &labels); err != nil {
				return err
			}

			var recommenders []*db.BookRecommenderReadOutput
			if err := tx.BookRecommenderRead(ctx, &db.Option{
				Filters: map[string][]string{
					"book_id": []string{b.ID},
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

			res = append(res, &BookListResponse{
				ID:           b.ID,
				CreatedAt:    b.CreatedAt,
				UpdatedAt:    b.UpdatedAt,
				Title:        b.Title,
				ISBN:         strValue(b.ISBN),
				Publisher:    strValue(b.Publisher),
				Pubdate:      strValue(b.Pubdate),
				Authors:      strArrayValue(b.Authors),
				Cover:        strValue(b.Cover),
				Labels:       convLabels(labels),
				Recommenders: convRecommenders(recommenders),
				Stocks:       convStocks(stocks),
			})
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"books": res,
	})
}

func bookListOpt(c *gin.Context, user *User) *db.Option {
	opt := &db.Option{
		Filters: map[string][]string{},
		Sort:    &db.Sort{},
		Range:   &db.Range{},
	}

	if vs, ok := c.GetQueryArray("title"); ok {
		opt.Filters["title"] = vs
	}
	if vs, ok := c.GetQueryArray("isbn"); ok {
		opt.Filters["isbn"] = vs
	}
	if vs, ok := c.GetQueryArray("publisher"); ok {
		opt.Filters["publisher"] = vs
	}
	if vs, ok := c.GetQueryArray("authors"); ok {
		opt.Filters["authors"] = vs
	}
	if vs, ok := c.GetQueryArray("labels"); ok {
		opt.Filters["labels"] = vs
	}
	if vs, ok := c.GetQueryArray("owner"); ok {
		opt.Filters["owner"] = vs
	}

	own := c.DefaultQuery("own", "false")
	switch own {
	case "true":
		opt.Filters["owner_id"] = []string{user.ID}
	case "false":
		// do nothing
	default:
		// do nothing
	}

	sort := c.DefaultQuery("sort", "title:asc")
	splited := strings.Split(sort, ":")
	if len(splited) != 2 {
		opt.Sort.By = "title"
		opt.Sort.IsAsc = true
	} else {
		switch splited[0] {
		case "created_at", "title", "isbn", "publisher", "pubdate", "authors", "owner":
			opt.Sort.By = splited[0]
		default:
			opt.Sort.By = "title"
		}

		switch splited[1] {
		case "ASC", "asc":
			opt.Sort.IsAsc = true
		case "DESC", "desc":
			opt.Sort.IsAsc = false
		default:
			opt.Sort.IsAsc = true
		}
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err == nil {
		opt.Range.Limit = limit
	}
	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.ParseUint(offsetStr, 10, 64)
	if err == nil {
		opt.Range.Offset = offset
	}

	return opt
}
