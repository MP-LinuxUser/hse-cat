package domain

// scooter
type Scooter struct {
	baseVehicle
	energyUser
}

func NewScooter(name, itemID string, watts, hours float64, simplicity, wear int) *Scooter {
	return &Scooter{
		baseVehicle: baseVehicle{baseItem{name, itemID}, simplicity, wear},
		energyUser:  energyUser{watts, hours},
	}
}

// bicycle

type Bicycle struct {
	baseVehicle
}

func NewBicycle(name, itemID string, simplicity, wear int) *Bicycle {
	return &Bicycle{
		baseVehicle: baseVehicle{baseItem{name, itemID}, simplicity, wear},
	}
}

// electro-bicycle

type ElecBicycle struct {
	baseVehicle
	energyUser
}

func NewElecBicycle(name, itemID string, watts, hours float64, simplicity, wear int) *ElecBicycle {
	return &ElecBicycle{
		baseVehicle: baseVehicle{baseItem{name, itemID}, simplicity, wear},
		energyUser:  energyUser{watts, hours},
	}
}
