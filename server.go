// File: server.go
// Creation: Wed May 29 16:35:29 2024
// Time-stamp: <2024-06-07 14:50:42>
// Copyright (C): 2024 Pierre Lecocq

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Feedback struct {
	Type    string
	Message string
}

func deleteHandler(w http.ResponseWriter, r *http.Request, db *Database) {
	var f Feedback

	err := r.ParseForm()

	if err != nil {
		f = Feedback{Type: "error", Message: "Invalid form data"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	id, err := strconv.Atoi(r.Form.Get("id"))

	if err != nil || id <= 0 {
		f = Feedback{Type: "error", Message: "Invalid ID"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	err = db.DeleteTask(int64(id))

	if err != nil {
		f = Feedback{Type: "error", Message: "Error"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	w.Header().Set("HX-Trigger", "reload-tasks")
	f = Feedback{Type: "success", Message: "Task deleted successfully"}
	template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
}

func refreshHandler(w http.ResponseWriter, r *http.Request, db *Database) {
	var f Feedback

	err := r.ParseForm()

	if err != nil {
		f = Feedback{Type: "error", Message: "Invalid form data"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	id, err := strconv.Atoi(r.Form.Get("id"))

	if err != nil || id <= 0 {
		f = Feedback{Type: "error", Message: "Invalid ID"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	err = db.RefreshTaskDueDate(int64(id))

	if err != nil {
		f = Feedback{Type: "error", Message: "Error"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	w.Header().Set("HX-Trigger", "reload-tasks")
	f = Feedback{Type: "success", Message: "Task updated successfully"}
	template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
}

func updateHandler(w http.ResponseWriter, r *http.Request, db *Database) {
	var f Feedback
	var t Task
	var u bool = false

	err := r.ParseForm()

	if err != nil {
		f = Feedback{Type: "error", Message: "Invalid form data"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	id, err := strconv.Atoi(r.Form.Get("id"))

	if err != nil || id <= 0 {
		f = Feedback{Type: "error", Message: "Invalid ID"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	t.ID = int64(id)

	if r.Form.Has("title") {
		title := r.Form.Get("title")

		if len(title) == 0 {
			f = Feedback{Type: "error", Message: "Invalid title"}
			template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
			return
		}

		t.Title = title
		u = true
	}

	if r.Form.Has("state") {
		state := r.Form.Get("state")

		if len(state) == 0 {
			f = Feedback{Type: "error", Message: "Invalid state"}
			template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
			return
		}

		t.State = state
		u = true
	}

	if r.Form.Has("due_date") {
		dstr := r.Form.Get("due_date")

		if len(dstr) == 0 {
			f = Feedback{Type: "error", Message: "Invalid due date"}
			template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
			return
		}

		ddate, err := time.Parse("2006-01-02 15:04:05", dstr)

		if err != nil {
			f = Feedback{Type: "error", Message: "Invalid due date format"}
			template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
			return
		}

		t.DueAt = &ddate
		u = true
	}

	if u == false {
		f = Feedback{Type: "error", Message: "No data to update"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	err = db.UpdateTask(t)

	if err != nil {
		f = Feedback{Type: "error", Message: "Error"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	w.Header().Set("HX-Trigger", "reload-tasks")
	f = Feedback{Type: "success", Message: "Task updated successfully"}
	template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
}

func createHandler(w http.ResponseWriter, r *http.Request, db *Database) {
	var f Feedback

	err := r.ParseForm()

	if err != nil {
		f = Feedback{Type: "error", Message: "Invalid form data"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	title := r.Form.Get("title")

	if len(title) == 0 {
		f = Feedback{Type: "error", Message: "Invalid title"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	err = db.CreateTask(Task{Title: title})

	if err != nil {
		f = Feedback{Type: "error", Message: "Error"}
		template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
		return
	}

	w.Header().Set("HX-Trigger", "reload-tasks")
	f = Feedback{Type: "success", Message: "Task created successfully"}
	template.Must(template.ParseFiles("templates/_feedback.html")).Execute(w, f)
}

func tasksHandler(w http.ResponseWriter, _ *http.Request, db *Database) {
	tmpl := template.Must(template.ParseFiles("templates/_tasks.html"))

	tmpl.Execute(w, struct {
		Tasks []Task
	}{
		Tasks: db.FetchTasks(),
	})
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))

	tmpl.Execute(w, struct {
		Title string
	}{
		Title: "Today or never!",
	})
}

func StartServer(port int, db *Database) {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		deleteHandler(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		refreshHandler(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		updateHandler(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		createHandler(w, r, db)
	}).Methods("POST")

	r.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		tasksHandler(w, r, db)
	}).Methods("GET")

	// Index
	r.HandleFunc("/", homeHandler).Methods("GET")

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Serve
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), r); err != nil {
		log.Printf("error listening: %v", err)
	}
}
