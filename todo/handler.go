package todo

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type createRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	Tags        []string   `json:"tags"`
	DueDate     *string    `json:"due_date"`
}

type updateRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
	Status      Status   `json:"status"`
	Tags        []string `json:"tags"`
}

// RegisterRoutes mounts all todo endpoints on the given RouterGroup.
func RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", listTodos)
	rg.POST("", createTodo)
	rg.GET("/:id", getTodo)
	rg.PUT("/:id", updateTodo)
	rg.DELETE("/:id", deleteTodo)
	rg.PATCH("/:id/complete", completeTodo)
}

func listTodos(c *gin.Context) {
	todos := GetAll()

	statusFilter := c.Query("status")
	priorityFilter := c.Query("priority")
	if statusFilter != "" || priorityFilter != "" {
		filtered := make([]*Todo, 0, len(todos))
		for _, t := range todos {
			if statusFilter != "" && string(t.Status) != statusFilter {
				continue
			}
			if priorityFilter != "" && string(t.Priority) != priorityFilter {
				continue
			}
			filtered = append(filtered, t)
		}
		todos = filtered
	}
	c.JSON(http.StatusOK, gin.H{"data": todos, "total": len(todos)})
}

func createTodo(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t := &Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Tags:        req.Tags,
	}
	created := Create(t)
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

func getTodo(c *gin.Context) {
	t := GetByID(c.Param("id"))
	if t == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

func updateTodo(c *gin.Context) {
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := &Todo{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		Tags:        req.Tags,
	}
	updated := Update(c.Param("id"), updates)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func deleteTodo(c *gin.Context) {
	if !Delete(c.Param("id")) {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func completeTodo(c *gin.Context) {
	updated := Update(c.Param("id"), &Todo{Status: StatusDone})
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}
