package net

import "time"

// Configuration Структура конфигурации TCP/IP или UDP сервера.
type Configuration struct {
	// ID Уникальный идентификатор сервера, может быть любым уникальным, в пределах приложения, строковым значением.
	// Если значение не указано, при запуске сервера, создаётся уникальное временное значение, меняющееся
	// при каждом запуске.
	// Default value: ""
	ID string `yaml:"ID" json:"id"`

	// Address Публичный адрес на котором сервер доступен извне.
	// Например, если сервер находится за прокси, тут указывается реальный адрес подключения к серверу.
	// Default value: "" - make automatically
	Address string `yaml:"Address" json:"address"`

	// Host IP адрес или имя хоста на котором поднимается сервер, можно указывать 0.0.0.0 для всех ip адресов.
	// Default value: "0.0.0.0"
	Host string `yaml:"Host" json:"host" default-value:"0.0.0.0"`

	// Port TCP/IP порт занимаемый сервером.
	// Default value: 0
	Port uint16 `yaml:"Port" json:"port" default-value:"-"`

	// Socket Unix socket, systemd socket на котором поднимается сервер, только для unix-like операционных
	// систем Linux, Unix, Mac.
	// Default value: ""
	Socket string `yaml:"Socket" json:"socket" default-value:"-"`

	// SocketMode Файловые разрешения доступа к юникс-сокету.
	// Значение задаётся в восьмеричной системе счисления и не должно превышать 32 бита.
	// Default value: "0666"
	SocketMode string `yaml:"SocketMode" json:"socket_mode" default-value:"0666"`

	// Mode Режим открытия сокета, возможные значения: tcp, tcp4, tcp6, unix, unixpacket, socket, systemd.
	// udp, udp4, udp6 - Сервер поднимается на указанном Host:Port;
	// tcp, tcp4, tcp6 - Сервер поднимается на указанном Host:Port;
	// unix, unixpacket - Сервер поднимается на указанном unix/unixpacket;
	// socket  - Сервер поднимается на socket, только для unix-like операционных систем. Параметры Host:Port
	//           игнорируются, используется только путь к сокету;
	// systemd - Порт или сокет открывает systemd и передаёт слушателя порта через файловый дескриптор сервису,
	//           запущенному от пользователя без права открытия привилегированных портов. Максимально удобный
	//           способ при использовании правильного безопасно настроенного linux сервера.
	//           Более подробно можно посмотреть в документации man systemd.socket(5);
	// Default value: "tcp"
	Mode string `yaml:"Mode" json:"mode" default-value:"tcp"`

	// TLSPublicKeyPEM Путь и имя файла содержащего публичный ключ (сертификат) в PEM формате, включая CA
	// сертификаты всех промежуточных центров сертификации, если ими подписан ключ.
	// Применяется только для TCP соединений, для UDP не используется.
	// Default value: ""
	TLSPublicKeyPEM string `yaml:"TLSPublicKeyPEM" json:"tls_public_key_pem"`

	// TLSPrivateKeyPEM Путь и имя файла содержащего секретный/приватный ключ в PEM формате.
	// Применяется только для TCP соединений, для UDP не используется.
	// Default value: ""
	TLSPrivateKeyPEM string `yaml:"TLSPrivateKeyPEM" json:"tls_private_key_pem"`

	// ProxyProtocol Включение прокси-протокола.
	// Прокси-протокол позволяет серверу получать информацию о подключении клиента, передаваемую через
	// прокси-серверы и средства балансировки нагрузки, такие как Nginx, HAProxy, Amazon Elastic Load
	// Balancer (ELB) и многие другие.
	// Поддерживаются запросы с реализацией прокси протокола версий 1 и 2.
	// PROXY protocol: https://www.haproxy.org/download/2.3/doc/proxy-protocol.txt.
	// Default value: false
	ProxyProtocol bool `yaml:"ProxyProtocol" json:"proxy_protocol"`

	// ProxyProtocolReadHeaderTimeout Максимальное время ожидания получения данных о клиенте через прокси-протокол.
	// Время ожидание используется только при включённом прокси-протоколе.
	// Default value: 0s - no timeout
	ProxyProtocolReadHeaderTimeout time.Duration `yaml:"ProxyProtocolReadHeaderTimeout" json:"proxy_protocol_read_header_timeout"`
}

/**

   Пример конфигурации YAML:


      ## ID Уникальный идентификатор сервера, может быть любым уникальным, в пределах приложения, строковым значением.
      ## Если значение не указано, при запуске сервера, создаётся уникальное временное значение, меняющееся
      ## при каждом запуске.
      ## Default value: ""
      ID: ""

      ## Публичный адрес по которому сервер доступен извне.
      ## Например, если сервер находится за прокси, тут указывается реальный адрес подключения к серверу.
      ## Default value: ""
      Address: !!str "http://localhost/"

      ## IP адрес или имя хоста на котором поднимается сервер, можно указывать 0.0.0.0 для всех ip адресов.
      ## Default value: "0.0.0.0".
      #Host: !!str "[2a03:e2c0:a32::2]"
      #Host: !!str "example.hostname.local"
      Host: !!str "0.0.0.0"

      ## Tcp/ip порт занимаемый сервером.
      ## Default value: 0
      Port: !!int 1080

      ## Юникс сокет, на котором поднимается сервер, только для unix-like операционных систем Linux, Unix, Mac.
      ## Default value: ""
      Socket: !!str "run/example.sock"

      ## Файловые разрешения доступа к юникс-сокету.
      ## Значение задаётся в восьмеричной системе счисления и не должно превышать 32 бита.
      ## Default value: "0666"
      SocketMode: !!str "0666"

      ## Режим открытия сокета, возможные значения: tcp, tcp4, tcp6, unix, unixpacket, socket, systemd.
      ## udp, udp4, udp6 - Сервер поднимается на указанном Host:Port;
      ## tcp, tcp4, tcp6 - Сервер поднимается на указанном Host:Port;
      ## unix, unixpacket - Сервер поднимается на указанном unix/unixpacket;
      ## socket  - Сервер поднимается на socket, только для unix-like операционных систем. Параметры Host:Port
      ##           игнорируются, используется только путь к сокету;
      ## systemd - Порт или сокет открывает systemd и передаёт слушателя порта через файловый дескриптор сервису,
      ##           запущенному от пользователя без права открытия привилегированных портов. Максимально удобный
      ##           способ при использовании правильного безопасно настроенного linux сервера.
      ##           Более подробно можно посмотреть в документации man systemd.socket(5);
      ## Default value: "tcp"
      Mode: !!str "tcp"

      ## Путь и имя файла содержащего публичный ключ (сертификат) в PEM формате, включая CA
      ## сертификаты всех промежуточных центров сертификации, если ими подписан ключ.
      ## Применяется только для TCP соединений, для UDP не используется.
      ## Default value: ""
      TLSPublicKeyPEM: !!str "/etc/application/certificate.pub"

      ## Путь и имя файла содержащего секретный/приватный ключ в PEM формате.
      ## Применяется только для TCP соединений, для UDP не используется.
      ## Default value: ""
      TLSPrivateKeyPEM: !!str "/etc/application/certificate.key"

      ## Включение прокси-протокола.
      ## Прокси-протокол позволяет серверу получать информацию о подключении клиента, передаваемую через
      ## прокси-серверы и средства балансировки нагрузки, такие как Nginx, HAProxy, Amazon Elastic Load
      ## Balancer (ELB) и многие другие.
      ## Поддерживаются запросы с реализацией прокси протокола версий 1 и 2.
      ## PROXY protocol: https://www.haproxy.org/download/2.3/doc/proxy-protocol.txt.
      ## Default value: false
      ProxyProtocol: !!bool false

      ## Максимальное время ожидания получения данных о клиенте через прокси-протокол.
      ## Время ожидание используется только при включённом прокси-протоколе.
      ## Default value: 0s - no timeout
      ProxyProtocolReadHeaderTimeout: 0s


**/
