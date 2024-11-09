// Interconversion between map list and slice
package mapslice

import (
	"cmp"
	"encoding/json"
	"io"
	"maps"
	"slices"
)

// MaplistToTable : Convert map list to table.
//
// e.g.:
// [{"key1": 1, "key2": 2}, {"key1": 3, "key2": 4}]  =>  ["key1", "key2"] and [[1, 2], [3, 4]]
//
// Empty elements has a zero value of type V.
//
// Order of the columns is such that the `fixedColumns` are placed first if they exist,
// followed by the other keys sorted in ascending order. `fixColumns` accepts nil.
func MaplistToTable[K cmp.Ordered, V comparable](maplist []map[K]V, fixedColumns []K) (tblheader []K, tbldata [][]V) {
	// 入力チェック
	if len(maplist) == 0 {
		return
	}

	// 渡されたマップリストからキーを抽出
	// 重複を除去するための集合としてマップのキーを利用
	set := map[K]struct{}{}
	for _, mrow := range maplist {
		for c := range maps.Keys(mrow) {
			set[c] = struct{}{}
		}
	}

	// スライスのヘッダ部を生成
	if fixedColumns != nil {
		// 固定列が指定されている場合は一旦除去
		for _, fc := range fixedColumns {
			delete(set, fc)
		}
		// 先頭に固定列を挿入、その後のキーはソート
		tblheader = slices.Concat(fixedColumns, slices.Sorted(maps.Keys(set)))
	} else {
		// キーをソート
		tblheader = slices.Sorted(maps.Keys(set))
	}
	if len(tblheader) == 0 {
		return
	}

	// スライスのデータ部を生成
	tbldata = make([][]V, 0, len(maplist))
	for _, mrow := range maplist {
		srow := make([]V, len(tblheader))
		for key, val := range mrow {
			pos := find(tblheader, key)
			if pos != -1 {
				srow[pos] = val
			}
		}
		tbldata = append(tbldata, srow)
	}

	return
}

// TableToMaplist : Convert table to map list.
//
// e.g.:
// ["key1", "key2"] and [[1, 2], [3, 4]]  =>  [{"key1": 1, "key2": 2}, {"key1": 3, "key2": 4}]
//
// If the size of the data row is larger than the size of the header, discard the value.
// If `ignoreZero` is true, do not store zero values.
func TableToMaplist[K cmp.Ordered, V comparable](tblheader []K, tbldata [][]V, ignoreZero bool) (maplist []map[K]V) {
	// 入力チェック
	if len(tblheader) == 0 || len(tbldata) == 0 {
		return
	}

	// マップリストを生成
	maplist = make([]map[K]V, 0, len(tbldata))
	for _, srow := range tbldata {
		mrow := map[K]V{}
		for pos, val := range srow {
			if pos >= len(tblheader) {
				continue
			}
			if ignoreZero && val == *new(V) {
				continue
			}
			mrow[tblheader[pos]] = val
		}
		if len(mrow) == 0 {
			continue
		}
		maplist = append(maplist, mrow)
	}
	if len(maplist) == 0 {
		maplist = nil
	}

	return
}

// ReadJson : Read JSON format data to map list.
func ReadJson[K comparable, V any](r io.Reader) ([]map[K]V, error) {
	// JSONを全文読み込む
	jsonBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	// json.Unmarshal を使用するためにラップする
	jsonBytes = slices.Concat([]byte(`{"body":`), jsonBytes, []byte("}"))

	// オブジェクトに変換
	var jsonObj struct {
		Body []map[K]V `json:"body"`
	}
	if err = json.Unmarshal(jsonBytes, &jsonObj); err != nil {
		return nil, err
	}
	return jsonObj.Body, nil
}

// WriteJson : Write JSON format data from map list.
func WriteJson[K comparable, V any](w io.Writer, maplist []map[K]V) error {
	// オブジェクトから変換
	jsonObj := struct {
		Body []map[K]V `json:"body"`
	}{
		Body: maplist,
	}
	jsonBytes, err := json.Marshal(&jsonObj)
	if err != nil {
		return err
	}

	// 前後の {"body": ... } を除去する
	jsonBytes = jsonBytes[len(`{"body":`) : len(jsonBytes)-len("}")]

	// JSONを出力
	w.Write(jsonBytes)
	return nil
}

// find : Find a element from the slice.
//
// Returns the position of the element if found, or -1 if not found.
func find[K comparable](t []K, v K) int {
	pos := -1
	for i := range t {
		if t[i] == v {
			pos = i
			break
		}
	}
	return pos
}
