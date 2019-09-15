package force

import (
	"reflect"

	"github.com/gravitational/trace"
)

// EmptyContext returns empty execution context
func EmptyContext() ExecutionContext {
	return WithRuntimeScope(nil)
}

// EvalInto evaluates variable in within the execution context
// into variable out
func EvalInto(ctx ExecutionContext, inRaw, out interface{}) error {
	if inRaw == nil {
		return nil
	}
	in, err := Eval(ctx, inRaw)
	if err != nil {
		return trace.Wrap(err)
	}
	if in == nil {
		return nil
	}
	inType := reflect.TypeOf(in)
	outType := reflect.TypeOf(out)
	if outType.Kind() != reflect.Ptr {
		return trace.BadParameter("out should be a pointer, got %T(%v)", out, outType.Kind())
	}
	outVal := reflect.ValueOf(out)

	switch inType.Kind() {
	case reflect.Struct:
		if outType.Elem().Kind() != reflect.Struct {
			return trace.BadParameter("in is %v then out should be pointer to struct, got %T", inType, out)
		}
		inVal := reflect.ValueOf(in)
		for i := 0; i < inType.NumField(); i++ {
			fieldVal := inVal.Field(i)
			fieldType := inType.Field(i)
			if fieldVal.Interface() == nil {
				continue
			}
			eval, err := Eval(ctx, fieldVal.Interface())
			if err != nil {
				return trace.Wrap(err)
			}
			if eval == nil {
				continue
			}
			if fieldType.Name == metadataFieldName || fieldType.Tag.Get(codeTag) == codeSkip {
				continue
			}
			outField := outVal.Elem().FieldByName(fieldType.Name)
			if !outField.IsValid() {
				return trace.NotFound("struct %T has no field %v", out, fieldType.Name)
			}
			// simple case, can assign two primitive evaluated types
			if outField.Type().AssignableTo(reflect.TypeOf(eval)) {
				outField.Set(reflect.ValueOf(eval))
			} else {
				evalType := reflect.TypeOf(eval)
				evalValue := reflect.ValueOf(eval)
				if evalType.Kind() == reflect.Ptr && !evalValue.Elem().IsValid() {
					continue
				}
				if evalType.Kind() == reflect.Ptr && evalType.Elem().Kind() == reflect.Struct {
					tempVal := reflect.New(OriginalType(evalType.Elem()))
					err := EvalInto(ctx, reflect.ValueOf(eval).Elem().Interface(), tempVal.Interface())
					if err != nil {
						return trace.Wrap(err)
					}
					outField.Set(tempVal)
				} else {
					if err := EvalInto(ctx, eval, outField.Addr().Interface()); err != nil {
						return trace.Wrap(err)
					}
				}
			}
		}
		return nil
	case reflect.Ptr:
		return trace.BadParameter("can't evaluate %v(%T) into %v(%T)", in, in, out, out)
	case reflect.Slice:
		inVal := reflect.ValueOf(in)
		outSlice := reflect.MakeSlice(outType.Elem(), inVal.Len(), inVal.Len())
		for i := 0; i < inVal.Len(); i++ {
			elem := inVal.Index(i).Interface()
			if err := EvalInto(ctx, elem, outSlice.Index(i).Addr().Interface()); err != nil {
				return trace.Wrap(err)
			}
		}
		outVal.Elem().Set(outSlice)
		return nil
	default:
		evaluated, err := Eval(ctx, in)
		if err != nil {
			return trace.Wrap(err)
		}
		outElem := outVal.Elem()
		if !outElem.CanSet() {
			return trace.BadParameter("can't set value of %v(%T) to %v(%T)", out, out, evaluated, evaluated)
		}
		outElem.Set(reflect.ValueOf(evaluated))
		return nil
	}
}

// Eval evaluates variable based on the execution context
func Eval(ctx ExecutionContext, variable interface{}) (interface{}, error) {
	switch v := variable.(type) {
	case []interface{}:
		outSlice := make([]interface{}, len(v))
		for i := range v {
			out, err := Eval(ctx, v[i])
			if err != nil {
				return nil, trace.Wrap(err)
			}
			outSlice[i] = out
		}
		return outSlice, nil
	case IntVar:
		return v.Eval(ctx)
	case BoolVar:
		return v.Eval(ctx)
	case StringVar:
		return v.Eval(ctx)
	case StringsVar:
		return v.Eval(ctx)
	default:
		return v, nil
	}
}

// EvalString evaluates empty or missing string into ""
func EvalString(ctx ExecutionContext, v StringVar) (string, error) {
	if v == nil {
		return "", nil
	}
	return v.Eval(ctx)
}

// EvalBool evaluates empty or unspecified bool to false
func EvalBool(ctx ExecutionContext, in BoolVar) (bool, error) {
	if in == nil {
		return false, nil
	}
	return in.Eval(ctx)
}

func EvalPInt64(ctx ExecutionContext, in IntVar) (*int64, error) {
	if in == nil {
		return nil, nil
	}
	out, err := in.Eval(ctx)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	val := int64(out)
	return &val, nil
}

func EvalPInt32(ctx ExecutionContext, in IntVar) (*int32, error) {
	if in == nil {
		return nil, nil
	}
	out, err := in.Eval(ctx)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	val := int32(out)
	return &val, nil
}

func PInt32(in int32) *int32 {
	return &in
}

func PInt64(in int64) *int64 {
	return &in
}