package main

import (
	"context"
	"log"
	"sync"

	rd "github.com/redis/go-redis/v9"
)

// Mutex wrapper/helper to handle concurrent access to the inner Redis client
type Redis struct {
	client *rd.Client
	mutex  sync.Mutex
}

func RedisDefault() Redis {
	return Redis{
		client: nil,
		mutex:  sync.Mutex{},
	}
}

func (r *Redis) set(client *rd.Client) {
	if r.client != nil {
		panic("Unexpected error: Trying to re-initialize a Redis instance")
	}

	r.client = client
}

// Important: Don't forget to defer .unlock() !
func (r *Redis) get() *rd.Client {
	r.mutex.Lock()

	return r.client
}

// Defer this function when accessing the wrapped inner value with .get() to unlock the Mutex guard
func (r *Redis) unlock() {
	r.mutex.Unlock()
}

func (r *Redis) subscribe(ctx context.Context, channel string) <-chan *rd.Message {
	client := r.get()
	defer r.unlock()
	sub := client.Subscribe(ctx, channel)

	iface, err := sub.Receive(ctx)
	if err != nil {
		// Since it's called once, it should't fail, else it's a Redis issue
		log.Fatalf("Failed to subscribe to channel %s\n", channel)
	}

	switch iface.(type) {
	case *rd.Subscription:
		break
	default: // Message, Pong or Error
		log.Fatalf("Unexpected result (%s) on channel (%s) subscription. It shouldn't return anything else than SUBSCRIBE\n", iface, channel)
	}

	return sub.Channel()
}

func (r *Redis) publish(ctx context.Context, channel string, message string) error {
	client := r.get()
	defer r.unlock()

	if err := client.Publish(ctx, channel, message).Err(); err != nil {
		return err
	}

	return nil
}
