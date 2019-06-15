package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/countsheep123/library/db"
	"github.com/countsheep123/library/obj"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type BookCreateRequest struct {
	Title     string   `json:"title"`
	ISBN      *string  `json:"isbn"`
	Publisher *string  `json:"publisher"`
	Pubdate   *string  `json:"pubdate"`
	Authors   []string `json:"authors"`
	Cover     *string  `json:"cover"`
}

func (r *BookCreateRequest) Validate() error {
	if len(r.Title) == 0 {
		return obj.InvalidRequest{
			Msg: "title is required",
		}
	}
	return nil
}

type BookCreateResponse struct {
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

// BookCreate : Create Book
func (s *Server) BookCreate(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	var req *BookCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse(c, err)
		return
	}

	if err := req.Validate(); err != nil {
		errResponse(c, err)
		return
	}

	id := xid.New().String()
	ts := fmt.Sprint(time.Now().UnixNano())

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		if err := tx.BookCreate(ctx, &db.BookCreateInput{
			ID:        id,
			CreatedAt: ts,
			UpdatedAt: ts,
			Title:     req.Title,
			ISBN:      req.ISBN,
			Publisher: req.Publisher,
			Pubdate:   req.Pubdate,
			Authors:   req.Authors,
			Cover:     req.Cover,
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		switch err.(type) {
		case obj.Duplicate:
			var books []*db.BookReadOutput
			if err := s.handler.BookRead(ctx, &db.Option{
				Filters: map[string][]string{
					"isbn": []string{*req.ISBN},
				},
			}, &books); err != nil {
				errResponse(c, err)
				return
			}

			if len(books) == 0 {
				errResponse(c, obj.Internal{})
				return
			}

			c.JSON(http.StatusConflict, gin.H{
				"id": books[0].ID,
			})
			return
		}
		errResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, &BookCreateResponse{
		ID:        id,
		CreatedAt: ts,
		UpdatedAt: ts,
		Title:     req.Title,
		ISBN:      strValue(req.ISBN),
		Publisher: strValue(req.Publisher),
		Pubdate:   strValue(req.Pubdate),
		Authors:   strArrayValue(req.Authors),
		Cover:     strValue(req.Cover),
	})
}
