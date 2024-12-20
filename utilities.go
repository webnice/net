package net

import (
	"fmt"
	"net"
	"os"
	runtimeDebug "runtime/debug"
	"strconv"
	"strings"
)

// Наполнение конфигурации значениями по умолчанию.
// Проверка и исправление значений.
func defaultConfiguration(conf *Configuration) {
	if conf.SocketMode == "" {
		conf.SocketMode = strconv.FormatUint(uint64(parseFileModeWithDefault("")), 8)
	}
	// Проверка Mode.
	switch strings.ToLower(conf.Mode) {
	case netUdp, netUdp4, netUdp6, netTcp, netTcp4, netTcp6, netUnix, netUnixPacket, netSystemd:
		conf.Mode = strings.ToLower(conf.Mode)
	case netSocket:
		conf.Mode = netUnix
	default:
		conf.Mode = netTcp
	}
	if conf.Mode == netUnix && conf.Socket == "" || conf.Mode == netUnixPacket && conf.Socket == "" {
		conf.Mode = netTcp
	}
	if conf.Address == "" && conf.Mode == netTcp {
		if conf.Port == 0 {
			conf.Address = conf.Host
		} else {
			conf.Address = conf.HostPort()
		}
	}
}

// Конвертация строки содержащей восьмеричное значение в 32 битное число os.FileMode.
func parseFileModeWithDefault(mode string) (ret os.FileMode) {
	var (
		err  error
		ui64 uint64
	)

	if ui64, err = strconv.ParseUint(mode, 8, 32); err != nil {
		ret = os.FileMode(defaultSocketFileMode)
		return
	}
	ret = os.FileMode(uint32(ui64))

	return
}

// Разбор адреса, определение порта через net.LookupPort, в том числе портов заданных через синонимы,
// например ":http".
func parseAddress(addr string, mode string) (ret *Configuration, err error) {
	const bColon = ":"
	var (
		sp []string
		n  int
	)

	addr = strings.TrimSpace(addr)
	ret, sp = new(Configuration), make([]string, 2)
	defer defaultConfiguration(ret)
	if mode != "" {
		ret.Mode = mode
	}
	if sp[0], sp[1], err = net.SplitHostPort(addr); err != nil {
		ret.Host, err = addr, nil
		return
	}
	if n, err = net.LookupPort(netTcp, strings.Join(sp[1:], bColon)); err != nil {
		return
	}
	ret.Host, ret.Port = sp[0], uint16(n)
	switch sp, err = net.LookupHost(sp[0]); err {
	case nil:
		ret.Host = sp[0]
	default:
		err = nil
	}

	return
}

// Выбор ошибки и добавление стека вызова к ошибке восстановления после паники.
func recoverErrorWithStack(e1 any, e2 error) (err error) {
	if err = e2; e1 != nil && err == nil {
		switch et := e1.(type) {
		case error:
			err = et
		default:
			err = fmt.Errorf("%v", e1)
		}
		err = fmt.Errorf("%s\n%s", err, string(runtimeDebug.Stack()))
	}

	return
}

// Функция закрытия файлового дескриптора, при тестировании должна отключаться.
func fileClose(file *os.File) error { return file.Close() }

// Ожидание сигнала из канала, закрытие канала получения сигнала.
func safeWait(ch chan struct{}) {
	<-ch
	safeClose(ch)
}

// Безопасное закрытие канала.
func safeClose(ch chan struct{}) {
	defer func() { _ = recover() }()
	close(ch)
}

// Безопасная отправка сигнала.
func safeSendSignal(ch chan struct{}) {
	defer func() { _ = recover() }()
	ch <- struct{}{}
}
