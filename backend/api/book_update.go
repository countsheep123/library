package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/countsheep123/library/db"
	"github.com/countsheep123/library/obj"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BookUpdateRequest struct {
	Title     *string  `json:"title"`
	ISBN      *string  `json:"isbn"`
	Publisher *string  `json:"publisher"`
	Pubdate   *string  `json:"pubdate"`
	Authors   []string `json:"authors"`
	Cover     *string  `json:"cover"`
}

func (r *BookUpdateRequest) Validate() error {
	if r.Title != nil && len(*r.Title) == 0 {
		return obj.InvalidRequest{
			Msg: "title is required",
		}
	}
	return nil
}

type BookUpdateResponse struct {
	ID        string   `json:"id"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Title     string   `json:"title"`
	ISBN      string   `json:"isbn"`
	Publisher string   `json:"publisher"`
	Pubdate   string   `json:"pubdate"`
	Authors   []string `json:"authors"`
	Cover     string   `json:"cover"`
}

// BookUpdate : Update Book
func (s *Server) BookUpdate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	id := c.Param("book_id")

	var req *BookUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		errResponse(c, err)
		return
	}

	ts := fmt.Sprint(time.Now().UnixNano())

	res := &BookUpdateResponse{
		ID: id,
	}

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		opt := &db.Option{
			Filters: map[string][]string{
				"id": []string{id},
			},
		}
		var books []*db.BookReadOutput
		if err := tx.BookRead(ctx, opt, &books); err != nil {
			return err
		}

		if len(books) == 0 {
			return obj.NotFound{}
		}

		res.CreatedAt = books[0].CreatedAt
		res.UpdatedAt = ts

		filters := map[string]string{
			"id": id,
		}

		if err := tx.BookUpdate(ctx, &db.BookUpdateInput{
			UpdatedAt: ts,
			Title:     req.Title,
			ISBN:      req.ISBN,
			Publisher: req.Publisher,
			Pubdate:   req.Pubdate,
			Authors:   req.Authors,
			Cover:     req.Cover,
		}, filters); err != nil {
			return err
		}

		if req.Title != nil {
			res.Title = *req.Title
		} else {
			res.Title = books[0].Title
		}

		if req.ISBN != nil {
			res.ISBN = *req.ISBN
		} else {
			res.ISBN = strValue(books[0].ISBN)
		}

		if req.Publisher != nil {
			res.Publisher = *req.Publisher
		} else {
			res.Publisher = strValue(books[0].Publisher)
		}

		if req.Pubdate != nil {
			res.Pubdate = *req.Pubdate
		} else {
			res.Pubdate = strValue(books[0].Pubdate)
		}

		if req.Authors != nil {
			res.Authors = req.Authors
		} else {
			res.Authors = strArrayValue(books[0].Authors)
		}

		if req.Cover != nil {
			res.Cover = *req.Cover
		} else {
			res.Cover = strValue(books[0].Cover)
		}

		return nil
	}); err != nil {
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
