package storage

import (
	"go-dev-task-lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ/internal/phone/model"
	"go-dev-task-lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ/internal/phone/storage/db_decorator"
)

type Storage interface {
	GetPhone(id int64) (phone *model.Phone, err error)
}

func New() Storage {
	return db_decorator.New()
}
