package models

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
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

// NOTE: different kinds of databases (e.g., Athena, Snowflake) will use different endpoints (lambda functions), same api.

// NOTE: if a json tag is used more than once it is factored into a struct to avoid inconsistencies

const (
	QuerySucceeded = "succeeded"
	QueryFailed    = "failed"
	QueryRunning   = "running"
	QueryCanceled  = "canceled"
)

// LambdaInput is the collection of all possible args to the Lambda function.
type LambdaInput struct {
	ExecuteAsyncQuery       *ExecuteAsyncQueryInput       `json:"executeAsyncQuery"`
	ExecuteAsyncQueryNotify *ExecuteAsyncQueryNotifyInput `json:"executeAsyncQueryNotify"` // uses Step functions
	ExecuteQuery            *ExecuteQueryInput            `json:"executeQuery"`
	GetDatabases            *GetDatabasesInput            `json:"getDatabases"`
	GetQueryResults         *GetQueryResultsInput         `json:"getQueryResults"`
	GetQueryStatus          *GetQueryStatusInput          `json:"getQueryStatus"`
	GetTables               *GetTablesInput               `json:"getTables"`
	GetTablesDetail         *GetTablesDetailInput         `json:"getTablesDetail"`
	InvokeNotifyLambda      *InvokeNotifyLambdaInput      `json:"invokeNotifyLambda"`
	NotifyAppSync           *NotifyAppSyncInput           `json:"notifyAppSync"`
	StopQuery               *StopQueryInput               `json:"stopQuery"`
}

type GetDatabasesInput struct {
	OptionalName // if nil get all databases
}

// NOTE: we will assume this is small an not paginate
type GetDatabasesOutput struct {
	Databases []*NameAndDescription `json:"databases,omitempty"`
}

type GetTablesInput struct {
	Database
	OnlyPopulated bool `json:"onlyPopulated,omitempty"` // if true, only return table containing data
}

// NOTE: we will assume this is small an not paginate
type GetTablesOutput struct {
	TablesDetail
}

type TablesDetail struct {
	Tables []*TableDetail `json:"tables"`
}

type TableDetail struct {
	TableDescription
	Columns []*TableColumn `json:"columns"`
}

type TableDescription struct {
	Database
	NameAndDescription
}

type GetTablesDetailInput struct {
	Database
	Names []string `json:"names" validate:"required"`
}

// NOTE: we will assume this is small an not paginate
type GetTablesDetailOutput struct {
	TablesDetail
}

type TableColumn struct {
	NameAndDescription
	Type string `json:"type" validate:"required"`
}

type ExecuteAsyncQueryNotifyInput struct {
	ExecuteAsyncQueryInput
	LambdaInvoke
	UserDataToken
	DelaySeconds int `json:"delaySeconds" validate:"omitempty,gt=0"` // wait this long before starting workflow (default 0)
}

type ExecuteAsyncQueryNotifyOutput struct {
	WorkflowIdentifier
}

// Blocking query
type ExecuteQueryInput = ExecuteAsyncQueryInput

type ExecuteQueryOutput = GetQueryResultsOutput // call GetQueryResults() to page thu results

type ExecuteAsyncQueryInput struct {
	Database
	SQLQuery
}

type ExecuteAsyncQueryOutput struct {
	QueryStatus
	QueryIdentifier
}

type GetQueryStatusInput = QueryIdentifier

type GetQueryStatusOutput struct {
	QueryStatus
	SQLQuery
	Stats *QueryResultsStats `json:"stats,omitempty"` // present only on successful queries
}

type GetQueryResultsInput struct {
	QueryIdentifier
	Pagination
	PageSize *int64 `json:"pageSize" validate:"omitempty,gt=0,lt=1000"` // only return this many rows per call
}

type GetQueryResultsOutput struct {
	GetQueryStatusOutput
	ResultsPage QueryResultsPage `json:"resultsPage" validate:"required"`
}

type QueryResultsPage struct {
	Pagination
	NumRows int    `json:"numRows"  validate:"required"` // number of rows in page of results, len(Rows)
	Rows    []*Row `json:"rows"  validate:"required"`
}

type QueryResultsStats struct {
	ExecutionTimeMilliseconds int64 `json:"executionTimeMilliseconds"  validate:"required"`
	DataScannedBytes          int64 `json:"dataScannedBytes"  validate:"required"`
}

type StopQueryInput = QueryIdentifier

type StopQueryOutput = GetQueryStatusOutput

type InvokeNotifyLambdaInput struct {
	LambdaInvoke
	QueryIdentifier
	WorkflowIdentifier
	UserDataToken
}

type InvokeNotifyLambdaOutput struct {
}

type NotifyAppSyncInput struct {
	NotifyInput
}

type NotifyAppSyncOutput struct {
	StatusCode int `json:"statusCode" validate:"required"` // the http status returned from POSTing callback to appsync
}

type NotifyInput struct { // notify lambdas need to have this as input
	GetQueryStatusInput
	ExecuteAsyncQueryNotifyOutput
	UserDataToken
}

type NameAndDescription struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type OptionalName struct {
	Name *string `json:"name,omitempty"`
}

type SQLQuery struct {
	SQL string `json:"sql" validate:"required"`
}

type QueryIdentifier struct {
	QueryID string `json:"queryId" validate:"required"`
}

type Database struct {
	DatabaseName string `json:"databaseName" validate:"required"`
}

type Row struct {
	Columns []*Column `json:"columns"`
}

type Column struct {
	Value string `json:"value"`
}

type Pagination struct {
	PaginationToken *string `json:"paginationToken,omitempty"`
}

type QueryStatus struct {
	Status   string `json:"status" validate:"required,oneof=running,succeeded,failed,canceled"`
	SQLError string `json:"sqlError,omitempty"`
}

type WorkflowIdentifier struct {
	WorkflowID string `json:"workflowId" validate:"required"`
}

type UserDataToken struct {
	UserData string `json:"userData" validate:"required,gt=0"` // token passed though to notifications (usually the userid)
}

type LambdaInvoke struct {
	LambdaName string `json:"lambdaName" validate:"required"` // the name of the lambda to call when done
	MethodName string `json:"methodName" validate:"required"` // the method to call on the lambda
}
