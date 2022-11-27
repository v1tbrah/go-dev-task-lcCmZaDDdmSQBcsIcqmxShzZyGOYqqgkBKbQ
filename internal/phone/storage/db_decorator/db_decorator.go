package db_decorator

import (
	"fmt"
	"strconv"
	"sync"

	dberr "go-dev-task-lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ/internal/phone/storage/error"

	"go-dev-task-lcCmZaDDdmSQBcsIcqmxShzZyGOYqqgkBKbQ/internal/phone/model"
)

type DBDecorator struct {
	phones map[int64]string
	mu     sync.RWMutex
}

func New() *DBDecorator {

	newDBDecorator := &DBDecorator{}

	newDBDecorator.mu = sync.RWMutex{}

	newPhones := newPhones()

	newDBDecorator.phones = newPhones

	return newDBDecorator

}

func newPhones() map[int64]string {

	newPhones := make(map[int64]string, 100)

	var currEndPhoneNumber uint8

	for i := 1; i < 100; i++ {
		if currEndPhoneNumber < 10 {
			newPhones[int64(i)] = "8-800-555-35-3" + strconv.Itoa(int(currEndPhoneNumber))
		} else {
			newPhones[int64(i)] = "8-800-555-35-" + strconv.Itoa(int(currEndPhoneNumber))
		}
		currEndPhoneNumber++
	}

	return newPhones

}

func (d *DBDecorator) GetPhone(id int64) (phone *model.Phone, err error) {

	d.mu.RLock()
	defer d.mu.RUnlock()

	number, isFound := d.phones[id]
	if !isFound {
		return nil, fmt.Errorf(`db_decorator: %w`, dberr.ErrPhoneIsNotFound)
	}

	phone = &model.Phone{
		ID:     id,
		Number: number,
	}

	return phone, nil

}
