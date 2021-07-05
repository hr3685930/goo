package pool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

type Pool struct {
	clients    chan *Client
	connCnt    int32
	cap        int32
	idleDur    time.Duration
	maxLifeDur time.Duration
	timeout    time.Duration
	factor     Factor
	lock       sync.RWMutex
	mode       int
}

type Client struct {
	*grpc.ClientConn
	timeUsed time.Time
	timeInit time.Time
	pool     *Pool
}

type Factor func() (*grpc.ClientConn, error)

var (
	ErrPoolInit     = errors.New("Pool init occurred error")
	ErrGetTimeout   = errors.New("Getting connection client timeout from pool")
	ErrDialConn     = errors.New("Dialing connection occurred error")
	ErrPoolIsClosed = errors.New("Pool is closed")
)

const (
	PoolGetModeStrict = iota
	PoolGetModeLoose
)

func Default(factor Factor, init, cap int32) (*Pool, error) {
	return Init(factor, init, cap, 10*time.Second, 60*time.Second, 10*time.Second, PoolGetModeLoose)
}

func Init(factor Factor, init, cap int32, idleDur, maxLifeDur, timeout time.Duration, mode int) (*Pool, error) {
	if factor == nil {
		return nil, ErrPoolInit
	}
	if init < 0 || cap <= 0 || idleDur < 0 || maxLifeDur < 0 {
		return nil, ErrPoolInit
	}
	// init pool
	if init > cap {
		init = cap
	}
	pool := &Pool{
		clients:    make(chan *Client, cap),
		cap:        cap,
		idleDur:    idleDur,
		maxLifeDur: maxLifeDur,
		timeout:    timeout,
		factor:     factor,
		mode:       mode,
	}
	// init client
	for i := int32(0); i < init; i++ {
		client, err := pool.createClient()
		if err != nil {
			return nil, ErrPoolInit
		}
		pool.clients <- client
	}
	return pool, nil
}

func (pool *Pool) createClient() (*Client, error) {
	conn, err := pool.factor()
	if err != nil {
		return nil, ErrPoolInit
	}
	now := time.Now()
	client := &Client{
		ClientConn: conn,
		timeUsed:   now,
		timeInit:   now,
		pool:       pool,
	}
	atomic.AddInt32(&pool.connCnt, 1)
	return client, nil
}

func (pool *Pool) Get(ctx context.Context) (*Client, error) {
	if pool.IsClose() {
		return nil, ErrPoolIsClosed
	}

	var client *Client
	now := time.Now()
	select {
	case <-ctx.Done():
		if pool.mode == PoolGetModeStrict {
			pool.lock.Lock()
			defer pool.lock.Unlock()

			var err error
			if pool.connCnt >= int32(pool.cap) {
				err = ErrGetTimeout
			} else {
				client, err = pool.createClient()
			}
			return client, err
		}
	case client = <-pool.clients:
		if client != nil && pool.idleDur > 0 && client.timeUsed.Add(pool.idleDur).After(now) {
			client.timeUsed = now
			return client, nil
		}
	}
	if client != nil {
		client.Destory()
	}
	client, err := pool.createClient()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (pool *Pool) Close() {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	if pool.IsClose() {
		return
	}

	clients := pool.clients
	pool.clients = nil

	go func() {
		for {
			select {
			case client := <-clients:
				if client != nil {
					client.Destory()
				}
			case <-time.Tick(pool.timeout):
				if len(clients) <= 0 {
					close(clients)
					break
				}
			}
		}
	}()
}

func (pool *Pool) IsClose() bool {
	return pool == nil || pool.clients == nil
}

func (pool *Pool) Size() int {
	pool.lock.RLock()
	defer pool.lock.RUnlock()

	return len(pool.clients)
}

func (pool *Pool) ConnCnt() int32 {
	return pool.connCnt
}

func (client *Client) Close() {
	go func() {
		pool := client.pool
		now := time.Now()
		if pool.IsClose() {
			client.Destory()
			return
		}
		if pool.maxLifeDur > 0 && client.timeInit.Add(pool.maxLifeDur).Before(now) {
			client.Destory()
			return
		}
		if client.ClientConn == nil {
			return
		}
		client.timeUsed = now
		client.pool.clients <- client
	}()
}

func (client *Client) Destory() {
	if client.ClientConn != nil {
        _ = client.ClientConn.Close()
		atomic.AddInt32(&client.pool.connCnt, -1)
	}
	client.ClientConn = nil
	client.pool = nil
}

func (client *Client) TimeInit() time.Time {
	return client.timeInit
}

func (client *Client) TimeUsed() time.Time {
	return client.timeUsed
}
