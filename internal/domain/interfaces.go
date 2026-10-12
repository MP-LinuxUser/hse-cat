package domain

type PowerConsumer interface {
	DailyKWh() float64
}

type Vehicle interface {
	Item
	Simplicity() int
	Wear() int
}

type Item interface {
	Name() string
	ItemID() string
}

type baseItem struct {
	name   string
	itemID string
}

func (b baseItem) Name() string   { return b.name }
func (b baseItem) ItemID() string { return b.itemID }

// транспорт

type baseVehicle struct {
	baseItem
	simplicity int
	wear       int
}

func (b baseVehicle) Simplicity() int { return b.simplicity }
func (b baseVehicle) Wear() int       { return b.wear }

// потребляет энергию

type energyUser struct {
	watts float64
	hours float64
}

func (e energyUser) DailyKWh() float64 {
	return e.watts * e.hours / 1000
}
