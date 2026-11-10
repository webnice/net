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
		fnFpc:      net.FilePacketConn,
		fnNf:       os.NewFile,
		fnFc:       fileClose,
	}

	nut.isRun.Store(false)
	nut.isShutdown.Store(false)

	return nut
}

// ID Уникальный идентификатор сервера.
// Если идентификатор не был указан в конфигурации, создаётся при запуске.
func (nut *impl) ID() string {
	if nut.conf == nil {
		return ""
	}
	return nut.conf.ID
}

// Handler Назначение основной функции TCP или сокет сервера. Функция должна назначаться до запуска сервера.
func (nut *impl) Handler(fn HandlerFn) Interface { nut.handler = fn; return nut }

// HandlerUdp Назначение основной функции UDP сервера. Функция должна назначаться до запуска сервера.
func (nut *impl) HandlerUdp(fn HandlerUdpFn) Interface { nut.handlerUdp = fn; return nut }

// Clean Очистка последней ошибки.
func (nut *impl) Clean() Interface { nut.setError(nil); return nut }

// Error Последняя ошибка сервера.
func (nut *impl) Error() error { return nut.getError() }

// Установка последней ошибки сервера.
func (nut *impl) setError(err error) {
	nut.errLck.Lock()
	defer nut.errLck.Unlock()
	nut.err = err
}

// Получение последней ошибки сервера.
func (nut *impl) getError() error {
	nut.errLck.Lock()
	defer nut.errLck.Unlock()
	return nut.err
}

// Errors Справочник ошибок.
func (nut *impl) Errors() *Error { return errSingleton }

// Wait Блокируемая функция ожидания завершения веб сервера, если он запущен.
// Если сервер не запущен, функция завершается немедленно.
func (nut *impl) Wait() Interface {
	nut.lck.Lock()
	ch, running := nut.onShutdown, nut.isRun.Load()
	nut.lck.Unlock()
	if !running {
		return nut
	}
	if ch != nil {
		// Ожидание закрытия канала потоком сервера при его завершении.
		<-ch
	}

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
	// Закрытие соединения. Ожидание завершения работы сервера выполняется функцией Wait().
	nut.setError(nut.listener.Close())

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
