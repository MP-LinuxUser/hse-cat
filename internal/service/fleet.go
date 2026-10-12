package service

import (
	"fmt"
	"hse-cat/internal/domain"
	"slices"
)

type Fleet struct {
	inspector ServiceCenter
	items     []domain.Item
}

func NewFleet(inspector ServiceCenter) *Fleet {
	return &Fleet{inspector: inspector}
}

func (f *Fleet) Items() []domain.Item {
	return slices.Clone(f.items)
}

func (f *Fleet) Add(item domain.Item) error {
	if err := f.inspector.Inspect(item); err != nil {
		return fmt.Errorf("object %q failed the check: %w", item.ItemID(), err)
	}
	if f.has(item.ItemID()) {
		return fmt.Errorf("obect %q already in the park", item.ItemID())
	}
	f.items = append(f.items, item)
	return nil
}

func (f *Fleet) has(id string) bool {
	for _, it := range f.items {
		if it.ItemID() == id {
			return true
		}
	}
	return false
}

func (f *Fleet) CountVehicles() int {
	cnt := 0
	for _, it := range f.items {
		if _, ok := it.(domain.Vehicle); ok {
			cnt += 1
		}
	}
	return cnt
}

func (f *Fleet) Count() int {
	return len(f.items)
}

func (f *Fleet) TotalDailyKWh() float64 {
	sum := 0.0
	for _, item := range f.items {
		if pc, ok := item.(domain.PowerConsumer); ok {
			sum += pc.DailyKWh()
		}
	}
	return sum
}

const beginnerMinSimplicity = 6

func (f *Fleet) BeginnerVehicles() []domain.Vehicle {
	var res []domain.Vehicle
	for _, item := range f.items {
		if v, ok := item.(domain.Vehicle); ok && v.Simplicity() >= beginnerMinSimplicity {
			res = append(res, v)
		}
	}
	return res
}
