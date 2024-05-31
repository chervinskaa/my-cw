package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/gin-gonic/gin"
)

type MeasurementController struct {
	service app.MeasurementService
}

func NewMeasurementController(service app.MeasurementService) *MeasurementController {
	return &MeasurementController{
		service: service,
	}
}

func (ctrl *MeasurementController) SaveMeasurement(c *gin.Context) {
	var measurement domain.Measurement
	if err := c.ShouldBindJSON(&measurement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	measurement, err := ctrl.service.Save(measurement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, measurement)
}

func (ctrl *MeasurementController) GetMeasurementsByDeviceAndDate(c *gin.Context) {
	deviceId, err := strconv.ParseUint(c.Query("device_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device_id"})
		return
	}

	startDate, err := time.Parse("2006-01-02", c.Query("start_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date"})
		return
	}

	endDate, err := time.Parse("2006-01-02", c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date"})
		return
	}

	measurements, err := ctrl.service.FindByDeviceAndDate(deviceId, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, measurements)
}

func (ctrl *MeasurementController) UpdateMeasurement(c *gin.Context) {
	var measurement domain.Measurement
	if err := c.ShouldBindJSON(&measurement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	measurement, err := ctrl.service.Update(measurement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, measurement)
}

func (ctrl *MeasurementController) DeleteMeasurement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid measurement ID"})
		return
	}

	err = ctrl.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Measurement deleted successfully"})
}
