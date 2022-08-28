package homework

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var errInvalidEntity = errors.New("invalid entity")

func InsertStmt(entity interface{}) (string, []interface{}, error) {

	val := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)

	if typ == nil || isInvalidType(val) {
		return "", nil, errInvalidEntity
	}

	fieldArray, valueArray, err := getFieldAndValue(val)
	tableName := getTableName(val)
	if err != nil {
		return "", nil, err
	}

	return buildSql(fieldArray, valueArray, tableName), valueArray, nil
}

func isInvalidType(val reflect.Value) bool {
	kind := val.Kind()
	switch kind {
	case reflect.Pointer:
		return val.IsNil()
	case reflect.Struct:
		return val.NumField() == 0
	default:
		return false
	}
}

func getFieldAndValue(val reflect.Value) ([]string, []any, error) {
	return processFieldAndValue(val)
}

func processFieldAndValue(val reflect.Value) ([]string, []any, error) {
	typ := val.Type()

	for typ.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = typ.Elem()
	}

	var (
		fieldArray = make([]string, 0)
		valueArray = make([]any, 0)
	)

	fmt.Printf("进入方法，name:%s  kind:%s  value:%v\n", typ.Name(), typ.Kind(), val)

	switch typ.Kind() {
	case reflect.Struct:

		for i := 0; i < val.NumField(); i++ {

			var thisErr error
			var thisFieldArrays []string
			var thisValueArrays []any

			if val.Field(i).Kind() == reflect.Struct && typ.Field(i).Anonymous {
				thisFieldArrays, thisValueArrays, thisErr = processFieldAndValue(val.Field(i))
			}

			if val.Field(i).Kind() == reflect.Struct && typ.Field(i).Anonymous && thisErr == nil {
				combineArray(&fieldArray, &thisFieldArrays, &valueArray, &thisValueArrays)
			} else {
				fieldName := typ.Field(i).Name
				value := val.Field(i).Interface()
				fieldArray = append(fieldArray, fieldName)
				valueArray = append(valueArray, value)
			}
		}
	default:
		return nil, nil, errInvalidEntity
	}

	return fieldArray, valueArray, nil
}

func getTableName(val reflect.Value) string {
	typ := val.Type()

	for typ.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = typ.Elem()
	}

	return typ.Name()
}

func buildSql(fieldArray []string, valueArray []any, tableName string) string {
	var (
		isFirst  = true
		bd       = strings.Builder{}
		valuesBD = strings.Builder{}
	)

	bd.WriteString("INSERT INTO `")
	bd.WriteString(tableName)
	bd.WriteString("`(")

	for _, k := range fieldArray {

		if !isFirst {
			bd.WriteString(",")
			valuesBD.WriteString(",")
		}

		bd.WriteString("`")
		bd.WriteString(k)
		bd.WriteString("`")

		valuesBD.WriteString("?")

		isFirst = false
	}

	bd.WriteString(") VALUES(")
	bd.WriteString(valuesBD.String())
	bd.WriteString(");")

	return bd.String()
}

func combineArray(oriField, curField *[]string, oriValue, curValue *[]any) {

	for i, curK := range *curField {
		var (
			isContain = false
		)

		for _, oriK := range *oriField {
			if curK == oriK {
				isContain = true
				break
			}
		}

		if isContain {
			continue
		}
		*oriField = append(*oriField, curK)
		*oriValue = append(*oriValue, (*curValue)[i])

	}
}
