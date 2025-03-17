package sync

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test that the retries are actually causing the expected
// delay and  additionally make sure they are a fixed delay
// and not the default backoff value.
func TestWithRetry_DelayBetweenRetries(t *testing.T) {
	t.Parallel()

	counter := 0
	start := time.Now()
	err := withRetry(func() error {
		counter++
		if counter < 3 {
			return errors.New("test error")
		}
		return nil
	}, 3, 1) // 3 attempts, 1-second delay

	elapsed := time.Since(start)

	assert.NoError(t, err, "Expected success before max attempts")
	assert.GreaterOrEqual(t, elapsed.Seconds(), 2.0, "Expected at least 2 seconds of delay between all retries")
	assert.LessOrEqual(t, elapsed.Seconds(), 2.5, "Expected at most 2.5 seconds of delay between all retries")
}

// Test that we do not retry on immediate success.
func TestWithRetry_NoRetriesOnImmediateSuccess(t *testing.T) {
	t.Parallel()

	counter := 0
	err := withRetry(func() error {
		counter++
		return nil
	}, 5, 2) // 5 attempts, 2-second delay

	assert.NoError(t, err, "Expected no error when function succeeds immediately")
	assert.Equal(t, 1, counter, "Expected function to run only once without retries")
}

// Test that we properly succeed after a few but not max retires.
func TestWithRetry_SuccessAfterRetries(t *testing.T) {
	t.Parallel()

	counter := 0
	err := withRetry(func() error {
		counter++
		if counter < 2 {
			return errors.New("test error")
		}
		return nil
	}, 3, 1) // 3 attempts, 1-second delay

	assert.NoError(t, err, "Expected success before max attempts")
	assert.Equal(t, 2, counter, "Expected function to retry once before success")
}

// Test to make sure we properly fail after max amount of retries.
func TestWithRetry_MaxAttemptsFailure(t *testing.T) {
	t.Parallel()

	counter := 0
	err := withRetry(func() error {
		counter++
		return errors.New("test error")
	}, 3, 1) // 3 attempts, 1-second delay

	assert.Error(t, err, "Expected an error after max attempts")
	assert.Equal(t, 3, counter, "Expected function to be retried 3 times")
}
