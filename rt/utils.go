package rt

import "reflect"

func getInstanceVal(anyData interface{}) reflect.Value {
	instanceTpe := reflect.TypeOf(anyData)
	instanceVal := reflect.ValueOf(anyData)

	if instanceTpe.Kind() == reflect.Pointer {
		instanceVal = instanceVal.Elem()
	}

	return instanceVal
}
