package net

import (
	"crypto/sha512"
	"crypto/tls"
	"errors"
	"fmt"
	"hash"
	"io"
	"net"
	"strings"
	"testing"
)

// Функция возвращает TLS конфиг.
func testTls(key, crt string) (ret *tls.Config) {
	var (
		err  error
		cert tls.Certificate
	)

	if cert, err = tls.LoadX509KeyPair(crt, key); err != nil {
		println(fmt.Sprintf("функция LoadX509KeyPair(), ошибка: %v, ожидалась: %v", err, nil))
		return nil
	}
	return &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true, // Сертификат вообще не подписан. :)
	}
}

// Основная функция TCP сервера.
func testTcpHandler(l net.Listener) (err error) {
	var (
		ne   net.Error
		conn net.Conn
	)

	for {
		if conn, err = l.Accept(); err != nil {
			if errors.As(err, &ne) && ne.Timeout() {
				continue
			}
			break
		}
		go testTcpIncoming(conn)
	}

	return
}

// Функция обработки входящих соединений TCP сервера.
func testTcpIncoming(conn net.Conn) {
	const bufLen = 1024
	var (
		err  error
		tlsc *tls.Conn
		buf  []byte
		n    int
		ok   bool
		h    hash.Hash
	)

	buf = make([]byte, bufLen)
	defer func() { _ = conn.Close() }()
	if tlsc, ok = conn.(*tls.Conn); !ok {
		return
	}
	if n, err = tlsc.Read(buf); n > 0 {
		h = sha512.New()
		_, _ = h.Write(buf[:n])
		if _, err = io.WriteString(tlsc, fmt.Sprintf("%x", h.Sum(nil))); err != nil {
			println(fmt.Sprintf("запись во входящий поток прервана ошибкой: %s", err))
		}
	}
	switch err {
	case nil:
	case io.EOF:
		err = nil
	default:
		println(fmt.Sprintf("чтение входящего потока прервано ошибкой: %s", err))
		return
	}
}

// Функция создаёт процесс клиента, ожидая запуска горутины, возвращает ожидаемый результат и канал с результатом.
func testTcpClient(nut Interface, tlsCfg *tls.Config, onStart chan<- struct{}, t *testing.T) (
	ret string,
	rsp chan []byte,
) {
	const (
		bufLen  = 1024
		content = `Подруга дней моих суровых,
Голубка дряхлая моя!
Одна в глуши лесов сосновых
Давно, давно ты ждёшь меня.
Ты под окном своей светлицы
Горюешь, будто на часах,
И медлят поминутно спицы
В твоих наморщенных руках.
Глядишь в забытые вороты
На чёрный отдалённый путь:
Тоска, предчувствия, заботы
Теснят твою всечасно грудь.
То чудится тебе...`
	)
	var h hash.Hash

	rsp = make(chan []byte, 10)
	h = sha512.New()
	_, _ = h.Write([]byte(content))
	ret = fmt.Sprintf("%x", h.Sum(nil))
	go func(ch chan<- struct{}) {
		var (
			err  error
			conn *tls.Conn
			buf  []byte
			n    int
		)

		onStart <- struct{}{}
		if conn, err = tls.Dial("tcp", "127.0.0.1:8001", tlsCfg); err != nil {
			t.Errorf("функция DialTCP(), ошибка: %v, ожидалось: %v", err, nil)
		}
		if n, err = conn.Write([]byte(content)); err != nil {
			t.Errorf("функция conn.Write(), ошибка: %v, ожидалось: %v", err, nil)
		}
		buf = make([]byte, bufLen)
		if n, err = conn.Read(buf); n > 0 {
			rsp <- buf[:n]
		}
		switch err {
		case nil:
		case io.EOF:
			err = nil
		default:
			t.Errorf("чтение входящего потока прервано ошибкой: %s", err)
		}
		if err = conn.Close(); err != nil {
			t.Errorf("функция conn.Close(), ошибка: %v, ожидалось: %v", err, nil)
		}
		nut.Stop()
		close(rsp)
	}(onStart)

	return
}

// Тестирование TCP/IP клиента и сервера, с обменом данными и шифрованием.
func TestClientServerTcp(t *testing.T) {
	var (
		err      error
		nut      Interface
		key      *tmpFile
		crt      *tmpFile
		onStart  chan struct{}
		response chan []byte
		tlsCfg   *tls.Config
		buf      []byte
		data     string
	)

	nut = New()
	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	// Запуск сервера.
	tlsCfg = testTls(key.Filename, crt.Filename)
	nut.
		Handler(testTcpHandler).
		ListenAndServeTLS("127.0.0.1:8001", crt.Filename, key.Filename, tlsCfg)
	// Контролируемый запуск клиента.
	onStart = make(chan struct{})
	data, response = testTcpClient(nut, tlsCfg, onStart, t)
	safeWait(onStart)
	// Ожидание завершения сервера.
	if err = nut.Wait().
		Error(); err != nil {
		t.Errorf("функция Wait(), ошибка: %v, ожидалась: %v", err, nil)
	}
	// Чтение результата из буферизированного канала, в него поступит контрольная сумма полученных сервером данных.
	buf = <-response
	if !strings.EqualFold(data, string(buf)) {

		t.Log(data)
		t.Log(string(buf))

		t.Errorf("тестирование сеанса связи между сервером и клиентом завершилось провалом")
	}
}
