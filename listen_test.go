package net

import (
	"errors"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func getTestHandlerFn(isPanic bool) HandlerFn {
	return func(l net.Listener) (err error) {
		var ne net.Error

		for {
			if isPanic {
				_ = ne.Timeout()
			}
			if _, err = l.Accept(); err != nil {
				if errors.As(err, &ne) && ne.Timeout() {
					continue
				}
				break
			}
		}

		return
	}
}

func TestInvalidPort(t *testing.T) {
	const invalidAddress = ":170000"
	var nut Interface

	nut = New().
		ListenAndServe(invalidAddress)
	if nut.Error() == nil {
		t.Errorf("функция ListenAndServe(), не корректная проверка адреса")
	}
}

func TestNoConfigurationError(t *testing.T) {
	var wsv = New()

	wsv.ListenAndServeWithConfig(nil)
	defer wsv.Stop()
	if wsv.Error() == nil {
		t.Errorf("функция ListenAndServe(), не корректная проверка адреса")
	}
	if !errors.Is(wsv.Error(), Errors().NoConfiguration()) {
		t.Errorf("функция ListenAndServe(), получена не корректная ошибка")
	}
}

func TestAlreadyRunningError(t *testing.T) {
	const (
		testAddress1 = "127.0.0.1:18080"
		testAddress2 = "127.0.0.1:18081"
	)
	var nut Interface

	nut = New().
		Handler(getTestHandlerFn(false)).
		ListenAndServe(testAddress1)
	defer nut.Stop()
	if nut.Error() != nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", nut.Error(), nil)
	}
	nut.ListenAndServe(testAddress2)
	if nut.Error() == nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", nut.Error(), Errors().AlreadyRunning())
	}
	if !errors.Is(nut.Error(), Errors().AlreadyRunning()) {
		t.Errorf("функция ListenAndServe(), не корректная ошибка")
	}
	if !errors.Is(nut.
		Clean().                      // Очистка последней ошибки.
		ListenAndServe(testAddress1). // Запуск сервера, который уже запущен.
		Error(), Errors().AlreadyRunning()) {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", nut.Error(), Errors().AlreadyRunning())
	}
}

func TestPortIsBusy(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var w1, w2 Interface

	w1 = New().
		Handler(getTestHandlerFn(false))
	w1.ListenAndServe(testAddress1)
	defer w1.Stop()
	if w1.Error() != nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", w1.Error(), nil)
	}

	w2 = New().
		Handler(getTestHandlerFn(false))
	w2.ListenAndServe(testAddress1)
	defer w2.Stop()
	if w2.Error() == nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", w2.Error(), Errors().AlreadyRunning())
	}
}

func TestUnixSocket(t *testing.T) {
	const (
		testAddress1         = ".test.socket"
		testAddress1FileMode = os.FileMode(0666)
	)
	var (
		err  error
		conf *Configuration
		w1   Interface
		fi   os.FileInfo
	)

	conf, _ = parseAddress("", "")
	conf.Mode = "socket"
	conf.Socket = testAddress1
	w1 = New().
		Handler(getTestHandlerFn(false))
	w1.ListenAndServeWithConfig(conf)
	if w1.Error() != nil {
		t.Errorf("функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v", w1.Error(), nil)
	}
	if fi, err = os.Stat(testAddress1); err != nil {
		t.Errorf("проверка юникс сокета завершилась ошибкой: %s", err)
	}
	if fi.Mode().Perm() != testAddress1FileMode.Perm() {
		t.Errorf(
			"разрешения доступа юникс сокета, Mode(): %v, ожидалось: %v",
			fi.Mode().Perm(), testAddress1FileMode.Perm(),
		)
	}
	if err = w1.
		Stop().
		Error(); err != nil {
		t.Errorf("функция Stop(), ошибка: %v, ожидалось: %v", err, nil)
	}
	if _, err = os.Stat(testAddress1); os.IsExist(err) {
		t.Errorf("юникс сокет не был удалён после остановки сервера")
	}
}

func TestUnixSocketSocketMode(t *testing.T) {
	const (
		tMode            = "0600"
		tAddress         = ".test.socket"
		tAddressFileMode = os.FileMode(0600)
	)
	var (
		err  error
		conf *Configuration
		w1   Interface
		fi   os.FileInfo
	)

	conf, _ = parseAddress("", "")
	conf.Mode = "socket"
	conf.Socket = tAddress
	conf.SocketMode = tMode
	w1 = New().
		Handler(getTestHandlerFn(false))
	w1.ListenAndServeWithConfig(conf)
	if w1.Error() != nil {
		t.Errorf("функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v", w1.Error(), nil)
	}
	if fi, err = os.Stat(tAddress); err != nil {
		t.Errorf("проверка юникс сокета завершилась ошибкой: %s", err)
	}
	if fi.Mode().Perm() != tAddressFileMode.Perm() {
		t.Errorf(
			"разрешения доступа юникс сокета, Mode(): %v, ожидалось: %v",
			fi.Mode().Perm(), tAddressFileMode.Perm(),
		)
	}
	if err = w1.
		Stop().
		Error(); err != nil {
		t.Errorf("функция Stop(), ошибка: %v, ожидалось: %v", err, nil)
	}
	if _, err = os.Stat(tAddress); os.IsExist(err) {
		t.Errorf("юникс сокет не был удалён после остановки сервера")
	}
}

func TestServe(t *testing.T) {
	const (
		testAddress1 = "127.0.0.1:18080"
		testAddress2 = "127.0.0.1:18080"
		errString    = "use of closed network connection"
	)
	var (
		err error
		ltn net.Listener
		w1  Interface
	)

	if ltn, err = net.Listen("tcp", testAddress1); err != nil {
		t.Errorf("функция Listen(%q, %q), прервана ошибкой: %s", "tcp", testAddress1, err)
	}
	w1 = New().
		Handler(getTestHandlerFn(false)).
		Serve(ltn)
	defer w1.Stop()
	if w1.(*impl).conf == nil {
		t.Errorf("ошибка создания конфигурации сервера")
	}
	if w1.(*impl).conf.Address != testAddress1 && w1.(*impl).conf.Address != testAddress2 {
		t.Errorf("ошибка создания конфигурации сервера из net.Listener. Address: %q, ожидался: %q",
			w1.(*impl).conf.Address,
			testAddress1,
		)
	}
	if err = w1.
		Stop().
		Error(); err != nil {
		t.Errorf("функция Stop(), ошибка: %v, ожидалось: %v", err, nil)
	}
	if err = ltn.Close(); err == nil {
		t.Errorf("функция Close(), ошибка: %v, ожидалась ошибка", err)
	}
	if !strings.Contains(err.Error(), errString) {
		t.Errorf("функция Close(), ошибка: %v, ожидалось: %v", err, errString)
	}
}

func TestServeAlreadyRunning(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var (
		err error
		ltn net.Listener
		w1  Interface
	)

	if ltn, err = net.Listen("tcp", testAddress1); err != nil {
		t.Errorf("функция Listen(%q, %q), прервана ошибкой: %s", "tcp", testAddress1, err)
	}
	w1 = New().
		Handler(getTestHandlerFn(false)).
		Serve(ltn)
	defer w1.Stop()
	if w1.(*impl).conf == nil {
		t.Errorf("ошибка создания конфигурации сервера")
	}
	if err = w1.
		Serve(ltn).
		Error(); err == nil {
		t.Errorf("функция Serve(), ошибка: %v, ожидалось: %v", err, Errors().AlreadyRunning())
	}
	if !errors.Is(err, Errors().AlreadyRunning()) {
		t.Errorf("функция Serve(), ошибка: %v, ожидалось: %v", err, Errors().AlreadyRunning())
	}
}

func TestServeErrServerHandlerIsNotSet(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var (
		err error
		ltn net.Listener
		w1  Interface
	)

	if ltn, err = net.Listen("tcp", testAddress1); err != nil {
		t.Errorf("функция Listen(%q, %q), прервана ошибкой: %s", "tcp", testAddress1, err)
	}
	defer func() { _ = ltn.Close() }()
	w1 = New().
		Serve(ltn)
	if err = w1.Error(); err == nil {
		t.Errorf("функция Serve(), ошибка: %v, ожидалось: %v", err, Errors().ServerHandlerIsNotSet())
	}
}

func TestServeErrServerHandlerUdpIsNotSet(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var (
		err error
		ltn net.PacketConn
		w1  Interface
	)

	if ltn, err = net.ListenPacket("udp", testAddress1); err != nil {
		t.Errorf("функция Listen(%q, %q), прервана ошибкой: %s", "tcp", testAddress1, err)
	}
	defer func() { _ = ltn.Close() }()
	w1 = New().
		ServeUdp(ltn)
	if err = w1.Error(); err == nil {
		t.Errorf("функция Serve(), ошибка: %v, ожидалось: %v", err, Errors().ServerHandlerIsNotSet())
	}
}

// Тестирование паники в основной функции сервера полученной от пользователя.
func TestServeHandlerPanic(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var (
		err error
		ltn net.Listener
		nut Interface
	)

	if ltn, err = net.Listen("tcp", testAddress1); err != nil {
		t.Errorf("функция Listen(%q, %q), прервана ошибкой: %s", "tcp", testAddress1, err)
	}
	nut = New().
		Handler(getTestHandlerFn(true)).
		Serve(ltn).
		Wait()
	_ = ltn.Close()
	if err = nut.Error(); err == nil {
		t.Errorf("функция Serve() функция повреждена")
	}
}

func TestWait(t *testing.T) {
	const (
		testAddress1 = "127.0.0.1:18080"
		testAddress2 = ".test.socket"
		ticTimeout   = time.Second / 4
	)
	var (
		tic  *time.Ticker
		cou  uint32
		w1   Interface
		conf *Configuration
	)

	w1 = New().
		Handler(getTestHandlerFn(false)).
		ListenAndServe(testAddress1)
	if w1.Error() != nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалась: %v", w1.Error(), nil)
	}
	go func() {
		tic = time.NewTicker(ticTimeout)
		defer tic.Stop()
		for {
			<-tic.C
			if cou++; cou > 4 {
				w1.Stop()
				break
			}
		}
	}()
	w1.Wait()
	if cou <= 4 {
		t.Errorf("функция Wait() повреждена")
	}
	w1 = New().
		Handler(getTestHandlerFn(false))
	conf, _ = parseAddress("", "")
	conf.Mode = "socket"
	conf.Socket = testAddress2
	w1.ListenAndServeWithConfig(conf)
	if w1.Error() != nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалась: %v", w1.Error(), nil)
	}
	go func() {
		tic = time.NewTicker(ticTimeout)
		defer tic.Stop()
		for {
			<-tic.C
			if cou++; cou > 4 {
				w1.Stop()
				break
			}
		}
	}()
	w1.Wait()
	if cou <= 4 {
		t.Errorf("функция Wait() повреждена")
	}
}

// Создание ID, если не указан.
func TestImpl_RunUuidGen(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var nut Interface

	nut = New().
		Handler(getTestHandlerFn(false)).
		ListenAndServe(testAddress1)
	defer nut.Stop()
	if nut.Error() != nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", nut.Error(), nil)
	}
	if nut.ID() == "" {
		t.Errorf("функция ID(), ошибка, ожидалось не пустое значение")
	}
}

// Проверка статического ID, если указан.
func TestImpl_ListenAndServeWithConfigUuidStatic(t *testing.T) {
	const testAddress1 = ".test.socket"
	var (
		err  error
		id   string
		conf *Configuration
		nut  Interface
	)

	id = uuid.NewString()
	conf, _ = parseAddress("", "")
	conf.ID, conf.Mode, conf.Socket = id, "socket", testAddress1
	nut = New().
		Handler(getTestHandlerFn(false))
	nut.ListenAndServeWithConfig(conf)
	if nut.Error() != nil {
		t.Errorf("функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v", nut.Error(), nil)
	}
	if nut.ID() != id {
		t.Errorf("функция ID(), вернулось: %q, ожидалось: %q", nut.ID(), id)
	}
	if err = nut.
		Stop().
		Error(); err != nil {
		t.Errorf("функция Stop(), ошибка: %v, ожидалось: %v", err, nil)
	}
}
