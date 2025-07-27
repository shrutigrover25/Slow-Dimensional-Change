package router

import (
	"github.com/gin-gonic/gin"
	"mercor/internal/db"
	job "mercor/internal/domain/jobs"
	timelog "mercor/internal/domain/timelog"
	payment "mercor/internal/domain/paymentLineItem"
)

func InitRoutes(r *gin.Engine) {
	database := db.Connect()

	// JOB
	jobHandler := job.NewHandler(job.NewService(job.NewRepository(database)))
	jobHandler.RegisterRoutes(r)

	// TIMELOG
	tlHandler := timelog.NewHandler(timelog.NewService(timelog.NewRepository(database)))
	tlHandler.RegisterRoutes(r)

	// PAYMENT
	plHandler := payment.NewHandler(payment.NewService(payment.NewRepository(database)))
	plHandler.RegisterRoutes(r)
}
