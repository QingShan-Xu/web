package rt

import (
	"log"
	"reflect"
)

func IsGroupRouter(router *Router) bool {
	if len(router.Children) == 0 {
		return false
	}
	return true
}

func checkItem(route Router) {
	var name string
	if route.Path != "" {
		name = route.Path
	}
	if route.Name != "" {
		name = route.Name
	}

	if route.Method == "" {
		log.Fatalf("route: %s Method is Enpty", name)
	}

	if route.Type != "" && route.MODEL == nil {
		log.Fatalf("route: %s has Type but not has Model", name)
	}
	if route.Type != "" && route.Handler != nil {
		log.Fatalf("route: %s Cannot have both attributes Type and Handler simultaneously", name)
	}

	if route.Type == "" && route.Handler == nil {
		log.Fatalf("route: %s not has func to handle", name)
	}

	if route.Type == "GET_LIST" && (route.LIMIT != 0 || route.OFFSET != 0) {
		log.Fatalf("route: %s Cannot have LIMIT or OFFSET when Type = GET_LIST", name)
	}

	if (route.Type == "GET_ONE" || route.Type == "GET_LIST") && route.Method != "GET" {
		log.Fatalf("route: %s Cannot have Type = %s and Method = %s", name, route.Type, route.Method)
	}

	// if strings.Contains(route.Path, ":") && route.Data != nil {
	// 	log.Fatalf("route: %s Cannot have both Query and Data", name)
	// }

	// if route.MODEL == nil {
	// 	allowedMap := map[string]struct{}{
	// 		"NoToken":     {},
	// 		"Handler":     {},
	// 		"Method":      {},
	// 		"Type":        {},
	// 		"Bind":        {},
	// 		"Path":        {},
	// 		"Middlewares": {},
	// 		"Children":    {},
	// 		"Name":        {},
	// 		"MODEL":       {},
	// 		"TABLE":       {},
	// 		"DISTINCT":    {},
	// 		"SELECT":      {},
	// 		"OMIT":        {},
	// 		"MAP_COLUMNS": {},
	// 		"WHERE":       {},
	// 		"NOT":         {},
	// 		"OR":          {},
	// 		"INNER_JOINS": {},
	// 		"JOINS":       {},
	// 		"GROUP":       {},
	// 		"HAVING":      {},
	// 		"ORDER":       {},
	// 		"LIMIT":       {},
	// 		"OFFSET":      {},
	// 		"PRELOAD":     {},
	// 		"RAW":         {},
	// 	}

	// 	filedNmae, isAllowed := check(allowedMap, route)
	// 	if !isAllowed {
	// 		log.Fatalf("route: %s not allowed field: %s", name, filedNmae)
	// 	}
	// }
}

func checkGroup(route Router) {
	// var name string
	// if route.Path != "" {
	// 	name = route.Path
	// }
	// if route.Name != "" {
	// 	name = route.Name
	// }

	// allowedMap := map[string]struct{}{
	// 	"NoToken":     {},
	// 	"Path":        {},
	// 	"Middlewares": {},
	// 	"Children":    {},
	// 	"Name":        {},
	// }
	// filedNmae, isAllowed := check(allowedMap, route)
	// if !isAllowed {
	// 	log.Fatalf("routeGroup: %s not allowed field: %s", name, filedNmae)
	// }
}

func check(allowedMap map[string]struct{}, beChecked interface{}) (string, bool) {
	tpe := reflect.TypeOf(beChecked)
	val := reflect.ValueOf(beChecked)
	for i := 0; i < val.NumField(); i++ {
		field := tpe.Field(i)
		fieldValue := val.Field(i)
		if _, allowed := allowedMap[field.Name]; !allowed {
			if !isZero(fieldValue) {
				return field.Name, false
			}
		}
	}
	return "", true
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Func, reflect.Ptr, reflect.Interface, reflect.Chan:
		return v.IsNil()
	default:
		return v.Interface() == reflect.Zero(v.Type()).Interface()
	}
}
