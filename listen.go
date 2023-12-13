package net

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path"
)

// ListenAndServe Открытие адреса или сокета без использования конфигурации сервера (конфигурация по
// умолчанию), после успешного открытия адреса, выполняется запуск сервера для обслуживания входящих соединений.
func (nut *impl) ListenAndServe(addr string) Interface {
	var conf *Configuration

	if conf, nut.err = parseAddress(addr, ""); nut.err != nil {
		return nut
	}

	return nut.ListenAndServeWithConfig(conf)
}

// ListenAndServeTLS Открытие адреса или сокета с использованием TLS, без использования конфигурации сервера
// (конфигурация по умолчанию), после успешного открытия адреса, выполняется запуск сервера для обслуживания
// входящих соединений.
func (nut *impl) ListenAndServeTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) Interface {
	var conf *Configuration

	if conf, nut.err = parseAddress(addr, ""); nut.err != nil {
		return nut
	}
	conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM = certFile, keyFile

	return nut.ListenAndServeTLSWithConfig(conf, tlsConfig)
}

// ListenAndServeWithConfig Настройка сервера с использованием переданной конфигурации, открытие адреса или сокета
// на прослушивание, после успешного открытия адреса, выполняется запуск сервера для
// обслуживания входящих соединений.
func (nut *impl) ListenAndServeWithConfig(conf *Configuration) Interface {
	if conf == nil {
		nut.err = Errors().NoConfiguration()
		return nut
	}
	nut.conf = conf

	return nut.Listen(nil)
}

// ListenAndServeTLSWithConfig Настройка сервера с использованием переданной конфигурации в режиме TLS, открытие
// адреса или сокета на прослушивание, после успешного открытия адреса, выполняется запуск сервера для
// обслуживания входящих соединений.
func (nut *impl) ListenAndServeTLSWithConfig(conf *Configuration, tlsConfig *tls.Config) Interface {
	if conf == nil {
		nut.err = Errors().NoConfiguration()
		return nut
	}
	nut.conf = conf
	if tlsConfig == nil {
		if tlsConfig, nut.err = nut.NewTLSConfigDefault(conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM); nut.err != nil {
			return nut
		}
	}

	return nut.Listen(tlsConfig)
}

// NewListener Создание нового слушателя соединений net.Listener на основе конфигурации сервера.
func (nut *impl) NewListener(conf *Configuration) (
	ret net.Listener,
	rpc net.PacketConn,
	err error,
) {
	var (
		lstWithNames map[string][]net.Listener
		listeners    []net.Listener
		ok           bool
	)

	defaultConfiguration(conf)
	switch conf.Mode {
	case netSystemd:
		switch conf.Socket {
		case "":
			// Название сокета(ов) не указано.
			if listeners, err = nut.ListenersSystemdWithoutNames(); err != nil {
				return
			}
		default:
			// Название сокета(ов) указано.
			if lstWithNames, err = nut.ListenersSystemdWithNames(); err != nil {
				return
			}
			// Выбор сокета по названию.
			if listeners, ok = lstWithNames[path.Base(conf.Socket)]; !ok {
				err = Errors().ListenSystemdNotFound()
				return
			}
		}
		if len(listeners) > 0 {
			ret = listeners[0]
		}
	case netUnix, netUnixPacket:
		_ = os.Remove(conf.Socket)
		ret, err = net.Listen(conf.Mode, conf.Socket)
		_ = os.Chmod(conf.Socket, os.FileMode(conf.SocketMode))
	case netUdp, netUdp4, netUdp6, netUnixgram:
		rpc, err = net.ListenPacket(conf.Mode, conf.HostPort())
	default:
		ret, err = net.Listen(conf.Mode, conf.HostPort())
	}

	return
}

// NewListenerTLS Создание нового слушателя соединений net.Listener в режиме TLS, на основе конфигурации
// сервера.
func (nut *impl) NewListenerTLS(conf *Configuration, tlsConfig *tls.Config) (
	ret net.Listener,
	rpc net.PacketConn,
	err error,
) {
	const errTemplate = "публичный ключ %q, секретный ключ %q, ошибка: %s"
	var (
		lst net.Listener
		ler error
	)

	if lst, rpc, ler = nut.NewListener(conf); tlsConfig == nil {
		if tlsConfig, err = nut.NewTLSConfigDefault(conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM); err != nil {
			err = fmt.Errorf(errTemplate, conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM, err)
			return
		}
	}
	if err = ler; ler != nil {
		return
	}
	ret = tls.NewListener(lst, tlsConfig)

	return
}

// NewTLSConfigDefault Создание TLS конфигурации по умолчанию, на основе секретного и публичного ключей.
func (nut *impl) NewTLSConfigDefault(tlsPublicFile string, tlsPrivateFile string) (ret *tls.Config, err error) {
	ret = &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: make([]tls.Certificate, 1),
	}
	if ret.Certificates[0], err = tls.LoadX509KeyPair(tlsPublicFile, tlsPrivateFile); err != nil {
		return
	}

	return
}

// Listen Создание слушателя по TLS конфигурации, запуск прослушивания входящих соединений и запуск сервера.
func (nut *impl) Listen(tlsConfig *tls.Config) Interface {
	var (
		lTcp net.Listener
		lUdp net.PacketConn
	)

	if nut.isRun.Load() {
		nut.err = Errors().AlreadyRunning()
		return nut
	}
	switch tlsConfig == nil {
	case true:
		lTcp, lUdp, nut.err = nut.NewListener(nut.conf)
	case false:
		lTcp, lUdp, nut.err = nut.NewListenerTLS(nut.conf, tlsConfig)
	}
	if nut.err != nil {
		return nut
	}
	switch {
	case lUdp != nil:
		return nut.ServeUdp(lUdp)
	default:
		return nut.Serve(lTcp)
	}
}

// Serve Запуск функции сервера для входящих соединений на основе переданного слушателя net.Listener.
func (nut *impl) Serve(ltn net.Listener) Interface { return nut.serve(netListenerTcp(ltn)) }

// ServeUdp Запуск функции сервера для входящих UDP пакетов на основе переданного слушателя net.PacketConn.
func (nut *impl) ServeUdp(lpc net.PacketConn) Interface { return nut.serve(netListenerUdp(lpc)) }

// Запуск функции сервера для входящих соединений на основе слушателя netListener, который может содержать
// соединения следующих типов: udp, tcp, socket.
func (nut *impl) serve(nl *netListener) Interface {
	var (
		conf *Configuration
		onUp chan struct{}
	)

	// Защита от возможной смертельной блокировки при остановке сервера из разных потоков.
	nut.lck.Lock()
	defer nut.lck.Unlock()
	// Выход, если сервер запущен или начато завершение работы сервера.
	if nut.isRun.Load() || nut.isShutdown.Load() {
		nut.err = Errors().AlreadyRunning()
		return nut
	}
	if nut.listener = nl; nut.conf == nil {
		conf, _ = parseAddress(nut.listener.Addr().String(), nut.listener.Addr().Network())
		defaultConfiguration(conf)
		nut.conf = conf
	}
	if nut.onShutdown == nil {
		nut.onShutdown = make(chan struct{})
	}
	// Обеспечение контролируемого синхронного запуска потока.
	onUp = make(chan struct{})
	go nut.run(onUp)
	safeWait(onUp) // Текущий поток ожидает обратную связь из запущенного потока.

	return nut
}

// Процесс веб сервера.
func (nut *impl) run(onUp chan struct{}) {
	var err error

	// Финализация сокетов.
	defer func() {
		if nut.conf.Socket == "" {
			return
		}
		switch nut.conf.Mode {
		case netUnix, netUnixPacket:
			_ = os.Remove(nut.conf.Socket)
		}
	}()
	// Обеспечение синхронного запуска потока.
	nut.isRun.Store(true)
	onUp <- struct{}{}
	switch nut.listener.isUdp() {
	case true:
		if nut.handlerUdp == nil {
			nut.err = Errors().ServerHandlerUdpIsNotSet()
			return
		}
	default:
		if nut.handler == nil {
			nut.err = Errors().ServerHandlerIsNotSet()
			return
		}
	}
	// Запуск основной функции сервера.
	if err = nut.safeHandlerRun(); !nut.isShutdown.Load() && err != nil {
		nut.err = err
	}
	// Финализация флагов.
	nut.isShutdown.Store(false)
	nut.isRun.Store(false)
	safeSendSignal(nut.onShutdown)
}

// Безопасный запуск пользовательской основной функции сервера.
func (nut *impl) safeHandlerRun() (err error) {
	defer func() { err = recoverErrorWithStack(recover(), err) }()
	switch nut.listener.isUdp() {
	case true:
		err = nut.handlerUdp(nut.listener.Udp())
	default:
		err = nut.handler(nut.listener.Tcp())
	}

	return
}
