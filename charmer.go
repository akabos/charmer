package charmer

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Charm(x interface{}, cfg *viper.Viper, cmd *cobra.Command) error {
	y := reflect.ValueOf(x)
	if y.Kind() != reflect.Ptr {
		return errors.New("pointer expected")
	}
	y = y.Elem()
	if y.Kind() != reflect.Struct {
		return errors.New("pointer to struct expected")
	}
	t := y.Type()
	for i := 0; i < y.NumField(); i++ {
		tf := t.Field(i)
		c, cobraOk := tf.Tag.Lookup("cobra")
		v, viperOk := tf.Tag.Lookup("viper")
		if !cobraOk && !viperOk {
			continue
		}
		if cobraOk && !viperOk {
			return errors.Errorf("cobra tag only allowed alongside viper tag: field=%v", tf.Name)
		}
		if cobraOk && cmd == nil {
			return errors.Errorf("command argument can't be nil alongside cobra tag: field=%v", tf.Name)
		}
		if cobraOk {
			flag := cmd.PersistentFlags().Lookup(c)
			if flag == nil {
				flag = cmd.Flags().Lookup(c)
			}
			if flag == nil {
				return errors.Errorf("flag not found: field=%v cobra=%v", tf.Name, c)
			}
			_ = cfg.BindPFlag(v, flag)
		}
		switch tf.Type.Kind() {
		case reflect.String:
			y.Field(i).SetString(cfg.GetString(v))
		case reflect.Bool:
			y.Field(i).SetBool(cfg.GetBool(v))
		case reflect.Int8:
			y.Field(i).SetInt(cfg.GetInt64(v))
		case reflect.Int16:
			y.Field(i).SetInt(cfg.GetInt64(v))
		case reflect.Int32:
			y.Field(i).SetInt(cfg.GetInt64(v))
		case reflect.Int64:
			y.Field(i).SetInt(cfg.GetInt64(v))
		case reflect.Int:
			y.Field(i).SetInt(cfg.GetInt64(v))
		case reflect.Uint8:
			y.Field(i).SetUint(cfg.GetUint64(v))
		case reflect.Uint16:
			y.Field(i).SetUint(cfg.GetUint64(v))
		case reflect.Uint32:
			y.Field(i).SetUint(cfg.GetUint64(v))
		case reflect.Uint64:
			y.Field(i).SetUint(cfg.GetUint64(v))
		case reflect.Uint:
			y.Field(i).SetUint(cfg.GetUint64(v))
		case reflect.Slice:
			switch tf.Type.Elem().Kind() {
			case reflect.String:
				for _, s := range cfg.GetStringSlice(v) {
					y.Field(i).Set(reflect.Append(y.Field(i), reflect.ValueOf(s)))
				}
			case reflect.Int:
				for _, s := range cfg.GetIntSlice(v) {
					y.Field(i).Set(reflect.Append(y.Field(i), reflect.ValueOf(s)))
				}
			default:
				return errors.Errorf("unsupported slice element: field=%v kind=%v", tf.Name, tf.Type.Elem().Kind().String())
			}
		case reflect.Map:
			return errors.Errorf("not implemented: field=%v kind=%v", tf.Name, tf.Type.Kind())
		default:
			return errors.Errorf("unsupported: field=%v kind=%v", tf.Name, tf.Type.Kind().String())
		}
	}
	return nil
}
