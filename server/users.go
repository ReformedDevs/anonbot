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

type editUserForm struct {
	Username string
	Email    string
	IsAdmin  bool
}

func (s *Server) editUser(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id   = mux.Vars(r)["id"]
			u    = &db.User{}
			form = &editUserForm{}
			ctx  = pongo2.Context{
				"form": form,
			}
		)
		if len(id) != 0 {
			if err := c.C.Find(u, id).Error; err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return nil
			}
			s.copyStruct(u, form)
			ctx["title"] = "Edit User"
			ctx["action"] = "Save"
		} else {
			ctx["title"] = "New User"
			ctx["action"] = "Create"
		}
		if r.Method == http.MethodPost {
			s.populateStruct(r.Form, form)
			s.copyStruct(form, u)
			if err := c.C.Save(u).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/users", http.StatusFound)
			return nil
		}
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
