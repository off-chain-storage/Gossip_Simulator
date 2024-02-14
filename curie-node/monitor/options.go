package monitor

import "net"

func (s *Service) buildUDPAddr() error {
	cfg := s.cfg

	udpServer, err := net.ResolveUDPAddr("udp4", cfg.UDPAddr)
	if err != nil {
		log.WithError(err).Error("Failed to resolve UDP address")
		return err
	}

	s.udpServer = udpServer
	return nil
}

func (s *Service) Conn() (*net.UDPConn, error) {
	conn, err := net.DialUDP("udp4", nil, s.udpServer)
	if err != nil {
		log.WithError(err).Error("Failed to start UDP listener")
		return nil, err
	}

	return conn, nil
}
