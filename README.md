# Time Queue
Time queue system for turn based games. Events are sorted by time as they enter the queue, and their time is decreased as other events are removed from the queue.

Based on the rogue basin article [A priority queue based turn scheduling system](http://www.roguebasin.com/index.php?title=A_priority_queue_based_turn_scheduling_system)

## Install

go get github.com/redarcanearcher/timequeue

## Sample

```
q := timequeue.Queue{}

q.Push(timequeue.Event{ ID: 1, Time: 3.3, Data: "attack"})
q.Push(timequeue.Event{ ID: 2, Time: 1, Data: "defend"})
q.Push(timequeue.Event{ ID: 3, Time: 4.5, Data: "rest"})

for {
    event, ok := q.Pop()
    if !ok {
        fmt.Println("end of queue")
        return
    }
    fmt.Printf("ID: %v, Delta: %v, Data: %v\n", event.ID, event.Time, event.Data)
}
```

```
Output:
ID: 2, Delta: 1, Data: defend
ID: 1, Delta: 2.3, Data: attack
ID: 3, Delta: 1.2, Data: rest
end of queue
```