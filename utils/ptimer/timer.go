// 定时器,为了方便后续统一维护和二次扩展，故此将所有的方法重新二次封装一遍
package ptimer

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/os/gtimer"
)

// Entry is the timing job.
type Entry = *gtimer.Entry

// Timer is the timer manager, which uses ticks to calculate the timing interval.
type Timer = gtimer.Timer

// JobFunc is the timing called job function in timer.
type JobFunc = gtimer.JobFunc

// DefaultOptions creates and returns a default options object for Timer creation.
func DefaultOptions() gtimer.TimerOptions {
	return gtimer.DefaultOptions()
}

// SetTimeout runs the job once after duration of `delay`.
// It is like the one in javascript.
func SetTimeout(ctx context.Context, delay time.Duration, job JobFunc) {
	gtimer.SetTimeout(ctx, delay, job)
}

// SetInterval runs the job every duration of `delay`.
// It is like the one in javascript.
func SetInterval(ctx context.Context, interval time.Duration, job JobFunc) {
	gtimer.SetInterval(ctx, interval, job)
}

// Add adds a timing job to the default timer, which runs in interval of `interval`.
func Add(ctx context.Context, interval time.Duration, job JobFunc) Entry {
	return gtimer.Add(ctx, interval, job)
}

// AddEntry adds a timing job to the default timer with detailed parameters.
//
// The parameter `interval` specifies the running interval of the job.
//
// The parameter `singleton` specifies whether the job running in singleton mode.
// There's only one of the same job is allowed running when its a singleton mode job.
//
// The parameter `times` specifies limit for the job running times, which means the job
// exits if its run times exceeds the `times`.
//
// The parameter `status` specifies the job status when it's firstly added to the timer.
func AddEntry(ctx context.Context, interval time.Duration, job JobFunc, isSingleton bool, times int, status int) Entry {
	return gtimer.AddEntry(ctx, interval, job, isSingleton, times, status)
}

// AddSingleton is a convenience function for add singleton mode job.
func AddSingleton(ctx context.Context, interval time.Duration, job JobFunc) Entry {
	return gtimer.AddSingleton(ctx, interval, job)
}

// AddOnce is a convenience function for adding a job which only runs once and then exits.
func AddOnce(ctx context.Context, interval time.Duration, job JobFunc) Entry {
	return gtimer.AddOnce(ctx, interval, job)
}

// AddTimes is a convenience function for adding a job which is limited running times.
func AddTimes(ctx context.Context, interval time.Duration, times int, job JobFunc) Entry {
	return gtimer.AddTimes(ctx, interval, times, job)
}

// DelayAdd adds a timing job after delay of `interval` duration.
// Also see Add.
func DelayAdd(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	gtimer.DelayAdd(ctx, delay, interval, job)
}

// DelayAddEntry adds a timing job after delay of `interval` duration.
// Also see AddEntry.
func DelayAddEntry(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc, isSingleton bool, times int, status int) {
	gtimer.DelayAddEntry(ctx, delay, interval, job, isSingleton, times, status)
}

// DelayAddSingleton adds a timing job after delay of `interval` duration.
// Also see AddSingleton.
func DelayAddSingleton(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	gtimer.DelayAddSingleton(ctx, delay, interval, job)
}

// DelayAddOnce adds a timing job after delay of `interval` duration.
// Also see AddOnce.
func DelayAddOnce(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	gtimer.DelayAddOnce(ctx, delay, interval, job)
}

// DelayAddTimes adds a timing job after delay of `interval` duration.
// Also see AddTimes.
func DelayAddTimes(ctx context.Context, delay time.Duration, interval time.Duration, times int, job JobFunc) {
	gtimer.DelayAddTimes(ctx, delay, interval, times, job)
}
