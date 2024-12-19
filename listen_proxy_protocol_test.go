package net

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	proxyproto "github.com/pires/go-proxyproto"
)

func TestImpl_ListenAndServe_ProxyProtocol(t *testing.T) {
	const (
		testAddress, testNetwork        = "127.0.0.1:8080", "tcp"
		uri                             = "http" + "://" + testAddress + "/"
		testProxyAddress, testProxyPort = "111.222.21.22", 43210
	)
	var (
		nut           Interface
		conf          *Configuration
		testTransport http.RoundTripper
		origTransport http.RoundTripper
	)

	// Конфигурация.
	conf, _ = parseAddress(testAddress, testNetwork)
	conf.ProxyProtocol = true
	// Сервер.
	nut = New().
		Handler(func(l net.Listener) (err error) {
			var (
				route *chi.Mux
				srv   *http.Server
			)
			route = chi.NewRouter()
			route.Get("/", func(wr http.ResponseWriter, rq *http.Request) {
				wr.Header().Set("Content-Type", "text/plain")
				_, _ = io.WriteString(wr, rq.RemoteAddr)
			})
			srv = &http.Server{
				Addr:    testAddress,
				Handler: route,
			}
			err = srv.Serve(l)

			return
		}).
		ListenAndServeWithConfig(conf)
	if nut.Error() != nil {
		t.Errorf("запуск сервера прерван ошибкой: %s", nut.Error())
		return
	}
	// Тестовый транспорт.
	testTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			var (
				target *net.TCPAddr
				header *proxyproto.Header
			)

			if target, err = net.ResolveTCPAddr(testNetwork, testAddress); err != nil {
				return
			}
			conn, err = net.DialTCP(testNetwork, nil, target)
			header = &proxyproto.Header{
				Version:           1,
				Command:           proxyproto.PROXY,
				TransportProtocol: proxyproto.TCPv4,
				SourceAddr: &net.TCPAddr{
					IP:   net.ParseIP(testProxyAddress),
					Port: testProxyPort,
				},
				DestinationAddr: &net.TCPAddr{
					IP:   net.ParseIP("127.0.0.1"),
					Port: 8080,
				},
			}
			if _, err = header.WriteTo(conn); err != nil {
				return
			}

			return
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	// Временная подмена транспорта.
	origTransport = http.DefaultTransport
	http.DefaultTransport = testTransport
	defer func() { http.DefaultTransport = origTransport }()
	// Клиент.
	go func() {
		defer func() { nut.Stop() }() // Остановка сервера.
		rq, e := http.Get(uri)
		if e != nil {
			t.Errorf("запрос к %q прерван ошибкой: %s", uri, e)
		}
		defer func() { _ = rq.Body.Close() }()
		buf, e := io.ReadAll(rq.Body)
		if e != nil {
			t.Errorf("чтение ответа на запрос к %q прервано ошибкой: %s", uri, e)
		}
		testProxyProto := fmt.Sprintf("%s:%d", testProxyAddress, testProxyPort)
		if rq.StatusCode != 200 || !strings.Contains(string(buf), testProxyProto) {
			t.Errorf("получен ответ сервера: %q, ожидался: %q", string(buf), testProxyProto)
		}
	}()
	// Ожидание завершения сервера.
	if err := nut.Wait().
		Error(); err != nil {
		t.Errorf("сервер завершился с ошибкой: %s", err)
		return
	}
	nut.Stop()
}
