package payment

import (
	"net/http"
	"github.com/gin-gonic/gin"
);

type Handler struct {
	svc Service
}

func NewHandler(s Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/payment-line-items", h.Create)
	r.GET("/payment-line-items/:uid", h.GetByUID)
	r.PUT("/payment-line-items/:uid", h.Update)
	r.DELETE("/payment-line-items/:uid", h.Delete)
	r.GET("/contractors/:id/payment-line-items", h.GetByContractor)
}

func (h *Handler) Create(c *gin.Context) {
	var req PaymentLineItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) GetByUID(c *gin.Context) {
	resp, err := h.svc.GetByUID(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Update(c *gin.Context) {
	var req PaymentLineItem
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.Update(c.Param("uid"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Delete(c *gin.Context) {
	err := h.svc.Delete(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetByContractor(c *gin.Context) {
	resp, err := h.svc.GetByContractor(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
