package monitor

func (s *Service) SetUDPConn() {
	conn, err := s.Conn()
	if err != nil {
		log.WithError(err).Error("Failed to start UDP listener")
		return
	}

	s.conn = conn
}

func (s *Service) SendUDPMessage(msg string) error {
	if s.conn == nil {
		s.SetUDPConn()
	}

	_, err := s.conn.Write([]byte(msg + "\n"))
	if err != nil {
		log.WithError(err).Error("Failed to send UDP message")
		return err
	}

	return nil
}
