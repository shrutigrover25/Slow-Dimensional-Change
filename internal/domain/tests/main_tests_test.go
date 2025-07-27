package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mercor/internal/domain/jobs"
	"mercor/internal/domain/router"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	router.InitRoutes(r)
	return r
}

func TestJobCRUD(t *testing.T) {
	r := setupRouter()

	// --- CREATE job
	jobCreate := map[string]any{
		"title":        "Backend Developer",
		"status":       "active",
		"rate":         42.5,
		"companyId":    uuid.New().String(),
		"contractorId": uuid.New().String(),
	}
	body, _ := json.Marshal(jobCreate)
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var createdJob jobs.Job
	err := json.Unmarshal(resp.Body.Bytes(), &createdJob)
	assert.Nil(t, err)
	assert.Equal(t, "Backend Developer", createdJob.Title)
	assert.Equal(t, 1, createdJob.Version)

	// --- GET job by UID
	req, _ = http.NewRequest("GET", "/jobs/"+createdJob.UID.String(), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var fetched jobs.Job
	json.Unmarshal(resp.Body.Bytes(), &fetched)
	assert.Equal(t, createdJob.UID, fetched.UID)

	// --- UPDATE (PUT) → creates new version
	update := map[string]any{
		"title":  "Backend Engineer",
		"status": "extended",
	}
	updateJSON, _ := json.Marshal(update)
	req, _ = http.NewRequest("PUT", "/jobs/"+createdJob.UID.String(), bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var updatedJob jobs.Job
	json.Unmarshal(resp.Body.Bytes(), &updatedJob)
	assert.Equal(t, "Backend Engineer", updatedJob.Title)
	assert.Equal(t, 2, updatedJob.Version)
}

func TestTimeLogCRUD(t *testing.T) {
	r := setupRouter()

	// Create Job First
	jobPayload := map[string]any{
		"title":        "TimeLogJob",
		"status":       "active",
		"rate":         55.5,
		"companyId":    uuid.New().String(),
		"contractorId": uuid.New().String(),
	}
	body, _ := json.Marshal(jobPayload)
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var createdJob jobs.Job
	json.Unmarshal(resp.Body.Bytes(), &createdJob)

	// --- CREATE Timelog
	tlog := map[string]any{
		"startTime":    time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		"endTime":      time.Now().Format(time.RFC3339),
		"contractorId": createdJob.ContractorID.String(),
		"jobUid":       createdJob.UID.String(),
	}
	tlogJSON, _ := json.Marshal(tlog)
	req, _ = http.NewRequest("POST", "/timelogs", bytes.NewBuffer(tlogJSON))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestPaymentLineItemFlow(t *testing.T) {
	r := setupRouter()

	// Create Job
	job := map[string]any{
		"title":        "PaymentJob",
		"status":       "active",
		"rate":         100,
		"companyId":    uuid.New().String(),
		"contractorId": uuid.New().String(),
	}
	body, _ := json.Marshal(job)
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var createdJob jobs.Job
	json.Unmarshal(resp.Body.Bytes(), &createdJob)

	// Create Timelog
	tlog := map[string]any{
		"startTime":    time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
		"endTime":      time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		"contractorId": createdJob.ContractorID.String(),
		"jobUid":       createdJob.UID.String(),
	}
	tlogJSON, _ := json.Marshal(tlog)
	req, _ = http.NewRequest("POST", "/timelogs", bytes.NewBuffer(tlogJSON))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var tlogResp map[string]any
	json.Unmarshal(resp.Body.Bytes(), &tlogResp)

	// Create PaymentLineItem
	payment := map[string]any{
		"contractorId": createdJob.ContractorID.String(),
		"amount":       100,
		"issuedAt":     time.Now().Format(time.RFC3339),
		"jobUid":       createdJob.UID.String(),
		"timelogUid":   tlogResp["uid"],
	}
	paymentJSON, _ := json.Marshal(payment)
	req, _ = http.NewRequest("POST", "/payment-line-items", bytes.NewBuffer(paymentJSON))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	fmt.Println("✅ Payment Line Item created and linked successfully.")
}
