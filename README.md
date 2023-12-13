# net

[![GoDoc](https://godoc.org/github.com/webnice/net?status.png)](http://godoc.org/github.com/webnice/net)
[![Go Report Card](https://goreportcard.com/badge/github.com/webnice/net)](https://goreportcard.com/report/github.com/webnice/net)
[![Coverage Status](https://coveralls.io/repos/github/webnice/net/badge.svg?branch=master%0Av1)](https://coveralls.io/github/webnice/net?branch=master%0Av1)

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

Пример веб сервера с роутингом через "github.com/go-chi/chi/v5".

```go
package main

import (
	"io"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	wns "github.com/webnice/net"
)

func main() {
	nut := wns.New().
		Handler(func(l net.Listener) error {
			route := chi.NewRouter()
			route.Get("/", func(wr http.ResponseWriter, rq *http.Request) {
				wr.Header().Set("Content-Type", "text/plain")
				_, _ = io.WriteString(wr, "Hello, World!")
			})
			srv := &http.Server{
				Addr:    "localhost:8080",
				Handler: route,
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
		log.Fatalf("сервер завершился с ошибкой: %s", err)
		return
	}
}
```
