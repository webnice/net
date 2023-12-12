package net

// Все ошибки определены как константы.
const (
	cAlreadyRunning                = "Сервер уже запущен."
	cNoConfiguration               = "Конфигурация сервера отсутствует либо равна nil."
	cListenSystemdPID              = "Переменная окружения LISTEN_PID пустая, либо содержит не верное значение."
	cListenSystemdFDS              = "Переменная окружения LISTEN_FDS пустая, либо содержит не верное значение."
	cListenSystemdNotFound         = "Получение сокета systemd по имени, имя не найдено."
	cListenSystemdQuantityNotMatch = "Полученное количество LISTEN_FDS не соответствует переданному LISTEN_FDNAMES."
	cTLSIsNil                      = "Конфигурация TLS сервера пустая."
	cServerHandlerIsNotSet         = "Не установлен обработчик основной функции TCP сервера."
	cServerHandlerUdpIsNotSet      = "Не установлен обработчик основной функции UDP сервера."
)

// Константы указываются в объектах в качестве фиксированного адреса на протяжении всего времени работы приложения.
// Ошибка с ошибкой могут сравниваться по содержимому, по адресу и т.д.
var (
	errSingleton                     = &Error{}
	errAlreadyRunning                = err(cAlreadyRunning)
	errNoConfiguration               = err(cNoConfiguration)
	errListenSystemdPID              = err(cListenSystemdPID)
	errListenSystemdFDS              = err(cListenSystemdFDS)
	errListenSystemdNotFound         = err(cListenSystemdNotFound)
	errListenSystemdQuantityNotMatch = err(cListenSystemdQuantityNotMatch)
	errTLSIsNil                      = err(cTLSIsNil)
	errServerHandlerIsNotSet         = err(cServerHandlerIsNotSet)
	errServerHandlerUdpIsNotSet      = err(cServerHandlerUdpIsNotSet)
)

type (
	// Error object of package
	Error struct{}
	err   string
)

// Error The error built-in interface implementation
func (e err) Error() string { return string(e) }

// Errors Справочник ошибок.
func Errors() *Error { return errSingleton }

// ОШИБКИ.

// AlreadyRunning Сервер уже запущен.
func (e *Error) AlreadyRunning() error { return &errAlreadyRunning }

// NoConfiguration Конфигурация сервера отсутствует либо равна nil.
func (e *Error) NoConfiguration() error { return &errNoConfiguration }

// ListenSystemdPID Переменная окружения LISTEN_PID пустая, либо содержит не верное значение.
func (e *Error) ListenSystemdPID() error { return &errListenSystemdPID }

// ListenSystemdFDS Переменная окружения LISTEN_FDS пустая, либо содержит не верное значение.
func (e *Error) ListenSystemdFDS() error { return &errListenSystemdFDS }

// ListenSystemdNotFound Получение сокета systemd по имени, имя не найдено.
func (e *Error) ListenSystemdNotFound() error { return &errListenSystemdNotFound }

// ListenSystemdQuantityNotMatch Полученное количество LISTEN_FDS не соответствует переданному LISTEN_FDNAMES.
func (e *Error) ListenSystemdQuantityNotMatch() error { return &errListenSystemdQuantityNotMatch }

// TLSIsNil Конфигурация TLS сервера пустая.
func (e *Error) TLSIsNil() error { return &errTLSIsNil }

// ServerHandlerIsNotSet Не установлен обработчик основной функции TCP сервера.
func (e *Error) ServerHandlerIsNotSet() error { return &errServerHandlerIsNotSet }

// ServerHandlerUdpIsNotSet Не установлен обработчик основной функции UDP сервера.
func (e *Error) ServerHandlerUdpIsNotSet() error { return &errServerHandlerUdpIsNotSet }
