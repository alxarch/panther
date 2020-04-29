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
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/panther-labs/panther/pkg/awsglue"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
)

type LogType struct {
	Name        string
	Description string
	Schema      interface{}
	NewParser   ParserFactory
}

type ParserFactory func() Parser

func (entry *LogType) GlueTableMetaData() *awsglue.GlueTableMetadata {
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

func Register(entries ...LogType) error {
	for _, entry := range entries {
		if err := defaultRegistry.Register(entry); err != nil {
			return err
		}
	}
	return nil
}

func AvailableLogTypes() []LogType {
	return defaultRegistry.LogTypes()
}

func NewParser(logType string) (Parser, error) {
	entry := defaultRegistry.Get(logType)
	if entry != nil {
		return entry.NewParser(), nil
	}
	return nil, errors.Errorf("unregistered LogType %q", logType)
}

func MustRegister(entries ...LogType) {
	if err := Register(entries...); err != nil {
		panic(err)
	}
}

// Parser represents a parser for a supported log type
type Parser interface {
	// Parse attempts to parse the provided log line
	// If the provided log is not of the supported type the method returns nil and an error
	Parse(log string) ([]*PantherLogJSON, error)
}

// Validator can be used to validate schemas of log fields
var Validator = validator.New()

func indexOf(values []string, search string) int {
	for i, value := range values {
		if value == search {
			return i
		}
	}
	return -1
}

func checkLogEntrySchema(logType string, schema interface{}) error {
	if schema == nil {
		return errors.Errorf("nil schema for log type %q", logType)
	}
	data, err := jsoniter.Marshal(schema)
	if err != nil {
		return errors.Errorf("invalid schema struct for log type %q: %s", logType, err)
	}
	var fields map[string]interface{}
	if err := jsoniter.Unmarshal(data, &fields); err != nil {
		return errors.Errorf("invalid schema struct for log type %q: %s", logType, err)
	}
	// TODO: [parsers] Use reflect to check provided schema struct for required panther fields
	return nil
}
