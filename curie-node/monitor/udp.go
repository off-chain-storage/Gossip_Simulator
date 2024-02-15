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

	log.Info("1 Sent UDP message: ", msg)

	_, err := s.conn.Write([]byte(msg + "\n"))
	if err != nil {
		log.WithError(err).Error("Failed to send UDP message")
		return err
	}

	log.Info("2 Sent UDP message: ", msg)

	return nil
}
