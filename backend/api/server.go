package api

import (
	"net/http"

	"github.com/countsheep123/library/db"
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	handler    *db.Handler
	staticPath string

	userMap map[string]*User
}

func New(h *db.Handler, sp string) (*Server, error) {
	userMap, err := loadUserMap(h)
	if err != nil {
		return nil, err
	}

	return &Server{
		handler:    h,
		staticPath: sp,
		userMap:    userMap,
	}, nil
}

func (s *Server) Handler() (http.Handler, error) {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.Use(static.Serve("/", static.LocalFile(s.staticPath, false)))

	api := router.Group("/api", s.checkUser())
	{
		api.GET("/context", s.Context)

		users := api.Group("/users")
		{
			users.GET("", s.UserList)
			users.PATCH("/:user_id", s.UserUpdate)
		}

		books := api.Group("/books")
		{
			books.GET("", s.BookList)
			books.POST("", s.BookCreate)
			books.GET("/:book_id", s.BookRead)
			books.PATCH("/:book_id", s.BookUpdate)
			books.DELETE("/:book_id", admin(), s.BookDelete)

			books.POST("/:book_id/labels", s.BookLabelCreate)
			books.DELETE("/:book_id/labels/:label_id", s.BookLabelDelete)
			books.POST("/:book_id/recommends", s.BookRecommenderCreate)
			books.DELETE("/:book_id/recommends", s.BookRecommenderDelete)

			books.POST("/:book_id/stocks", s.StockCreate)
			books.PATCH("/:book_id/stocks/:stock_id", s.StockUpdate)
			books.DELETE("/:book_id/stocks/:stock_id", s.StockDelete)
		}

		marks := api.Group("/marks")
		{
			marks.GET("", s.MarkList)
			marks.POST("", s.MarkCreate)
			marks.DELETE("/:mark_id", s.MarkDelete)
		}

		locations := api.Group("/locations")
		{
			locations.GET("", s.LocationList)
			locations.POST("", admin(), s.LocationCreate)
			locations.DELETE("/:location_id", admin(), s.LocationDelete)
		}

		records := api.Group("/records")
		{
			records.GET("", s.RecordList)
			records.POST("", s.RecordCreate)
			records.DELETE("/:record_id", s.RecordDelete)
		}
	}

	return router, nil
}

func admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		user, err := getUser(ctx)
		if err != nil {
			raven.CaptureError(err, nil)
			zap.S().Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if !user.IsAdmin {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
