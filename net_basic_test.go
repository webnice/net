package net

import (
	"io"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestBasicWebServer(t *testing.T) {
	const addr = "localhost:8080"
	const uri = "http://" + addr + "/"

	// Сервер.
	nut := New().
		Handler(func(l net.Listener) (err error) {
			var (
				ehr *echo.Echo
				srv *http.Server
			)

			ehr = echo.New()
			ehr.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, World!")
			})
			srv = &http.Server{
				Addr:    addr,
				Handler: ehr,
			}
			err = srv.Serve(l)

			return
		}).
		ListenAndServe(addr)
	if nut.Error() != nil {
		t.Errorf("запуск сервера прерван ошибкой: %s", nut.Error())
		return
	}
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
		if rq.StatusCode != 200 || !strings.Contains(string(buf), "Hello, World") {
			t.Errorf("получен не корректный ответ сервера: %q", string(buf))
		}
	}()
	// Ожидание завершения сервера.
	if err := nut.Wait().
		Error(); err != nil {
		t.Errorf("сервер завершился с ошибкой: %s", err)
		return
	}
}
