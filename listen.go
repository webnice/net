package net

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/pires/go-proxyproto"
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
		nut.setError(Errors().NoConfiguration())
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
		nut.setError(Errors().NoConfiguration())
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
	var all []listenSystemdItem

	defaultConfiguration(conf)
	switch conf.Mode {
	case netSystemd:
		// Загрузка слушателей соединений из файловых дескрипторов, переданных службой systemd.
		if all, err = nut.listenSystemdLoadAll(); err != nil {
			return
		}
		for n := range all {
			// Выбор сокета по названию, если название указано.
			if conf.Socket != "" && all[n].name != path.Base(conf.Socket) {
				continue
			}
			switch all[n].nl.isUdp() {
			case true:
				// Пакетный сокет используется только при отсутствии потокового сокета.
				if rpc == nil && ret == nil {
					rpc = all[n].nl.Udp()
				}
			default:
				if ret == nil {
					ret = all[n].nl.Tcp()
				}
			}
		}
		// Слушатель с указанным названием не найден.
		if conf.Socket != "" && ret == nil && rpc == nil {
			err = Errors().ListenSystemdNotFound()
			return
		}
	case netUnix, netUnixPacket:
		_ = os.Remove(conf.Socket)
		if conf.ListenConfig != nil {
			ret, err = conf.ListenConfig.Listen(context.Background(), conf.Mode, conf.Socket)
		} else {
			ret, err = net.Listen(conf.Mode, conf.Socket)
		}
		_ = os.Chmod(conf.Socket, parseFileModeWithDefault(conf.SocketMode))
	case netUdp, netUdp4, netUdp6, netUnixgram:
		if conf.ListenConfig != nil {
			rpc, err = conf.ListenConfig.ListenPacket(context.Background(), conf.Mode, conf.HostPort())
		} else {
			rpc, err = net.ListenPacket(conf.Mode, conf.HostPort())
		}
	default:
		if conf.ListenConfig != nil {
			ret, err = conf.ListenConfig.Listen(context.Background(), conf.Mode, conf.HostPort())
		} else {
			ret, err = net.Listen(conf.Mode, conf.HostPort())
		}
	}
	// Включение ProxyProtocol.
	if ret != nil && conf.ProxyProtocol {
		//
		// TODO: Необходимо реализовать функции ConnPolicy и ValidateHeader.
		// TODO: ВНИМАНИЕ!!! Архитектура функции ConnPolicy придумана плохо, "одарённым" автором, возможна паника!
		// TODO: ВНИМАНИЕ!!! Архитектура функции ValidateHeader придумана плохо, "одарённым" автором, возможна паника!
		//
		ret = &proxyproto.Listener{
			Listener:          ret,
			ReadHeaderTimeout: conf.ProxyProtocolReadHeaderTimeout,
			ConnPolicy:        nil,
			ValidateHeader:    nil,
		}
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
			if lst != nil {
				_ = lst.Close()
			}
			if rpc != nil {
				_ = rpc.Close()
				rpc = nil
			}
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
		err  error
	)

	if nut.isRun.Load() {
		nut.setError(Errors().AlreadyRunning())
		return nut
	}
	switch tlsConfig == nil {
	case true:
		lTcp, lUdp, err = nut.NewListener(nut.conf)
	case false:
		lTcp, lUdp, err = nut.NewListenerTLS(nut.conf, tlsConfig)
	}
	if err != nil {
		nut.setError(err)
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
func (nut *impl) Serve(ltn net.Listener) Interface { return nut.ServeWithId(ltn, "") }

// ServeWithId Запуск функции сервера для входящих соединений на основе переданного слушателя net.Listener с
// указанием ID сервера.
func (nut *impl) ServeWithId(ltn net.Listener, id string) Interface {
	//goland:noinspection ALL
	return nut.serve(netListenerTcp(ltn), id)
}

// ServeUdp Запуск функции сервера для входящих UDP пакетов на основе переданного слушателя net.PacketConn.
func (nut *impl) ServeUdp(lpc net.PacketConn) Interface { return nut.ServeUdpWithId(lpc, "") }

// ServeUdpWithId Запуск функции сервера для входящих UDP пакетов на основе переданного слушателя net.PacketConn с
// указанием ID сервера.
func (nut *impl) ServeUdpWithId(lpc net.PacketConn, id string) Interface {
	//goland:noinspection ALL
	return nut.serve(netListenerUdp(lpc), id)
}

// Запуск функции сервера для входящих соединений на основе слушателя netListener, который может содержать
// соединения следующих типов: udp, tcp, socket.
func (nut *impl) serve(nl *netListener, id string) Interface {
	var (
		conf *Configuration
		onUp chan struct{}
	)

	// Защита от возможной смертельной блокировки при остановке сервера из разных потоков.
	nut.lck.Lock()
	defer nut.lck.Unlock()
	// Выход, если сервер запущен или начато завершение работы сервера.
	if nut.isRun.Load() || nut.isShutdown.Load() {
		nut.setError(Errors().AlreadyRunning())
		return nut
	}
	if nut.listener = nl; nut.conf == nil {
		conf, _ = parseAddress(nut.listener.Addr().String(), nut.listener.Addr().Network())
		defaultConfiguration(conf)
		conf.ID = id
		nut.conf = conf
	}
	// Новый канал завершения работы сервера. Канал создаётся всегда, так как после
	// завершения работы сервера канал находится в закрытом состоянии.
	nut.onShutdown = make(chan struct{})
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
	// Проверка и создание уникального ID сервера.
	if nut.conf.ID == "" {
		nut.conf.ID = uuid.NewString()
	}
	// Проверка назначения основной функции сервера, выполняется до сигнала обратной связи о запуске сервера.
	switch {
	case nut.listener.isUdp() && nut.handlerUdp == nil:
		nut.setError(Errors().ServerHandlerUdpIsNotSet())
		onUp <- struct{}{}
		return
	case !nut.listener.isUdp() && nut.handler == nil:
		nut.setError(Errors().ServerHandlerIsNotSet())
		onUp <- struct{}{}
		return
	}
	// Обеспечение синхронного запуска потока.
	nut.isRun.Store(true)
	onUp <- struct{}{}
	// Запуск основной функции сервера.
	if err = nut.safeHandlerRun(); !nut.isShutdown.Load() && err != nil {
		nut.setError(err)
	}
	// Финализация флагов. Выполняется всегда, независимо от результата запуска сервера.
	nut.isShutdown.Store(false)
	nut.isRun.Store(false)
	safeClose(nut.onShutdown)
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
