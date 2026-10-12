package service

import (
	"testing"

	"hse-cat/internal/domain"
)

func TestBasicServiceCenter_Inspect(t *testing.T) {
	sc := NewBasicServiceCenter(70) // порог износа 70%

	tests := []struct {
		name    string
		item    domain.Item
		wantErr bool // true = должен отклонить
	}{
		// транспорт: износ
		{"новый самокат (износ 0)", domain.NewScooter("S", "VK-1", 350, 6, 7, 0), false},
		{"износ ровно на пороге", domain.NewScooter("S", "VK-1", 350, 6, 7, 70), false},
		{"износ выше порога", domain.NewScooter("S", "VK-1", 350, 6, 7, 71), true},
		{"отрицательный износ", domain.NewBicycle("B", "VK-1", 5, -1), true},

		// транспорт: простота
		{"простота 1 (минимум)", domain.NewBicycle("B", "VK-1", 1, 0), false},
		{"простота 10 (максимум)", domain.NewBicycle("B", "VK-1", 10, 0), false},
		{"простота 0", domain.NewBicycle("B", "VK-1", 0, 0), true},
		{"простота 11", domain.NewBicycle("B", "VK-1", 11, 0), true},

		// общие правила
		{"пустое название", domain.NewScooter("", "VK-1", 350, 6, 7, 0), true},
		{"пустой инвентарный номер", domain.NewBicycle("B", "", 5, 0), true},

		// вещи
		{"шлем принят", domain.NewHelmet("H", "VK-2"), false},
		{"шлем без номера", domain.NewHelmet("H", ""), true},
		{"шкаф принят", domain.NewPowerContainer("P", "VK-3", 1200, 10), false},
		{"шкаф с отрицательной мощностью", domain.NewPowerContainer("P", "VK-3", -100, 5), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sc.Inspect(tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("Inspect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
