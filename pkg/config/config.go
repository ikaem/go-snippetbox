// pkg/config/config.go
package config

import "log"

type Application struct {
	InfoLog  log.Logger
	ErrorLog log.Logger
}
