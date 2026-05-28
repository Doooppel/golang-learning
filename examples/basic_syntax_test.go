package examples

import (
	"fmt"
	"io"
	"testing"
)

func TestBasicSyntax(t *testing.T) {
	//array test start
	testArr := make([]int, 10)
	testArr[0] = 1

	testArr = append(testArr, 2, 3, 4)
	fmt.Println(testArr)
	fmt.Println(len(testArr))

	testArr2 := make([]int, 0, 10)
	testArr2 = append(testArr2, 2)
	fmt.Println(testArr2)

	var list = []int{1, 2, 2, 3, 5, 8, 13}
	fmt.Println(list)

	for index, value := range list {
		fmt.Printf("index: %d, value: %d\n", index, value)
	}
	//array test end

	//map test start
	testMap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	fmt.Println(testMap)

	myMap := map[string]int{
		"a": 1,
	}
	fmt.Println(myMap)

	testMap2 := make(map[string]int)
	testMap2["key1"] = 9999
	fmt.Println(testMap2)

	for key, value := range testMap {
		key := key + "_suffix"
		value := value * 10
		fmt.Printf("key: %s, value: %d\n", key, value)
	}

	for i, j := 1, 10; i <= 10; i, j = i+1, j-1 {
		fmt.Println(i, j)
	}

	var personalMap = map[string]string{
		"name": "Doppel",
		"age":  "28",
	}

	if value, exists := personalMap["age"]; exists {
		fmt.Println("age", value, exists)
	}

	// BasicSyntaxExamples()
}

func TestPolymorphsim(t *testing.T) {
	myReader := FakeReader{data: "test polymorphic reader"}
	n, err := ProcessData(myReader)
	fmt.Printf("Read %d bytes, error: %v\n", n, err)
}

type FakeReader struct {
	data string
}

func (fr FakeReader) Read(p []byte) (n int, err error) {
	fmt.Println("Reading data")
	fmt.Println("Data:", fr.data)
	return len(fr.data), io.EOF
}
func ProcessData(r io.Reader) (n int, err error) {
	buf := make([]byte, 1024)
	return r.Read(buf)
}
