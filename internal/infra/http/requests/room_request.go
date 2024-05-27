package requests

import "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"

type RoomRequest struct {
	OrganizationId uint64 `json:"organization_id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

func (r RoomRequest) ToDomainModel() (interface{}, error) {
	return domain.Room{
		OrganizationId: r.OrganizationId,
		Name:           r.Name,
		Description:    r.Description,
	}, nil
}
