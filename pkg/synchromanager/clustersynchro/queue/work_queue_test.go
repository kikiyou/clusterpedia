package queue

import (
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/clock"
)

func newTestBasicWorkQueue() (*basicWorkQueue, *clock.FakeClock) {
	fakeClock := clock.NewFakeClock(time.Now())
	wq := &basicWorkQueue{
		clock: fakeClock,
		queue: make(map[string]time.Time),
	}
	return wq, fakeClock
}

func compareResults(t *testing.T, expected, actual []string) {
	expectedSet := sets.NewString()
	for _, u := range expected {
		expectedSet.Insert(u)
	}
	actualSet := sets.NewString()
	for _, u := range actual {
		actualSet.Insert(u)
	}
	if !expectedSet.Equal(actualSet) {
		t.Errorf("Expected %#v, got %#v", expectedSet.List(), actualSet.List())
	}
}

func TestGetWork(t *testing.T) {
	q, clock := newTestBasicWorkQueue()
	q.Enqueue("foo1", -1*time.Minute)
	q.Enqueue("foo2", -1*time.Minute)
	q.Enqueue("foo3", 1*time.Minute)
	q.Enqueue("foo4", 1*time.Minute)
	expected := []string{"foo1", "foo2"}
	compareResults(t, expected, q.GetWork())
	compareResults(t, []string{}, q.GetWork())
	// Dial the time to 1 hour ahead.
	clock.Step(time.Hour)
	expected = []string{"foo3", "foo4"}
	compareResults(t, expected, q.GetWork())
	compareResults(t, []string{}, q.GetWork())
}
