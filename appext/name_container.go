// Copyright 2025 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package appext

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"buf.build/go/app"
)

type nameContainer struct {
	app.Container

	appName string

	configDirPath     string
	configDirPathOnce sync.Once
	cacheDirPath      string
	cacheDirPathOnce  sync.Once
	dataDirPath       string
	dataDirPathOnce   sync.Once
	port              uint16
	portErr           error
	portOnce          sync.Once
}

func newNameContainer(baseContainer app.Container, appName string) (*nameContainer, error) {
	if err := validateAppName(appName); err != nil {
		return nil, err
	}
	return &nameContainer{
		Container: baseContainer,
		appName:   appName,
	}, nil
}

func (c *nameContainer) AppName() string {
	return c.appName
}

func (c *nameContainer) ConfigDirPath() string {
	c.configDirPathOnce.Do(c.setConfigDirPath)
	return c.configDirPath
}

func (c *nameContainer) CacheDirPath() string {
	c.cacheDirPathOnce.Do(c.setCacheDirPath)
	return c.cacheDirPath
}

func (c *nameContainer) DataDirPath() string {
	c.dataDirPathOnce.Do(c.setDataDirPath)
	return c.dataDirPath
}

func (c *nameContainer) Port() (uint16, error) {
	c.portOnce.Do(c.setPort)
	return c.port, c.portErr
}

func (c *nameContainer) setConfigDirPath() {
	c.configDirPath = c.getDirPath("CONFIG_DIR", app.ConfigDirPath)
}

func (c *nameContainer) setCacheDirPath() {
	c.cacheDirPath = c.getDirPath("CACHE_DIR", app.CacheDirPath)
}

func (c *nameContainer) setDataDirPath() {
	c.dataDirPath = c.getDirPath("DATA_DIR", app.DataDirPath)
}

func (c *nameContainer) setPort() {
	c.port, c.portErr = c.getPort()
}

func (c *nameContainer) getDirPath(envSuffix string, getBaseDirPath func(app.EnvContainer) (string, error)) string {
	dirPath := c.Container.Env(getAppNameEnvPrefix(c.appName) + envSuffix)
	if dirPath == "" {
		baseDirPath, err := getBaseDirPath(c.Container)
		if err == nil {
			dirPath = filepath.Join(baseDirPath, c.appName)
		}
	}
	return dirPath
}

func (c *nameContainer) getPort() (uint16, error) {
	portString := c.Container.Env(getAppNameEnvPrefix(c.appName) + "PORT")
	if portString == "" {
		portString = c.Container.Env("PORT")
		if portString == "" {
			return 0, nil
		}
	}
	port, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		return 0, fmt.Errorf("could not parse port %q to uint16: %w", portString, err)
	}
	return uint16(port), nil
}

func getAppNameEnvPrefix(appName string) string {
	return strings.ToUpper(strings.ReplaceAll(appName, "-", "_")) + "_"
}

func validateAppName(appName string) error {
	if appName == "" {
		return errors.New("empty application name")
	}
	for _, c := range appName {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return fmt.Errorf("invalid application name: %s", appName)
		}
	}
	return nil
}
