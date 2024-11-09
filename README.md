# go-mapslice

Interconversion between map list and slice

## Input/Output Examples

- Functions using generics.
- Map keys and slice header elements must satisfy `cmp.Ordered`.
- Map values and slice data elements must satisfy `comparable`.


Map list (JSON-like)

```
[
    {"key1": 1, "key2": 2,},
    {"key1": 3, "key2": 4,},
    ...
]
```

Slice (CSV-like)

```
[
    ["key1, "key2"],
    [1, 2],
    [3, 4],
]
```

## Usage

### Map list to slice

```go
import "github.com/bugph0bia/go-mapslice"

func main() {

	maplist := []map[string]int{
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

	header, data := mapslice.MapToSlice(input, nil)
	// header = []string{"key1", "key2", "key3", "key4"}
	// data   = [][]int{{1, 2, 3, 0}, {4, 0, 5, 6}}

	header, data = mapslice.MapToSlice(input, []string{"key3"})
	// header = []string{"key3", "key1", "key2", "key4"}
	// data   = [][]int{{3, 1, 2, 0}, {5, 4, 0, 6}}
}
```

- Empty elements has a zero value of type V.
- Order of the columns is such that the `fixColumns` (second argument) are placed first if they exist, followed by the other keys sorted in ascending order. `fixColumns` accepts `nil`.

### Slice to Map list

```go
import "github.com/bugph0bia/go-mapslice"

func main() {

	header := []string{"key1", "key2", "key3", "key4"}
	data := [][]int{{1, 2, 3, 0}, {4, 0, 5, 6}}

	maplist := SliceToMap(header, data, false)
	// maplist = []map[string]int{
	// 	{
	// 		"key1": 1,
	// 		"key2": 2,
	// 		"key3": 3,
	// 		"key4": 0,
	// 	},
	// 	{
	// 		"key1": 4,
	// 		"key2": 0,
	// 		"key3": 5,
	// 		"key4": 6,
	// 	},
	// }

	maplist := SliceToMap(header, data, true)
	// maplist = []map[string]int{
	// 	{
	// 		"key1": 1,
	// 		"key2": 2,
	// 		"key3": 3,
	// 	},
	// 	{
	// 		"key1": 4,
	// 		"key3": 5,
	// 		"key4": 6,
	// 	},
	// }

}
```

- If the size of the data row is larger than the size of the header, discard the value.
- If `ignoreZero` (second argument) is true, do not store zero values.
