package mapslice

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaplistToTable(t *testing.T) {
	input := []map[string]int{
		{
			"key1": 1,
			"key2": 2,
			"key3": 3,
		},
		{
			"key1": 4,
			"key3": 5,
			"key4": 6,
		},
	}

	// 固定列を指定しない場合
	wanth := []string{"key1", "key2", "key3", "key4"}
	wantd := [][]int{{1, 2, 3, 0}, {4, 0, 5, 6}}
	acth, actd := MaplistToTable(input, nil)
	assert.Equal(t, wanth, acth)
	assert.Equal(t, wantd, actd)

	// 固定列を指定する場合
	wanth = []string{"key3", "key1", "key2", "key4"}
	wantd = [][]int{{3, 1, 2, 0}, {5, 4, 0, 6}}
	acth, actd = MaplistToTable(input, []string{"key3"})
	assert.Equal(t, wanth, acth)
	assert.Equal(t, wantd, actd)

	// 入力が空の場合
	input = []map[string]int{}
	wanth = nil
	wantd = nil
	acth, actd = MaplistToTable(input, nil)
	assert.Equal(t, wanth, acth)
	assert.Equal(t, wantd, actd)

	// 入力が空の場合
	input = []map[string]int{make(map[string]int)}
	wanth = nil
	wantd = nil
	acth, actd = MaplistToTable(input, nil)
	assert.Equal(t, wanth, acth)
	assert.Equal(t, wantd, actd)
}

func TestTableToMaplist(t *testing.T) {

	inputh := []string{"key1", "key2", "key3", "key4"}
	inputd := [][]int{{1, 2, 3, 0}, {4, 0, 5, 6}}

	// ゼロ値を格納する場合
	want := []map[string]int{
		{
			"key1": 1,
			"key2": 2,
			"key3": 3,
			"key4": 0,
		},
		{
			"key1": 4,
			"key2": 0,
			"key3": 5,
			"key4": 6,
		},
	}
	act := TableToMaplist(inputh, inputd, false)
	assert.Equal(t, want, act)

	// ゼロ値を格納しない場合
	want = []map[string]int{
		{
			"key1": 1,
			"key2": 2,
			"key3": 3,
		},
		{
			"key1": 4,
			"key3": 5,
			"key4": 6,
		},
	}
	act = TableToMaplist(inputh, inputd, true)
	assert.Equal(t, want, act)

	// ヘッダよりも長いデータを保つ場合
	inputd = [][]int{{1, 2, 3, 4, 5, 6}}
	want = []map[string]int{
		{
			"key1": 1,
			"key2": 2,
			"key3": 3,
			"key4": 4,
		},
	}
	act = TableToMaplist(inputh, inputd, false)
	assert.Equal(t, want, act)

	// ゼロ値により出力が空になる場合
	inputd = [][]int{{0, 0, 0, 0}}
	want = nil
	act = TableToMaplist(inputh, inputd, true)
	assert.Equal(t, want, act)

	// 入力が空の場合
	inputh = []string{}
	inputd = [][]int{}
	want = nil
	act = TableToMaplist(inputh, inputd, false)
	assert.Equal(t, want, act)
}

func TestReadJson(t *testing.T) {
	json := `
		[
			{
				"key1": 1,
				"key2": 2,
				"key3": 3
			},
			{
				"key1": 4,
				"key3": 5,
				"key4": 6
			}
		]
	`
	r := strings.NewReader(json)
	want := []map[string]int{
		{
			"key1": 1,
			"key2": 2,
			"key3": 3,
		},
		{
			"key1": 4,
			"key3": 5,
			"key4": 6,
		},
	}
	act, err := ReadJson[string, int](r)
	assert.Equal(t, want, act)
	assert.NoError(t, err)
}

func TestWriteJson(t *testing.T) {
	w := new(bytes.Buffer)
	input := []map[string]int{
		{
			"key1": 1,
			"key2": 2,
			"key3": 3,
		},
		{
			"key1": 4,
			"key3": 5,
			"key4": 6,
		},
	}
	want := `[{"key1":1,"key2":2,"key3":3},{"key1":4,"key3":5,"key4":6}]`
	err := WriteJson(w, input)
	act := w.String()
	assert.Equal(t, want, act)
	assert.NoError(t, err)
}
