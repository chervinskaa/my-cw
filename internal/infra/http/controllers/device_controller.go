package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type DeviceController struct {
	DeviceService app.DeviceService
}

func NewDeviceController(ds app.DeviceService) DeviceController {
	return DeviceController{
		DeviceService: ds,
	}
}

func (c DeviceController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.DeviceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dev, err := req.ToDomainModel()
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		savedDevice, err := c.DeviceService.Save(dev)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		deviceDto := resources.DeviceDto{}.DomainToDto(savedDevice)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(deviceDto)
	}
}

func (c DeviceController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		devices, err := c.DeviceService.FindAll()
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		deviceDtos := make([]resources.DeviceDto, len(devices))
		for i, device := range devices {
			deviceDtos[i] = resources.DeviceDto{}.DomainToDto(device)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deviceDtos)
	}
}

func (c DeviceController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		device, err := c.DeviceService.Find(id)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		deviceDto := resources.DeviceDto{}.DomainToDto(device)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deviceDto)
	}
}

func (c DeviceController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.DeviceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		device, err := req.ToDomainModel()
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		updatedDevice, err := c.DeviceService.Update(device)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		deviceDto := resources.DeviceDto{}.DomainToDto(updatedDevice)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deviceDto)
	}
}

func (c DeviceController) InstallDevice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceId, err := strconv.ParseUint(r.URL.Query().Get("deviceId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid device ID", http.StatusBadRequest)
			return
		}

		roomId, err := strconv.ParseUint(r.URL.Query().Get("roomId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid room ID", http.StatusBadRequest)
			return
		}

		err = c.DeviceService.InstallDevice(deviceId, roomId)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (c DeviceController) UninstallDevice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceId, err := strconv.ParseUint(r.URL.Query().Get("deviceId"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid device ID", http.StatusBadRequest)
			return
		}

		err = c.DeviceService.UninstallDevice(deviceId)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (c DeviceController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := r.Context().Value(DevKey).(domain.Device)

		err := c.DeviceService.Delete(d.Id)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
