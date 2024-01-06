package rate

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type SlidingLog struct {
	conn redis.UniversalClient

	// PipelineFn is exposed for black box tests in the same package.
	PipelineFn func(context.Context, func(redis.Pipeliner) error) ([]redis.Cmder, error)
}

type StorageHandler interface {
	As(interface{}) error
}

// NewSlidingLog creates a new SlidingLog instance with a storage.Handler. In case
// the storage is offline, it's expected to return nil and an error to handle.
func NewSlidingLog(cluster StorageHandler, pipeline bool) (*SlidingLog, error) {
	conn := new(redis.UniversalClient)
	if err := cluster.As(conn); err != nil {
		return nil, err
	}

	return NewSlidingLogRedis(*conn, pipeline), nil
}

// NewSlidingLogRedis creates a new SlidingLog instance with a redis.UniversalClient.
func NewSlidingLogRedis(conn redis.UniversalClient, pipeline bool) *SlidingLog {
	r := &SlidingLog{
		conn: conn,
	}
	r.PipelineFn = r.conn.TxPipelined
	if pipeline {
		r.PipelineFn = r.conn.Pipelined
	}
	return r
}

// SetCount will trim the rolling window log, add an item and return the count of the items in a window before the add.
func (r *SlidingLog) SetCount(ctx context.Context, now time.Time, keyName string, per int64) (int64, error) {
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Second)

	var res *redis.IntCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(ctx, keyName, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		res = pipe.ZCard(ctx, keyName)

		element := redis.Z{
			Score: float64(now.UnixNano()),
		}

		element.Member = strconv.Itoa(int(now.UnixNano()))

		pipe.ZAdd(ctx, keyName, &element)
		pipe.Expire(ctx, keyName, time.Duration(per)*time.Second)

		return nil
	}

	if _, err := r.PipelineFn(ctx, pipeFn); err != nil {
		return 0, err
	}

	return res.Result()
}

// GetCount will trim the rolling window log and return the count of items remaining.
func (r *SlidingLog) GetCount(ctx context.Context, now time.Time, keyName string, per int64) (int64, error) {
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Second)

	var res *redis.IntCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(ctx, keyName, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		res = pipe.ZCard(ctx, keyName)

		return nil
	}

	if _, err := r.PipelineFn(ctx, pipeFn); err != nil {
		return 0, err
	}

	return res.Result()
}

// Set will append to a sorted set in redis and return the contents of the window as a slice.
func (r *SlidingLog) Set(ctx context.Context, now time.Time, keyName string, per int64) ([]string, error) {
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Second)

	var res *redis.StringSliceCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(ctx, keyName, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		res = pipe.ZRange(ctx, keyName, 0, -1)

		element := redis.Z{
			Score: float64(now.UnixNano()),
		}

		element.Member = strconv.Itoa(int(now.UnixNano()))

		pipe.ZAdd(ctx, keyName, &element)
		pipe.Expire(ctx, keyName, time.Duration(per)*time.Second)

		return nil
	}

	if _, err := r.PipelineFn(ctx, pipeFn); err != nil {
		return nil, err
	}

	return res.Result()
}

// Get will trim the rolling window log and return the contents of the window as a slice.
func (r *SlidingLog) Get(ctx context.Context, now time.Time, keyName string, per int64) ([]string, error) {
	onePeriodAgo := now.Add(time.Duration(-1*per) * time.Second)

	var res *redis.StringSliceCmd

	pipeFn := func(pipe redis.Pipeliner) error {
		pipe.ZRemRangeByScore(ctx, keyName, "-inf", strconv.Itoa(int(onePeriodAgo.UnixNano())))
		res = pipe.ZRange(ctx, keyName, 0, -1)

		return nil
	}

	if _, err := r.PipelineFn(ctx, pipeFn); err != nil {
		return nil, err
	}

	return res.Result()
}
