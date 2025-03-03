package roles

import (
	"github.com/meidum/dns/db"
	"github.com/meidum/dns/util"
	bolt "go.etcd.io/bbolt"
	"net/http"
)

func deleteRole(w http.ResponseWriter, r *http.Request, path string, database *bolt.DB) {
	// Validate initial request with type, path, and header
	if r.Method != "DELETE" {
		util.Responses.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	} else if len(r.URL.Path[len(path):]) == 0 {
		util.Responses.Error(w, http.StatusBadRequest, "role must be specified in path")
		return
	} else if r.Header.Get("Authorization") == "" {
		util.Responses.Error(w, http.StatusUnauthorized, "header 'Authorization' is required")
		return
	}

	// Verify JWT in headers
	token, err := db.TokenFromString(r.Header.Get("Authorization"), database)
	if err != nil {
		util.Responses.Error(w, http.StatusUnauthorized, "failed to authenticate: "+err.Error())
		return
	}

	// Get user from token
	u, err := db.UserFromToken(token, database)
	if err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check user role
	if u.Role != "admin" {
		util.Responses.Error(w, http.StatusForbidden, "user must be of role 'admin'")
		return
	}

	// Delete role
	if err := db.DeleteRole(r.URL.Path[len(path):], database); err != nil {
		util.Responses.Error(w, http.StatusInternalServerError, "failed to delete role: "+err.Error())
		return
	}

	util.Responses.Success(w)
}
