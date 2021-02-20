package refl

import (
	"fmt"
	"reflect"

	"github.com/jakubDoka/sterr"
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

// Overwrite takes target and default struct, all fields of target
// that has default value will be overwritten by default struct, this is deep
// copy only in case of raw values, pointers will point to same values
func Overwrite(target, defaults interface{}, onlyZeroValues bool) {

	tv, dv := reflect.ValueOf(target), reflect.ValueOf(defaults)
	if tv.Kind() != reflect.Ptr || dv.Kind() != reflect.Ptr {
		panic("target or defaults is not a pointer")
	}

	tv, dv = tv.Elem(), dv.Elem()
	if tv.Kind() != reflect.Struct || dv.Kind() != reflect.Struct {
		panic("target or defaults is not pointer to struct")
	}

	tt, dt := tv.Type(), dv.Type()

	if tt != dt {
		panic("type of target and defaults does not match")
	}

	for i := 0; i < tv.NumField(); i++ {
		tf, df := tv.Field(i), dv.Field(i)

		if df.IsZero() {
			continue
		}

		if !tf.CanInterface() || !tf.CanSet() {
			continue
		}

		if tf.Kind() == reflect.Struct {
			Overwrite(tf.Addr().Interface(), df.Addr().Interface(), onlyZeroValues)
			continue
		}

		if onlyZeroValues && !tf.IsZero() {
			continue
		}

		tf.Set(df)
	}
}

// Convert related errors
var (
	ErrNested = sterr.New("when parsing field '%s'")
)

// Convert converts one value to another by method, method have to
// be defined for scr and has to have structure func() (dest, error) with value receiver
// if scr does not implement it but it is struct, its values will be
// converted instead if possible, scr and dest has to be passed as pointer
func Convert(scr, dest interface{}, method string) error {
	sv := reflect.ValueOf(scr).Elem()
	st := sv.Type()
	dv := reflect.ValueOf(dest).Elem()
	dt := dv.Type()

	_, ok := st.MethodByName(method)
	if !ok {
		switch st.Kind() {
		case reflect.Struct:
			for i := 0; i < sv.NumField(); i++ {
				f := st.Field(i)
				if _, ok := dt.FieldByName(f.Name); !ok {
					continue
				}

				fs := sv.FieldByName(f.Name)
				fd := dv.FieldByName(f.Name)
				if !fs.CanAddr() || !fs.CanInterface() || !fd.CanAddr() || !fd.CanInterface() {
					continue
				}

				err := Convert(fs.Addr().Interface(), fd.Addr().Interface(), method)
				if err != nil {
					return ErrNested.Args(f.Name).Wrap(err)
				}
			}
			return nil
		}

		if dv.CanSet() {
			dv.Set(sv)
		}

		return nil
	}

	m := sv.MethodByName(method)
	vals := m.Call(nil)
	if !vals[1].IsNil() {
		return vals[1].Interface().(error)
	}

	if dv.CanSet() {
		dv.Set(vals[0])
	}

	return nil
}
