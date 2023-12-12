package net

// Configuration Структура конфигурации TCP/IP или UDP сервера.
type Configuration struct {
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

	// SocketMode Файловые разрешения доступа к юникс сокету, при его использовании.
	// Default value: "640"
	SocketMode uint32 `yaml:"SocketMode" json:"socket_mode" default-value:"0640"`

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
}

/**

   Пример конфигурации YAML:

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

      ## Файловые разрешения доступа к юникс сокету, при его использовании.
      ## Default value: "640"
      SocketMode: !!int 640

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


**/
