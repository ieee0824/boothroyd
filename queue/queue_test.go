package queue

import (
	"testing"
	"time"
	"os"
	"reflect"
)

func TestNewInnerQueue(t *testing.T) {
	{
		q := newInnnerQueue()
		if q.delay != 5*time.Second {
			t.Error("delay time is miss match")
		}

	}
	{
		os.Setenv("DELAY_TIME", "10")
		q := newInnnerQueue()
		if q.delay != 10*time.Second {
			t.Error("delay time is miss match")
		}
		os.Unsetenv("DELAY_TIME")
	}
}

func TestInnerQueueEnqueue(t *testing.T) {
	{
		q := newInnnerQueue()
		tests := []struct{
			input interface{}
			want int
		}{
			{"5000兆円ほしい", 1},
			{1, 2},
			{100, 3},
			{500000000000, 4},
		}

		for _, test := range tests {
			q.enqueue(test.input)
			if len(q.status) != test.want {
				t.Fatalf("want %q, but %q:", test.want, len(q.status))
			}
		}
	}
}

func TestInnnerQueueDequeue(t *testing.T) {
	defer func(){
		os.Unsetenv("DELAY_TIME")
	}()
	{
		os.Setenv("DELAY_TIME", "0")
		q := newInnnerQueue()
		tests := []struct{
			input interface{}
		}{
			{"5000兆円ほしい"},
			{1},
			{100},
			{500000000000},
			{"にゃんぱすー"},
		}
		for _, test := range tests {
			q.enqueue(test.input)
		}
		for _, test := range tests {
			select {
			case d := <- q.dequeue():
				if !reflect.DeepEqual(d, test.input) {
					t.Fatalf("want %v, but %v:", test.input, d)
				}
			}
		}

	}
}