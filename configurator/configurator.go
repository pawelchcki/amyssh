package configurator

import (
	// "fmt"
	"io/ioutil"
	"reflect"

	"launchpad.net/goyaml"
)

func meldStructs(parent, target interface{}) {
	meldValueStructs(reflect.ValueOf(parent), reflect.ValueOf(target))
}

func meldValueStructs(parent, target reflect.Value) {
	if target.Kind() == reflect.Ptr {
		target = target.Elem()
	}
	if parent.Kind() == reflect.Ptr {
		parent = parent.Elem()
	}
	// fmt.Printf("format", ...)
	for i, n := 0, parent.NumField(); i < n; i++ {
		if target.Field(i).CanSet() {
			switch parent.Field(i).Kind() {
			case reflect.Struct:
				meldValueStructs(parent.Field(i), target.Field(i))

			case reflect.Array, reflect.Slice, reflect.String:
				if parent.Field(i).Len() > 0 {
					target.Field(i).Set(parent.Field(i))
				}

			default:
				if reflect.Zero(target.Field(i).Type()).Interface() != parent.Field(i).Interface() {
					target.Field(i).Set(parent.Field(i))
				}
			}
		}
	}
}

func NewConfig(cfg interface{}) *ConfigHolder {
	holder := ConfigHolder{}
	holder.defaultCfg = reflect.ValueOf(cfg)
	if holder.defaultCfg.Kind() != reflect.Struct {
		panic("NewConfig requires struct type passed to it.")
	}
	return &holder
}

func (c *ConfigHolder) LoadFromFile(filePath string) {
	c.fileCfg = reflect.New(c.defaultCfg.Type())

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		// panic(err)
		return
	}
	err = goyaml.Unmarshal([]byte(contents), c.fileCfg.Interface())
	if err != nil {
		panic(err)
	}
}

func (c *ConfigHolder) Config() interface{} {
	if !c.cachedConfig.IsValid() {
		c.cachedConfig = reflect.New(c.defaultCfg.Type())
		meldValueStructs(c.defaultCfg, c.cachedConfig)

		if c.optionCfg.IsValid() {
			meldValueStructs(c.optionCfg, c.cachedConfig)
		}

		if len(c.ConfigFilePath) > 0 {
			c.LoadFromFile(c.ConfigFilePath)
		}

		if c.fileCfg.IsValid() {
			meldValueStructs(c.fileCfg, c.cachedConfig)
		}
	}
	return c.cachedConfig.Elem().Interface()
}

func (c *ConfigHolder) Options() interface{} {
	if !c.optionCfg.IsValid() {
		c.optionCfg = reflect.New(c.defaultCfg.Type())
	}
	return c.optionCfg.Interface()
}

type ConfigHolder struct {
	defaultCfg   reflect.Value
	optionCfg    reflect.Value
	fileCfg      reflect.Value
	cachedConfig reflect.Value

	ConfigFilePath string
}

var globalConfig *ConfigHolder

func Config() interface{} {
	if globalConfig == nil {
		panic("config.Global not initialized")
	}
	return globalConfig.Config()
}

func Options() interface{} {
	if globalConfig == nil {
		panic("config.Global not initialized")
	}
	return globalConfig.Options()
}

func Initialize(cfg interface{}) *ConfigHolder {
	globalConfig = NewConfig(cfg)
	return globalConfig
}

func ConfigFilePath() *string {
	if globalConfig == nil {
		panic("config.Global not initialized")
	}
	return &globalConfig.ConfigFilePath
}
