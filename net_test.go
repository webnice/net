package net

import (
	"errors"
	"net"
	"reflect"
	"testing"
)

// Тестирование конструктора.
func TestNew(t *testing.T) {
	var (
		nut Interface
		nio *impl
	)

	nut = New()
	nio = nut.(*impl)
	if nut == nil {
		t.Errorf("фунция New(), вернулся nil, ожидался интерфейс")
	}
	if nio.err != nil {
		t.Errorf("фунция New(), ошибка не равна nil")
	}
	if nio.isRun.Load() {
		t.Errorf("фунция New(), isRun=%t, ожидалось: %t", nio.isRun.Load(), false)
	}
	if nio.handler != nil {
		t.Errorf("фунция New(), handler=%v, ожидался nil", nio.handler)
	}
	if nio.isShutdown.Load() {
		t.Errorf("фунция New(), isShutdown=%t, ожидалось: %t", nio.isShutdown.Load(), false)
	}
	if nio.onShutdown != nil {
		t.Errorf("фунция New(), onShutdown=%v, ожидался nil", nio.onShutdown)
	}
	if nio.conf != nil {
		t.Errorf("фунция New(), conf=%v, ожидался nil", nio.conf)
	}
}

// Тестирование назначения основной функции сервера.
func TestHandler(t *testing.T) {
	var (
		her HandlerFn
		nut Interface
		nio *impl
	)

	her = func(l net.Listener) error { return nil }
	nut = New()
	nio = nut.(*impl)
	if nio.handler != nil {
		t.Errorf("не верно создан объект Interface")
		return
	}
	nut.Handler(her)
	if nio.handler == nil {
		t.Errorf("не корректно работает функция Handler(), ожидалось назначение функции сервера")
	}
	if reflect.ValueOf(nio.handler).Pointer() != reflect.ValueOf(her).Pointer() {
		t.Errorf("не корректно работает функция Handler(), ожидалось назначение переданной функции сервера")
	}
}

// Тестирование функции возвращения последней ошибки.
func TestError(t *testing.T) {
	const _TestString = "m7SqTD9K2FEstVjD2QR9"
	var (
		err error
		nio = New().(*impl)
	)

	err = errors.New(_TestString)
	if nio.err = err; !errors.Is(err, nio.Error()) {
		t.Errorf("фунция Error(), не корректный результат")
	}
}

func TestWaitNotRun(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var (
		nut Interface
		err error
	)

	nut = New().
		Handler(getTestHandlerFn(false))
	if err = nut.
		Wait().
		Error(); err != nil {
		t.Errorf("фунция Wait(), ошибка: %v, ожидалось: %v", err, nil)
	}
}

func TestImpl_IsRunning(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var (
		err error
		ltn net.Listener
		w1  Interface
		ok  bool
	)

	if ltn, err = net.Listen("tcp", testAddress1); err != nil {
		t.Errorf("функция Listen(%q, %q), прервана ошибкой: %s", "tcp", testAddress1, err)
	}
	w1 = New().
		Handler(getTestHandlerFn(false))
	if ok = w1.IsRunning(); ok {
		t.Errorf("функция IsRunning(), вернулось: %t, ожидалось: %t", ok, false)
	}
	w1.Serve(ltn)
	defer w1.Stop()
	if w1.(*impl).conf == nil {
		t.Errorf("ошибка создания конфигурации сервера")
	}
	if ok = w1.IsRunning(); !ok {
		t.Errorf("функция IsRunning(), вернулось: %t, ожидалось: %t", ok, true)
	}
}
