package service

func (b *Broker) DisconnectWorker(sender string, msg string) {
	b.Socket.SendMessage(
		sender, COMMAND_DISCONNECT, "broker", []string{msg},
	)
}
