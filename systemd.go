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
		errPIDTpl     = `получение PID из переменной окружения %q прервано ошибкой: %s`
		errFDSTpl     = `получение файлового дескриптора из переменной окружения %q прервано ошибкой: %s`
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
		listenFdPrefix = `LISTEN_FD_`
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

// ListenersSystemdWithoutNames Возвращает срез net.Listener сокетов переданных в процесс сервера
// из службы linux - systemd.
func (nut *impl) ListenersSystemdWithoutNames() (ret []net.Listener, err error) {
	var (
		file  *os.File
		files []*os.File
		n     int
		pc    net.Listener
	)

	files, err = nut.ListenLoadFilesFdWithNames()
	ret = make([]net.Listener, len(files))
	for n, file = range files {
		if pc, err = nut.fnFl(file); err == nil {
			ret[n], _ = pc, nut.fnFc(file)
		}
	}

	return
}

// ListenersSystemdWithNames Возвращает карту срезов net.Listener сокетов переданных в процесс сервера
// из службы linux - systemd.
func (nut *impl) ListenersSystemdWithNames() (ret map[string][]net.Listener, err error) {
	var (
		file    *os.File
		files   []*os.File
		current []net.Listener
		pc      net.Listener
		ok      bool
	)

	files, err = nut.ListenLoadFilesFdWithNames()
	ret = make(map[string][]net.Listener)
	for _, file = range files {
		if pc, err = nut.fnFl(file); err == nil {
			if current, ok = ret[file.Name()]; !ok {
				ret[file.Name()] = make([]net.Listener, 0, 1)
			}
			ret[file.Name()] = append(current, pc)
			_ = nut.fnFc(file)
		}
	}

	return
}

// ListenersSystemdTLSWithoutNames Возвращает срез net.nnlistener для TLS сокетов переданных в процесс сервера
// из службы linux - systemd.
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

// ListenersSystemdTLSWithNames Возвращает карту срезов net.nnlistener для TLS сокетов переданных в процесс сервера
// из службы linux - systemd.
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

	return
}
