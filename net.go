package net

import (
	"net"
	"os"
	"sync"
	"sync/atomic"
)

// New Конструктор объекта сущности пакета, возвращается интерфейс пакета.
func New() Interface {
	var nut = &impl{
		lck:        new(sync.Mutex),
		isRun:      new(atomic.Bool),
		isShutdown: new(atomic.Bool),
		fnFl:       net.FileListener,
		fnNf:       os.NewFile,
		fnFc:       fileClose,
	}

	nut.isRun.Store(false)
	nut.isShutdown.Store(false)

	return nut
}

// Handler Назначение основной функции TCP или сокет сервера. Функция должна назначаться до запуска сервера.
func (nut *impl) Handler(fn HandlerFn) Interface { nut.handler = fn; return nut }

// HandlerUdp Назначение основной функции UDP сервера. Функция должна назначаться до запуска сервера.
func (nut *impl) HandlerUdp(fn HandlerUdpFn) Interface { nut.handlerUdp = fn; return nut }

// Clean Очистка последней ошибки.
func (nut *impl) Clean() Interface { nut.err = nil; return nut }

// Error Последняя ошибка сервера.
func (nut *impl) Error() error { return nut.err }

// Errors Справочник ошибок.
func (nut *impl) Errors() *Error { return errSingleton }

// Wait Блокируемая функция ожидания завершения веб сервера, если он запущен.
// Если сервер не запущен, функция завершается немедленно.
func (nut *impl) Wait() Interface {
	if !nut.isRun.Load() {
		//nut.err = ErrNotRunning()
		return nut
	}
	safeWait(nut.onShutdown)
	nut.onShutdown = nil

	return nut
}

// Stop Завершение работы сервера/функции сервера.
func (nut *impl) Stop() Interface {
	// Защита от возможной смертельной блокировки при остановке сервера из разных потоков.
	nut.lck.Lock()
	defer nut.lck.Unlock()
	// Выход, если сервер не запущен или уже начато завершение работы сервера.
	if !nut.isRun.Load() || nut.isShutdown.Load() {
		return nut
	}
	// Флаг начала завершения работы сервера.
	nut.isShutdown.Store(true)
	// Закрытие соединения.
	_ = nut.listener.Close()
	safeClose(nut.onShutdown)
	nut.onShutdown = nil

	return nut
}

// IsRunning Статус выполнения сервера.
// Вернётся истина, если сервер запущен.
func (nut *impl) IsRunning() (ret bool) {
	// Защита от возможной смертельной блокировки при остановке сервера из разных потоков.
	nut.lck.Lock()
	defer nut.lck.Unlock()
	// Выход, если сервер запущен или начато завершение работы сервера.
	if nut.isRun.Load() || nut.isShutdown.Load() {
		ret = true
	}

	return
}
