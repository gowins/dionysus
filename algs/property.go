package algs

import "reflect"

func Comparable(i interface{}) bool {
	//only nil without rtype
	//if i == nil {
	//	return false
	//}
	if t := reflect.TypeOf(i); t != nil {
		return t.Comparable()
	}
	return false
}
