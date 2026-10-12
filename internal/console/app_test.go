package console

import (
	"bytes"
	"strings"
	"testing"

	"hse-cat/internal/service"
)

// runApp запускает консоль с заранее заготовленным вводом и возвращает весь вывод.
func runApp(t *testing.T, input string) string {
	t.Helper() // если тест упадёт, ошибка укажет на вызывающую строку, а не сюда
	fleet := service.NewFleet(service.NewBasicServiceCenter(70))
	var out bytes.Buffer
	NewApp(fleet, strings.NewReader(input), &out).Run()
	return out.String()
}

func assertContains(t *testing.T, output string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(output, w) {
			t.Errorf("в выводе нет %q", w)
		}
	}
}

func TestApp_Demo(t *testing.T) {
	out := runApp(t, "9\n0\n")
	assertContains(t, out, "Принято:", "Отклонено:", "16.60", "VK-006")
}

func TestApp_AddScooterManually(t *testing.T) {
	// 1 → название → номер → мощность → часы → простота → износ, потом 8 (инвентарь), 0 (выход)
	out := runApp(t, "1\nСамокат\nVK-1\n350\n6\n7\n10\n8\n0\n")
	assertContains(t, out, "Принято: Самокат", "VK-1")
}

func TestApp_RejectedByWear(t *testing.T) {
	out := runApp(t, "2\nВелосипед\nVK-1\n5\n90\n0\n")
	assertContains(t, out, "Отклонено:")
}

func TestApp_DotInFloat(t *testing.T) {
	out := runApp(t, "5\nШкаф\nVK-9\n1000.5\n2\n0\n")
	assertContains(t, out, "Принято: Шкаф")
}

func TestApp_UnknownMenuItem(t *testing.T) {
	out := runApp(t, "42\n0\n")
	assertContains(t, out, "Нет такого пункта")
}

// Ввод закончился без "0" — программа не должна зависнуть.
func TestApp_EndOfInput(t *testing.T) {
	out := runApp(t, "6\n")
	assertContains(t, out, "Всего единиц")
}
