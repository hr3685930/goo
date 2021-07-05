package rabbitmq

import (
    "goo/pkg/queue"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/golang-module/carbon"
    "github.com/streadway/amqp"
    "reflect"
    "strings"
    "time"
)

type RabbitMQ struct {
    conn         *amqp.Connection
    channel      *amqp.Channel
    QueueName    string
    Exchange     string
    Key          string
    DLKQueueName string
    DLKExchange  string
    DLKKey       string
    MQUrl        string
    Prefix       string
}

func NewRabbitMQ(user, pass, host, port, vhost, prefix string) queue.Queue {
    mqUrl := "amqp://" + user + ":" + pass + "@" + host + ":" + port + vhost
    return &RabbitMQ{MQUrl: mqUrl, Prefix: prefix}
}

// new conn
func (r *RabbitMQ) Connect() error {
    conn, err := amqp.Dial(r.MQUrl)
    if err != nil {
        return errors.New(fmt.Sprintf("ampq connect error %s", err))
    }
    channel, err := conn.Channel()
    if err != nil {
        return errors.New(fmt.Sprintf("ampq channel error %s", err))
    }

    r.conn = conn
    r.channel = channel
    return nil
}

// producer connect
func (r *RabbitMQ) ProducerConnect() queue.Queue {
    channel, err := r.conn.Channel()
    if err != nil {
        amqperr := err.(*amqp.Error)
        if amqperr.Code == amqp.ChannelError {
            _ = r.Connect()
            return &RabbitMQ{MQUrl: r.MQUrl, channel: r.channel, conn: r.conn, Prefix: r.Prefix}
        }
        panic(fmt.Sprintf("ampq channel error %s", err))
    }

    r.channel = channel
    return &RabbitMQ{MQUrl: r.MQUrl, channel: channel, conn: r.conn, Prefix: r.Prefix}
}

// producer connect
func (r *RabbitMQ) ConsumerConnect() queue.Queue {
    channel, err := r.conn.Channel()
    if err != nil {
        amqperr := err.(*amqp.Error)
        if amqperr.Code == amqp.ChannelError {
            _ = r.Connect()
            return &RabbitMQ{MQUrl: r.MQUrl, channel: r.channel, conn: r.conn, Prefix: r.Prefix}
        }
        panic(fmt.Sprintf("ampq channel error %s", err))
    }

    r.channel = channel
    return &RabbitMQ{MQUrl: r.MQUrl, channel: channel, conn: r.conn, Prefix: r.Prefix}
}

func (r *RabbitMQ) Close() {
    _ = r.channel.Close()
    _ = r.conn.Close()
}

//error handler
func (r *RabbitMQ) failOnErr(err error, message string) {
    if err != nil {
        panic(fmt.Sprintf("%s:%s", message, err))
    }
}

// abstract topic
func (r *RabbitMQ) Topic(topic string) {
    r.SetExchange(topic)
}

//exchange
func (r *RabbitMQ) SetExchange(exchangeName string) {
    r.Exchange = r.Prefix + "." + exchangeName
}

//queue and key
func (r *RabbitMQ) SetQueueKey(queue string) {
    r.QueueName = strings.ToLower(strings.Replace(r.Prefix+"."+queue, ".", "_", -1))
    r.Key = strings.ToLower(r.Prefix + "." + queue)
}

//DLK Exchange
func (r *RabbitMQ) SetDLKExchange(DLKexchange string) {
    r.DLKExchange = DLKexchange
}

//DLK queue and key
func (r *RabbitMQ) SetDLKQueueKey(queue string) {
    r.DLKQueueName = strings.ToLower(strings.Replace(queue, ".", "_", -1))
    r.DLKKey = strings.ToLower(queue)
}

//Producer
func (r *RabbitMQ) Producer(job queue.JobBase) {
    if r.Exchange == "" {
        panic("not exchange")
    }

    queueName := reflect.TypeOf(job).String()
    r.SetQueueKey(queueName[1:])

    message, err := json.Marshal(job)
    r.failOnErr(err, "Umarshal failed")

    //create Exchange
    err = r.channel.ExchangeDeclare(r.Exchange, "topic", true, false, false, false, nil)
    r.failOnErr(err, "Failed to declare an exchange")

    //create queue
    q, err := r.channel.QueueDeclare(r.QueueName, true, false, false, false, nil)
    r.failOnErr(err, "Failed to declare a queue")

    //bind
    err = r.channel.QueueBind(q.Name, r.Key, r.Exchange, false, nil)
    r.failOnErr(err, "Failed to declare a queue bind")

    //pub
    err = r.channel.Publish(r.Exchange, r.Key, false, false,
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         []byte(message),
            DeliveryMode: 2,
            Timestamp:    time.Now(),
        })
    //
    _ = r.channel.Close()
}

//DLK
func (r *RabbitMQ) DLK(base queue.JobBase, sleep, retry int32) {
    DLKexchange := "delay." + r.Exchange
    DLKkey := "delay." + r.Key

    //create DLK Exchange
    if r.DLKExchange == "" {
        r.SetDLKExchange(DLKexchange)
        err := r.channel.ExchangeDeclare(r.DLKExchange, "topic", true, false, false, false, nil)
        r.failOnErr(err, "Failed to declare an delay_exchange")
    }

    //create DLK queue
    if r.DLKQueueName == "" {
        r.SetDLKQueueKey(DLKkey)
        args := make(amqp.Table)
        args["x-dead-letter-exchange"] = r.Exchange
        args["x-dead-letter-routing-key"] = r.Key
        args["x-message-ttl"] = 1000 * sleep
        q, err := r.channel.QueueDeclare(r.DLKQueueName, true, false, false, false, args)
        r.failOnErr(err, "Failed to declare a delay_queue")
        //bind DLK
        err = r.channel.QueueBind(q.Name, r.DLKKey, r.DLKExchange, false, nil)
        r.failOnErr(err, "Failed to declare a delay_queue bind")
    }

    message, err := json.Marshal(base)
    r.failOnErr(err, "Umarshal failed")
    //pub
    header := make(map[string]interface{}, 1)
    header["retry_num"] = retry + int32(1)
    err = r.channel.Publish(r.DLKExchange, r.DLKKey, false, false,
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         []byte(message),
            DeliveryMode: 1,
            Timestamp:    time.Now(),
            Headers:      header,
        })

}

// Consumer
func (r *RabbitMQ) Consumer(base queue.JobBase, sleep, retry int32) {
    queueName := reflect.TypeOf(base).String()
    r.SetQueueKey(queueName[1:])

    q, err := r.channel.QueueDeclare(r.QueueName, true, false, false, false, nil)
    r.failOnErr(err, "Failed to declare a queue")
    //Qos
    err = r.channel.Qos(1, 0, false)
    //Consumer
    messages, err := r.channel.Consume(q.Name, "", false, false, false, false, nil)
    forever := make(chan bool)

    go func() {
        for d := range messages {
            err := json.Unmarshal(d.Body, base)
            if err != nil {
                r.failOnErr(err, "Unmarshal error")
            }

            err = base.Handler()
            queueErr := err.(*queue.Error)
            if queueErr != nil {
                retryNum, ok := d.Headers["retry_num"].(int32)
                if !ok {
                    retryNum = int32(0)
                }

                if retryNum < retry {
                    r.DLK(base, sleep, retryNum)
                } else {
                    r.ExportErr(queueErr, d)
                }
            }

            err = d.Ack(false)
            if err != nil {
                r.ExportErr(queue.Err(err), d)
            }
        }
    }()

    <-forever
}

func (r *RabbitMQ) ExportErr(err error, d amqp.Delivery) {
    e := err.(*queue.Error)
    r.Err(queue.FailedJobs{
        Connection: "rabbitmq",
        Queue:      r.QueueName,
        Message:    string(d.Body),
        Exception:  err.Error(),
        Stack:      e.Stack(),
        FiledAt:    carbon.Now(),
    })
}

func (r *RabbitMQ) Err(failed queue.FailedJobs) {
    queue.ErrJob <- failed
}
