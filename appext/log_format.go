// Copyright 2025-2026 Buf Technologies, Inc.
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
	"fmt"
	"strings"
)

// LogFormat is a format to print logs in.
type LogFormat string

const (
	// LogFormatText is the text log format.
	LogFormatText LogFormat = "text"
	// LogFormatColor is the colored text log format.
	//
	// This is the default value when parsing LogFormats. However, unless BuilderWithLoggerProvider
	// is used, there is no difference between LogFormatText and LogFormatColor.
	LogFormatColor LogFormat = "color"
	// LogFormatJSON is the JSON log format.
	LogFormatJSON LogFormat = "json"
)

// AllLogFormatStrings contains all valid values for the --log-format flag.
var AllLogFormatStrings = []string{
	string(LogFormatText),
	string(LogFormatColor),
	string(LogFormatJSON),
}

// ParseLogFormat parses the log format for the string.
//
// If logFormatString is empty, this returns LogFormatColor.
func ParseLogFormat(logFormatString string) (LogFormat, error) {
	switch LogFormat(strings.TrimSpace(strings.ToLower(logFormatString))) {
	case LogFormatText:
		return LogFormatText, nil
	case LogFormatColor, "":
		return LogFormatColor, nil
	case LogFormatJSON:
		return LogFormatJSON, nil
	default:
		return "", fmt.Errorf("unknown log format [text,color,json]: %q", logFormatString)
	}
}
