package interfaces

type IServer interface {
	// GetAddr returns the server address.
	GetAddr() string
	// GetPort returns the server port.
	GetPort() string
}
