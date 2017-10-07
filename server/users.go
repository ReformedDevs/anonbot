package server

import (
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	users := []*db.User{}
	if err := s.database.C.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "users.html", pongo2.Context{
		"title": "Users",
		"users": users,
	})
}

func (s *Server) editUser(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id  = mux.Vars(r)["id"]
			u   = &db.User{}
			ctx = pongo2.Context{
				"title": "Edit User",
				"user":  u,
			}
		)
		if err := c.C.Find(u, id).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		var (
			isAdmin = u.IsAdmin
		)
		if r.Method == http.MethodPost {
			isAdmin = len(r.Form.Get("is_admin")) != 0
			u.IsAdmin = isAdmin
			if err := c.C.Save(u).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/users", http.StatusFound)
			return nil
		}
		ctx["is_admin"] = isAdmin
		s.render(w, r, "edituser.html", ctx)
		return nil
	})
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id = mux.Vars(r)["id"]
			u  = &db.User{}
		)
		if err := c.C.Find(u, id).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		if r.Method == http.MethodPost {
			if err := c.C.Delete(u).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/users", http.StatusFound)
			return nil
		}
		s.render(w, r, "delete.html", pongo2.Context{
			"title": "Delete User",
			"name":  u.Username,
		})
		return nil
	})
}
