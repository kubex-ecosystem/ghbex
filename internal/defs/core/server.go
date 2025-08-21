package core

import "github.com/rafa-mori/ghbex/internal/defs/interfaces"

type Server struct {
	Addr string `yaml:"addr" json:"addr"`
	Port string `yaml:"port" json:"port"`
}

func NewServerType(addr, port string) *Server {
	return &Server{
		Addr: addr,
		Port: port,
	}
}

func NewServer(addr, port string) interfaces.IServer { return NewServerType(addr, port) }
func (s *Server) GetAddr() string                    { return s.Addr }
func (s *Server) GetPort() string                    { return s.Port }
