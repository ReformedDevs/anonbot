package server

import (
	"net/http"

	"github.com/ReformedDevs/anonbot/db"
	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

func (s *Server) schedules(w http.ResponseWriter, r *http.Request) {
	schedules := []*db.Schedule{}
	if err := s.database.C.
		Preload("Account").
		Order("account_id").
		Find(&schedules).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.render(w, r, "schedules.html", pongo2.Context{
		"title":     "Schedules",
		"schedules": schedules,
	})
}

type editScheduleForm struct {
	Cron      string
	AccountID int64
}

func (s *Server) editSchedule(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id   = mux.Vars(r)["id"]
			sc   = &db.Schedule{}
			form = &editScheduleForm{}
			ctx  = pongo2.Context{
				"form": form,
			}
		)
		if len(id) != 0 {
			if err := c.C.Find(sc, id).Error; err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return nil
			}
			s.copyStruct(sc, form)
			ctx["title"] = "Edit Schedule"
			ctx["action"] = "Save"
		} else {
			ctx["title"] = "New Schedule"
			ctx["action"] = "Create"
		}
		if r.Method == http.MethodPost {
			s.populateStruct(r.Form, form)
			s.copyStruct(form, sc)
			if err := sc.Calculate(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			if err := c.C.Save(sc).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			s.tweeter.Trigger(nil)
			http.Redirect(w, r, "/schedules", http.StatusFound)
			return nil
		}
		accounts := []*db.Account{}
		if err := c.C.Find(&accounts).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}
		ctx["accounts"] = accounts
		s.render(w, r, "editschedule.html", ctx)
		return nil
	})
}

func (s *Server) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	s.database.Transaction(func(c *db.Connection) error {
		var (
			id = mux.Vars(r)["id"]
			sc = &db.Schedule{}
		)
		if err := c.C.Preload("Account").Find(sc, id).Error; err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return nil
		}
		if r.Method == http.MethodPost {
			if err := c.C.Delete(sc).Error; err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			http.Redirect(w, r, "/schedules", http.StatusFound)
			return nil
		}
		s.render(w, r, "delete.html", pongo2.Context{
			"title": "Delete Schedule",
			"name":  sc.Account.Name,
		})
		return nil
	})
}
