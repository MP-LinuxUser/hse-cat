package service

import (
	"errors"
	"fmt"
	"hse-cat/internal/domain"
)

type ServiceCenter interface {
	Inspect(v domain.Item) error
}

type BasicServiceCenter struct {
	maxWear int
}

func NewBasicServiceCenter(maxWear int) *BasicServiceCenter {
	return &BasicServiceCenter{maxWear: maxWear}
}

func (sc *BasicServiceCenter) Inspect(item domain.Item) error {
	if item.Name() == "" {
		return errors.New("name missing")
	}
	if item.ItemID() == "" {
		return errors.New("id missing")
	}
	if pc, ok := item.(domain.PowerConsumer); ok && pc.DailyKWh() < 0 {
		return errors.New("negative value of power consumer")
	}
	if v, ok := item.(domain.Vehicle); ok {
		return sc.inspectVehicle(v)
	}
	return nil
}

func (sc *BasicServiceCenter) inspectVehicle(v domain.Vehicle) error {
	if s := v.Simplicity(); s > 10 || s < 1 {
		return fmt.Errorf("simplicity %d out of range 1-10", s)
	}
	if w := v.Wear(); w > 100 || w < 0 {
		return fmt.Errorf("wear %d out of range 0-100", w)
	}
	if v.Wear() > sc.maxWear {
		return fmt.Errorf("wear %d exceeds the permissible limit %d", v.Wear(), sc.maxWear)
	}
	return nil
}
