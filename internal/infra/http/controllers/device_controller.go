package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
)

type DeviceController struct {
	DeviceService app.DeviceService
	RoomService   app.RoomService
}

func NewDeviceController(ds app.DeviceService, rs app.RoomService) DeviceController {
	return DeviceController{
		DeviceService: ds,
		RoomService:   rs,
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

		if deviceRequest.RoomId == nil {
			err := errors.New("roomId is required")
			log.Printf("DeviceController: %s", err)
			BadRequest(w, err)
			return
		}

		_, err = c.RoomService.Find(*deviceRequest.RoomId)
		if err != nil {
			log.Printf("DeviceController: Error finding room: %s", err)
			BadRequest(w, errors.New("room not found"))
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

func (c *DeviceController) FindByRoomId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomIdParam := chi.URLParam(r, "roomId")
		roomId, err := strconv.ParseUint(roomIdParam, 10, 64)
		if err != nil {
			log.Printf("DeviceController: Invalid room ID: %s", err)
			BadRequest(w, errors.New("invalid room ID"))
			return
		}

		log.Printf("DeviceController: Looking for devices with room ID: %d", roomId)

		devices, err := c.DeviceService.FindByRoomId(roomId)
		if err != nil {
			log.Printf("DeviceController: Error finding devices: %s", err)
			InternalServerError(w, errors.New("failed to retrieve devices"))
			return
		}

		if len(devices) == 0 {
			log.Printf("DeviceController: No devices found for room ID: %d", roomId)
		}

		devicesDto := resources.DevicesDto{}.DomainToDto(devices)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(devicesDto)
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
		device := r.Context().Value(DevKey).(domain.Device)

		device.RoomId = nil
		updatedDevice, err := c.DeviceService.Update(device)
		if err != nil {
			log.Printf("DeviceController: Error uninstalling device: %s", err)
			InternalServerError(w, errors.New("failed to uninstall device"))
			return
		}

		deviceDto := resources.DeviceDto{}.DomainToDto(updatedDevice)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deviceDto)
	}
}
