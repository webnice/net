package net

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestDefaultConfiguration(t *testing.T) {
	const (
		_TestHost    = "5BjSMzDCuHJWHKx2JqbW.Dd5Vr8ytnD968dDrNc3s"
		_TestAddress = "https://abc.pw"
		_TestSocket  = "/var/run/test.socket"
	)
	var conf = new(Configuration)

	defaultConfiguration(conf)
	if conf.Port != 0 {
		t.Errorf("ошибка конфигурации по умолчанию, Port: %d, ожидался порт: %d", conf.Port, 0)
	}
	if conf.Mode != "tcp" {
		t.Errorf("ошибка конфигурации по умолчанию, Mode: %q, ожидался: %q", conf.Mode, "tcp")
	}

	conf.Host = _TestHost
	conf.Mode = "socket"
	defaultConfiguration(conf)
	if conf.Mode != "tcp" {
		t.Errorf("ошибка конфигурации по умолчанию, Mode: %q, ожидался: %q", conf.Mode, "tcp")
	}
	if conf.Address != _TestHost {
		t.Errorf("ошибка конфигурации по умолчанию, Address: %q, ожидался: %q", conf.Address, _TestHost)
	}

	conf = new(Configuration)
	conf.Host = _TestHost
	conf.Port = 1234
	conf.Mode = _TestHost
	defaultConfiguration(conf)
	if conf.Mode != "tcp" {
		t.Errorf("ошибка конфигурации по умолчанию, Mode: %q, ожидался: %q", conf.Mode, "tcp")
	}
	if conf.Address != fmt.Sprintf("%s:%d", _TestHost, conf.Port) {
		t.Errorf(
			"ошибка конфигурации по умолчанию, Address: %q, ожидался: %q",
			conf.Address,
			fmt.Sprintf("%s:%d", _TestHost, conf.Port),
		)
	}

	conf = new(Configuration)
	conf.Address = _TestAddress
	conf.Host = _TestHost
	conf.Port = 3210
	conf.Mode = "unixpacket"
	defaultConfiguration(conf)
	if conf.Mode != "tcp" {
		t.Errorf("ошибка конфигурации по умолчанию, Mode: %q, ожидался: %q", conf.Mode, "tcp")
	}
	if conf.Address != _TestAddress {
		t.Errorf("ошибка конфигурации по умолчанию, Address: %q, ожидался: %q", conf.Address, _TestAddress)
	}

	conf = new(Configuration)
	conf.Host = _TestHost
	conf.Mode = "socket"
	conf.Socket = _TestSocket
	defaultConfiguration(conf)
	if conf.Mode != "unix" {
		t.Errorf("ошибка конфигурации по умолчанию, Mode: %q, ожидался: %q", conf.Mode, "unix")
	}
	if conf.Address != "" {
		t.Errorf("ошибка конфигурации по умолчанию, Address: %q, ожидался: %q", conf.Address, _TestAddress)
	}
	if conf.Mode != "unix" {
		t.Errorf("ошибка конфигурации по умолчанию, Mode: %q, ожидался: %q", conf.Mode, "unix")
	}
	if conf.HostPort() != fmt.Sprintf("unix:%s", _TestSocket) {
		t.Errorf(
			"ошибка конфигурации по умолчанию, HostPort: %q, ожидался: %q",
			conf.HostPort(),
			fmt.Sprintf("unix:%s", _TestSocket),
		)
	}
}

func TestParseAddress(t *testing.T) {
	var err error
	var conf *Configuration

	conf, err = parseAddress("", "")
	if conf == nil && err == nil {
		t.Errorf("функция parseAddress(), конфигурация равна nil, ожидалось не nil")
	}
	conf, err = parseAddress(":https", "")
	if conf.Port != 443 || err != nil {
		t.Errorf("функция parseAddress(), Port: %d, ожидался: %d, ошибка: %v", conf.Port, 443, err)
	}
	conf, err = parseAddress(":http", "")
	if conf.Port != 80 || err != nil {
		t.Errorf("функция parseAddress(), Port: %d, ожидался: %d, ошибка: %v", conf.Port, 80, err)
	}
	conf, err = parseAddress("abcd:9080", "")
	if conf.Port != 9080 || err != nil {
		t.Errorf("функция parseAddress(), Port: %d, ожидался: %d, ошибка: %v", conf.Port, 9080, err)
	}
	if conf.Host != "abcd" {
		t.Errorf("функция parseAddress(), Host: %q, ожидался: %q", conf.Host, "abcd")
	}
	conf, err = parseAddress("127.0.0.1:1080", "udp")
	if conf.Port != 1080 || conf.Mode != "udp" || err != nil {
		t.Errorf(
			"функция parseAddress(), Port: %q, Mode: %q, ожидалось: Port: %q, Mode: %q, ошибка: %v",
			conf.Port, conf.Mode,
			1080, "udp",
			err,
		)
	}
	// Проверка ошибки.
	_, err = parseAddress("abcd:abcd", "")
	if err == nil {
		t.Errorf("функция parseAddress() повреждена")
	}
}

func TestRecoverErrorWithStack(t *testing.T) {
	var (
		e1  any
		e2  any
		e3  error
		err error
	)

	if err = recoverErrorWithStack(e1, err); err != nil {
		t.Errorf("функция recoverErrorWithStack() повреждена")
	}
	e2, err = "test 321", nil
	if err = recoverErrorWithStack(e2, err); err == nil {
		t.Errorf("функция recoverErrorWithStack(), ошибка: %v, ожидалось ошибка", err)
	}
	e3, err = errors.New("test 123"), nil
	if err = recoverErrorWithStack(e3, err); err == nil {
		t.Errorf("функция recoverErrorWithStack(), ошибка: %v, ожидалось ошибка", err)
	}
	e3, err = errors.New("K2brD4iCOhIm63S7j1d1k1z"), nil
	if err = recoverErrorWithStack(nil, e3); !errors.Is(err, e3) {
		t.Errorf("функция recoverErrorWithStack(), ошибка: %v, ожидалось: %v", err, e3)
	}
	if err = recoverErrorWithStack(e2, e3); !errors.Is(err, e3) {
		t.Errorf("функция recoverErrorWithStack(), ошибка: %v, ожидалось: %v", err, e3)
	}
	e3 = nil
	if err = recoverErrorWithStack(e2, e3); errors.Is(err, e3) {
		t.Errorf("функция recoverErrorWithStack(), ошибка: %v, ожидалось ошибка со стеком", err)
	}
}

func TestFileClose(t *testing.T) {
	var (
		err error
		fh  *os.File
		nm  string
	)

	if fh, err = os.CreateTemp(os.TempDir(), ""); err != nil {
		t.Errorf("функция os.CreateTemp() повреждена")
	}
	nm = fh.Name()
	defer func(name string) { _ = os.RemoveAll(name) }(nm)
	if err = fileClose(fh); err != nil {
		t.Errorf("функция fileClose(), ошибка: %v, ожидалась: %v", err, nil)
	}
	if err = fh.Close(); err == nil {
		t.Errorf("функция fh.Close(), ошибка: %v, ожидалась ошибка", err)
	}
}

// Тестирование конвертации строки с восьмеричным значением в os.FileMode.
func TestParseFileModeWithDefault(t *testing.T) {
	var ret os.FileMode

	if ret = parseFileModeWithDefault("0755"); ret != os.FileMode(0755) {
		t.Errorf("функция parseFileModeWithDefault(), результат: %o, ожидалось: %o", ret, 0755)
	}
	if ret = parseFileModeWithDefault("0666"); ret != os.FileMode(0666) {
		t.Errorf("функция parseFileModeWithDefault(), результат: %o, ожидалось: %o", ret, 0666)
	}
	if ret = parseFileModeWithDefault("0"); ret != os.FileMode(0) {
		t.Errorf("функция parseFileModeWithDefault(), результат: %o, ожидалось: %o", ret, 0)
	}
	if ret = parseFileModeWithDefault(""); ret != os.FileMode(defaultSocketFileMode) {
		t.Errorf("функция parseFileModeWithDefault(), результат: %o, ожидалось: %o", ret, defaultSocketFileMode)
	}
	if ret = parseFileModeWithDefault("не число"); ret != os.FileMode(defaultSocketFileMode) {
		t.Errorf("функция parseFileModeWithDefault(), результат: %o, ожидалось: %o", ret, defaultSocketFileMode)
	}
}

// Тестирование режима unixgram в конфигурации по умолчанию.
func TestDefaultConfigurationUnixgram(t *testing.T) {
	const testSocket = ".test-unixgram.socket"
	var conf = &Configuration{Mode: netUnixgram, Socket: testSocket}

	defaultConfiguration(conf)
	if conf.Mode != netUnixgram {
		t.Errorf("режим unixgram с указанным сокетом, Mode: %q, ожидался: %q", conf.Mode, netUnixgram)
	}
	if conf.HostPort() != testSocket {
		t.Errorf("режим unixgram с указанным сокетом, HostPort: %q, ожидался: %q", conf.HostPort(), testSocket)
	}
	// Без указания сокета режим unixgram должен быть заменён на TCP.
	conf = &Configuration{Mode: netUnixgram}
	defaultConfiguration(conf)
	if conf.Mode != netTcp {
		t.Errorf("режим unixgram без сокета, Mode: %q, ожидался: %q", conf.Mode, netTcp)
	}
	// Режим unix без сокета.
	conf = &Configuration{Mode: netUnix}
	defaultConfiguration(conf)
	if conf.Mode != netTcp {
		t.Errorf("режим unix без сокета, Mode: %q, ожидался: %q", conf.Mode, netTcp)
	}
	// Режим unixpacket без сокета.
	conf = &Configuration{Mode: netUnixPacket}
	defaultConfiguration(conf)
	if conf.Mode != netTcp {
		t.Errorf("режим unixpacket без сокета, Mode: %q, ожидался: %q", conf.Mode, netTcp)
	}
	// Конфигурация SocketMode по умолчанию.
	conf = new(Configuration)
	defaultConfiguration(conf)
	if conf.SocketMode == "" {
		t.Errorf("конфигурация по умолчанию, SocketMode пустой, ожидалось значение по умолчанию")
	}
}

// Тестирование безопасного закрытия канала.
func TestSafeClose(t *testing.T) {
	var ch = make(chan struct{})

	safeClose(ch)
	safeClose(ch) // Повторное закрытие не должно приводить к панике.
}

// Тестирование безопасного ожидания сигнала из канала.
func TestSafeWait(t *testing.T) {
	var ch = make(chan struct{}, 1)

	ch <- struct{}{}
	safeWait(ch)
	safeWait(ch) // Ожидание из закрытого канала не должно блокироваться.
}
