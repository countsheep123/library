package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/countsheep123/library/db"
	"github.com/countsheep123/library/obj"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

type userKey struct {
}

var (
	uKey = userKey{}
)

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

func newUser(u *db.UserReadOutput) *User {
	return &User{
		ID:      u.ID,
		Name:    u.Name,
		Company: strValue(u.Company),
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}

func (s *Server) Context(c *gin.Context) {
	ctx := c.Request.Context()
	user, err := getUser(ctx)
	if err != nil {
		errResponse(c, err)
		return
	}
	zap.S().Debugf("user_id = %s", user.ID)

	c.JSON(http.StatusOK, user)
}

func (s *Server) checkUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetHeader("X-Forwarded-Email")
		zap.S().Debugf("email = %s", email)

		name := c.GetHeader("X-Forwarded-User")
		zap.S().Debugf("name = %s", name)

		if len(email) == 0 || len(name) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		zap.S().Debugf("header = %v", c.Request.Header)

		ctx := c.Request.Context()

		user, err := s.getUserByEmail(ctx, email)
		if err != nil {
			switch err.(type) {
			case obj.NotFound:
				isAdmin := false

				count, err := s.handler.UserCount(ctx, &db.Option{})
				if err != nil {
					raven.CaptureError(err, nil)
					zap.S().Error(err)
					c.AbortWithStatus(http.StatusInternalServerError)
				}
				if count == 0 {
					isAdmin = true
				}

				if err := s.createUser(ctx, email, name, isAdmin); err != nil {
					raven.CaptureError(err, nil)
					zap.S().Error(err)
					c.AbortWithStatus(http.StatusInternalServerError)
				}
				user, err = s.getUserByEmail(ctx, email)
				if err != nil {
					raven.CaptureError(err, nil)
					zap.S().Error(err)
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
			default:
				raven.CaptureError(err, nil)
				zap.S().Error(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		}

		zap.S().Debugf("user_id = %s", user.ID)

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), uKey, user))

		c.Next()
	}
}

func (s *Server) getUserByEmail(ctx context.Context, email string) (*User, error) {
	user, ok := s.userMap[email]
	if ok {
		return user, nil
	}

	opt := &db.Option{
		Filters: map[string][]string{
			"email": []string{email},
		},
	}

	var us []*db.UserReadOutput
	if err := s.handler.UserRead(ctx, opt, &us); err != nil {
		return nil, err
	}

	if len(us) == 0 {
		return nil, obj.NotFound{}
	}

	user = newUser(us[0])

	s.userMap[email] = user

	return user, nil
}

func (s *Server) createUser(ctx context.Context, email, name string, isAdmin bool) error {
	id := xid.New().String()
	ts := fmt.Sprint(time.Now().UnixNano())

	if err := s.handler.Transact(ctx, func(tx *db.Transaction) error {
		if err := tx.UserCreate(ctx, &db.UserCreateInput{
			ID:        id,
			CreatedAt: ts,
			UpdatedAt: ts,
			Email:     email,
			Name:      name,
			// Company:   "",
			IsAdmin: isAdmin,
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func getUser(ctx context.Context) (*User, error) {
	v := ctx.Value(uKey)
	if v == nil {
		return nil, obj.Internal{}
	}
	u, ok := v.(*User)
	if !ok {
		return nil, obj.Internal{}
	}
	return u, nil
}

func loadUserMap(h *db.Handler) (map[string]*User, error) {
	ctx := context.Background()

	opt := &db.Option{}

	var users []*db.UserReadOutput
	if err := h.UserRead(ctx, opt, &users); err != nil {
		return nil, err
	}

	userMap := map[string]*User{}
	for _, u := range users {
		userMap[u.Email] = newUser(u)
	}

	return userMap, nil
}
