package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
)

type MeasurementController struct {
	MeasurementService app.MeasurementService
	DeviceService      app.DeviceService
}

func NewMeasurementController(ms app.MeasurementService, ds app.DeviceService) *MeasurementController {
	return &MeasurementController{
		DeviceService:      ds,
		MeasurementService: ms,
	}
}

func (c *MeasurementController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var measurementRequest requests.MeasurementRequest
		err := json.NewDecoder(r.Body).Decode(&measurementRequest)
		if err != nil {
			log.Printf("MeasurementController: Error decoding request body: %s", err)
			BadRequest(w, errors.New("invalid request payload"))
			return
		}

		if measurementRequest.DeviceId == 0 {
			err := errors.New("DeviceId is required")
			log.Printf("MeasurementController: %s", err)
			BadRequest(w, err)
			return
		}

		device, err := c.DeviceService.Find(measurementRequest.DeviceId)
		if err != nil {
			log.Printf("MeasurementController: Error finding device: %s", err)
			BadRequest(w, errors.New("device not found"))
			return
		}

		deviceDomain, ok := device.(domain.Device)
		if !ok {
			log.Printf("MeasurementController: Error asserting device to domain.Device")
			InternalServerError(w, errors.New("failed to assert device type"))
			return
		}

		measurement, err := measurementRequest.ToDomainModel()
		if err != nil {
			log.Printf("MeasurementController: Error converting to domain model: %s", err)
			BadRequest(w, err)
			return
		}

		measurement.RoomId = deviceDomain.RoomId

		createdMeasurement, err := c.MeasurementService.Save(measurement)
		if err != nil {
			log.Printf("MeasurementController: %s", err)
			InternalServerError(w, errors.New("failed to save measurement"))
			return
		}

		measurementDto := resources.MeasurementDto{}.DomainToDto(createdMeasurement)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(measurementDto)
	}
}

func (c *MeasurementController) FindByDeviceAndDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceID := chi.URLParam(r, "deviceId")
		if deviceID == "" {
			http.Error(w, "Invalid device_id", http.StatusBadRequest)
			return
		}

		deviceId, err := strconv.ParseUint(deviceID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid device_id", http.StatusBadRequest)
			return
		}

		startDateStr := chi.URLParam(r, "startDate")
		if startDateStr == "" {
			http.Error(w, "Invalid start_date", http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date", http.StatusBadRequest)
			return
		}

		endDateStr := chi.URLParam(r, "endDate")
		log.Printf("Received endDateStr: %s", endDateStr) // Add log here
		if endDateStr == "" {
			http.Error(w, "Invalid end_date", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date", http.StatusBadRequest)
			return
		}

		measurements, err := c.MeasurementService.FindByDeviceAndDate(deviceId, startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(measurements); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func (c *MeasurementController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		measurement := r.Context().Value(MeasurementKey).(domain.Measurement)

		measurementDto := resources.MeasurementDto{}.DomainToDto(measurement)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(measurementDto)
	}
}

func (c *MeasurementController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		measurement := r.Context().Value(MeasurementKey).(domain.Measurement)

		var measurementRequest requests.MeasurementRequest
		err := json.NewDecoder(r.Body).Decode(&measurementRequest)
		if err != nil {
			log.Printf("MeasurementController: Error decoding request body: %s", err)
			BadRequest(w, errors.New("invalid request payload"))
			return
		}

		measurement.Value = measurementRequest.Value
		measurement.DeviceId = measurementRequest.DeviceId
		updatedMeasurement, err := c.MeasurementService.Update(measurement)
		if err != nil {
			log.Printf("MeasurementController: Error updating measurement: %s", err)
			InternalServerError(w, errors.New("failed to update measurement"))
			return
		}

		measurementDto := resources.MeasurementDto{}.DomainToDto(updatedMeasurement)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(measurementDto)
	}
}

func (c *MeasurementController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		measurement := r.Context().Value(MeasurementKey).(domain.Measurement)

		err := c.MeasurementService.Delete(measurement.Id)
		if err != nil {
			log.Printf("MeasurementController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}