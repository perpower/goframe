// 定时任务,为了方便后续统一维护和二次扩展，故此将所有的方法重新二次封装一遍
package pcron

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gtimer"
)

type Entry = *gcron.Entry

// JobFunc is the timing called job function in timer.
type JobFunc = gtimer.JobFunc

// Add adds a timed task to default cron object.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func Add(ctx context.Context, pattern string, job JobFunc, name ...string) (Entry, error) {
	return gcron.Add(ctx, pattern, job, name...)
}

// AddSingleton adds a singleton timed task, to default cron object.
// A singleton timed task is that can only be running one single instance at the same time.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func AddSingleton(ctx context.Context, pattern string, job JobFunc, name ...string) (Entry, error) {
	return gcron.AddSingleton(ctx, pattern, job, name...)
}

// AddOnce adds a timed task which can be run only once, to default cron object.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func AddOnce(ctx context.Context, pattern string, job JobFunc, name ...string) (Entry, error) {
	return gcron.AddOnce(ctx, pattern, job, name...)
}

// AddTimes adds a timed task which can be run specified times, to default cron object.
// A unique `name` can be bound with the timed task.
// It returns and error if the `name` is already used.
func AddTimes(ctx context.Context, pattern string, times int, job JobFunc, name ...string) (Entry, error) {
	return gcron.AddTimes(ctx, pattern, times, job, name...)
}

// DelayAdd adds a timed task to default cron object after `delay` time.
func DelayAdd(ctx context.Context, delay time.Duration, pattern string, job JobFunc, name ...string) {
	gcron.DelayAdd(ctx, delay, pattern, job, name...)
}

// DelayAddSingleton adds a singleton timed task after `delay` time to default cron object.
func DelayAddSingleton(ctx context.Context, delay time.Duration, pattern string, job JobFunc, name ...string) {
	gcron.DelayAddSingleton(ctx, delay, pattern, job, name...)
}

// DelayAddOnce adds a timed task after `delay` time to default cron object.
// This timed task can be run only once.
func DelayAddOnce(ctx context.Context, delay time.Duration, pattern string, job JobFunc, name ...string) {
	gcron.DelayAddOnce(ctx, delay, pattern, job, name...)
}

// DelayAddTimes adds a timed task after `delay` time to default cron object.
// This timed task can be run specified times.
func DelayAddTimes(ctx context.Context, delay time.Duration, pattern string, times int, job JobFunc, name ...string) {
	gcron.DelayAddTimes(ctx, delay, pattern, times, job, name...)
}

// Search returns a scheduled task with the specified `name`.
// It returns nil if no found.
func Search(name string) Entry {
	return gcron.Search(name)
}

// Remove deletes scheduled task which named `name`.
func Remove(name string) {
	gcron.Remove(name)
}

// Size returns the size of the timed tasks of default cron.
func Size() int {
	return gcron.Size()
}

// Entries return all timed tasks as slice.
func Entries() []Entry {
	return gcron.Entries()
}

// Start starts running the specified timed task named `name`.
// If no`name` specified, it starts the entire cron.
func Start(name ...string) {
	gcron.Start(name...)
}

// Stop stops running the specified timed task named `name`.
// If no`name` specified, it stops the entire cron.
func Stop(name ...string) {
	gcron.Stop(name...)
}
