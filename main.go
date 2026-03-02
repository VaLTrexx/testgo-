package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *pgxpool.Pool

func main() {
	dsn := "postgres://postgres:abhinav@localhost:5432/test"

	db, _ = pgxpool.New(context.Background(), dsn)

	r := gin.Default()

	r.POST("/tasks", createTask)
	r.GET("/tasks", getTasks)
	r.GET("/tasks/:id", getTask)
	r.PUT("/tasks/:id", updateTask)

	r.Run(":8080")
}

func createTask(c *gin.Context) {
	var task Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(
		context.Background(),
		"INSERT INTO tasks (title, completed) VALUES ($1, $2) RETURNING id",
		task.Title,
		task.Completed,
	).Scan(&task.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func getTasks(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT id, title, completed FROM tasks")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.Title, &t.Completed)
		tasks = append(tasks, t)
	}

	c.JSON(http.StatusOK, tasks)
}

func getTask(c *gin.Context) {
	id := c.Param("id")

	var task Task
	err := db.QueryRow(
		context.Background(),
		"SELECT id, title, completed FROM tasks WHERE id=$1",
		id,
	).Scan(&task.ID, &task.Title, &task.Completed)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(
		context.Background(),
		"UPDATE tasks SET title=$1, completed=$2 WHERE id=$3",
		task.Title,
		task.Completed,
		id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	task.ID = 0
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
