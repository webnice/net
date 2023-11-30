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
* systemd - Сервера, запускаемые через systemd с использованием технологии файловый сокет, когда прослушиваемый порт открывает systemd от пользователя root, затем, открытый порт передаёт процессу запущенному без прав, через файловый дескриптор (документация: man systemd.socket(5) ).

#### Установка
```bash
go get github.com/webnice/net
```
