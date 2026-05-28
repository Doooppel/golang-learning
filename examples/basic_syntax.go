package examples

// 这个文件展示了Go语言的基础语法
// 虽然这个文件不会被编译到主程序中，但可以作为学习参考

import (
	"fmt"
	"time"
)

// 基础语法示例函数
func BasicSyntaxExamples() {
	// ========== 1. 变量声明 ==========
	// 方式1: var 关键字
	var name string = "Go语言"

	// 方式2: 类型推断
	var age = 25

	// 方式3: 短变量声明（最常用）
	email := "example@example.com"

	// 方式4: 多变量声明
	var x, y int = 1, 2
	a, b := 3, 4

	fmt.Println(name, age, email, x, y, a, b)

	// ========== 2. 常量 ==========
	const pi = 3.14159
	const (
		StatusOK       = 200
		StatusNotFound = 404
	)

	// ========== 3. 基本数据类型 ==========
	var (
		// 整数类型
		intVal   int   = 42
		int8Val  int8  = 127
		int64Val int64 = 9223372036854775807

		// 浮点数
		float32Val float32 = 3.14
		float64Val float64 = 3.141592653589793

		// 布尔值
		isTrue bool = true

		// 字符串
		str string = "Hello, 世界"

		// 字节（uint8的别名）
		byteVal byte = 'A'

		// rune（int32的别名，用于Unicode字符）
		runeVal rune = '中'
	)

	fmt.Println(intVal, int8Val, int64Val, float32Val, float64Val,
		isTrue, str, byteVal, runeVal)

	// ========== 4. 数组和切片 ==========
	// 数组（固定长度）
	var arr [5]int = [5]int{1, 2, 3, 4, 5}

	// 切片（动态长度，更常用）
	slice := []int{1, 2, 3, 4, 5}
	slice = append(slice, 6) // 追加元素

	// 使用make创建切片
	slice2 := make([]int, 0, 10) // 长度0，容量10

	fmt.Println(arr, slice, slice2)

	// ========== 5. Map（字典） ==========
	// 创建map
	userMap := make(map[string]int)
	userMap["age"] = 25
	userMap["score"] = 100

	// 字面量创建
	userMap2 := map[string]string{
		"name":  "张三",
		"email": "zhangsan@example.com",
	}

	// 访问和检查
	if age, exists := userMap["age"]; exists {
		fmt.Println("年龄:", age)
	}

	fmt.Println(userMap, userMap2)

	// ========== 6. 结构体 ==========
	type Person struct {
		Name  string
		Age   int
		Email string
	}

	// 创建结构体实例
	person1 := Person{
		Name:  "李四",
		Age:   30,
		Email: "lisi@example.com",
	}

	// 指针
	person2 := &Person{
		Name: "王五",
		Age:  28,
	}

	fmt.Println(person1, person2)

	// ========== 7. 控制流 ==========
	// if-else
	if age > 18 {
		fmt.Println("成年人")
	} else {
		fmt.Println("未成年人")
	}

	// switch
	switch age {
	case 18:
		fmt.Println("刚成年")
	case 25, 26, 27:
		fmt.Println("二十多岁")
	default:
		fmt.Println("其他年龄")
	}

	// for循环
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// range遍历
	numbers := []int{1, 2, 3, 4, 5}
	for index, value := range numbers {
		fmt.Printf("索引: %d, 值: %d\n", index, value)
	}

	// 遍历map
	for key, value := range userMap {
		fmt.Printf("键: %s, 值: %d\n", key, value)
	}

	// ========== 8. 函数 ==========
	result := add(10, 20)
	fmt.Println("加法结果:", result)

	sum, product := calculate(5, 6)
	fmt.Println("和:", sum, "积:", product)

	// ========== 9. 错误处理 ==========
	if err := mayFail(); err != nil {
		fmt.Println("错误:", err)
	}

	// ========== 10. 接口 ==========
	var writer Writer = &ConsoleWriter{}
	writer.Write([]byte("Hello, Interface!"))

	// ========== 11. Goroutine（并发） ==========
	go func() {
		fmt.Println("这是在一个goroutine中执行")
	}()

	time.Sleep(100 * time.Millisecond)

	// ========== 12. Channel（通道） ==========
	ch := make(chan string)

	go func() {
		ch <- "Hello from channel"
	}()

	message := <-ch
	fmt.Println(message)

	// ========== 13. defer（延迟执行） ==========
	deferExample()
}

// 普通函数
func add(a, b int) int {
	return a + b
}

// 多返回值函数
func calculate(a, b int) (int, int) {
	return a + b, a * b
}

// 命名返回值
func divide(a, b float64) (result float64, err error) {
	if b == 0 {
		err = fmt.Errorf("除数不能为0")
		return
	}
	result = a / b
	return
}

// 错误处理示例
func mayFail() error {
	return fmt.Errorf("这是一个错误示例")
}

// 接口定义
type Writer interface {
	Write([]byte) error
}

// 接口实现
type ConsoleWriter struct{}

func (cw *ConsoleWriter) Write(data []byte) error {
	fmt.Println(string(data))
	return nil
}

// defer示例
func deferExample() {
	fmt.Println("开始")
	defer fmt.Println("这会最后执行（defer）")
	defer fmt.Println("这会倒数第二个执行")
	fmt.Println("结束")
	// 输出顺序：开始 -> 结束 -> 这会倒数第二个执行 -> 这会最后执行（defer）
}
