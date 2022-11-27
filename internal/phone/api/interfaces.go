package api

import (
	"go-dev-task-lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ/internal/phone/model"
)

type Storage interface {
	GetPhone(id int64) (phone *model.Phone, err error)
}
