package service

type Request struct {
	Addresse string
	Message  *Message
}

func NewRequest(addressee string, msg *Message) *Request {
	return &Request{
		Addresse: addressee,
		Message:  msg,
	}
}

func (r *Request) Serialize() []string {
	res := []string{
		r.Addresse,
		COMMAND_REQUEST,
		r.Message.Sender,
		r.Message.CorrelationId,
		r.Message.ServiceType,
		r.Message.ServiceInstance,
	}
	res = append(res, r.Message.Payload...)

	return res
}
