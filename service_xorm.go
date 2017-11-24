package raweb

import (
	"errors"
	"github.com/go-xorm/xorm"
	"regexp"
	"fmt"
	"github.com/go-xorm/core"
)

type XOrmService struct {
	databases map[string]*xorm.Engine
}

var DefaultXOrmService XOrmService

var driverStringRegex = regexp.MustCompile("^([^:]+)://(.*)$")

func (service *XOrmService) Start(config Config) error {
	if xormRawConfig, ok := config["xorm"]; ok {
		xormConfig, ok := xormRawConfig.(map[interface{}]interface{})
		if !ok {
			return errors.New("xorm configuration is invalid")
		}
		service.databases = make(map[string]*xorm.Engine)
		for xormConfigRawKey, xormConfigRawValue := range xormConfig {
			xormConfigKey, ok := xormConfigRawKey.(string)
			if !ok {
				return errors.New(fmt.Sprintf("the key %s must be string", xormConfigRawKey))
			}
			switch xormConfigValue := xormConfigRawValue.(type) {
			case string:
				matcher := driverStringRegex.FindStringSubmatch(xormConfigValue)
				if matcher != nil {
					if db, err := xorm.NewEngine(matcher[1], matcher[2]); err == nil {
						service.databases[xormConfigKey] = db
					} else {
						return err
					}
				} else {
					return errors.New(fmt.Sprintf("1 database %s does not has a valid URL", xormConfigKey))
				}
			case map[interface{}]interface{}:
				var url string
				if rawUrl, ok := xormConfigValue["url"]; !ok {
					return errors.New(fmt.Sprintf("%s.url is not present", xormConfigKey))
				} else if url, ok = rawUrl.(string); !ok {
					return errors.New(fmt.Sprintf("%s.url must be string", xormConfigKey))
				}
				matcher := driverStringRegex.FindStringSubmatch(url)
				var engine *xorm.Engine
				if len(matcher) > 0 {
					if db, err := xorm.NewEngine(matcher[1], matcher[2]); err == nil {
						service.databases[xormConfigKey] = db
						engine = db
					} else {
						return err
					}
				}
				for k, v := range xormConfigValue {
					switch k {
					case "url":
						break // Just ignore it
					case "show_sql":
						showSQL, ok := v.(bool)
						if !ok {
							return errors.New(fmt.Sprintf("%s.%s must be boolean", xormConfigKey, k))
						}
						engine.ShowSQL(showSQL)
						break // Just ignore it
					case "log_level":
						logLevel, ok := v.(string)
						if !ok {
							return errors.New(fmt.Sprintf("%s.%s must be string", xormConfigKey, k))
						}
						switch logLevel {
						case "debug":
							engine.SetLogLevel(core.LOG_DEBUG)
						case "info":
							engine.SetLogLevel(core.LOG_INFO)
						case "warning":
							engine.SetLogLevel(core.LOG_WARNING)
						case "err":
							engine.SetLogLevel(core.LOG_ERR)
						case "off":
							engine.SetLogLevel(core.LOG_OFF)
						case "unknown":
							engine.SetLogLevel(core.LOG_UNKNOWN)
						default:
							return errors.New(fmt.Sprintf("%s is not a valid value for %s.%s", logLevel, xormConfigKey, k))
						}
						break // Just ignore it
					default:
						return errors.New(fmt.Sprintf("%s.%s xorm configuration is not supported", xormConfigKey, k))
					}
				}
			default:
				return errors.New(fmt.Sprintf("%s xorm configuration is not supported", xormConfigKey))
			}
		}
	} else {
		return errors.New("xorm configuration was not found")
	}
	return nil
}
