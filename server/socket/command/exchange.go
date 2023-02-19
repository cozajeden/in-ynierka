package command

const defaultChannelBuffer = 10

type Exchange struct {
	observers []chan *Command
}

func (e *Exchange) Dispatch(command *Command) {
	for _, o := range e.observers {
		o <- command
	}
}

func (e *Exchange) GetSubscriberChannel() chan *Command {
	ch := make(chan *Command, defaultChannelBuffer)
	e.observers = append(e.observers, ch)
	return ch
}

func (e *Exchange) Unsubscribe(subscriber chan *Command) {
	for i := range e.observers {
		if e.observers[i] == subscriber {
			e.observers = append(e.observers[:i], e.observers[i+1:]...)
			return
		}
	}
}

func NewExchange() *Exchange {
	return &Exchange{
		make([]chan *Command, 0),
	}
}
