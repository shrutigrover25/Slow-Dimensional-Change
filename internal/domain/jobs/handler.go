package jobs

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(s Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/jobs", h.Create)
	r.GET("/jobs/:uid", h.GetByUID)
	r.PUT("/jobs/:uid", h.Update)
	r.PUT("/jobs/:uid/status", h.UpdateStatus)
	r.GET("/companies/:id/jobs", h.GetByCompany)
}

func (h *Handler) Create(c *gin.Context) {
	var job Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, job)
}

func (h *Handler) GetByUID(c *gin.Context) {
	job, err := h.svc.GetByUID(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *Handler) Update(c *gin.Context) {
	var updated Job
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	job, err := h.svc.Update(c.Param("uid"), updated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	status := c.Query("status")
	job, err := h.svc.UpdateStatus(c.Param("uid"), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *Handler) GetByCompany(c *gin.Context) {
	jobs, err := h.svc.GetActiveJobsByCompany(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}
