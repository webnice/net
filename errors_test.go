package net

import (
	"strings"
	"testing"
)

func TestErrors(t *testing.T) {
	var nut Interface

	if nut = New(); nut.Errors() != errSingleton {
		t.Fatalf("фунция Errors(), функция повреждена")
	}
	if Errors() != errSingleton {
		t.Fatalf("фунция Errors(), функция повреждена")
	}
}

func TestErrAlreadyRunning(t *testing.T) {
	var v any

	if Errors().AlreadyRunning() != &errAlreadyRunning {
		t.Errorf("фунция ErrAlreadyRunning(), функция повреждена")
	}
	switch v = Errors().AlreadyRunning().Error(); s := v.(type) {
	case string:
		if !strings.EqualFold(s, cAlreadyRunning) {
			t.Fatalf("фунция ErrAlreadyRunning(), функция повреждена")
		}
	default:
		t.Fatalf("функции ошибок пакета повреждены")
	}
}

func TestErrNoConfiguration(t *testing.T) {
	var v any

	if Errors().NoConfiguration() != &errNoConfiguration {
		t.Errorf("фунция ErrAlreadyRunning(), функция повреждена")
	}
	switch v = Errors().NoConfiguration().Error(); s := v.(type) {
	case string:
		if !strings.EqualFold(s, cNoConfiguration) {
			t.Fatalf("фунция ErrNoConfiguration(), функция повреждена")
		}
	default:
		t.Fatalf("функции ошибок пакета повреждены")
	}
}

// Тестирование всех ошибок пакета: возврат константного адреса ошибки и корректность текста ошибки.
func TestErrorsAll(t *testing.T) {
	var items = []struct {
		name string
		got  error
		ref  error
		text string
	}{
		{"AlreadyRunning", Errors().AlreadyRunning(), error(&errAlreadyRunning), cAlreadyRunning},
		{"NoConfiguration", Errors().NoConfiguration(), error(&errNoConfiguration), cNoConfiguration},
		{"ListenSystemdPID", Errors().ListenSystemdPID(), error(&errListenSystemdPID), cListenSystemdPID},
		{"ListenSystemdFDS", Errors().ListenSystemdFDS(), error(&errListenSystemdFDS), cListenSystemdFDS},
		{"ListenSystemdNotFound", Errors().ListenSystemdNotFound(), error(&errListenSystemdNotFound), cListenSystemdNotFound},
		{"ListenSystemdQuantityNotMatch", Errors().ListenSystemdQuantityNotMatch(), error(&errListenSystemdQuantityNotMatch), cListenSystemdQuantityNotMatch},
		{"TLSIsNil", Errors().TLSIsNil(), error(&errTLSIsNil), cTLSIsNil},
		{"ServerHandlerIsNotSet", Errors().ServerHandlerIsNotSet(), error(&errServerHandlerIsNotSet), cServerHandlerIsNotSet},
		{"ServerHandlerUdpIsNotSet", Errors().ServerHandlerUdpIsNotSet(), error(&errServerHandlerUdpIsNotSet), cServerHandlerUdpIsNotSet},
		{"ListenerIsNotCreated", Errors().ListenerIsNotCreated(), error(&errListenerIsNotCreated), cListenerIsNotCreated},
	}

	for n := range items {
		if items[n].got != items[n].ref {
			t.Errorf("функция ошибки %q вернула не константный адрес ошибки", items[n].name)
		}
		if !strings.EqualFold(items[n].got.Error(), items[n].text) {
			t.Errorf("функция ошибки %q, текст ошибки: %q, ожидался: %q", items[n].name, items[n].got.Error(), items[n].text)
		}
	}
}
