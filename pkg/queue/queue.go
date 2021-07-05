package queue

import "runtime/debug"

type JobBase interface {
	Handler() (queueErr *Error)
}

type Queue interface {
	// connect
	Connect() error
	// Producer connect
	ProducerConnect() Queue
	// Consumer connect
	ConsumerConnect() Queue
	// abstract topic
	Topic(topic string)

	Producer(job JobBase)
	// sleep retry
	Consumer(job JobBase, sleep, retry int32)
	// report
	Err(failed FailedJobs)

	Close()
}

type Error struct {
	s     string
	stack string
}

func (qe *Error) Error() string {
	return qe.s
}

func Err(err error) *Error {
	return &Error{
		s:     err.Error(),
		stack: string(debug.Stack()),
	}
}

func (qe *Error) Stack() string {
	return qe.stack
}
