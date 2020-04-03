package api

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

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/aws/aws-sdk-go/service/glue"
)

var (
	awsSession          *session.Session
	glueClient          *glue.Glue
	athenaClient        *athena.Athena
	athenaS3ResultsPath *string
)

func SessionInit() {
	awsSession = session.Must(session.NewSession())
	glueClient = glue.New(awsSession)
	athenaClient = athena.New(awsSession)

	if os.Getenv("ATHENA_BUCKET") != "" {
		results := "s3://" + os.Getenv("ATHENA_BUCKET") + "/athena_api/"
		athenaS3ResultsPath = &results
	}
}

// API provides receiver methods for each route handler.
type API struct{}
