package connections

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// adaptMiddleware converts a standard http.Handler middleware into a Gin-compatible middleware
func adaptMiddleware(auth func(http.Handler) http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		called := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			c.Request = r
			c.Next()
		})
		auth(next).ServeHTTP(c.Writer, c.Request)
		if !called {
			c.Abort()
		}
	}
}

// RegisterRoutes initializes the connections module dependencies and mounts its routes on the Gin engine
func RegisterRoutes(r *gin.Engine, db *sql.DB, auth func(http.Handler) http.Handler, bus EventPublisher) {
	repo := NewRepository(db)
	svc := NewService(db, repo, bus)
	h := NewHandler(svc)

	authGin := adaptMiddleware(auth)

	// Register on both prefixes to support standard v1 proxy routing and flat /api/ connections endpoint requests
	for _, prefix := range []string{"/api/v1/connections", "/api/connections"} {
		g := r.Group(prefix)
		g.Use(authGin)

		g.POST("/request", h.RequestConnection)
		g.POST("/:id/accept", h.AcceptConnection)
		g.POST("/:id/decline", h.DeclineConnection)
		g.DELETE("/:id", h.RemoveConnection)
		g.POST("/users/:user_id/block", h.BlockUser)
		g.DELETE("/users/:user_id/block", h.UnblockUser)
		g.GET("", h.GetConnections)
		g.GET("/pending", h.GetPendingRequests)
		g.GET("/mutual/:user_id", h.GetMutualConnections)
		g.GET("/suggestions", h.GetSuggestions)
	}
}
