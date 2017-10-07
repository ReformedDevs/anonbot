package server

import (
	"net/url"
	"reflect"
	"strconv"
)

// copyStruct enumerates the exported fields in src and copies their values to
// the equivalent fields in dest. It is assumed that fields with the same name
// share the same type, otherwise a panic may result.
func (s *Server) copyStruct(src interface{}, dest interface{}) {
	var (
		srcType = reflect.TypeOf(src).Elem()
		srcVal  = reflect.ValueOf(src).Elem()
		destVal = reflect.ValueOf(dest).Elem()
	)
	for i := 0; i < srcType.NumField(); i++ {
		fDestVal := destVal.FieldByName(srcType.Field(i).Name)
		if fDestVal.Kind() != reflect.Invalid {
			fDestVal.Set(srcVal.Field(i))
		}
	}
}

// populateStruct loads the values from the form into the specified struct.
func (s *Server) populateStruct(form url.Values, v interface{}) {
	var (
		vType = reflect.TypeOf(v).Elem()
		vVal  = reflect.ValueOf(v).Elem()
	)
	for i := 0; i < vType.NumField(); i++ {
		var (
			fVal = vVal.Field(i)
			s    = form.Get(vType.Field(i).Name)
		)
		switch fVal.Kind() {
		case reflect.String:
			fVal.SetString(s)
		case reflect.Int64:
			iVal, _ := strconv.ParseInt(s, 10, 64)
			fVal.SetInt(iVal)
		}
	}
}
