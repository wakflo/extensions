// Copyright 2022-present Wakflo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

func NewDefaultLogger(service string) zerolog.Logger {
	return NewStdErr("debug", service)
}

func NewStdErr(level string, service string) zerolog.Logger {
	var err error

	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	var out io.Writer
	out = zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02T15:04:05.999Z07:00",
	}

	l := zerolog.New(out).Level(lvl)
	l = l.With().Timestamp().Logger()
	if service != "" {
		l = l.With().Str("service", service).Logger()
	}

	return l
}
