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

// ErrAlreadyRunning Сервер уже запущен.
func ErrAlreadyRunning() error { return &errAlreadyRunning }

// ErrNoConfiguration Конфигурация сервера отсутствует либо равна nil.
func ErrNoConfiguration() error { return &errNoConfiguration }

// ErrListenSystemdPID Переменная окружения LISTEN_PID пустая, либо содержит не верное значение.
func ErrListenSystemdPID() error { return &errListenSystemdPID }

// ErrListenSystemdFDS Переменная окружения LISTEN_FDS пустая, либо содержит не верное значение.
func ErrListenSystemdFDS() error { return &errListenSystemdFDS }

// ErrListenSystemdNotFound Получение сокета systemd по имени, имя не найдено.
func ErrListenSystemdNotFound() error { return &errListenSystemdNotFound }

// ErrListenSystemdQuantityNotMatch Полученное количество LISTEN_FDS не соответствует переданному LISTEN_FDNAMES.
func ErrListenSystemdQuantityNotMatch() error { return &errListenSystemdQuantityNotMatch }

// ErrTLSIsNil Конфигурация TLS сервера пустая.
func ErrTLSIsNil() error { return &errTLSIsNil }

// ErrServerHandlerIsNotSet Не установлен обработчик основной функции TCP сервера.
func ErrServerHandlerIsNotSet() error { return &errServerHandlerIsNotSet }

// ErrServerHandlerUdpIsNotSet Не установлен обработчик основной функции UDP сервера.
func ErrServerHandlerUdpIsNotSet() error { return &errServerHandlerUdpIsNotSet }
