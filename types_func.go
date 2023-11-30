package net

import "fmt"

// HostPort Формирование строки HOST:PORT.
func (uco *Configuration) HostPort() (ret string) {
	switch uco.Mode {
	case netUnix, netUnixPacket:
		ret = fmt.Sprintf("%s:%s", uco.Mode, uco.Socket)
	case netSystemd:
		ret = uco.Mode
	default:
		ret = fmt.Sprintf("%s:%d", uco.Host, uco.Port)
	}

	return
}
