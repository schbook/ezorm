package ezorm

import (
	"reflect"
	"testing"
	"time"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/5/25
 *
 */

func TestSnakeString(t *testing.T) {
	output := snakeString("HelloEzorm")
	if output != "hello_ezorm" {
		t.Errorf("SnakeString Wrong Output With:%v", output)
	}
}

func TestCamelString(t *testing.T) {
	output := camelString("hello_ezorm")
	if output != "HelloEzorm" {
		t.Errorf("CamelString Wrong Output With:%v", output)
	}
}

func TestSetFieldValue(t *testing.T) {
	type tp = struct {
		Integer  int
		Str      string
		Bl       bool
		Slice    []int
		Uinteger uint
		Flt      float32
		Ts       time.Time
	}

	var output tp
	input := tp{
		Integer:  1,
		Str:      "hello",
		Bl:       true,
		Slice:    []int{1, 2, 3},
		Uinteger: 2,
		Flt:      3.14,
		Ts:       time.Now(),
	}

	inputVal := reflect.ValueOf(&input).Elem()
	outputVal := reflect.ValueOf(&output).Elem()

	for i := 0; i < inputVal.NumField(); i++ {
		setFieldValue(outputVal.Field(i), inputVal.Field(i).Interface())

		if reflect.DeepEqual(inputVal.Field(i).Interface(),outputVal.Field(i).Interface())==false {
			t.Errorf("setFieldValue error, inputVal:%v, outputVal:%v", inputVal.Field(i).Interface(), outputVal.Field(i).Interface())
		}
	}
}
