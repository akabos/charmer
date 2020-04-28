package charmer

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleEnv() {
	var args struct {
		Example string `viper:"example"`
	}

	_ = os.Setenv("EXAMPLE", "example")

	cfg := viper.New()
	_ = cfg.BindEnv("example")

	if err := Charm(&args, cfg, nil); err != nil {
		panic(err)
	}

	fmt.Println(args.Example)

	// Output:
	// example
}

func ExampleFlag() {

	var args struct {
		Flag string `cobra:"flag" viper:"env"`
	}

	cmd := &cobra.Command{}
	cmd.Flags().String("flag", "default", "")
	_ = cmd.ParseFlags([]string{"--flag=example"})

	cfg := viper.New()

	if err := Charm(&args, cfg, cmd); err != nil {
		panic(err)
	}

	fmt.Println(args.Flag)

	// Output:
	// example
}

func TestCharm_Cobra(t *testing.T) {
	cmd := &cobra.Command{}
	cfg := viper.New()

	t.Run("flag", func(t *testing.T) {
		var args struct {
			Value string `cobra:"flag" viper:"config"`
		}
		cmd.ResetFlags()
		cmd.Flags().String("flag", "flag-default", "")
		_ = cmd.ParseFlags([]string{"--flag=flag-value"})

		err := Charm(&args, cfg, cmd)
		require.NoError(t, err)
		require.Equal(t, "flag-value", args.Value)
	})

	t.Run("flag default", func(t *testing.T) {
		var args struct {
			Value string `cobra:"flag" viper:"config"`
		}
		cmd.ResetFlags()
		cmd.Flags().String("flag", "flag-default", "")

		err := Charm(&args, cfg, cmd)
		require.NoError(t, err)
		require.Equal(t, "flag-default", args.Value)
	})

	t.Run("persistent flag", func(t *testing.T) {
		var args struct {
			Value string `cobra:"flag" viper:"config"`
		}
		cmd.ResetFlags()
		cmd.PersistentFlags().String("flag", "flag-default", "")
		_ = cmd.ParseFlags([]string{"--flag=flag-value"})

		err := Charm(&args, cfg, cmd)
		require.NoError(t, err)
		require.Equal(t, "flag-value", args.Value)
	})

	t.Run("persistent flag default", func(t *testing.T) {
		var args struct {
			Value string `cobra:"flag" viper:"config"`
		}
		cmd.ResetFlags()
		cmd.PersistentFlags().String("flag", "flag-default", "")

		err := Charm(&args, cfg, cmd)
		require.NoError(t, err)
		require.Equal(t, "flag-default", args.Value)
	})
}

func TestCharm_Viper(t *testing.T) {
	var args struct {
		Value string `viper:"config"`
	}

	cfg := viper.New()
	cfg.Set("config", "config")

	err := Charm(&args, cfg, nil)
	require.NoError(t, err)

	require.Equal(t, "config", args.Value)
}

func TestCharm_Errors(t *testing.T) {
	t.Run("not a pointer", func(t *testing.T) {
		err := Charm(struct{}{}, nil, nil)
		if assert.Error(t, err) {
			assert.Equal(t, "pointer expected", err.Error())
		}
	})
	t.Run("not a pointer to struct", func(t *testing.T) {
		v := 0
		err := Charm(&v, nil, nil)
		if assert.Error(t, err) {
			assert.Equal(t, "pointer to struct expected", err.Error())
		}
	})
	t.Run("incomplete spec", func(t *testing.T) {
		var args struct {
			Value string `cobra:"value"`
		}
		err := Charm(&args, nil, nil)
		if assert.Error(t, err) {
			assert.Equal(t, "cobra tag only allowed alongside viper tag: field=Value", err.Error())
		}
	})
	t.Run("flag not found", func(t *testing.T) {
		var args struct {
			Value string `cobra:"value" viper:"value"`
		}
		err := Charm(&args, viper.New(), &cobra.Command{})
		if assert.Error(t, err) {
			assert.Equal(t, "flag not found: field=Value cobra=value", err.Error())
		}
	})
	t.Run("unsupported field type", func(t *testing.T) {
		var args struct {
			Value *struct{} `cobra:"value" viper:"value"`
		}
		cmd := &cobra.Command{}
		cmd.Flags().String("value", "", "")

		err := Charm(&args, viper.New(), cmd)
		if assert.Error(t, err) {
			assert.Equal(t, "unsupported: field=Value kind=ptr", err.Error())
		}
	})
}

func TestCharm_Types(t *testing.T) {
	var args struct {
		String      string   `viper:"string"`
		Bool        bool     `viper:"bool"`
		Int8        int8     `viper:"int8"`
		Int16       int16    `viper:"int16"`
		Int32       int32    `viper:"int32"`
		Int64       int64    `viper:"int64"`
		Int         int      `viper:"int"`
		Uint8       uint8    `viper:"uint8"`
		Uint16      uint16   `viper:"uint16"`
		Uint32      uint32   `viper:"uint32"`
		Uint64      uint64   `viper:"uint64"`
		Uint        uint     `viper:"uint"`
		StringSlice []string `viper:"string-slice"`
		IntSlice    []int    `viper:"int-slice"`
	}

	cfg := viper.New()
	cfg.Set("string", "string")
	cfg.Set("bool", true)
	cfg.Set("int", 1)
	cfg.Set("int8", int8(1))
	cfg.Set("int16", int16(1))
	cfg.Set("int32", int32(1))
	cfg.Set("int64", int64(1))
	cfg.Set("uint", uint(1))
	cfg.Set("uint8", uint8(1))
	cfg.Set("uint16", uint16(1))
	cfg.Set("uint32", uint32(1))
	cfg.Set("uint64", uint64(1))
	cfg.Set("string-slice", []string{"value"})
	cfg.Set("int-slice", []int{1})

	err := Charm(&args, cfg, nil)
	require.NoError(t, err)

	assert.Equal(t, "string", args.String)
	assert.Equal(t, true, args.Bool)
	assert.Equal(t, int(1), args.Int)
	assert.Equal(t, int8(1), args.Int8)
	assert.Equal(t, int16(1), args.Int16)
	assert.Equal(t, int32(1), args.Int32)
	assert.Equal(t, int64(1), args.Int64)
	assert.Equal(t, uint(1), args.Uint)
	assert.Equal(t, uint8(1), args.Uint8)
	assert.Equal(t, uint16(1), args.Uint16)
	assert.Equal(t, uint32(1), args.Uint32)
	assert.Equal(t, uint64(1), args.Uint64)
	assert.Equal(t, []string{"value"}, args.StringSlice)
	assert.Equal(t, []int{1}, args.IntSlice)
}
