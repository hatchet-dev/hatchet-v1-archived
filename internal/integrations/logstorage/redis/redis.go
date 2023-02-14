package redis

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/hatchet-dev/hatchet/internal/integrations/logstorage"
)

type RedisLogStorageManager struct {
	client *redis.Client
}

type InitOpts struct {
	RedisHost     string
	RedisPort     string
	RedisUsername string
	RedisPassword string
	RedisDB       int
}

func NewRedisLogStorageManager(opts *InitOpts) (*RedisLogStorageManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", opts.RedisHost, opts.RedisPort),
		Username: opts.RedisUsername,
		Password: opts.RedisPassword,
		DB:       opts.RedisDB,
	})

	_, err := client.Ping(context.Background()).Result()
	return &RedisLogStorageManager{client}, err
}

func (c *RedisLogStorageManager) GetID() string {
	return "redis"
}

func (c *RedisLogStorageManager) PushLogLine(ctx context.Context, path string, log []byte) error {
	_, err := c.client.XAdd(context.TODO(), &redis.XAddArgs{
		Stream: path,
		ID:     "*",
		Values: map[string]interface{}{
			"log": log,
		},
	}).Result()

	return err
}

func (c *RedisLogStorageManager) StreamLogs(ctx context.Context, opts *logstorage.LogGetOpts, w io.WriteCloser) error {
	if doesExist, err := c.doesStreamExist(opts.Path); !doesExist || err != nil {
		return err
	}

	errorchan := make(chan error)
	redisCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		wg.Wait()
		close(errorchan)
	}()

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			errorchan <- nil
		case <-redisCtx.Done():
			errorchan <- nil
		}
	}()

	go func() {
		defer wg.Done()

		// check intermittently that the stream still exists -- it may have been
		// cleaned up automatically
		failedCount := 0
		for {
			x, err := c.client.Exists(
				context.Background(),
				opts.Path,
			).Result()

			// if the stream does not exist, increment the failed counter
			if x == 0 || err != nil {
				failedCount++
			} else {
				failedCount = 0
			}

			if failedCount >= 2 {
				errorchan <- nil
				return
			}

			// wait 5 seconds in between pings
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		defer wg.Done()

		lastID := "0-0"

		for {
			if redisCtx.Err() != nil {
				errorchan <- nil
				return
			}

			xstream, err := c.client.XRead(
				redisCtx,
				&redis.XReadArgs{
					Streams: []string{opts.Path, lastID},
					Block:   0,
					Count:   int64(opts.Count),
				},
			).Result()

			if err != nil {
				errorchan <- err
				return
			}

			messages := xstream[0].Messages

			lastID, err = c.parseAndWriteXMessages(messages, w)

			if err != nil {
				errorchan <- err
				return
			}
		}
	}()

	var err error

	for err = range errorchan {
		cancel()
	}

	return err
}

func (c *RedisLogStorageManager) ReadLogs(ctx context.Context, opts *logstorage.LogGetOpts, w io.WriteCloser) error {
	if doesExist, err := c.doesStreamExist(opts.Path); !doesExist || err != nil {
		return err
	}

	xstream, err := c.client.XRead(
		ctx,
		&redis.XReadArgs{
			Streams: []string{opts.Path, "0-0"},
			Block:   0,
			Count:   int64(opts.Count),
		},
	).Result()

	if err != nil {
		return err
	}

	messages := xstream[0].Messages

	_, err = c.parseAndWriteXMessages(messages, w)

	if err != nil {
		return err
	}

	return nil
}

func (c *RedisLogStorageManager) ClearLogs(ctx context.Context, path string) error {
	_, err := c.client.Del(
		context.Background(),
		path,
	).Result()

	if err != nil {
		return err
	}

	return nil
}

func (c *RedisLogStorageManager) doesStreamExist(path string) (bool, error) {
	x, err := c.client.Exists(
		context.Background(),
		path,
	).Result()

	if err != nil {
		return false, err
	}

	return x != 0, nil
}

func (c *RedisLogStorageManager) parseAndWriteXMessages(messages []redis.XMessage, w io.WriteCloser) (string, error) {
	lastID := messages[len(messages)-1].ID

	for _, msg := range messages {
		dataInter, ok := msg.Values["log"]

		if !ok {
			continue
		}

		dataString, ok := dataInter.(string)

		if !ok {
			continue
		}

		_, err := w.Write([]byte(dataString))

		if err != nil {
			return "", err
		}
	}

	return lastID, nil
}
