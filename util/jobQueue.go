package main

import "fmt"
import "time"

func main() {
	const length = 1000
	queue := newQueue(length)
	myChannel := make(chan string)

	for i := 0; i < length; i++ {
		queue.enqueue(&Email{recepient: "recepient", sender: "sender", title: "title", body: "body"})
	}

	fmt.Println("Emails loaded")

	for i := 0; i < 6; i++ {
		email := queue.dequeue()
		go email.sendInBackground(myChannel)
	}

	time.Sleep(time.Second * 60)
	fmt.Println("Emails sent")
}

type Queue struct {
	front      int
	rear       int
	size       int
	emailQueue []Email
}

type Email struct {
	recepient string
	sender    string
	title     string
	body      string
}

func newQueue(size int) Queue {
	q := Queue{front: -1, rear: -1, emailQueue: make([]Email, size), size: size}
	return q
}

func (q *Queue) enqueue(e *Email) {
	if ((q.rear + 1) % q.size) == q.front {
		panic("Full")
	} else if q.isEmpty() {
		q.front = 0
		q.rear = 0
	} else {
		q.rear = (q.rear + 1) % q.size
	}
	q.emailQueue[q.rear] = *e
}

func (q *Queue) dequeue() Email {
	if q.isEmpty() {
		//throw an error
		panic("Empty List")
	} else if q.front == q.rear {
		e := q.emailQueue[q.front]
		q.front = -1
		q.rear = -1
		return e
	} else {
		e := q.emailQueue[q.front]
		q.front = (q.front + 1) % q.size
		return e
	}
}

func (q *Queue) isEmpty() bool {
	return q.front == -1 && q.rear == -1
}

func (q *Queue) isFull() bool {
	return false
}

func (e *Email) sendInBackground(stringChan chan string) {
	fmt.Println("email send")
	time.Sleep(time.Millisecond * 10)
	fmt.Println("email sent")
}
