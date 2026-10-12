package domain

// helmet

type Helmet struct {
	baseItem
}

func NewHelmet(name, itemID string) *Helmet {
	return &Helmet{baseItem: baseItem{name, itemID}}
}

// powerContainer

type PowerContainer struct {
	baseItem
	energyUser
}

func NewPowerContainer(name, itemID string, watts, hours float64) *PowerContainer {
	return &PowerContainer{
		baseItem:   baseItem{name, itemID},
		energyUser: energyUser{watts, hours},
	}
}
