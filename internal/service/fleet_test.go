package service

import (
	"errors"
	"math"
	"testing"

	"hse-cat/internal/domain"
)

// ===== Подставной техосмотр =====

// fakeServiceCenter реализует интерфейс ServiceCenter, но без настоящих правил:
// он просто возвращает то, что ему сказали, и считает вызовы.
type fakeServiceCenter struct {
	err   error // что вернуть из Inspect: nil = принять, ошибка = отклонить
	calls int   // сколько раз Fleet к нему обратился
}

func (f *fakeServiceCenter) Inspect(item domain.Item) error {
	f.calls++
	return f.err
}

var errRejected = errors.New("отклонено фейковым техосмотром")

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

// ===== Приём техники =====

func TestFleet_Add_RejectedByServiceCenter(t *testing.T) {
	fake := &fakeServiceCenter{err: errRejected}
	fleet := NewFleet(fake) // ← подмена: вместо BasicServiceCenter передаём фейк

	err := fleet.Add(domain.NewBicycle("B", "VK-1", 5, 0))

	if !errors.Is(err, errRejected) {
		t.Fatalf("ожидали ошибку техосмотра, получили %v", err)
	}
	if fleet.Count() != 0 {
		t.Errorf("отклонённое устройство попало в парк: Count() = %d", fleet.Count())
	}
	if fake.calls != 1 {
		t.Errorf("техосмотр вызван %d раз, ожидали 1", fake.calls)
	}
}

func TestFleet_Add_AcceptedByServiceCenter(t *testing.T) {
	fleet := NewFleet(&fakeServiceCenter{}) // err == nil → всё принимает

	// Износ 99% — настоящий техосмотр отклонил бы. Но решение принимает
	// ServiceCenter, а не Fleet, поэтому с фейком устройство принимается.
	err := fleet.Add(domain.NewScooter("S", "VK-1", 350, 6, 7, 99))

	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if fleet.Count() != 1 {
		t.Errorf("Count() = %d, ожидали 1", fleet.Count())
	}
}

func TestFleet_Add_DuplicateID(t *testing.T) {
	fleet := NewFleet(&fakeServiceCenter{})

	if err := fleet.Add(domain.NewBicycle("Первый", "VK-1", 5, 0)); err != nil {
		t.Fatalf("первое добавление: %v", err)
	}
	err := fleet.Add(domain.NewHelmet("Второй", "VK-1"))

	if err == nil {
		t.Fatal("дубликат номера должен быть отклонён")
	}
	if fleet.Count() != 1 {
		t.Errorf("Count() = %d, ожидали 1", fleet.Count())
	}
}

// ===== Отчёты =====

func TestFleet_Counts(t *testing.T) {
	fleet := NewFleet(&fakeServiceCenter{})
	fleet.Add(domain.NewScooter("S", "VK-1", 350, 6, 7, 0))
	fleet.Add(domain.NewBicycle("B", "VK-2", 5, 0))
	fleet.Add(domain.NewHelmet("H", "VK-3"))

	if fleet.Count() != 3 {
		t.Errorf("Count() = %d, ожидали 3", fleet.Count())
	}
	if fleet.CountVehicles() != 2 {
		t.Errorf("CountVehicles() = %d, ожидали 2", fleet.CountVehicles())
	}
}

func TestFleet_TotalDailyKWh(t *testing.T) {
	fleet := NewFleet(&fakeServiceCenter{})
	fleet.Add(domain.NewScooter("S", "VK-1", 1000, 1, 7, 0))  // 1 кВт·ч
	fleet.Add(domain.NewPowerContainer("P", "VK-2", 2000, 1)) // 2 кВт·ч
	fleet.Add(domain.NewBicycle("B", "VK-3", 5, 0))           // не потребляет
	fleet.Add(domain.NewHelmet("H", "VK-4"))                  // не потребляет

	if got := fleet.TotalDailyKWh(); !almostEqual(got, 3.0) {
		t.Errorf("TotalDailyKWh() = %v, ожидали 3.0", got)
	}
}

func TestFleet_TotalDailyKWh_EmptyFleet(t *testing.T) {
	fleet := NewFleet(&fakeServiceCenter{})
	if got := fleet.TotalDailyKWh(); got != 0 {
		t.Errorf("пустой парк: TotalDailyKWh() = %v, ожидали 0", got)
	}
}

// Правило «для новичков»: простота >= 6. Проверяем границу 5/6.
func TestFleet_BeginnerVehicles(t *testing.T) {
	fleet := NewFleet(&fakeServiceCenter{})
	fleet.Add(domain.NewBicycle("B", "VK-5", 5, 0))           // 5 — нет
	fleet.Add(domain.NewBicycle("B", "VK-6", 6, 0))           // 6 — да (граница)
	fleet.Add(domain.NewScooter("S", "VK-10", 350, 6, 10, 0)) // 10 — да
	fleet.Add(domain.NewHelmet("H", "VK-H"))                  // не транспорт — нет

	got := fleet.BeginnerVehicles()

	want := []string{"VK-6", "VK-10"}
	if len(got) != len(want) {
		t.Fatalf("получили %d устройств, ожидали %d", len(got), len(want))
	}
	for i, v := range got {
		if v.ItemID() != want[i] {
			t.Errorf("[%d] = %s, ожидали %s", i, v.ItemID(), want[i])
		}
	}
}

func TestFleet_Items_ReturnsCopy(t *testing.T) {
	fleet := NewFleet(&fakeServiceCenter{})
	fleet.Add(domain.NewHelmet("H", "VK-1"))

	items := fleet.Items()
	items[0] = nil // портим копию

	if fleet.Items()[0] == nil {
		t.Error("Items() вернул внутренний слайс, а не копию")
	}
}
