package refl

import (
	"fmt"
	"reflect"
)

// AssertHomogeneity returns whether struct contains only sup datatype
func AssertHomogeneity(val, sup reflect.Type) error {
	if sup == val {
		return nil
	} else if val.Kind() == reflect.Slice {
		return AssertHomogeneity(val.Elem(), sup)
	} else if val.Kind() != reflect.Struct {
		return fmt.Errorf("type contains %v when only %v is expected", val.Name(), sup.Name())
	}

	fn := val.NumField()

	for i := 0; i < fn; i++ {
		if err := AssertHomogeneity(val.Field(i).Type, sup); err != nil {
			return err
		}
	}

	return nil
}

// OverwriteDefault takes target and default struct, all fields of target
// that has default value will be overwritten by default struct, this is deep
// copy only in case of raw values, pointers will point to same values
func OverwriteDefault(target, defaults interface{}) {
	tv, dv := reflect.ValueOf(target).Elem(), reflect.ValueOf(defaults).Elem()
	tt, dt := tv.Type(), dv.Type()

	if tt != dt {
		panic("type of target and defaults does not match")
	}

	for i := 0; i < tv.NumField(); i++ {
		tf, df := tv.Field(i), dv.Field(i)

		if !tf.CanInterface() || !tf.CanSet() {
			continue
		}

		if tf.Kind() == reflect.Struct {
			OverwriteDefault(tf.Addr().Interface(), df.Addr().Interface())
			continue
		}

		if !tf.IsZero() {
			continue
		}

		tf.Set(df)
	}
}
