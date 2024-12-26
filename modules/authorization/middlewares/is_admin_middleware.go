package middlewares

import (
	"encoding/json"
	"github.com/google/uuid"
	"gopher_tix/modules/authorization/repositories"
	"net/http"
)

type IsAdminMiddleware struct {
	authorizeRepo    *repositories.AuthorizeRepository
	getCurrentUserID func(r *http.Request) (uuid.UUID, error)
}

func NewIsAdminMiddleware(
	authorizeRepo *repositories.AuthorizeRepository,
	getCurrentUserID func(r *http.Request) (uuid.UUID, error),
) *IsAdminMiddleware {
	return &IsAdminMiddleware{
		authorizeRepo:    authorizeRepo,
		getCurrentUserID: getCurrentUserID,
	}
}

func (m *IsAdminMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := m.getCurrentUserID(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized: failed to get current user",
			})
			return
		}

		isAdmin, err := m.authorizeRepo.IsAdmin(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "internal server error: failed to check admin status",
			})
			return
		}

		if !isAdmin {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "forbidden: admin access required",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// HandlerFunc is a convenience wrapper that returns a http.HandlerFunc directly
func (m *IsAdminMiddleware) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.Handler(next).ServeHTTP(w, r)
	}
}
