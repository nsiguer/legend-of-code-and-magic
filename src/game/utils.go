package game

import (
	"reflect"
)

func in_array(v interface{}, in interface{}) (ok bool, i int) {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for ; i < val.Len(); i++ {
			if ok = v == val.Index(i).Interface(); ok {
				return
			}
		}
	}
	return
}


func iterate_combinations(elems []*Move) [][]*Move{
	moves := make([][]*Move, 0)
	n := len(elems)
	for num:=0;num < (1 << uint(n));num++ {
		combination := []*Move{}
		for ndx:=0;ndx<n;ndx++ {
			// (is the bit "on" in this number?)
			if num & (1 << uint(ndx)) != 0 {
				// (then add it to the combination)
				combination = append(combination, elems[ndx])
			}
		}
		if len(combination) > 0 {
			moves = append(moves, combination)
		
		}
	}
	return moves
}

func permutations(arr []*Move) ([][]*Move){
    var helper func([]*Move, int)
    res := make([][]*Move, 0)

    helper = func(arr []*Move, n int){
        if n == 1{
            tmp := make([]*Move, len(arr))
            copy(tmp, arr)
            res = append(res, tmp)
        } else {
			
            for i := 0; i < n; i++{
                helper(arr, n - 1)
                if n % 2 == 1{
                    tmp := arr[i]
                    arr[i] = arr[n - 1]
                    arr[n - 1] = tmp
                } else {
                    tmp := arr[0]
                    arr[0] = arr[n - 1]
                    arr[n - 1] = tmp
                }
            }
        }
    }
    helper(arr, len(arr))
    return res
}

func del_in_int_array(v int, a []int) (ok bool, r []int) {
    idx := -1
    for i, c := range(a) { if v == c { idx = i ; break } }
    if idx != -1 {
        r = append(a[:idx], a[idx+1:]...)
    }
    return
}