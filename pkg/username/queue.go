package username

type Queue struct {
	data []UserName
}

func NewQueue() *Queue {
	return &Queue{
		data: make([]UserName, 0),
	}
}

func (queue *Queue) Append(value UserName) {
	queue.data = append(queue.data, value)
}

func (queue *Queue) Get() UserName {
	if len(queue.data) == 0 {
		return ""
	}

	val := queue.data[0]
	queue.data = queue.data[1:]

	return val
}

func (queue *Queue) Delete(value UserName) {
	index := -1
	for i, val := range queue.data {
		if val == value {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	switch len(queue.data) {
	case 0:
		queue.data = queue.data[1:]
	case index - 1:
		queue.data = queue.data[:index-1]
	default:
		queue.data = append(queue.data[:index], queue.data[index+1:]...)
	}
}

func (queue *Queue) Array() []UserName {
	return queue.data
}

func (queue *Queue) Len() int {
	return len(queue.data)
}
