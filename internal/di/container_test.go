package di

import (
	"errors"
	"io"
	"strings"
	"testing"

	"hse-cat/internal/domain"
)

// Техосмотр, который всем отказывает
type rejectAll struct{}

func (rejectAll) Inspect(domain.Item) error { return errors.New("отказ") }

func newTestContainer() *Container {
	return NewContainer(strings.NewReader(""), io.Discard)
}

func TestContainer_SubstituteServiceCenter(t *testing.T) {
	c := newTestContainer()
	c.serviceCenter = rejectAll{} // подменяем ДО того, как создан парк

	err := c.Fleet().Add(domain.NewHelmet("Шлем", "VK-1"))

	if err == nil {
		t.Error("парк должен был использовать подменённый техосмотр и отказать")
	}
}
func TestContainer_DefaultServiceCenter(t *testing.T) {
	c := newTestContainer()
	// настоящий техосмотр не пропустит самокат с износом 90%
	err := c.Fleet().Add(domain.NewScooter("S", "VK-1", 350, 6, 7, 90))
	if err == nil {
		t.Error("по умолчанию должен работать настоящий техосмотр")
	}
}
