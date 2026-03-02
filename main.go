package main

import (
	"context"

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
