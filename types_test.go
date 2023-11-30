package net

import (
	"fmt"
	"testing"
)

func TestHostPort(t *testing.T) {
	const (
		vHost, vPort, vSocket               = "w2KkARz5ZsjblT0Rg1rq", 64513, "hj7XDG3pzOxpWm6TyCJv.socket"
		vModeTcp, vModeSocket, vModeSystemd = netTcp, netUnix, netSystemd
	)
	var (
		uco *Configuration
		hpo string
		hpe string
	)

	uco = new(Configuration)
	defaultConfiguration(uco)
	if hpo = uco.HostPort(); hpo != ":0" {
		t.Errorf("функция HostPort, результат: %q, ожидалось: %q", hpo, ":0")
	}
	// TCP
	uco = &Configuration{
		Host:   vHost,
		Port:   vPort,
		Socket: vSocket,
		Mode:   vModeTcp,
	}
	hpe = fmt.Sprintf("%s:%d", vHost, vPort)
	if hpo = uco.HostPort(); hpo != hpe {
		t.Errorf("функция HostPort, результат: %q, ожидалось: %q", hpo, hpe)
	}
	// Socket
	uco.Mode = vModeSocket
	hpe = fmt.Sprintf("%s:%s", netUnix, vSocket)
	if hpo = uco.HostPort(); hpo != hpe {
		t.Errorf("функция HostPort, результат: %q, ожидалось: %q", hpo, hpe)
	}
	// Systemd
	uco.Mode = vModeSystemd
	hpe = vModeSystemd
	if hpo = uco.HostPort(); hpo != hpe {
		t.Errorf("функция HostPort, результат: %q, ожидалось: %q", hpo, hpe)
	}
}
