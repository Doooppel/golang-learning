package examples

// 这个文件展示了Go语言的高级特性
// 包括：接口、泛型、反射、并发等

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// ========== 1. 接口和类型断言 ==========

// Animal 接口定义
type Animal interface {
	Speak() string
	Move() string
}

// Dog 实现Animal接口
type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return "汪汪"
}

func (d Dog) Move() string {
	return "跑步"
}

// Cat 实现Animal接口
type Cat struct {
	Name string
}

func (c Cat) Speak() string {
	return "喵喵"
}

func (c Cat) Move() string {
	return "跳跃"
}

// 接口使用示例
func InterfaceExample() {
	var animal Animal
	
	animal = Dog{Name: "旺财"}
	fmt.Println(animal.Speak(), animal.Move())
	
	animal = Cat{Name: "小花"}
	fmt.Println(animal.Speak(), animal.Move())
	
	// 类型断言
	if dog, ok := animal.(Dog); ok {
		fmt.Println("这是一只狗:", dog.Name)
	}
	
	// 类型switch
	switch v := animal.(type) {
	case Dog:
		fmt.Println("是狗:", v.Name)
	case Cat:
		fmt.Println("是猫:", v.Name)
	default:
		fmt.Println("未知动物")
	}
}

// ========== 2. 泛型（Go 1.18+） ==========

// Stack 泛型栈实现
type Stack[T any] struct {
	items []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// 泛型函数示例
func GenericExample() {
	// 整数栈
	intStack := NewStack[int]()
	intStack.Push(1)
	intStack.Push(2)
	intStack.Push(3)
	
	if item, ok := intStack.Pop(); ok {
		fmt.Println("弹出:", item)
	}
	
	// 字符串栈
	stringStack := NewStack[string]()
	stringStack.Push("Hello")
	stringStack.Push("World")
}

// ========== 3. 反射（Reflection） ==========

func ReflectionExample() {
	type User struct {
		Name  string `json:"name" db:"user_name"`
		Age   int    `json:"age" db:"user_age"`
		Email string `json:"email" db:"user_email"`
	}
	
	user := User{
		Name:  "张三",
		Age:   25,
		Email: "zhangsan@example.com",
	}
	
	// 获取类型信息
	t := reflect.TypeOf(user)
	fmt.Println("类型名称:", t.Name())
	
	// 遍历字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("字段: %s, 类型: %s\n", field.Name, field.Type)
		
		// 获取tag
		jsonTag := field.Tag.Get("json")
		dbTag := field.Tag.Get("db")
		fmt.Printf("  JSON tag: %s, DB tag: %s\n", jsonTag, dbTag)
	}
	
	// 获取值信息
	v := reflect.ValueOf(user)
	fmt.Println("值:", v)
}

// ========== 4. 并发：Goroutine和Channel ==========

func ConcurrencyExample() {
	// 无缓冲channel
	ch := make(chan int)
	
	// 启动goroutine发送数据
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch) // 关闭channel
	}()
	
	// 接收数据
	for value := range ch {
		fmt.Println("接收:", value)
	}
	
	// 有缓冲channel
	bufferedCh := make(chan string, 3)
	bufferedCh <- "消息1"
	bufferedCh <- "消息2"
	bufferedCh <- "消息3"
	
	fmt.Println("缓冲channel长度:", len(bufferedCh))
}

// ========== 5. Select语句（多路复用） ==========

func SelectExample() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- "来自ch1"
	}()
	
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- "来自ch2"
	}()
	
	// select会等待任意一个case就绪
	select {
	case msg1 := <-ch1:
		fmt.Println("收到:", msg1)
	case msg2 := <-ch2:
		fmt.Println("收到:", msg2)
	case <-time.After(3 * time.Second):
		fmt.Println("超时")
	}
}

// ========== 6. Context（上下文） ==========

func ContextExample() {
	// 创建带超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	// 在goroutine中使用context
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context已取消:", ctx.Err())
				return
			default:
				fmt.Println("工作中...")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()
	
	time.Sleep(3 * time.Second)
}

// ========== 7. sync包：互斥锁和等待组 ==========

func SyncExample() {
	var mu sync.Mutex
	var wg sync.WaitGroup
	counter := 0
	
	// 启动多个goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock() // 加锁
			counter++
			mu.Unlock() // 解锁
		}()
	}
	
	wg.Wait() // 等待所有goroutine完成
	fmt.Println("计数器值:", counter)
}

// ========== 8. sync.RWMutex（读写锁） ==========

type SafeMap struct {
	mu    sync.RWMutex
	items map[string]int
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		items: make(map[string]int),
	}
}

func (sm *SafeMap) Get(key string) (int, bool) {
	sm.mu.RLock() // 读锁
	defer sm.mu.RUnlock()
	value, exists := sm.items[key]
	return value, exists
}

func (sm *SafeMap) Set(key string, value int) {
	sm.mu.Lock() // 写锁
	defer sm.mu.Unlock()
	sm.items[key] = value
}

// ========== 9. 错误处理和自定义错误 ==========

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("错误代码: %d, 消息: %s", e.Code, e.Message)
}

func CustomErrorExample() error {
	return &CustomError{
		Code:    404,
		Message: "资源未找到",
	}
}

// ========== 10. 函数作为一等公民 ==========

func FunctionAsValueExample() {
	// 函数类型
	type Operation func(int, int) int
	
	add := func(a, b int) int { return a + b }
	multiply := func(a, b int) int { return a * b }
	
	// 函数作为参数
	calculate := func(op Operation, a, b int) int {
		return op(a, b)
	}
	
	fmt.Println("加法:", calculate(add, 5, 3))
	fmt.Println("乘法:", calculate(multiply, 5, 3))
	
	// 返回函数
	getOperation := func(op string) Operation {
		switch op {
		case "add":
			return add
		case "multiply":
			return multiply
		default:
			return add
		}
	}
	
	op := getOperation("multiply")
	fmt.Println("结果:", op(4, 5))
}

