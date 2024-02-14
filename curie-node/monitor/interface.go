package monitor

type Monitor interface {
	SendUDPMessage(msg string) error
}
