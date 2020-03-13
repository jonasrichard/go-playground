
## Worker pool

### First solution

Let us see how to implement a worker pool in a non-blocking way. We have several
worker goroutines which spend some time with computation. The problem is that
once the router which distributes the work, sends new work items to a worker
may be blocked if the worker is busy.

### Second solution

We are collecting the pending items to a queue and every time we get data we
try to resend the pending list. What we managed to send we remove from the
pending list.

There are two problems, the first is latency: if we don't get any input, we
don't try to resend the data. The second is: if only one worker is slow, we
try to resend the same work items to the same busy worker.

### Third solution

The router manages the status of the workers (available or busy). If we got
a message but the proper worker is busy, we append to the pending list.
If the worker is available we send it to that worker. The workers also can
send back a bool message to sign that they become available. When a worker
is available again, we look for a work item, which is supposed to be handled
by that worker, and remove that item from the pending list.
