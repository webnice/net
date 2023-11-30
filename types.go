package net

import (
	"net"
	"os"
	"sync"
	"sync/atomic"
)

const (
	netUdp        = "udp"
	netUdp4       = "udp4"
	netUdp6       = "udp6"
	netUnixgram   = "unixgram"
	netTcp        = "tcp"
	netTcp4       = "tcp4"
	netTcp6       = "tcp6"
	netUnix       = "unix"
	netUnixPacket = "unixpacket"
	netSocket     = "socket"
	netSystemd    = "systemd"
)

// Объект сущности, реализующий интерфейс Interface.
type impl struct {
	lck        *sync.Mutex                          // Защита от гонки.
	err        error                                // Сохранение последней ошибки.
	isRun      *atomic.Bool                         // Состояние выполнения сервера, =истина - запущен, =ложь - остановлен.
	handler    HandlerFn                            // Основная функция TCP сервера.
	handlerUdp HandlerUdpFn                         // Основная функция UDP сервера.
	listener   *netListener                         // Слушатель сокета сервера содержащий либо UDP либо TCP соединение.
	isShutdown *atomic.Bool                         // Флаг начала завершения работы сервера.
	onShutdown chan struct{}                        // Канал передачи сигнала об окончании завершения работы сервера.
	conf       *Configuration                       // Конфигурация сервера.
	fnFl       func(*os.File) (net.Listener, error) // Функция net.FileListener, подменяемая при тестировании.
	fnNf       func(uintptr, string) *os.File       // Функция os.NewFile, подменяемая при тестировании.
	fnFc       func(*os.File) error                 // Функция закрытия файлового дескриптора, подменяемая при тестировании.
}

// HandlerFn Описание типа функции TCP или сокет сервера.
type HandlerFn func(net.Listener) error

// HandlerUdpFn Описание типа функции UDP или сервера пакетов.
type HandlerUdpFn func(net.PacketConn) error

// Значения переменных systemd, загружаемые из окружения.
type listenEnv struct {
	pid   int      // ID процесса, получающего открытые соединения через файловые дескрипторы.
	fds   int      // Количество передаваемых соединений.
	names []string // Имена файловых дескрипторов.
}
