package timequeue

import "testing"


func TestQueueEvents(t *testing.T) {

	q := Queue{}

	if e := q.Events(); len(e) != 0 {
		t.Fail()
		t.Logf("expected zero length events from fresh queue.")
	}

	q.Push(Event{Time: 3,ID: 1, Data: nil})
	q.Push(Event{Time: 2, ID: 2, Data: nil})
	q.Push(Event{Time: 1, ID: 3, Data: nil})

	e := q.Events()
	if len(e) != 3 {
		t.Fail()
		t.Logf("expected 3 length events from used queue.")
	}

	if e[0].ID != 3 {
		t.Fail()
		t.Logf("expected ID==3 at 0 index.")
	}

	if e[2].ID != 1 {
		t.Fail()
		t.Logf("expected ID==1 at 1 index.")
	}
}

func TestQueuePushPop(t *testing.T) {
	q := Queue{}

	if _, ok := q.Pop(); ok {
		t.Fail()
		t.Log("ok==true on fresh queue")
	}

	q.Push(Event{Time: 3,ID: 1})

	if _, ok := q.Pop(); !ok {
		t.Fail()
		t.Log("ok==false after a single push")
	}

	q.Push(Event{Time: 4,ID: 2})
	q.Push(Event{Time: 2,ID: 3})
	q.Push(Event{Time: 5,ID: 5})
	e, _ := q.Pop()
	if e.ID != 3 {
		t.Fail()
		t.Log("expected id==3 event popped")
	}

	if q.events[q.start].Time != 2 {
		t.Fail()
		t.Log("time was not decreased")
	}

	e, _ = q.Pop()
	if e.ID != 2 {
		t.Fail()
		t.Log("expected id==2 event popped")
	}

	q.Pop()
}

func TestQueueExpand(t *testing.T) {
	q := Queue{}

	for i := 0; i < initsize+2; i++ {
		q.Push(Event{Time: 0, ID: i})
	}

	if len(q.events) == initsize {
		t.Fail()
		t.Log("expand should have increased buffer size")
	}
}

func TestQueueRingBuffer(t *testing.T) {
	q := Queue{}

	for i := 0; i < initsize; i++ {
		q.Push(Event{Time: float32(i), ID: i+1})
	}

	for i := 0; i < initsize-1; i++ {
		q.Pop()
	}

	for i := 0; i < initsize-1; i++ {
		q.Push(Event{Time: float32(i), ID: i+1})
	}

	if len(q.events) != initsize {
		t.Fail()
		t.Log("buffer should have reused old indices instead of expanding")
	}

	q.Push(Event{Time: 1, ID: 1})

	if len(q.events) == initsize {
		t.Fail()
		t.Log("buffer should have expanded at this point")
	}
}

func TestQueueCancel(t *testing.T) {
	q := Queue{}

	if _, ok := q.Cancel(1); ok {
		t.Fail()
		t.Log("ok==true from cancel on fresh queue")
	}

	q.Push(Event{Time: 4,ID: 2})
	q.Push(Event{Time: 2,ID: 3})
	q.Push(Event{Time: 5,ID: 5})
	q.Push(Event{Time: 6,ID: 6})

	if _, ok := q.Cancel(2); !ok {
		t.Fail()
		t.Log("ok==false on valid cancel")
	}

	if _, ok := q.Cancel(2); ok {
		t.Fail()
		t.Log("ok==true on invalid cancel")
	}

	q.Pop()

	if e, _ := q.Pop(); e.ID == 2 {
		t.Fail()
		t.Log("cancelled event was popped, should've been skipped")
	}

	q.Cancel(6);

	if _, ok := q.Pop(); ok {
		t.Fail()
		t.Log("ok==true but last event was cancelled")
	}
}