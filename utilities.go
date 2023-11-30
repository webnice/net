package net

import (
	"fmt"
	"net"
	"os"
	runtimeDebug "runtime/debug"
	"strings"
)

// Наполнение конфигурации значениями по умолчанию.
// Проверка и исправление значений.
func defaultConfiguration(conf *Configuration) {
	if conf.SocketMode == 0 {
		conf.SocketMode = uint32(os.FileMode(0640))
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

// Разбор адреса, определение порта через net.LookupPort, в том числе портов заданных через синонимы,
// например ":http".
func parseAddress(addr string, mode string) (ret *Configuration, err error) {
	const addrSep = ":"
	var (
		sp []string
		n  int
		e  error
	)

	ret = new(Configuration)
	defer defaultConfiguration(ret)
	if mode != "" {
		ret.Mode = mode
	}
	if sp = strings.Split(addr, addrSep); len(sp) <= 1 {
		ret.Host = sp[0]
		return
	}
	n, e = net.LookupPort(netTcp, strings.Join(sp[1:], addrSep))
	if err = e; err != nil {
		return
	}
	ret.Host, ret.Port = sp[0], uint16(n)

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
