package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type DeviceController struct {
	DeviceService       app.DeviceService
	RoomService         app.RoomService
	OrganizationService app.OrganizationService
}

func NewDeviceController(ds app.DeviceService, rs app.RoomService, os app.OrganizationService) DeviceController {
	return DeviceController{
		DeviceService:       ds,
		RoomService:         rs,
		OrganizationService: os,
	}
}

func (c *DeviceController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var deviceRequest requests.DeviceRequest
		err := json.NewDecoder(r.Body).Decode(&deviceRequest)
		if err != nil {
			log.Printf("DeviceController: Error decoding request body: %s", err)
			BadRequest(w, errors.New("invalid request payload"))
			return
		}

		if deviceRequest.OrganizationId == 0 {
			err := errors.New("OrganizationId is required")
			log.Printf("DeviceController: %s", err)
			BadRequest(w, err)
			return
		}

		_, err = c.OrganizationService.Find(deviceRequest.OrganizationId)
		if err != nil {
			log.Printf("DeviceController: Error finding organization: %s", err)
			BadRequest(w, errors.New("organization not found"))
			return
		}

		device, err := deviceRequest.ToDomainModel()
		if err != nil {
			log.Printf("DeviceController: Error converting to domain model: %s", err)
			BadRequest(w, err)
			return
		}

		createdDevice, err := c.DeviceService.Save(device)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			InternalServerError(w, errors.New("failed to save device"))
			return
		}

		deviceDto := resources.DeviceDto{}.DomainToDto(createdDevice)
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

func (c *DeviceController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		device := r.Context().Value(DevKey).(domain.Device)

		deviceDto := resources.DeviceDto{}.DomainToDto(device)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deviceDto)
	}
}

func (c *DeviceController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		device := r.Context().Value(DevKey).(domain.Device)

		var deviceRequest requests.DeviceRequest
		err := json.NewDecoder(r.Body).Decode(&deviceRequest)
		if err != nil {
			log.Printf("DeviceController: Error decoding request body: %s", err)
			BadRequest(w, errors.New("invalid request payload"))
			return
		}

		device.Characteristics = deviceRequest.Characteristics
		device.PowerConsumption = deviceRequest.PowerConsumption
		device.Units = deviceRequest.Units

		updatedDevice, err := c.DeviceService.Update(device)
		if err != nil {
			log.Printf("DeviceController: Error updating device: %s", err)
			InternalServerError(w, errors.New("failed to update device"))
			return
		}

		deviceDto := resources.DeviceDto{}.DomainToDto(updatedDevice)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deviceDto)
	}
}

func (c *DeviceController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		device := r.Context().Value(DevKey).(domain.Device)

		err := c.DeviceService.Delete(device.Id)
		if err != nil {
			log.Printf("DeviceController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c *DeviceController) Install() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		device := r.Context().Value(DevKey).(domain.Device)

		var req requests.DeviceRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("DeviceController: Error decoding request body: %s", err)
			BadRequest(w, errors.New("invalid request payload"))
			return
		}

		if req.RoomId == nil {
			BadRequest(w, errors.New("roomId is required"))
			return
		}

		_, err = c.RoomService.Find(*req.RoomId)
		if err != nil {
			log.Printf("DeviceController: Room not found: %s", err)
			BadRequest(w, errors.New("room not found"))
			return
		}

		device.RoomId = req.RoomId
		updatedDevice, err := c.DeviceService.Update(device)
		if err != nil {
			log.Printf("DeviceController: Error installing device: %s", err)
			InternalServerError(w, errors.New("failed to install device"))
			return
		}

		deviceDto := resources.DeviceDto{}.DomainToDto(updatedDevice)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deviceDto)
	}
}

func (c *DeviceController) Uninstall() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отримання пристрою з контексту
		device, ok := r.Context().Value(DevKey).(domain.Device)
		if !ok {
			log.Printf("DeviceController: Error getting device from context")
			InternalServerError(w, errors.New("failed to get device from context"))
			return
		}

		// Виклик сервісу для від'єднання пристрою
		uninstalledDevice, err := c.DeviceService.UninstallDevice(device)
		if err != nil {
			log.Printf("DeviceController: Error uninstalling device: %s", err)
			InternalServerError(w, errors.New("failed to uninstall device"))
			return
		}

		// Створення DTO і повернення відповіді
		deviceDto := resources.DeviceDto{}.DomainToDto(uninstalledDevice)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deviceDto)
	}
}
