package utils

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"reflect"
	"sort"
)

func In(strArray []string, target string) bool {
	sort.Strings(strArray)
	index := sort.SearchStrings(strArray, target)
	if index < len(strArray) && strArray[index] == target {
		return true
	}
	return false
}
func Index[T any](v T, array []T) int {
	if n := len(array); array != nil && n != 0 {
		i := 0
		for !reflect.DeepEqual(v, array[i]) {
			i++
		}
		if i != n {
			return i
		}
	}
	return -1
}

func DeleteSlice(slice interface{}, index int) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	length := sliceValue.Len()
	if slice == nil || length == 0 || (length-1) < index {
		return nil, errors.New("error")
	}
	if length-1 == index {
		return sliceValue.Slice(0, index).Interface(), nil
	} else if (length - 1) >= index {
		return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface(), nil
	}
	return nil, errors.New("error")
}

func MakeEmptyList(accentColor color.Color) *fyne.Container {
	text := canvas.NewText("无数据 ...", accentColor)
	text.TextSize = 20
	cEmpty := container.NewCenter(text)
	return cEmpty
}
