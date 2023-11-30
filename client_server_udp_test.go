package net

import (
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"net"
	"strings"
	"testing"
)

// Основная функция UDP сервера.
func testUdpHandler(udpServer net.PacketConn) (err error) {
	const bufLen, errConnClosed = 1024*64 - 8, "use of closed network connection"
	var (
		buf  []byte
		addr net.Addr
		n    int
	)

	for {
		buf = make([]byte, bufLen)
		if n, addr, err = udpServer.ReadFrom(buf); err != nil {
			if strings.Contains(err.Error(), errConnClosed) {
				err = nil
				break
			}

			continue
		}
		go testUdpIncoming(udpServer, addr, buf[:n])
	}

	return
}

// Функция обработки входящих UDP пакетов.
func testUdpIncoming(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	var (
		err error
		h   hash.Hash
	)

	h = sha512.New()
	_, err = h.Write(buf)
	if _, err = udpServer.WriteTo([]byte(fmt.Sprintf("%x", h.Sum(nil))), addr); err != nil {
		println(fmt.Sprintf("запись во входящее соединение прервана ошибкой: %s", err))
	}
}

// Функция создаёт процесс клиента, ожидая запуска горутины, возвращает ожидаемый результат и канал с результатом.
func testUdpClient(nut Interface, onStart chan<- struct{}, t *testing.T) (
	ret string,
	rsp chan []byte,
) {
	const (
		bufLen  = 1024*64 - 8
		content = `Предвижу всё: вас оскорбит
Печальной тайны объясненье.
Какое горькое презренье
Ваш гордый взгляд изобразит!
Чего хочу? с какою целью
Открою душу вам свою?
Какому злобному веселью,
Быть может, повод подаю!

Случайно вас когда-то встретя,
В вас искру нежности заметя,
Я ей поверить не посмел:
Привычке милой не дал ходу;
Свою постылую свободу
Я потерять не захотел.
Ещё одно нас разлучило...
Несчастной жертвой Ленской пал...
Ото всего, что сердцу мило,
Тогда я сердце оторвал;
Чужой для всех, ничем не связан,
Я думал: вольность и покой
Замена счастью. Боже мой!
Как я ошибся, как наказан!

Нет, поминутно видеть вас,
Повсюду следовать за вами,
Улыбку уст, движенье глаз
Ловить влюблёнными глазами,
Внимать вам долго, понимать
Душой всё ваше совершенство,
Пред вами в муках замирать,
Бледнеть и гаснуть... вот блаженство!

И я лишён того: для вас
Тащусь повсюду наудачу;
Мне дорог день, мне дорог час:
А я в напрасной скуке трачу
Судьбой отсчитанные дни.
И так уж тягостны они.
Я знаю: век уж мой измерен;
Но чтоб продлилась жизнь моя,
Я утром должен быть уверен,
Что с вами днём увижусь я...

Боюсь: в мольбе моей смиренной
Увидит ваш суровый взор
Затеи хитрости презренной —
И слышу гневный ваш укор.
Когда б вы знали, как ужасно
Томиться жаждою любви,
Пылать — и разумом всечасно
Смирять волнение в крови;
Желать обнять у вас колени,
И, зарыдав, у ваших ног
Излить мольбы, признанья, пени,
Всё, всё, что выразить бы мог,
А между тем притворным хладом
Вооружать и речь и взор,
Вести спокойный разговор,
Глядеть на вас весёлым взглядом!..

Но так и быть: я сам себе
Противиться не в силах боле;
Всё решено: я в вашей воле,
И предаюсь моей судьбе.`
	)
	var h hash.Hash

	rsp = make(chan []byte, 10)
	h = sha512.New()
	_, _ = h.Write([]byte(content))
	ret = fmt.Sprintf("%x", h.Sum(nil))
	go func(ch chan<- struct{}) {
		var (
			err  error
			addr *net.UDPAddr
			conn *net.UDPConn
			buf  []byte
			n    int
		)

		onStart <- struct{}{}
		if addr, err = net.ResolveUDPAddr("udp", "127.0.0.1:8001"); err != nil {
			t.Errorf("функция ResolveUDPAddr(), ошибка: %v, ожидалось: %v", err, nil)
		}
		if conn, err = net.DialUDP("udp", nil, addr); err != nil {
			t.Errorf("функция DialUDP(), ошибка: %v, ожидалось: %v", err, nil)
		}
		if n, err = io.WriteString(conn, content); err != nil {
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
			t.Errorf("чтение входящего пакета прервано ошибкой: %s", err)
		}
		if err = conn.Close(); err != nil {
			t.Errorf("функция conn.Close(), ошибка: %v, ожидалось: %v", err, nil)
		}
		nut.Stop()
		close(rsp)
	}(onStart)

	return
}

// Тестирование UDP клиента и сервера, с обменом данными.
func TestClientServerUdp(t *testing.T) {
	var (
		err      error
		nut      Interface
		onStart  chan struct{}
		data     string
		response chan []byte
		buf      []byte
	)

	nut = New().
		HandlerUdp(testUdpHandler).
		ListenAndServeWithConfig(&Configuration{
			Host: "127.0.0.1",
			Port: 8001,
			Mode: "udp",
		})

	// Контролируемый запуск клиента.
	onStart = make(chan struct{})
	data, response = testUdpClient(nut, onStart, t)
	safeWait(onStart)
	// Ожидание завершения сервера.
	if err = nut.Wait().
		Error(); err != nil {
		t.Errorf("функция Wait(), ошибка: %v, ожидалась: %v", err, nil)
	}
	// Чтение результата из буферизированного канала, в него поступит контрольная сумма полученных сервером данных.
	buf = <-response
	if !strings.EqualFold(data, string(buf)) {
		t.Errorf("тестирование сеанса связи между сервером и клиентом завершилось провалом")
	}
}
