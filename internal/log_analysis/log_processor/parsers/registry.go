package parsers

/**
 * Panther is a Cloud-Native SIEM for the Modern Security Team.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/panther-labs/panther/pkg/awsglue"
)

type LogType struct {
	Name        string
	Description string
	Schema      interface{}
	NewParser   ParserFactory
}

// // PantherLogFactory creates a serializable struct from a PantherEvent
// // This is not optimal in terms of performance
// type PantherLogFactory func(logType string, tm time.Time, fields ...pantherlog.Field) interface{}

func (entry *LogType) GlueTableMetadata() *awsglue.GlueTableMetadata {
	return awsglue.LogDataHourlyTableMetadata(entry.Name, entry.Description, entry.Schema)
}

func (entry *LogType) Check() error {
	if entry == nil {
		return errors.Errorf("nil log type entry")
	}
	if entry.Name == "" {
		return errors.Errorf("missing entry log type")
	}
	if entry.Description == "" {
		return errors.Errorf("missing description for log type %q", entry.Name)
	}
	// describes Glue table over processed data in S3
	// assert it does not panic here until some validation method is provided
	// TODO: [awsglue] Add some validation for the metadata in `awsglue` package
	_ = awsglue.LogDataHourlyTableMetadata(entry.Name, entry.Description, entry.Schema)

	return checkLogEntrySchema(entry.Name, entry.Schema)
}

type Registry struct {
	mu      sync.RWMutex
	entries map[string]*LogType
}

func NewRegistry(logTypes ...LogType) (*Registry, error) {
	r := &Registry{}
	for _, logType := range logTypes {
		if err := r.Register(logType); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// MustGet gets a registered LogType or panics
func (r *Registry) MustGet(name string) *LogType {
	if logType := r.Get(name); logType != nil {
		return logType
	}
	panic(fmt.Sprintf("unregistered log type %q", name))
}

func (r *Registry) Get(name string) *LogType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.entries[name]
}

func (r *Registry) LogTypes() (logTypes []LogType) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, logType := range r.entries {
		logTypes = append(logTypes, *logType)
	}
	return
}

func (r *Registry) Register(entry LogType) error {
	if err := entry.Check(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, duplicate := r.entries[entry.Name]; duplicate {
		return errors.Errorf("duplicate log type entry %q", entry.Name)
	}
	if r.entries == nil {
		r.entries = make(map[string]*LogType)
	}
	r.entries[entry.Name] = &entry
	return nil
}

var defaultRegistry Registry

func Get(logType string) *LogType {
	return defaultRegistry.Get(logType)
}

func MustGet(logType string) *LogType {
	return defaultRegistry.MustGet(logType)
}

func Register(entries ...LogType) error {
	for _, entry := range entries {
		if err := defaultRegistry.Register(entry); err != nil {
			return err
		}
	}
	return nil
}

func MustRegister(entries ...LogType) {
	if err := Register(entries...); err != nil {
		panic(err)
	}
}

func AvailableLogTypes() []LogType {
	return defaultRegistry.LogTypes()
}

func NewParser(logType string) (Interface, error) {
	entry := defaultRegistry.Get(logType)
	if entry != nil {
		return entry.NewParser(), nil
	}
	return nil, errors.Errorf("unregistered LogType %q", logType)
}

// func defaultPantherLogFactory(logType string, tm time.Time, fields ...PantherField) interface{} {
// 	return NewPantherLog(logType, tm, fields...)
// }
