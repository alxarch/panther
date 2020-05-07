package jsonutil

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
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
)

func TestAthenaRewrite(t *testing.T) {
	field := "@name"
	mapped := RewriteFieldNameAthena(field)
	require.Equal(t, "_at_sign_name", mapped)
}

func TestJSONIterExtension(t *testing.T) {
	RegisterAthenaRewrite()

	type S struct {
		Type string `json:"@type"`
	}
	var value S
	err := jsoniter.UnmarshalFromString(`{"@type":"foo"}`, &value)
	require.NoError(t, err)
	data, err := jsoniter.MarshalToString(&value)
	require.NoError(t, err)
	require.Equal(t, `{"_at_sign_type":"foo"}`, data)
}
