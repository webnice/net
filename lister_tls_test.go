package net

import (
	"errors"
	"net"
	"os"
	"testing"
)

type tmpFile struct {
	Filename string
}

func (tfo *tmpFile) Clean() { _ = os.RemoveAll(tfo.Filename) }

func newTmpFile(content []byte) (ret *tmpFile) {
	var (
		err error
		fh  *os.File
	)

	ret = new(tmpFile)
	if fh, err = os.CreateTemp(os.TempDir(), ""); err != nil {
		ret = nil
		return
	}
	defer func() { _ = fh.Close() }()
	ret.Filename = fh.Name()
	_, _ = fh.Write(content)

	return
}

func getKeyEcdsa() []byte {
	return []byte(`-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDCJGaDfRSjCg2zYopdy
M7SqBKeIpcEriH3GWTtwy3hlQSiloiyGOk25Ekpt/Ha04PahZANiAARu/6BxP3/t
kYuOdvDeAKD9fsC2m3pLEOzM+ZY8phS7qg4CTFT7Yej8UTaEX1WSd4Sq5F/zmLto
BE2ulX63u0MqdUd/GU6XIpn31kDt2MVqKgprixw7Ow3zIH47KDvdwa0=
-----END PRIVATE KEY-----`)
}

func getCrtEcdsa() []byte {
	return []byte(`-----BEGIN CERTIFICATE-----
MIICHjCCAaSgAwIBAgIUCUXV3GKvC6699j1lRVwby9bvbVQwCgYIKoZIzj0EAwIw
RTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAgFw0yMzExMjcxNjI4MDFaGA8yMTIzMTEw
MzE2MjgwMVowRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAf
BgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDB2MBAGByqGSM49AgEGBSuB
BAAiA2IABG7/oHE/f+2Ri4528N4AoP1+wLabeksQ7Mz5ljymFLuqDgJMVPth6PxR
NoRfVZJ3hKrkX/OYu2gETa6Vfre7Qyp1R38ZTpcimffWQO3YxWoqCmuLHDs7DfMg
fjsoO93BraNTMFEwHQYDVR0OBBYEFB+2RT/hFPvsibVkD5YixFyDSMAmMB8GA1Ud
IwQYMBaAFB+2RT/hFPvsibVkD5YixFyDSMAmMA8GA1UdEwEB/wQFMAMBAf8wCgYI
KoZIzj0EAwIDaAAwZQIwUDdaraaLyrL2+Lmj0xTvPI4+zUJ2qVPcVMgzKiElDCk+
75dcyLpKakC/CIgcIzhOAjEAouUFwLdcMQClcToO2jH0n3KWIJXv/y/X3yLN/t+W
4R/MDRGuSjdIm80Rtr3DbnB7
-----END CERTIFICATE-----`)
}

func getKeyRsa() []byte {
	return []byte(`-----BEGIN PRIVATE KEY-----
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
x/SaXicC5Rf0CmEgbxujpV8Yte1tNA==
-----END PRIVATE KEY-----`)
}

func getCrtRsa() []byte {
	return []byte(`-----BEGIN CERTIFICATE-----
MIIFbTCCA1WgAwIBAgIUBn/fuvXa3lP5L2IrsKMAktV4abwwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAgFw0yMzExMjcxNjMxMTNaGA8yMTIz
MTEwMzE2MzExM1owRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUx
ITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDCCAiIwDQYJKoZIhvcN
AQEBBQADggIPADCCAgoCggIBALWoExrMhDz+jm9m0MivX2zTmooH4qKLsqFRm27+
26BGESHqXmZ6/9B6Lfj4I1xgU5CexH/THzBXvgQQjf7hpDtZqlbWXip35bubDvaZ
tl9+OSduCNNnn3kXU4QvEeLBGdcR8H2Qev/5RSROxTBq5cz0smLmEKl0rHsT6vzM
RAM/5z1YwtLQdriy3+POv+qjJbzgfEt0467f/aaANYXmNBiu8otCRlqHsEqrIHbb
0hUInxoGMeLTKTvYMKsX1dnNRNLkBukaRPrcPlJK72/VbVqpb2MtjTcHbG1qiu5v
0tzjpkwDrK2kMtyRBrxokaMMdMmehkyH3MRt59ja+Am9uiZOp4sQxZmjOpJCkqkR
DlZTVu17te4p664sIXe/YvvFmX7nrLDysDWLF84UI9HAllrAy6zNU+dJby5vuRes
xXgtcmCAYXBB483oSt7/7Yqxw/99nOZ8NBWggsPp52GpcXS1EsGWLRBLD7DuLMmY
+uEgFwep9OR1OzbNdHgbnxs00gvd+oCbGL+2PsY9jek3qgup/7iTKanvPlWqLwEh
EXt/18mm2i4eTUTqhaF9ISJpfeUTZ6zPIwYy1Jw/W10UxcCG+PMzvjEbM2ElXNk1
jlBcC3kHSt74AX9nRSuqiagAc8mDSZgcv2awQntK/WGMV+gst3KBujVtXn+f+eou
OmsjAgMBAAGjUzBRMB0GA1UdDgQWBBRHMSiOsP9RMcxccbNufie8PaJz6zAfBgNV
HSMEGDAWgBRHMSiOsP9RMcxccbNufie8PaJz6zAPBgNVHRMBAf8EBTADAQH/MA0G
CSqGSIb3DQEBCwUAA4ICAQCuLlvHWgbQrhVj2E+kyUgbhE4nQoomxdx1AZFbO22+
fx6RN8/fR8lCMADlDGQK291pZvD6XqIJiAIrKxKnw1MzynwxXZH+6QbdWgwHR+ft
hdLypNu4/AJ8zGBXWyFam8fMMWr/WtfiVDXY7EqmIqHqD5NUdFRbSZ74y04l8WS3
DqJz8ZyiEevKaCn9ClPKOYRyTikSrelBAAkb/41WcXwEDqZQ21aZOSCpK3g/4Dwy
6p7J4E0ySFUliSGrDVXmPc70H96g4pwBBfMjwhuZIM9rnoDnH2Mavw6PPYSVS1Rl
Yl8u1rK8AAzCuppK+dQoLx/AU4RgrZygSsPa6dp7giG5lEjq/DXC5F2+LO2bSpXV
aqJNpPZ2+eO1G+30zKSfI9FShyapdi0KlVDR/1RD7f/ubQRF5imZe+9lDXaA+ZKy
axmY68Y40LEYRNH9qoqaCTmxxLwkJU4YNjK4ZRUGtRXVmjOqZVCpqF41JIIDqi/w
W0EqVBaRDYmLHNY3UJWcGSy0aaJh96Du6oojuo3MgLoJun367a8Uw42RMvCLu9GE
gf5z+NriNZFUQDACfnDynuYaa8Fj6IPXy/Y18WVsNqiIg8vYCcoNE8xpV2lP1uqj
1vtLWXn39KuQvOcFWuG0zWR6Z/YbSQaY6He4G8Tc3M2JRRmVZAuA2xRFjXHthaCL
6Q==
-----END CERTIFICATE-----`)
}

func TestLoadX509KeyPair(t *testing.T) {
	const testAddress1 = `localhost:18088`
	var (
		err error
		key *tmpFile
		crt *tmpFile
		nut Interface
	)

	key, crt = newTmpFile([]byte("123")), newTmpFile([]byte("321"))
	defer func() { key.Clean(); crt.Clean() }()
	nut = New().
		Handler(getTestHandlerFn(false)).
		ListenAndServeTLS(testAddress1, crt.Filename, key.Filename, nil)
	if err = nut.Error(); err == nil {
		t.Errorf("функция ListenAndServeTLS() повреждена")
	}
}

func TestListenAndServeTLSOk(t *testing.T) {
	const testAddress1 = `localhost:18088`
	var (
		err error
		key *tmpFile
		crt *tmpFile
		nut Interface
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	nut = New().
		Handler(getTestHandlerFn(false)).
		ListenAndServeTLS(testAddress1, crt.Filename, key.Filename, nil)
	if err = nut.Error(); err != nil {
		t.Errorf("функция ListenAndServeTLS(), ошибка: %v, ожидалось: %v", err, nil)
	}
	if err = nut.
		Stop().
		Error(); err != nil {
		t.Errorf("функция Stop(), ошибка: %v, ожидалось: %v", err, nil)
	}
}

func TestInvalidPortTLS(t *testing.T) {
	const invalidAddress = `:170000`
	var (
		key *tmpFile
		crt *tmpFile
		nut Interface
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	nut = New().
		ListenAndServeTLS(invalidAddress, crt.Filename, key.Filename, nil)
	if nut.Error() == nil {
		t.Errorf("функция ListenAndServeTLS(), не корректная проверка адреса")
	}
}

func TestTLSNoConfigurationError(t *testing.T) {
	var wsv = New()

	wsv.ListenAndServeTLSWithConfig(nil, nil)
	defer wsv.Stop()
	if wsv.Error() == nil {
		t.Errorf("функция ListenAndServeTLSWithConfig(), не корректная проверка адреса")
	}
	if !errors.Is(wsv.Error(), ErrNoConfiguration()) {
		t.Errorf("функция ListenAndServeTLSWithConfig(), получена не корректная ошибка")
	}
}

func TestNewListenerTLS(t *testing.T) {
	var (
		err error
		key *tmpFile
		crt *tmpFile
		nut Interface
		cfg *Configuration
		lst net.Listener
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	cfg = new(Configuration)
	defaultConfiguration(cfg)
	nut = New().
		Handler(getTestHandlerFn(false))
	if lst, _, err = nut.
		NewListenerTLS(cfg, nil); err == nil {
		t.Errorf("функция NewListenerTLS(), функция повреждена")
	}
	cfg.TLSPrivateKeyPEM, cfg.TLSPublicKeyPEM = key.Filename, crt.Filename
	if lst, _, err = nut.
		NewListenerTLS(cfg, nil); err != nil {
		t.Errorf("функция NewListenerTLS(), ошибка: %v, ожидалось: %v", err, nil)
	}
	nut = nut.Serve(lst)
	if err = nut.Error(); err != nil {
		t.Errorf("функция Serve(), функция повреждена")
	}
	if err = nut.Stop().
		Error(); err != nil {
		t.Errorf("функция Stop(), ошибка: %v, ожидалось: %v", err, nil)
	}
}
