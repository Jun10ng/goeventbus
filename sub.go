package goeventbus

// Sub contains subscriber's informations ,like channel etc.
type Sub struct {
	out chan interface{}
}

func NewSub() Sub {
	return Sub{
		out: make(chan interface{}),
	}
}

func (s *Sub) receive(msg interface{}) {
	s.out <- msg
}

// Out return a msg which out from Sub.out channel
func (s *Sub) Out() (msg interface{}) {
	msg = <-s.out
	return
}
