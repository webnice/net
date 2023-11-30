# net

[![GoDoc](https://godoc.org/github.com/webnice/net?status.png)](http://godoc.org/github.com/webnice/net)
[![Go Report Card](https://goreportcard.com/badge/github.com/webnice/net)](https://goreportcard.com/report/github.com/webnice/net)

#### Описание

Библиотека, надстройка над "net", для создания управляемого сервера на основе стандартной библиотеки net.
Предназначена для создания серверов:

* UDP - Сервера принимающие и отвечающие на UDP пакеты.
* TCP/IP - Сервера принимающие TCP/IP запросы (как чистые TCP/IP, так и http, rpc или gRPC и другие).
* TLS - Сервера на основе TCP/IP запросов с использованием TLS шифрования (те же сервера, что TCP/IP, но с использованием TLS шифрования, например https).
* socket - Сервера поднимающие unix socket и полностью работающие через него.
* systemd - Сервера, запускаемые через systemd с использованием технологии передачи соединения через файловый сокет, когда прослушиваемый порт открывает systemd от пользователя root, затем, открытый порт передаёт процессу запущенному без прав, через файловый дескриптор (документация: man systemd.socket(5)).

#### Подключение
```bash
go get github.com/webnice/net
```

### Использование в приложении

```go
package main

import (
	"log"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	wns "github.com/webnice/net"
)

func main() {
	nut := wns.New().
		Handler(func(l net.Listener) error {
			ehr := echo.New()
			ehr.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, World!")
			})
			srv := &http.Server{
				Addr:    "localhost:8080",
				Handler: ehr,
			}

			return srv.Serve(l)
		}).
		ListenAndServe("localhost:8080")
	if nut.Error() != nil {
		log.Fatalf("запуск сервера прерван ошибкой: %s", nut.Error())
		return
	}
	// Ожидание завершения сервера.
	if err := nut.Wait().
		Error(); err != nil {
		log.Fatalf("сервер завершился с ошибкой: %s", nut.Error())
		return
	}
}
```
