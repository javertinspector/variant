package step

import (
	"github.com/Sirupsen/logrus"
	log "github.com/Sirupsen/logrus"
	"github.com/mumoshu/variant/pkg/util/maputil"
	"reflect"
)

type stepDefImpl struct {
	raw map[string]interface{}
}

type StepDef interface {
	GetName() string
	Raw() map[string]interface{}
	Get(key string) interface{}
	GetStringMapOrEmpty(key string) map[string]interface{}
	Silent() bool
}

func (c stepDefImpl) GetName() string {
	str, ok := c.raw["name"].(string)

	if !ok {
		logrus.Panicf("name wasn't string! name=%s raw=%v", reflect.TypeOf(c.raw["name"]), c.raw)
	}

	return str
}

func (c stepDefImpl) Raw() map[string]interface{} {
	return c.raw
}

func (c stepDefImpl) Get(key string) interface{} {
	return c.raw[key]
}

func (c stepDefImpl) GetStringMapOrEmpty(key string) map[string]interface{} {
	ctx := log.WithField("raw", c.raw[key]).WithField("key", key).WithField("type", reflect.TypeOf(c.raw[key]))

	rawMap, expected := c.Get(key).(map[interface{}]interface{})

	if !expected {
		ctx.Debugf("step config ignored field with unepected type")
		return map[string]interface{}{}
	} else {
		result, err := maputil.CastKeysToStrings(rawMap)

		if err != nil {
			ctx.WithField("error", err).Debugf("step config failed to cast keys to strings")
			return map[string]interface{}{}
		}

		return result
	}
}

func (c stepDefImpl) Silent() bool {
	silent, ok := c.raw["silent"].(bool)
	if !ok {
		silent = false
	}
	return silent
}

func NewStepConfig(raw map[string]interface{}) StepDef {
	return stepDefImpl{
		raw: raw,
	}
}
