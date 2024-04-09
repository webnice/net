package net

import (
	"crypto/tls"
	"net"
)

// Interface Интерфейс пакета.
type Interface interface {
	// ФУНКЦИИ СЕРВЕРА

	// ID Уникальный идентификатор сервера.
	// Если идентификатор не был указан в конфигурации, создаётся при запуске.
	ID() string

	// Handler Назначение основной функции TCP или сокет сервера. Функция должна назначаться до запуска сервера.
	Handler(fn HandlerFn) Interface

	// HandlerUdp Назначение основной функции UDP сервера. Функция должна назначаться до запуска сервера.
	HandlerUdp(fn HandlerUdpFn) Interface

	// ПРОСЛУШИВАНИЕ СЕТЕВОГО СОЕДИНЕНИЯ

	// ListenAndServe Открытие адреса или сокета без использования конфигурации сервера (конфигурация по
	// умолчанию), после успешного открытия адреса, выполняется запуск сервера для обслуживания входящих соединений.
	ListenAndServe(addr string) Interface

	// ListenAndServeTLS Открытие адреса или сокета с использованием TLS, без использования конфигурации сервера
	// (конфигурация по умолчанию), после успешного открытия адреса, выполняется запуск сервера для обслуживания
	// входящих соединений.
	ListenAndServeTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) Interface

	// ListenAndServeWithConfig Настройка сервера с использованием переданной конфигурации, открытие адреса или сокета
	// на прослушивание, после успешного открытия адреса, выполняется запуск сервера для
	// обслуживания входящих соединений.
	ListenAndServeWithConfig(conf *Configuration) Interface

	// ListenAndServeTLSWithConfig Настройка сервера с использованием переданной конфигурации в режиме TLS, открытие
	// адреса или сокета на прослушивание, после успешного открытия адреса, выполняется запуск сервера для
	// обслуживания входящих соединений.
	ListenAndServeTLSWithConfig(conf *Configuration, tlsConfig *tls.Config) Interface

	// ListenersSystemdWithoutNames Возвращает срез net.Listener сокетов переданных в процесс сервера
	// из службы linux - systemd.
	ListenersSystemdWithoutNames() (ret []net.Listener, err error)

	// ListenersSystemdWithNames Возвращает карту срезов net.Listener сокетов переданных в процесс сервера
	// из службы linux - systemd.
	ListenersSystemdWithNames() (ret map[string][]net.Listener, err error)

	// ListenersSystemdTLSWithoutNames Возвращает срез net.nnlistener для TLS сокетов переданных в процесс сервера
	// из службы linux - systemd.
	ListenersSystemdTLSWithoutNames(tlsConfig *tls.Config) (ret []net.Listener, err error)

	// ListenersSystemdTLSWithNames Возвращает карту срезов net.listener для TLS сокетов переданных в процесс сервера
	// из службы linux - systemd.
	ListenersSystemdTLSWithNames(tlsConfig *tls.Config) (ret map[string][]net.Listener, err error)

	// NewListener Создание нового слушателя соединений net.Listener на основе конфигурации сервера.
	NewListener(conf *Configuration) (ret net.Listener, rpc net.PacketConn, err error)

	// NewListenerTLS Создание нового слушателя соединений net.Listener в режиме TLS, на основе конфигурации
	// сервера.
	NewListenerTLS(conf *Configuration, tlsConfig *tls.Config) (ret net.Listener, rpc net.PacketConn, err error)

	// NewTLSConfigDefault Создание TLS конфигурации по умолчанию, на основе секретного и публичного ключей.
	NewTLSConfigDefault(tlsPublicFile string, tlsPrivateFile string) (ret *tls.Config, err error)

	// СЕРВЕР

	// Serve Запуск функции сервера для входящих соединений на основе переданного слушателя net.Listener.
	Serve(net.Listener) Interface

	// ServeWithId Запуск функции сервера для входящих соединений на основе переданного слушателя net.Listener с
	// указанием ID сервера.
	ServeWithId(ltn net.Listener, id string) Interface

	// ServeUdp Запуск функции сервера для входящих UDP пакетов на основе переданного слушателя net.PacketConn.
	ServeUdp(lpc net.PacketConn) Interface

	// ServeUdpWithId Запуск функции сервера для входящих UDP пакетов на основе переданного слушателя net.PacketConn с
	// указанием ID сервера.
	ServeUdpWithId(lpc net.PacketConn, id string) Interface

	// Wait Блокируемая функция ожидания завершения веб сервера, если он запущен.
	// Если сервер не запущен, функция завершается немедленно.
	Wait() Interface

	// Stop Завершение работы сервера/функции сервера.
	Stop() Interface

	// IsRunning Статус выполнения сервера.
	// Вернётся истина, если сервер запущен.
	IsRunning() (ret bool)

	// ОШИБКИ

	// Clean Очистка последней ошибки.
	Clean() Interface

	// Errors Справочник ошибок.
	Errors() *Error

	// Error Последняя ошибка сервера.
	Error() error
}
