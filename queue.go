package timequeue

const initsize = 16 //must be power of 2

//Queue represents a queue of events sorted by time.
type Queue struct {
	events []Event
	start, end, length int
}

//Event represents something that is ready to be executed in the future.
//ID of zero is considered cancelled.
type Event struct {
	Time float32
	ID int
	Data interface{}
}

func (q Queue) index(i int) int {
	return i & (len(q.events)-1)
}

//Push schedules the event to be returned later with Pop().
//event.ID must not be zero, else it will be considered cancelled.
func (q *Queue) Push(event Event) {
	if q.length == len(q.events) {
		q.expand()
	}

	i := q.end
	for i != q.start {
		j := q.index(i - 1)
		if event.Time < q.events[j].Time {
			q.events[i] = q.events[j]
			i = j
		} else {
			break
		}
	}
	q.events[i] = event

	q.length++
	q.end = q.index(q.end + 1)
} 

//Pop advances time and returns the next event.
//If there are no remaining events, false is returned.
func (q *Queue) Pop() (Event, bool) {
	if q.length == 0 {
		return Event{}, false
	}

	event := Event{}
	for event.ID == 0 && q.length > 0 {
		event = q.events[q.start]
		q.events[q.start] = Event{}
		q.start = q.index(q.start+1)
		q.length--
	}

	if event.ID == 0 {
		return Event{}, false
	}

	for i := q.start; i != q.end; i = q.index(i+1) {
		q.events[i].Time -= event.Time 
	}

	return event, true
}

//Cancel removes the event from the queue.
//If the event is found, it returns the remaining time and true.
func (q *Queue) Cancel(id int) (float32, bool) {
	for i := range q.events {
		if q.events[i].ID == id {
			q.events[i].ID = 0
			return q.events[i].Time, true
		}
	}
	return 0, false
}

//Events returns a copy of all the events in the queue.
func (q Queue) Events() []Event {
	e := make([]Event, q.length)
	n := q.copy(e)
	return e[:n]
}

//Clear resets the queue to zero events.
func (q *Queue) Clear() {
	for i := range q.events {
		q.events[i] = Event{}
	}
	q.start = 0
	q.end = 0
	q.length = 0
}

func (q *Queue) expand() {
	if len(q.events) == 0 {
		q.events = make([]Event, initsize)
		return
	}

	e := make([]Event, len(q.events) << 1)
	q.length = q.copy(e)
	q.start = 0
	q.end = q.length
	
	q.events = e
}

func (q Queue) copy(e []Event) int {
	if q.length == 0 {
		return 0
	}

	l := q.length
	j := q.start
	for i := 0; i < l; i++ {
		if q.events[j].ID != 0 {
			e[i] = q.events[j]
		} else {
			i--
			l--
		}
		j = q.index(j+1)
	}

	return l
}