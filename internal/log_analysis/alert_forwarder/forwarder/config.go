package forwarder

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
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/kelseyhightower/envconfig"

	policiesclient "github.com/panther-labs/panther/api/gateway/analysis/client"
	"github.com/panther-labs/panther/pkg/gatewayapi"
)

var (
	env        envConfig
	awsSession *session.Session
	ddbClient  dynamodbiface.DynamoDBAPI
	sqsClient  sqsiface.SQSAPI

	httpClient   *http.Client
	policyConfig *policiesclient.TransportConfig
	policyClient *policiesclient.PantherAnalysis
)

type envConfig struct {
	AlertsTable      string `required:"true" split_words:"true"`
	AlertingQueueURL string `required:"true" split_words:"true"`
	AnalysisAPIHost  string `required:"true" split_words:"true"`
	AnalysisAPIPath  string `required:"true" split_words:"true"`
}

// Setup parses the environment and builds the AWS and http clients.
func Setup() {
	envconfig.MustProcess("", &env)

	awsSession = session.Must(session.NewSession())
	ddbClient = dynamodb.New(awsSession)
	sqsClient = sqs.New(awsSession)
	httpClient = gatewayapi.GatewayClient(awsSession)
	policyConfig = policiesclient.DefaultTransportConfig().
		WithHost(env.AnalysisAPIHost).
		WithBasePath(env.AnalysisAPIPath)
	policyClient = policiesclient.NewHTTPClientWithConfig(nil, policyConfig)
}
