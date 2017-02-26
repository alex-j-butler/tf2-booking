package retry

import (
	"math"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

var MaxRetries = 10
var errMaxRetriesReached = errors.New("exceeded retry limit")

type Func func(attempt int) (retry bool, err error)

func Do(fn Func) error {
	return DoWithOptions(fn, &Options{DefaultBackoff, MaxRetries})
}

func DoN(fn Func, maxRetries int) error {
	return DoWithOptions(fn, &Options{DefaultBackoff, maxRetries})
}

func DoWithOptions(fn Func, options *Options) error {
	var err error
	var cont bool

	if options.BackoffStrategy == nil {
		options.BackoffStrategy = DefaultBackoff
	}

	if options.MaxRetries == 0 {
		options.MaxRetries = MaxRetries
	}

	attempt := 1
	for {
		cont, err = fn(attempt)
		if !cont || err == nil {
			break
		}
		attempt++
		if attempt > options.MaxRetries {
			return errors.Wrap(err, errMaxRetriesReached.Error())
		}

		time.Sleep(options.BackoffStrategy(attempt))
	}

	return err
}

type BackoffStrategy func(retry int) time.Duration

type Options struct {
	BackoffStrategy BackoffStrategy
	MaxRetries      int
}

func IsMaxRetries(err error) bool {
	return err == errMaxRetriesReached
}

func DefaultBackoff(_ int) time.Duration {
	return 0 * time.Second
}

func ExponentialJitterBackoff(i int) time.Duration {
	return jitter(int(math.Pow(2, float64(i))))
}

func jitter(i int) time.Duration {
	ms := i * 1000

	maxJitter := ms / 3

	rand.Seed(time.Now().UnixNano())
	jitter := rand.Intn(maxJitter + 1)

	if rand.Intn(2) == 1 {
		ms = ms + jitter
	} else {
		ms = ms - jitter
	}

	if ms <= 0 {
		ms = 1
	}

	return time.Duration(ms) * time.Millisecond
}
