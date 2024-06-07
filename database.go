// File: database.go
// Creation: Fri May 24 07:41:56 2024
// Time-stamp: <2024-06-07 14:51:06>
// Copyright (C): 2024 Pierre Lecocq

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	handler *sql.DB
}

type Task struct {
	ID        int64      `sql:"id"`
	Title     string     `sql:"title"`
	State     string     `sql:"state"`
	DueAt     *time.Time `sql:"due_at"`
	CreatedAt *time.Time `sql:"created_at"`
	UpdatedAt *time.Time `sql:"updated_at"`
	Overdue   bool       `sql:"-"`
}

func NewDatabase(filename string) *Database {
	path, created := CreateDatabaseFile(filename)

	db, err := sql.Open("sqlite3", path)

	if err != nil {
		panic(err)
	}

	if created {
		PopulateDatabase(db)
		log.Println("Database populated successfully")
	}

	return &Database{
		handler: db,
	}
}

func CreateDatabaseFile(filename string) (string, bool) {
	created := false
	root := os.Getenv("HOME") + "/Library/Application Support/todayornever" // @see https://github.com/shibukawa/configdir

	if _, err := os.Stat(root); os.IsNotExist(err) {
		os.MkdirAll(root, os.ModePerm)
	}

	path := filepath.Join(root, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_CREATE, 0644)

		if err != nil {
			panic(err)
		}

		defer file.Close()
		created = true
	}

	return path, created
}

func PopulateDatabase(db *sql.DB) {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS task (
    id INTEGER,
    title TEXT,
    state TEXT,
    due_at TIMESTAMP DEFAULT (datetime('now', '+1 day', 'start of day')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id AUTOINCREMENT)
)`)
	if err != nil {
		panic(err)
	}
}

func (db *Database) DeleteTask(id int64) error {
	stmt, err := db.handler.Prepare(
		"DELETE FROM task WHERE id = ?",
	)

	stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RefreshTaskDueDate(id int64) error {
	stmt, err := db.handler.Prepare(
		"UPDATE task SET due_at = datetime('now', '+1 day', 'start of day') WHERE id = ?",
	)

	stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) SwitchTaskState(id int64) error {
	stmt, err := db.handler.Prepare(
		"UPDATE task SET state = CASE state WHEN \"todo\" THEN \"done\" ELSE \"todo\" END WHERE id = ?",
	)

	stmt.Exec(id)
	defer stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) UpdateTask(t Task) error {
	tdb, err := db.FetchTask(t.ID)

	if err != nil {
		return err
	}

	var clauses []string
	var args []interface{}

	if t.Title != "" && t.Title != tdb.Title {
		clauses = append(clauses, "title = ?")
		args = append(args, t.Title)
	}

	if t.State != "" && t.State != tdb.State {
		clauses = append(clauses, "state = ?")
		args = append(args, t.State)
	}

	if t.DueAt != nil && t.DueAt != tdb.DueAt {
		clauses = append(clauses, "due_at = ?")
		args = append(args, t.DueAt)
	}

	args = append(args, t.ID)

	stmt, err := db.handler.Prepare(
		fmt.Sprintf("UPDATE task SET %s WHERE id = ?", strings.Join(clauses, ", ")),
	)

	stmt.Exec(args...)

	defer stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) CreateTask(t Task) error {
	stmt, err := db.handler.Prepare(
		"INSERT INTO task (title, state, due_at) VALUES (?, 'todo', datetime('now', '+1 day', 'start of day'))",
	)

	stmt.Exec(t.Title)
	defer stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) FetchTask(id int64) (Task, error) {
	var t Task

	err := db.handler.QueryRow(
		"SELECT id, title, state, due_at, created_at, updated_at FROM task WHERE id = ?",
		id,
	).Scan(
		&t.ID,
		&t.Title,
		&t.State,
		&t.DueAt,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return t, err
	}

	return t, nil
}

func (db *Database) FetchTasks() []Task {
	var tasks []Task

	rows, err := db.handler.Query(
		"SELECT id, title, state, due_at, created_at, updated_at FROM task ORDER by due_at desc, (case state when 'todo' then 0 when 'done' then 1 end),  created_at desc",
	)
	defer rows.Close()

	err = rows.Err()

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		t := Task{}

		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.State,
			&t.DueAt,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			panic(err)
		}

		t.Overdue = t.DueAt.Before(time.Now())

		tasks = append(tasks, t)
	}

	err = rows.Err()

	if err != nil {
		panic(err)
	}

	return tasks
}
