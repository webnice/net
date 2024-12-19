package net

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
)

// Проверка функции загрузки переменных окружения выставляемых systemd.
func TestListenEnv(t *testing.T) {
	const (
		envListenPid, envListenFds = "LISTEN_PID", "LISTEN_FDS"
		envListenFdnames           = "LISTEN_FDNAMES"
		s0, s1, s2, s3             = "service0.socket", "service1.socket", "service2.socket", "service3.socket"
		sepColon                   = ":"
	)
	var (
		err   error
		obj   *impl
		env   *listenEnv
		sar   []string
		n, s  int
		found bool
	)

	obj = New().(*impl)
	if env, err = obj.ListenEnv(); err == nil {
		t.Errorf("функция ListenEnv() повреждена, ожидалась ошибка")
	}
	// Должна вернуться ошибка о не корректном PID.
	if !strings.Contains(err.Error(), "PID") {
		t.Errorf("функция ListenEnv() повреждена, ожидалась ошибка получения PID")
	}
	// Установка переменной окружения LISTEN_PID.
	if err = os.Setenv(envListenPid, "1"); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenPid)
	}
	// Должна вернуться ошибка о не корректном LISTEN_FDS.
	env, err = obj.ListenEnv()
	if !strings.Contains(err.Error(), envListenFds) {
		t.Errorf("функция ListenEnv() повреждена, ожидалась ошибка получения %q", envListenFds)
	}
	// Установка переменной окружения LISTEN_FDS.
	if err = os.Setenv(envListenFds, "0"); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenFds)
	}
	// Проверка соответствия PID.
	if env, err = obj.ListenEnv(); !errors.Is(err, Errors().ListenSystemdPID()) {
		t.Errorf("функция ListenEnv(), ошибка: %v, ожидалось: %v", err, Errors().ListenSystemdPID())
	}
	// Установка корректного значения.
	if err = os.Setenv(envListenPid, fmt.Sprint(os.Getpid())); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenPid)
	}
	defer func() { _ = os.Unsetenv(envListenPid) }()
	// Проверка ошибки не корректного значения переменной окружения LISTEN_FDS.
	if env, err = obj.ListenEnv(); !errors.Is(err, Errors().ListenSystemdFDS()) {
		t.Errorf("функция ListenEnv(), ошибка: %v, ожидалось: %v", err, Errors().ListenSystemdFDS())
	}
	// Установка корректного значения переменной LISTEN_FDS.
	sar = []string{s0, s1, s2, s3}
	if err = os.Setenv(envListenFds, fmt.Sprint(len(sar))); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenFds)
	}
	defer func() { _ = os.Unsetenv(envListenFds) }()
	// Проверка ошибки не соответствия количества заявленному значению.
	if env, err = obj.ListenEnv(); !errors.Is(err, Errors().ListenSystemdQuantityNotMatch()) {
		t.Errorf("функция ListenEnv(), ошибка: %v, ожидалось: %v", err, Errors().ListenSystemdQuantityNotMatch())
	}
	// Установка переменной окружения .
	if err = os.Setenv(envListenFdnames, strings.Join(sar, sepColon)); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenFdnames)
	}
	defer func() { _ = os.Unsetenv(envListenFdnames) }()
	if env, err = obj.ListenEnv(); err != nil {
		t.Errorf("функция ListenEnv() повреждена, ошибка не ожидалась")
	}
	// Проверка загруженных значений.
	if env.fds != len(sar) {
		t.Errorf("функция ListenEnv(), fds: %d, ожидалось: %d", env.fds, len(sar))
	}
	if len(env.names) != len(sar) {
		t.Errorf("функция ListenEnv(), names: %d, ожидалось: %d", len(env.names), len(sar))
	}
	for n = range sar {
		found = false
		for s = range env.names {
			if env.names[s] == sar[n] {
				found = true
				break
			}
		}
		if !found {
			t.Errorf(
				"функция ListenEnv() повреждена, не найдено значение %q в переменной %q",
				sar[n], envListenFdnames,
			)
		}
	}
}

func TestListenLoadFilesFdWithNames(t *testing.T) {
	const s0, s1, s2, s3 = "service0.socket", "service1.socket", "service2.socket", "service3.socket"
	var (
		err  error
		sar  []string
		obj  *impl
		fla  []*os.File
		n, s int
		ok   bool
	)

	obj = New().(*impl)
	obj.fnFl = func(_ *os.File) (ret net.Listener, err error) {
		if ret, err = net.Listen("tcp", "127.0.0.1:18080"); err == nil {
			_ = ret.Close()
		}
		return
	}
	obj.fnNf = func(_ uintptr, name string) *os.File { return os.NewFile(0, name) }
	obj.fnFc = func(_ *os.File) error { return nil }
	if _, err = obj.ListenLoadFilesFdWithNames(); err == nil {
		t.Errorf("функция ListenersSystemdWithoutNames(), ожидалась ошибка")
	}
	// Выставление переменных окружения.
	sar = []string{s0, s1, s2, s3}
	if err = os.Setenv("LISTEN_FDNAMES", strings.Join(sar, ":")); err != nil {
		t.Fatalf("невозможно установить переменные окружения LISTEN_FDNAMES")
	}
	defer func() { _ = os.Unsetenv("LISTEN_FDNAMES") }()
	if err = os.Setenv("LISTEN_PID", fmt.Sprint(os.Getpid())); err != nil {
		t.Fatalf("невозможно установить переменные окружения LISTEN_PID")
	}
	defer func() { _ = os.Unsetenv("LISTEN_PID") }()
	if err = os.Setenv("LISTEN_FDS", fmt.Sprint(len(sar))); err != nil {
		t.Fatalf("невозможно установить переменные окружения LISTEN_FDS")
	}
	defer func() { _ = os.Unsetenv("LISTEN_FDS") }()
	if fla, err = obj.ListenLoadFilesFdWithNames(); err != nil {
		t.Errorf("функция ListenLoadFilesFdWithNames(), ошибка: %v, ожидалось: %v", err, nil)
	}
	for s = range sar {
		ok = false
		for n = range fla {
			if sar[s] == fla[n].Name() {
				ok = true
				break
			}
		}
		if !ok {
			t.Fatalf("не найден сокет %q", sar[s])
		}
	}
}

func TestNewListener(t *testing.T) {
	const (
		keySystemd, s0, s1         = "systemd", "service0.socket", "service1.socket"
		envListenPid, envListenFds = "LISTEN_PID", "LISTEN_FDS"
		envListenFdnames           = "LISTEN_FDNAMES"
		sepColon                   = ":"
	)
	var (
		err  error
		nut  Interface
		sar  []string
		okFn func(*os.File) (net.Listener, error)
		erFn func(*os.File) (net.Listener, error)
	)

	// Переменные окружения для сокета systemd.
	sar = []string{s0}
	if err = os.Setenv(envListenFdnames, strings.Join(sar, sepColon)); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenFdnames)
	}
	defer func() { _ = os.Unsetenv(envListenFdnames) }()
	if err = os.Setenv(envListenPid, fmt.Sprint(os.Getpid())); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenPid)
	}
	defer func() { _ = os.Unsetenv(envListenPid) }()
	if err = os.Setenv(envListenFds, fmt.Sprint(len(sar))); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenFds)
	}
	defer func() { _ = os.Unsetenv(envListenFds) }()
	// Подготовка.
	okFn = func(_ *os.File) (ret net.Listener, err error) {
		ret, err = net.Listen("tcp", "127.0.0.1:18080")
		_ = ret.Close()
		return
	}
	erFn = func(_ *os.File) (ret net.Listener, err error) { err = errors.New("ошибка"); return }
	nut = New()
	nut.(*impl).fnNf = func(_ uintptr, _ string) *os.File { return os.NewFile(0, s0) }
	nut.(*impl).fnFc = func(_ *os.File) error { return nil }
	// Название сокета не указано.
	nut.(*impl).fnFl = erFn
	if _, _, err = nut.NewListener(&Configuration{Mode: keySystemd}); err == nil {
		t.Errorf("функция NewListener() повреждена, ожидалась ошибка")
	}
	nut.(*impl).fnFl = okFn
	if _, _, err = nut.NewListener(&Configuration{Mode: keySystemd}); err != nil {
		t.Errorf("функция NewListener() повреждена, ошибка не ожидалась")
	}
	// Название сокета указано.
	nut.(*impl).fnFl = erFn
	if _, _, err = nut.NewListener(&Configuration{Mode: keySystemd, Socket: s0}); err == nil {
		t.Errorf("функция NewListener() повреждена, ожидалась ошибка")
	}
	nut.(*impl).fnFl = okFn
	if _, _, err = nut.NewListener(&Configuration{Mode: keySystemd, Socket: s0}); err != nil {
		t.Errorf("функция NewListener() повреждена, ошибка не ожидалась")
	}
	// Тестирование ошибки ErrListenSystemdNotFound().
	if _, _, err = nut.NewListener(&Configuration{Mode: keySystemd, Socket: s1}); err == nil {
		t.Errorf("функция NewListener() повреждена, ожидалась ошибка")
	}
	if !errors.Is(err, Errors().ListenSystemdNotFound()) {
		t.Errorf("функция NewListener(), ошибка: %v, ожидалась: %v", err, Errors().ListenSystemdNotFound())
	}
	// Без ошибок.
	if _, _, err = nut.NewListener(&Configuration{Mode: keySystemd, Socket: s0}); err != nil {
		t.Errorf("функция NewListener() повреждена, ошибка не ожидалась")
	}
}

func TestListenersSystemdTLSWithoutNames(t *testing.T) {
	const (
		s0, s1                     = "service0.socket", "service1.socket"
		envListenPid, envListenFds = "LISTEN_PID", "LISTEN_FDS"
		envListenFdnames           = "LISTEN_FDNAMES"
		sepColon                   = ":"
	)
	var (
		err  error
		obj  *impl
		sar  []string
		okFn func(*os.File) (net.Listener, error)
		erFn func(*os.File) (net.Listener, error)
		key  *tmpFile
		crt  *tmpFile
		tcg  *tls.Config
	)

	// Переменные окружения для сокета systemd.
	sar = []string{s0, s1}
	if err = os.Setenv(envListenFdnames, strings.Join(sar, sepColon)); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenFdnames)
	}
	defer func() { _ = os.Unsetenv(envListenFdnames) }()
	if err = os.Setenv(envListenPid, fmt.Sprint(os.Getpid())); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenPid)
	}
	defer func() { _ = os.Unsetenv(envListenPid) }()
	if err = os.Setenv(envListenFds, fmt.Sprint(len(sar))); err != nil {
		t.Fatalf("невозможно установить переменные окружения %q", envListenFds)
	}
	defer func() { _ = os.Unsetenv(envListenFds) }()
	// Подготовка.
	okFn = func(_ *os.File) (ret net.Listener, err error) {
		if ret, err = net.Listen("tcp", "127.0.0.1:18080"); err == nil {
			_ = ret.Close()
		}
		return
	}
	erFn = func(_ *os.File) (ret net.Listener, err error) { err = errors.New("ошибка"); return }
	obj = New().(*impl)
	obj.fnNf = func(_ uintptr, _ string) *os.File { return os.NewFile(0, s0) }
	obj.fnFc = func(_ *os.File) error { return nil }
	obj.fnFl = erFn
	if _, err = obj.ListenersSystemdTLSWithoutNames(nil); err == nil {
		t.Errorf("функция ListenersSystemdWithoutNames(), ожидалась ошибка")
	}
	if _, err = obj.ListenersSystemdTLSWithNames(nil); err == nil {
		t.Errorf("функция ListenersSystemdWithoutNames(), ожидалась ошибка")
	}
	obj.fnFl = okFn
	// Проверка ошибки "Конфигурация TLS сервера пустая".
	if _, err = obj.ListenersSystemdTLSWithoutNames(nil); err == nil {
		t.Errorf("функция ListenersSystemdTLSWithoutNames() повреждена, ожидалась ошибка")
	}
	if _, err = obj.ListenersSystemdTLSWithNames(nil); err == nil {
		t.Errorf("функция ListenersSystemdTLSWithNames() повреждена, ожидалась ошибка")
	}
	// Тестирование с конфигурацией TLS.
	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	if tcg, err = obj.NewTLSConfigDefault(crt.Filename, key.Filename); err != nil {
		t.Fatalf("не удалось создать *tls.Config, ошибка: %s", err)
	}
	if _, err = obj.ListenersSystemdTLSWithoutNames(tcg); err != nil {
		t.Errorf("функция ListenersSystemdTLSWithoutNames(), ошибка: %v, ожидалось: %v", err, nil)
	}
	if _, err = obj.ListenersSystemdTLSWithNames(tcg); err != nil {
		t.Errorf("функция ListenersSystemdTLSWithNames(), ошибка: %v, ожидалось: %v", err, nil)
	}
}
