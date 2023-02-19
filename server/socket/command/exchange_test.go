package command

import "testing"

func TestExchange_GetSubscriberChannel(t *testing.T) {
	e := NewExchange()
	commandToSend := &Command{Command: 1, SensorType: 1, ForId: 1}

	subscriber := e.GetSubscriberChannel()
	e.Dispatch(commandToSend)

	receivedCommand := <-subscriber
	if commandToSend != receivedCommand {
		t.Errorf("Expected %v, got %v", commandToSend, receivedCommand)
	}
}

func TestExchange_Unsubscribe(t *testing.T) {
	e := NewExchange()
	commandToSend := &Command{Command: 1, SensorType: 1, ForId: 1}

	subscriber := e.GetSubscriberChannel()
	subscriber2 := e.GetSubscriberChannel()
	subscriber3 := e.GetSubscriberChannel()
	e.Unsubscribe(subscriber)
	e.Dispatch(commandToSend)

	if len(e.observers) != 2 {
		t.Errorf("Subscriber has not been closed")
	}

	if commandToSend != <-subscriber3 || commandToSend != <-subscriber2 {
		t.Errorf("wrong subscriber has been removed from subscriber list")
	}
}
