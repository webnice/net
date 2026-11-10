package net

import (
	"errors"
	"net"
	"testing"
)

// Тестирование внутренней структуры слушателя TCP соединений.
func TestNetListenerTcp(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var (
		err error
		ltn net.Listener
		nl  *netListener
	)

	if ltn, err = net.Listen("tcp", testAddress1); err != nil {
		t.Fatalf("функция Listen(%q, %q), прервана ошибкой: %s", "tcp", testAddress1, err)
	}
	nl = netListenerTcp(ltn)
	if nl == nil {
		t.Fatalf("функция netListenerTcp(), вернулся nil")
	}
	if nl.isUdp() {
		t.Errorf("функция isUdp(), вернулось: %t, ожидалось: %t", true, false)
	}
	if nl.Udp() != nil {
		t.Errorf("функция Udp(), вернулось не пустое значение, ожидался nil")
	}
	if nl.Tcp() != ltn {
		t.Errorf("функция Tcp(), вернулся не корректный слушатель")
	}
	if nl.Addr() == nil || nl.Addr().String() != ltn.Addr().String() {
		t.Errorf("функция Addr(), вернулся не корректный адрес")
	}
	if err = nl.Close(); err != nil {
		t.Errorf("функция Close(), ошибка: %v, ожидалось: %v", err, nil)
	}
	if err = nl.Close(); err == nil {
		t.Errorf("функция Close(), повторное закрытие, ожидалась ошибка")
	}
}

// Тестирование внутренней структуры слушателя UDP соединений.
func TestNetListenerUdp(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var (
		err error
		lpc net.PacketConn
		nl  *netListener
	)

	if lpc, err = net.ListenPacket("udp", testAddress1); err != nil {
		t.Fatalf("функция ListenPacket(%q, %q), прервана ошибкой: %s", "udp", testAddress1, err)
	}
	nl = netListenerUdp(lpc)
	if nl == nil {
		t.Fatalf("функция netListenerUdp(), вернулся nil")
	}
	if !nl.isUdp() {
		t.Errorf("функция isUdp(), вернулось: %t, ожидалось: %t", false, true)
	}
	if nl.Udp() != lpc {
		t.Errorf("функция Udp(), вернулся не корректный слушатель")
	}
	if nl.Tcp() != nil {
		t.Errorf("функция Tcp(), вернулось не пустое значение, ожидался nil")
	}
	if nl.Addr() == nil || nl.Addr().String() != lpc.LocalAddr().String() {
		t.Errorf("функция Addr(), вернулся не корректный адрес")
	}
	if err = nl.Close(); err != nil {
		t.Errorf("функция Close(), ошибка: %v, ожидалось: %v", err, nil)
	}
	if err = nl.Close(); err == nil {
		t.Errorf("функция Close(), повторное закрытие, ожидалась ошибка")
	}
}

// Тестирование ошибки UDP слушателя при закрытом соединении.
func TestNetListenerUdpCloseError(t *testing.T) {
	const testAddress1 = "127.0.0.1:18080"
	var (
		err error
		lpc net.PacketConn
		nl  *netListener
	)

	if lpc, err = net.ListenPacket("udp", testAddress1); err != nil {
		t.Fatalf("функция ListenPacket(%q, %q), прервана ошибкой: %s", "udp", testAddress1, err)
	}
	nl = netListenerUdp(lpc)
	_ = lpc.Close()
	if err = nl.Close(); err == nil {
		t.Errorf("функция Close(), ожидалась ошибка закрытого соединения")
	} else if !errors.Is(err, net.ErrClosed) {
		t.Errorf("функция Close(), ошибка: %v, ожидалась: %v", err, net.ErrClosed)
	}
}
