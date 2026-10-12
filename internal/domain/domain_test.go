package domain

import (
	"math"
	"testing"
)

// Дробные числа нельзя сравнивать через ==: 0.1+0.2 != 0.3 в компьютере.
// Сравниваем с допуском.
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestScooter_Getters(t *testing.T) {
	s := NewScooter("Самокат", "VK-1", 350, 6, 7, 10)

	if s.Name() != "Самокат" {
		t.Errorf("Name() = %q, want %q", s.Name(), "Самокат")
	}
	if s.ItemID() != "VK-1" {
		t.Errorf("ItemID() = %q, want %q", s.ItemID(), "VK-1")
	}
	if s.Simplicity() != 7 {
		t.Errorf("Simplicity() = %d, want 7", s.Simplicity())
	}
	if s.Wear() != 10 {
		t.Errorf("Wear() = %d, want 10", s.Wear())
	}
}

func TestDailyKWh(t *testing.T) {
	tests := []struct {
		name     string
		consumer PowerConsumer
		want     float64
	}{
		{"самокат 350 Вт × 6 ч", NewScooter("S", "VK-1", 350, 6, 7, 0), 2.1},
		{"электровелосипед 250 Вт × 4 ч", NewElecBicycle("E", "VK-2", 250, 4, 6, 0), 1.0},
		{"зарядный шкаф 1200 Вт × 10 ч", NewPowerContainer("P", "VK-3", 1200, 10), 12.0},
		{"ноль часов работы", NewScooter("S", "VK-4", 350, 0, 7, 0), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.consumer.DailyKWh()
			if !almostEqual(got, tt.want) {
				t.Errorf("DailyKWh() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Фиксирует модель: кто транспорт, кто потребляет энергию.
// Если кто-то случайно добавит шлему DailyKWh() — тест упадёт.
func TestCapabilities(t *testing.T) {
	tests := []struct {
		name          string
		item          Item
		isVehicle     bool
		consumesPower bool
	}{
		{"самокат", NewScooter("S", "VK-1", 350, 6, 7, 0), true, true},
		{"велосипед", NewBicycle("B", "VK-2", 9, 0), true, false},
		{"электровелосипед", NewElecBicycle("E", "VK-3", 250, 4, 6, 0), true, true},
		{"шлем", NewHelmet("H", "VK-4"), false, false},
		{"зарядный шкаф", NewPowerContainer("P", "VK-5", 1200, 10), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, isVehicle := tt.item.(Vehicle)
			_, consumesPower := tt.item.(PowerConsumer)

			if isVehicle != tt.isVehicle {
				t.Errorf("Vehicle: got %v, want %v", isVehicle, tt.isVehicle)
			}
			if consumesPower != tt.consumesPower {
				t.Errorf("PowerConsumer: got %v, want %v", consumesPower, tt.consumesPower)
			}
		})
	}
}
