package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/countsheep123/library/db"
	"github.com/countsheep123/library/obj"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RecordListResponse struct {
	ID         string  `json:"id"`
	LentAt     string  `json:"lent_at"`
	ReturnedAt *string `json:"returned_at"`
	UserID     string  `json:"user_id"`
	StockID    string  `json:"stock_id"`

	OwnerID   string   `json:"owner_id"`
	OwnerName string   `json:"owner_name"`
	BookID    string   `json:"book_id"`
	Title     string   `json:"title"`
	ISBN      string   `json:"isbn"`
	Publisher string   `json:"publisher"`
	Pubdate   string   `json:"pubdate"`
	Authors   []string `json:"authors"`
	Cover     string   `json:"cover"`

	Labels       []*label       `json:"labels"`
	Recommenders []*recommender `json:"recommenders"`
}

func (s *Server) RecordList(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	opt := recordListOpt(c, user)

	res := []*RecordListResponse{}
	var total uint64

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		var records []*db.RecordReadOutput
		if err := tx.RecordRead(ctx, opt, &records); err != nil {
			return err
		}

		count, err := tx.RecordCount(ctx, opt)
		if err != nil {
			return err
		}
		total = count

		for _, r := range records {
			var stocks []*db.StockReadOutput
			if err := tx.StockRead(ctx, &db.Option{
				Filters: map[string][]string{
					"id": []string{r.StockID},
				},
			}, &stocks); err != nil {
				return err
			}

			if len(stocks) == 0 {
				return obj.Internal{}
			}

			var books []*db.BookReadOutput
			if err := tx.BookRead(ctx, &db.Option{
				Filters: map[string][]string{
					"id": []string{stocks[0].BookID},
				},
			}, &books); err != nil {
				return err
			}

			if len(books) == 0 {
				return obj.Internal{}
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

			var users []*db.UserReadOutput
			if err := tx.UserRead(ctx, &db.Option{
				Filters: map[string][]string{
					"id": []string{stocks[0].UserID},
				},
			}, &users); err != nil {
				return err
			}

			if len(users) == 0 {
				return obj.Internal{}
			}

			res = append(res, &RecordListResponse{
				ID:           r.ID,
				LentAt:       r.LentAt,
				ReturnedAt:   r.ReturnedAt,
				UserID:       r.UserID,
				StockID:      r.StockID,
				OwnerID:      users[0].ID,
				OwnerName:    users[0].Name,
				BookID:       books[0].ID,
				Title:        books[0].Title,
				ISBN:         strValue(books[0].ISBN),
				Publisher:    strValue(books[0].Publisher),
				Pubdate:      strValue(books[0].Pubdate),
				Authors:      strArrayValue(books[0].Authors),
				Cover:        strValue(books[0].Cover),
				Labels:       convLabels(labels),
				Recommenders: convRecommenders(recommenders),
			})
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   total,
		"records": res,
	})
}

func recordListOpt(c *gin.Context, user *User) *db.Option {
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
	if vs, ok := c.GetQueryArray("book_id"); ok {
		opt.Filters["book_id"] = vs
	}
	if vs, ok := c.GetQueryArray("stock_id"); ok {
		opt.Filters["stock_id"] = vs
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

	returned := c.DefaultQuery("borrowed", "false")
	switch returned {
	case "true":
		opt.Filters["returned"] = []string{"false"}
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
		case "lent_at", "returned_at", "title", "isbn", "publisher", "pubdate", "authors", "owner":
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
