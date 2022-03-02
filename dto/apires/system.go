package apires

import (
	"casbin-auth-go/dto/model"
	"gopkg.in/guregu/null.v4"
)

type ListSystem struct {
	List        []*model.System `json:"list"`
	Total       int             `json:"total"`
	CurrentPage int             `json:"current_page"`
	PerPage     int             `json:"per_page"`
	NextPage    null.Int        `json:"next_page" swaggertype:"string"`
}
