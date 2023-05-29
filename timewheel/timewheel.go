package timewheel

/*
来源：https://github.com/ouqiang/timewheel
开源协议：MIT License

改动：
	1. key 的类型由 interface{} 改为 string
	2. 增加 updateTimer 接口
	3. 由 秒 级别改为 毫秒 级别
*/

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// Job 延时任务回调函数
type Job func(interface{})

// TaskData 回调函数参数类型

// TimeWheel 时间轮
type TimeWheel struct {
	interval time.Duration // 指针每隔多久往前移动一格
	ticker   *time.Ticker
	slots    []*list.List // 时间轮槽
	// key: 定时器唯一标识 value: 定时器所在的槽, 主要用于删除定时器, 不会出现并发读写，不加锁直接访问
	timer             map[string]int
	currentPos        int         // 当前指针指向哪一个槽
	slotNum           int         // 槽数量
	job               Job         // 定时器回调函数
	addTaskChannel    chan Task   // 新增任务channel
	removeTaskChannel chan string // 删除任务channel
	updateTaskChannel chan Task   // 更新任务channel
	stopChannel       chan bool   // 停止定时器channel

	tasks map[string]bool // task map
	mu    *sync.Mutex     // 锁
}

// Task 延时任务
type Task struct {
	delay  time.Duration // 延迟时间
	circle int           // 时间轮需要转动几圈
	key    string        // 定时器唯一标识, 用于删除定时器
	data   interface{}   // 回调函数参数
}

// New 创建时间轮
func New(interval time.Duration, slotNum int, job Job) (*TimeWheel, error) {
	if interval <= 0 || slotNum <= 0 || job == nil {
		return nil, errors.New("invalid interval or slotNum or job")
	}
	tw := &TimeWheel{
		interval:          interval,
		slots:             make([]*list.List, slotNum),
		timer:             make(map[string]int),
		currentPos:        0,
		job:               job,
		slotNum:           slotNum,
		addTaskChannel:    make(chan Task),
		removeTaskChannel: make(chan string),
		updateTaskChannel: make(chan Task),
		stopChannel:       make(chan bool),
		tasks:             make(map[string]bool),
		mu:                &sync.Mutex{},
	}

	tw.initSlots()

	return tw, nil
}

// 初始化槽，每个槽指向一个双向链表
func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

// Start 启动时间轮
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

// Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

// AddTimer 添加定时器 key为定时器唯一标识
func (tw *TimeWheel) AddTimer(delay time.Duration, key string, data interface{}) {
	if delay < 0 {
		return
	}
	tw.mu.Lock()
	tw.tasks[key] = true
	tw.mu.Unlock()
	tw.addTaskChannel <- Task{delay: delay, key: key, data: data}
}

// RemoveTimer 删除定时器 key为添加定时器时传递的定时器唯一标识
func (tw *TimeWheel) RemoveTimer(key string) {
	if key == "" {
		return
	}
	tw.mu.Lock()
	if _, ok := tw.tasks[key]; !ok {
		return
	}
	delete(tw.tasks, key)
	tw.mu.Unlock()
	tw.removeTaskChannel <- key
}

func (tw *TimeWheel) UpdateTimer(delay time.Duration, key string, data interface{}) {
	if delay < 0 {
		return
	}
	// 先删除, 防止回调触发
	tw.mu.Lock()
	delete(tw.tasks, key)
	tw.mu.Unlock()
	tw.updateTaskChannel <- Task{delay: delay, key: key, data: data}
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(&task)
		case key := <-tw.removeTaskChannel:
			tw.removeTask(key)
		case task := <-tw.updateTaskChannel:
			tw.updateTask(&task)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAndRunTask(l)
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

// 扫描链表中过期定时器, 并执行回调函数
func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		tw.mu.Lock()
		_, ok := tw.tasks[task.key]
		tw.mu.Unlock()
		// 已经 remove 的任务
		if !ok {
			next := e.Next()
			l.Remove(e)
			if task.key != "" {
				delete(tw.timer, task.key)
			}
			e = next
			continue
		}

		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}

		go tw.job(task.data)
		next := e.Next()
		l.Remove(e)
		if task.key != "" {
			delete(tw.timer, task.key)
		}
		e = next
	}
}

// 新增任务到链表中
func (tw *TimeWheel) addTask(task *Task) {
	pos, circle := tw.getPositionAndCircle(task.delay)
	task.circle = circle

	tw.slots[pos].PushBack(task)

	if task.key != "" {
		tw.timer[task.key] = pos
	}
}

// 更新 timer
func (tw *TimeWheel) updateTask(task *Task) {
	// 删除
	tw.removeTask(task.key)

	// 增加
	tw.addTask(task)
}

// 获取定时器在槽中的位置, 时间轮需要转动的圈数
func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos, circle int) {
	delaySeconds := int(d.Milliseconds())
	intervalSeconds := int(tw.interval.Milliseconds())
	circle = delaySeconds / intervalSeconds / tw.slotNum
	pos = (tw.currentPos + delaySeconds/intervalSeconds) % tw.slotNum

	return
}

// 从链表中删除任务
func (tw *TimeWheel) removeTask(key string) {
	// 获取定时器所在的槽
	position, ok := tw.timer[key]
	if !ok {
		return
	}
	// 获取槽指向的链表
	l := tw.slots[position]
	for e := l.Back(); e != nil; {
		task := e.Value.(*Task)
		if task.key == key {
			delete(tw.timer, task.key)
			l.Remove(e)
			break
		}

		e = e.Next()
	}
}
