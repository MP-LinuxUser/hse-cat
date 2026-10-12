package console

import (
	"bufio"
	"fmt"
	"hse-cat/internal/domain"
	"hse-cat/internal/service"
	"io"
	"strconv"
	"strings"
)

type App struct {
	fleet *service.Fleet
	in    *bufio.Scanner
	out   io.Writer
	done  bool
}

func NewApp(fleet *service.Fleet, in io.Reader, out io.Writer) *App {
	return &App{fleet: fleet, in: bufio.NewScanner(in), out: out}
}

func (a *App) printMenu() {
	fmt.Fprintln(a.out, `
=== ВышКат: учёт парка ===
1. Добавить самокат
2. Добавить велосипед
3. Добавить электровелосипед
4. Добавить шлем
5. Добавить зарядный шкаф
6. отчёт
7. для новичков
8. инвентарные номера
9. демо
0. Выход`)
}

func (a *App) Run() {
	for !a.done {
		a.printMenu()
		switch a.readLine("Выберите пункт: ") {
		case "1":
			a.addScooter()
		case "2":
			a.addBicycle()
		case "3":
			a.addElecBicycle()
		case "4":
			a.addHelmet()
		case "5":
			a.addPowerContainer()
		case "6":
			a.printReport()
		case "7":
			a.printBeginners()
		case "8":
			a.printInventory()
		case "9":
			a.runDemo()
		case "0":
			fmt.Fprintln(a.out, "Конец программы")
			return
		case "":
		default:
			fmt.Fprintln(a.out, "Нет такого пункта")
		}
	}
}

func (a *App) readLine(prompt string) string {
	fmt.Fprint(a.out, prompt)
	if !a.in.Scan() {
		a.done = true
		return ""
	}
	return strings.TrimSpace(a.in.Text())
}

func (a *App) readInt(prompt string) (int, bool) {
	s := a.readLine(prompt)
	n, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(a.out, "Ожидалось целое число, получено %q\n", s)
		return 0, false
	}
	return n, true
}

func (a *App) readFloat(prompt string) (float64, bool) {
	s := a.readLine(prompt)
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Fprintf(a.out, "Ожидалось float число, получено %q\n", s)
		return 0.0, false
	}
	return n, true
}

func (a *App) add(item domain.Item) {
	if err := a.fleet.Add(item); err != nil {
		fmt.Fprintln(a.out, "Отклонено:", err)
		return
	}
	fmt.Fprintln(a.out, "Принято:", item.Name())
}

func (a *App) addScooter() {
	name := a.readLine("Название: ")
	id := a.readLine("ID: ")
	watts, ok := a.readFloat("Мощность, Вт: ")
	if !ok {
		return
	}
	hours, ok := a.readFloat("Часы работы в сутки: ")
	if !ok {
		return
	}
	simplicity, ok := a.readInt("Простота для новичка (1-10): ")
	if !ok {
		return
	}
	wear, ok := a.readInt("Износ, %: ")
	if !ok {
		return
	}
	a.add(domain.NewScooter(name, id, watts, hours, simplicity, wear))
}

func (a *App) addBicycle() {
	name := a.readLine("Название: ")
	id := a.readLine("ID: ")
	simplicity, ok := a.readInt("Простота для новичка (1-10): ")
	if !ok {
		return
	}
	wear, ok := a.readInt("Износ, %: ")
	if !ok {
		return
	}
	a.add(domain.NewBicycle(name, id, simplicity, wear))
}

func (a *App) addElecBicycle() {
	name := a.readLine("Название: ")
	id := a.readLine("ID: ")
	watts, ok := a.readFloat("Мощность, Вт: ")
	if !ok {
		return
	}
	hours, ok := a.readFloat("Часы работы в сутки: ")
	if !ok {
		return
	}
	simplicity, ok := a.readInt("Простота для новичка (1-10): ")
	if !ok {
		return
	}
	wear, ok := a.readInt("Износ, %: ")
	if !ok {
		return
	}
	a.add(domain.NewElecBicycle(name, id, watts, hours, simplicity, wear))
}

func (a *App) addHelmet() {
	name := a.readLine("Название: ")
	id := a.readLine("ID: ")
	a.add(domain.NewHelmet(name, id))
}

func (a *App) addPowerContainer() {
	name := a.readLine("Название: ")
	id := a.readLine("ID: ")
	watts, ok := a.readFloat("Мощность, Вт: ")
	if !ok {
		return
	}
	hours, ok := a.readFloat("Часы работы в сутки: ")
	if !ok {
		return
	}
	a.add(domain.NewPowerContainer(name, id, watts, hours))
}

func (a *App) printReport() {
	total := a.fleet.Count()
	vehicles := a.fleet.CountVehicles()
	fmt.Fprintln(a.out, "\n--- Отчёт по парку ---")
	fmt.Fprintf(a.out, "Всего единиц:     %d\n", total)
	fmt.Fprintf(a.out, "  транспорт:      %d\n", vehicles)
	fmt.Fprintf(a.out, "  вещи:           %d\n", total-vehicles)
	fmt.Fprintf(a.out, "Энергия в сутки:  %.2f кВт·ч\n", a.fleet.TotalDailyKWh())
}

func (a *App) printBeginners() {
	list := a.fleet.BeginnerVehicles()
	fmt.Fprintln(a.out, "\n--- Подходит новичкам ---")
	if len(list) == 0 {
		fmt.Fprintln(a.out, "(нет подходящей техники)")
		return
	}
	for _, v := range list {
		fmt.Fprintf(a.out, "%-20s %-8s простота %d\n", v.Name(), v.ItemID(), v.Simplicity())
	}
}

func (a *App) printInventory() {
	items := a.fleet.Items()
	fmt.Fprintln(a.out, "\n--- Инвентарь ---")
	if len(items) == 0 {
		fmt.Fprintln(a.out, "(парк пуст)")
		return
	}
	for _, it := range items {
		fmt.Fprintf(a.out, "%-20s %s\n", it.Name(), it.ItemID())
	}
}

func (a *App) runDemo() {
	fmt.Fprintln(a.out, "\n=== Демо-сценарий ===")
	demo := []domain.Item{
		domain.NewScooter("Самокат City", "VK-001", 350, 6, 7, 10),
		domain.NewScooter("Самокат Sport", "VK-002", 500, 5, 4, 30), // не для новичков
		domain.NewBicycle("Велосипед Classic", "VK-003", 9, 5),
		domain.NewElecBicycle("Электровелосипед", "VK-004", 250, 4, 6, 85), // износ: отклонят
		domain.NewHelmet("Шлем", "VK-005"),
		domain.NewPowerContainer("Зарядный шкаф", "VK-006", 1200, 10),
		domain.NewBicycle("Велосипед Old", "VK-003", 8, 0), // дубликат номера
		domain.NewScooter("", "VK-007", 300, 5, 7, 0),      // без названия
	}
	for _, item := range demo {
		a.add(item)
	}
	a.printReport()
	a.printBeginners()
	a.printInventory()
}
