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
	var v interface{}

	if ErrAlreadyRunning() != &errAlreadyRunning {
		t.Errorf("фунция ErrAlreadyRunning(), функция повреждена")
	}
	switch v = ErrAlreadyRunning().Error(); s := v.(type) {
	case string:
		if !strings.EqualFold(s, cAlreadyRunning) {
			t.Fatalf("фунция ErrAlreadyRunning(), функция повреждена")
		}
	default:
		t.Fatalf("функции ошибок пакета повреждены")
	}
}

func TestErrNoConfiguration(t *testing.T) {
	var v interface{}

	if ErrNoConfiguration() != &errNoConfiguration {
		t.Errorf("фунция ErrAlreadyRunning(), функция повреждена")
	}
	switch v = ErrNoConfiguration().Error(); s := v.(type) {
	case string:
		if !strings.EqualFold(s, cNoConfiguration) {
			t.Fatalf("фунция ErrNoConfiguration(), функция повреждена")
		}
	default:
		t.Fatalf("функции ошибок пакета повреждены")
	}
}
