package net

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// ListenEnv Загрузка значений переменных окружения от systemd для работы через сокет.
func (nut *impl) ListenEnv() (ret *listenEnv, err error) {
	const (
		sepColon      = ":"
		listenPID     = "LISTEN_PID"
		listenFds     = "LISTEN_FDS"
		listenFdNames = "LISTEN_FDNAMES"
		errPIDTpl     = "получение PID из переменной окружения %q прервано ошибкой: %s"
		errFDSTpl     = "получение файлового дескриптора из переменной окружения %q прервано ошибкой: %s"
	)

	ret = new(listenEnv)
	if ret.pid, err = strconv.Atoi(os.Getenv(listenPID)); err != nil {
		err = fmt.Errorf(errPIDTpl, listenPID, err)
		return
	}
	if ret.fds, err = strconv.Atoi(os.Getenv(listenFds)); err != nil {
		err = fmt.Errorf(errFDSTpl, listenFds, err)
		return
	}
	ret.names = strings.Split(os.Getenv(listenFdNames), sepColon)
	// Проверка ID процесса.
	if ret.pid != os.Getpid() {
		err = Errors().ListenSystemdPID()
		return
	}
	// Проверка количества передаваемых соединений.
	if ret.fds == 0 {
		err = Errors().ListenSystemdFDS()
		return
	}
	// Проверка заявленного количества с фактически переданным количеством соединений.
	if ret.fds != len(ret.names) {
		err = Errors().ListenSystemdQuantityNotMatch()
		return
	}

	return
}

// ListenLoadFilesFdWithNames Загрузка файловых дескрипторов на основе переменных окружения
func (nut *impl) ListenLoadFilesFdWithNames() (ret []*os.File, err error) {
	const (
		listenFdBegin  = 3
		listenFdPrefix = "LISTEN_FD_"
	)
	var (
		env        *listenEnv
		name       string
		fd, offset int
	)

	if env, err = nut.ListenEnv(); err != nil {
		return
	}
	ret = make([]*os.File, 0, env.fds)
	for fd = listenFdBegin; fd < listenFdBegin+env.fds; fd++ {
		//syscall.CloseOnExec(fd)
		name = listenFdPrefix + strconv.Itoa(fd)
		if offset = fd - listenFdBegin; offset < len(env.names) && len(env.names[offset]) > 0 {
			name = env.names[offset]
		}
		ret = append(ret, nut.fnNf(uintptr(fd), name))
	}

	return
}

// ListenersSystemdWithoutNames Возвращает срез net.Listener сокетов, переданных в процесс сервера
// из службы Linux - systemd.
func (nut *impl) ListenersSystemdWithoutNames() (ret []net.Listener, err error) {
	var (
		all  []listenSystemdItem
		item listenSystemdItem
	)

	if all, err = nut.listenSystemdLoadAll(); err != nil {
		return
	}
	ret = make([]net.Listener, 0, len(all))
	for _, item = range all {
		if !item.nl.isUdp() {
			ret = append(ret, item.nl.Tcp())
		}
	}

	return
}

// ListenersSystemdPacketWithoutNames Возвращает срез net.PacketConn пакетных сокетов (udp, unixgram),
// переданных в процесс сервера из службы Linux - systemd.
func (nut *impl) ListenersSystemdPacketWithoutNames() (ret []net.PacketConn, err error) {
	var (
		all  []listenSystemdItem
		item listenSystemdItem
	)

	if all, err = nut.listenSystemdLoadAll(); err != nil {
		return
	}
	ret = make([]net.PacketConn, 0, len(all))
	for _, item = range all {
		if item.nl.isUdp() {
			ret = append(ret, item.nl.Udp())
		}
	}

	return
}

// ListenersSystemdWithNames Возвращает карту срезов net.Listener сокетов, переданных в процесс сервера
// из службы Linux - systemd.
func (nut *impl) ListenersSystemdWithNames() (ret map[string][]net.Listener, err error) {
	var (
		all  []listenSystemdItem
		item listenSystemdItem
	)

	if all, err = nut.listenSystemdLoadAll(); err != nil {
		return
	}
	ret = make(map[string][]net.Listener)
	for _, item = range all {
		if item.nl.isUdp() {
			continue
		}
		ret[item.name] = append(ret[item.name], item.nl.Tcp())
	}

	return
}

// ListenersSystemdPacketWithNames Возвращает карту срезов net.PacketConn пакетных сокетов (udp, unixgram),
// переданных в процесс сервера из службы Linux - systemd.
func (nut *impl) ListenersSystemdPacketWithNames() (ret map[string][]net.PacketConn, err error) {
	var (
		all  []listenSystemdItem
		item listenSystemdItem
	)

	if all, err = nut.listenSystemdLoadAll(); err != nil {
		return
	}
	ret = make(map[string][]net.PacketConn)
	for _, item = range all {
		if item.nl.isUdp() {
			ret[item.name] = append(ret[item.name], item.nl.Udp())
		}
	}

	return
}

// listenSystemdLoadAll Загрузка всех файловых дескрипторов, переданных службой systemd, и преобразование
// их в слушателей соединений. Каждый файловый дескриптор пробуется как потоковый (tcp, unix) через
// net.FileListener, при ошибке дескриптор пробуется как пакетный (udp, unixgram) через net.FilePacketConn.
func (nut *impl) listenSystemdLoadAll() (ret []listenSystemdItem, err error) {
	var (
		files []*os.File
		file  *os.File
		lst   net.Listener
		pck   net.PacketConn
	)

	if files, err = nut.ListenLoadFilesFdWithNames(); err != nil {
		return
	}
	ret = make([]listenSystemdItem, 0, len(files))
	for _, file = range files {
		if lst, err = nut.fnFl(file); err == nil {
			ret = append(ret, listenSystemdItem{name: file.Name(), nl: netListenerTcp(lst)})
			_ = nut.fnFc(file)
			continue
		}
		if pck, err = nut.fnFpc(file); err == nil {
			ret = append(ret, listenSystemdItem{name: file.Name(), nl: netListenerUdp(pck)})
			_ = nut.fnFc(file)
			continue
		}
	}
	// Если ни один из переданных дескрипторов не был преобразован в слушателя.
	if len(ret) == 0 {
		err = Errors().ListenerIsNotCreated()
	}

	return
}

// ListenersSystemdTLSWithoutNames Возвращает срез net.Listener для TLS сокетов, переданных в процесс сервера
// из службы Linux - systemd.
func (nut *impl) ListenersSystemdTLSWithoutNames(tlsConfig *tls.Config) (ret []net.Listener, err error) {
	var (
		listeners []net.Listener
		l         net.Listener
	)

	if listeners, err = nut.ListenersSystemdWithoutNames(); len(listeners) == 0 || err != nil {
		return
	}
	if tlsConfig == nil {
		err = Errors().TLSIsNil()
		return
	}
	ret = make([]net.Listener, 0, len(listeners))
	for _, l = range listeners {
		ret = append(ret, tls.NewListener(l, tlsConfig))
	}

	return
}

// ListenersSystemdTLSWithNames Возвращает карту срезов net.Listener для TLS сокетов, переданных в процесс
// сервера из службы Linux - systemd.
func (nut *impl) ListenersSystemdTLSWithNames(tlsConfig *tls.Config) (ret map[string][]net.Listener, err error) {
	var (
		listeners map[string][]net.Listener
		ll        []net.Listener
		l         net.Listener
		n         int
	)

	if listeners, err = nut.ListenersSystemdWithNames(); listeners == nil || err != nil {
		return
	}
	if tlsConfig == nil {
		err = Errors().TLSIsNil()
		return
	}
	for _, ll = range listeners {
		for n, l = range ll {
			ll[n] = tls.NewListener(l, tlsConfig)
		}
	}
	ret = listeners

	return
}
