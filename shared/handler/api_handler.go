package handler

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/esmailemami/chess/shared/errs"
	"github.com/esmailemami/chess/shared/util"
	"github.com/gin-gonic/gin"
)

func HandleAPI(fn interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			args       = make([]reflect.Value, 1)
			params     = make([]string, len(ctx.Params))
			paramsCall = 0
			fnType     = reflect.TypeOf(fn)
		)
		args[0] = reflect.ValueOf(ctx)

		// get the parameters in order from the request
		for i, param := range ctx.Params {
			params[i] = param.Value
		}

		for i := 1; i < fnType.NumIn(); i++ {
			var (
				argType  = fnType.In(i)
				argValue = reflect.New(argType).Elem()
				err      error
			)

			// if the argument is struct type in must be a query params from GET request or body from POST request
			if argType.Kind() == reflect.Struct && !argType.AssignableTo(reflect.TypeOf(time.Time{})) {
				switch ctx.Request.Method {
				case http.MethodGet:
					err = fillStructFromQuery(argValue, ctx)
				case http.MethodPost, http.MethodPut:
					err = ctx.ShouldBindJSON(argValue.Addr().Interface())
				default:
					return
				}
			} else {
				// the argument value is in the route
				err = util.SetReflectValueFromString(argValue, params[paramsCall])
				paramsCall++
			}

			if err != nil {
				errs.ErrorHandler(ctx.Writer, errs.BadRequestErr().WithError(err))
				return
			}

			args = append(args, argValue)
		}

		// call the function and handle the error
		result := reflect.ValueOf(fn).Call(args)

		err := result[1].Interface()

		if err != nil {
			errs.ErrorHandler(ctx.Writer, err.(error))
			return
		}

		resp := result[0].Elem()
		ctx.JSON(resp.FieldByName("Status").Interface().(int), resp.Interface())
	}
}

func fillStructFromQuery(strct reflect.Value, ctx *gin.Context) error {
	for i := 0; i < strct.NumField(); i++ {
		var (
			field       = strct.Field(i)
			structField = strct.Type().Field(i)
			jsonTag     = structField.Tag.Get("json")
			defaultTag  = structField.Tag.Get("default")
			fieldName   string
		)

		if jsonTag == "-" {
			continue
		}

		if jsonTag != "" {
			fieldName = strings.Split(jsonTag, ",")[0]
		} else {
			fieldName = field.Type().Name()
		}

		if paramValue, ok := ctx.GetQuery(fieldName); ok {
			if err := util.SetReflectValueFromString(field, paramValue); err != nil {
				return err
			}
		} else if defaultTag != "" {
			if err := util.SetReflectValueFromString(field, defaultTag); err != nil {
				return err
			}
		}
	}

	return nil
}
