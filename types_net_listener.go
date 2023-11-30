package net

import "net"

// Внутренняя структура хранения интерфейса слушателя соединения.
type netListener struct {
	udp net.PacketConn
	tcp net.Listener
}

// Возвращается истина для объекта содержащего UDP соединение.
func (l *netListener) isUdp() bool { return l.udp != nil }

// Close Закрытие соединения.
func (l *netListener) Close() error {
	switch l.isUdp() {
	case true:
		return l.udp.Close()
	default:
		return l.tcp.Close()
	}
}

// Addr Возвращается локальный сетевой адрес или прослушиваемый адрес, если таковой есть.
func (l *netListener) Addr() net.Addr {
	switch l.isUdp() {
	case true:
		return l.udp.LocalAddr()
	default:
		return l.tcp.Addr()
	}
}

// Udp Возвращается интерфейс UDP соединения.
func (l *netListener) Udp() net.PacketConn { return l.udp }

// Tcp Возвращается интерфейс TCP/IP соединения или сокета.
func (l *netListener) Tcp() net.Listener { return l.tcp }

// Конструктор объекта для UDP соединения.
func netListenerUdp(l net.PacketConn) (ret *netListener) {
	return &netListener{udp: l}
}

// Конструктор объекта для TCP соединения.
func netListenerTcp(l net.Listener) (ret *netListener) {
	return &netListener{tcp: l}
}
