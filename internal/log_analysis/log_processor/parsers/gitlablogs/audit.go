package gitlablogs

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
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/logs"
	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/timestamp"
)

const (
	// TypeAudit is the log type of Audit log records
	TypeAudit = PantherPrefix + ".Audit"
	// AuditDesc describes the Git log record
)

var LogTypeAudit = parsers.LogType{
	Name: TypeAudit,
	Description: `GitLab log file containing changes to group or project settings 
Reference: https://docs.gitlab.com/ee/administration/logs.html#audit_jsonlog`,
	Schema: struct {
		Audit
		logs.Meta
	}{},
	NewParser: NewAuditParser,
}

// Audit is a a GitLab log line from a failed interaction with git
// nolint: lll
type Audit struct {
	Severity      *string            `json:"severity" validate:"required" description:"The log level"`
	Time          *timestamp.RFC3339 `json:"time" validate:"required" description:"The event timestamp"`
	AuthorID      *int64             `json:"author_id" validate:"required" description:"User id that made the change"`
	EntityID      *int64             `json:"entity_id" validate:"required" description:"Id of the entity that was modified"`
	EntityType    *string            `json:"entity_type" validate:"required" description:"Type of the modified entity"`
	Change        *string            `json:"change" validate:"required" description:"Type of change to the settings"`
	From          *string            `json:"from" validate:"required" description:"Old setting value"`
	To            *string            `json:"to" validate:"required" description:"New setting value"`
	AuthorName    *string            `json:"author_name" validate:"required" description:"Name of the user that made the change"`
	TargetID      *int64             `json:"target_id" validate:"required" description:"Target id of the modified setting"`
	TargetType    *string            `json:"target_type" validate:"required" description:"Target type of the modified setting"`
	TargetDetails *string            `json:"target_details" validate:"required" description:"Details of the target of the modified setting"`
}

func (event *Audit) PantherEvent() *logs.Event {
	return logs.NewEvent(TypeAudit, event.Time.UTC())
}

// AuditParser parses gitlab rails logs
type AuditParser struct{}

var _ parsers.PantherEventer = (*Audit)(nil)

var _ parsers.Interface = (*AuditParser)(nil)

// New creates a new parser
func NewAuditParser() parsers.Interface {
	return &AuditParser{}
}

// Parse returns the parsed events or nil if parsing failed
func (p *AuditParser) Parse(log string) ([]*parsers.Result, error) {
	return parsers.QuickParseJSON(&Audit{}, log)
}
