package rt

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"gitee.com/be_clear/xjh/bm"

	"gorm.io/gorm"
)

func genDB(TX *gorm.DB, route *Route) {

	val := reflect.ValueOf(route).Elem()
	tpe := reflect.TypeOf(route).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := tpe.Field(i)
		fieldValue := val.Field(i)
		if fieldValue.IsZero() {
			continue
		}

		switch field.Name {
		case "TABLE":
			TX.Table(route.TABLE)
		case "SELECT":
			TX.Select(route.SELECT)
		case "OMIT":
			TX.Omit(route.OMIT...)
		case "WHERE":
			result := make(map[string]interface{}, 0)
			replaceFieldValues(route.Bind, route.WHERE, result)
			sql, data := map2SqlStr(result)
			TX.Where(sql, data...)
		case "ORDER":
			result := make(map[string]interface{}, 0)
			bindValue := reflect.ValueOf(route.Bind)
			if bindValue.Kind() == reflect.Ptr {
				bindValue = bindValue.Elem()
			}
			for key, val := range route.ORDER {
				var fieldV reflect.Value
				var fieldK reflect.Value

				if bindValue.Kind() == reflect.Struct {
					fieldV = bindValue.FieldByName(val)
					fieldK = bindValue.FieldByName(key)
				}

				if fieldV.IsValid() && fieldK.IsValid() && isSimpleType(fieldK.Interface()) && isSimpleType(fieldV.Interface()) {
					result[fieldK.String()] = fieldV.Interface()
				}
			}
			for key, val := range result {
				if reflect.TypeOf(val).Kind() == reflect.Bool {
					if reflect.ValueOf(val).Bool() {
						TX.Order(key + " ASC")
					} else {
						TX.Order(key + " DESC")
					}
				}
			}
		case "NOT":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			result := make(map[string]interface{}, 0)
			replaceFieldValues(route.Bind, route.NOT, result)
			sql, data := map2SqlStr(result)
			TX.Not(sql, data...)
		case "OR":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			result := make(map[string]interface{}, 0)
			replaceFieldValues(route.Bind, route.OR, result)
			sql, data := map2SqlStr(result)
			TX.Or(sql, data...)
		case "JOINS":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			TX.Joins(route.JOINS)
		case "INNER_JOINS":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			TX.InnerJoins(route.INNER_JOINS)
		case "PRELOAD":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			TX.Preload(route.PRELOAD)
		case "GROUP":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			TX.Group(route.GROUP)
		case "DISTINCT":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			for _, distinct := range route.DISTINCT {
				TX.Distinct(distinct)
			}
		case "Limit":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			TX.Limit(route.LIMIT)
		case "MAP_COLUMNS":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			TX.MapColumns(route.MAP_COLUMNS)
		case "HAVING":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			result := make(map[string]interface{}, 0)
			replaceFieldValues(route.Bind, route.HAVING, result)
			sql, data := map2SqlStr(result)
			TX.Having(sql, data...)
		case "OFFSET":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			TX.Offset(route.OFFSET)
		case "RAW":
			fmt.Printf("%s 是未测试功能, 请联系管理员", route.Type)
			result := make(map[string]interface{}, 0)
			replaceFieldValues(route.Bind, route.RAW, result)
			sql, data := map2SqlStr(result)
			TX.Raw(sql, data...)
		}
	}
}

func genType(TX *gorm.DB, route *Route) *bm.Response {
	noVModelVal, _ := createDataInstance(route.MODEL, false)
	model := noVModelVal.Interface()

	switch route.Type {
	case TYPE.GET_LIST:
		// 提取所需的字段
		paginationField := reflect.ValueOf(route.Bind).Elem().FieldByName("Pagination")
		if !paginationField.IsValid() {
			return new(bm.Response).FailFront("无效的分页数据")
		}
		pageSize := int(paginationField.FieldByName("PageSize").Int())
		current := int(paginationField.FieldByName("Current").Int())

		if pageSize == 0 || current == 0 {
			return new(bm.Response).FailFront("无效的分页数据")
		}
		var total int64

		// 处理数据库操作时的错误
		if err := TX.Count(&total).Error; err != nil {
			return new(bm.Response).FailBackend(fmt.Sprintf("计数失败 %s", err.Error()))
		}

		result := reflect.New(reflect.SliceOf(noVModelVal.Type())).Interface()

		if err := TX.Limit(pageSize).Offset((current - 1) * pageSize).Find(result).Error; err != nil {
			return new(bm.Response).FailBackend(fmt.Sprintf("数据查询失败 %s", err.Error()))
		}

		return new(bm.Response).Suc(bm.Pagination{
			PageSize: pageSize,
			Current:  current,
			Data:     result,
			Total:    int(total),
		}, "")

	case TYPE.GET_ONE:
		result := TX.First(model)

		if result.Error == gorm.ErrRecordNotFound {
			return new(bm.Response).FailFront("没有该数据")
		}
		if result.Error != nil {
			return new(bm.Response).FailBackend(fmt.Sprintf("数据查询失败 %s", result.Error))
		}
		return new(bm.Response).Suc(model, "")

	case TYPE.UPDATE_ONE:
		result := TX.First(model).Updates(route.Bind)
		if result.Error == gorm.ErrRecordNotFound {
			return new(bm.Response).FailFront("没有该数据")
		}
		if result.Error != nil {
			return new(bm.Response).FailBackend(fmt.Sprintf("更新失败 %s", result.Error))
		}
		return new(bm.Response).Suc(model, "")

	case TYPE.UPDATE_LIST:
		log.Fatal("未实现")

	case TYPE.CREATE_ONE:
		err := AssignMatchingFields(route.Bind, route.MODEL)
		if err != nil {
			return new(bm.Response).FailBackend(fmt.Sprintf("创建失败 %s", err.Error()))
		}
		result := TX.Create(route.MODEL)
		if result.Error != nil {
			return new(bm.Response).FailBackend(fmt.Sprintf("创建失败 %s", result.Error))
		}

		return new(bm.Response).Suc(route.MODEL, "")

	case TYPE.CREATE_LIST:
		result := TX.Create(route.Bind)
		if result.Error != nil {
			return new(bm.Response).FailBackend(fmt.Sprintf("创建失败 %s", result.Error))
		}
		return new(bm.Response).Suc(nil, "")
	case TYPE.DELETE_LIST:
		log.Fatal("未实现")

	case TYPE.DELETE:
		err := AssignMatchingFields(route.Bind, route.MODEL)
		if err != nil {
			return new(bm.Response).FailBackend(fmt.Sprintf("删除失败 %s", err.Error()))
		}

		result := TX.First(route.MODEL).Delete(route.MODEL)
		if result.Error == gorm.ErrRecordNotFound {
			return new(bm.Response).FailFront("没有该数据")
		}
		if result.Error != nil {
			return new(bm.Response).FailBackend(fmt.Sprintf("删除失败 %s", result.Error))
		}
		return new(bm.Response).Suc(route.MODEL, "")

	default:
		log.Fatal("无效的操作类型")
	}

	return nil
}

// 判断是否为简单类型
func isSimpleType(value interface{}) bool {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func replaceFieldValues(data interface{}, fieldValue map[string]string, result map[string]interface{}) {
	// 将 data 转换为 reflect.Value
	dataValue := reflect.ValueOf(data)

	// 如果 data 是指针类型，获取指向的值
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	// 确保 data 是一个结构体或者 map
	if dataValue.Kind() != reflect.Struct {
		return
	}

	for key, val := range fieldValue {
		var field reflect.Value

		// 如果 data 是一个结构体，从结构体字段获取值
		if dataValue.Kind() == reflect.Struct {
			field = dataValue.FieldByName(val)
		}

		// 检查字段是否存在并且是简单类型
		if field.IsValid() && isSimpleType(field.Interface()) {
			result[key] = field.Interface()
		}
	}
}

func map2SqlStr(data map[string]interface{}) (string, []interface{}) {
	sql := make([]string, 0)
	vals := make([]interface{}, 0)
	for key, val := range data {
		if isSimpleType(val) {
			sql = append(sql, fmt.Sprintf("%s = ?", key))
			vals = append(vals, val)
		}
	}
	return strings.Join(sql, "AND"), vals
}

// AssignMatchingFields 将 data 的字段值赋值给 model 中相应的字段
func AssignMatchingFields(data interface{}, model interface{}) error {
	// 通过反射获取指针指向的实际值
	dataVal := reflect.ValueOf(data).Elem()
	modelVal := reflect.ValueOf(model).Elem()

	// 确保 data 和 model 都是结构体
	if dataVal.Kind() != reflect.Struct || modelVal.Kind() != reflect.Struct {
		return fmt.Errorf("data 和 model 都必须是指向结构体的指针")
	}

	for i := 0; i < dataVal.NumField(); i++ {
		dataField := dataVal.Type().Field(i)
		modelField := modelVal.FieldByName(dataField.Name)

		// 检查 model 中是否存在对应的字段，并且类型相同
		if modelField.IsValid() && modelField.Type() == dataField.Type {
			if modelField.CanSet() {
				modelField.Set(dataVal.Field(i))
			}
		}
	}

	return nil
}
