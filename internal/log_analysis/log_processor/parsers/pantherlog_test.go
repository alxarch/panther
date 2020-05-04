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
	"time"

	"github.com/panther-labs/panther/internal/log_analysis/log_processor/parsers/logs"
)

type TestEvent struct {
	Ts       time.Time `json:"ts"`
	RemoteIP string    `json:"remote_ip,omitempty"`
}

var _ PantherEventer = (*TestEvent)(nil)

func (e *TestEvent) PantherEvent() *logs.Event {
	return logs.NewEvent("test", e.Ts,
		logs.IPAddress(e.RemoteIP))
}

type testParser struct{}

func (p *testParser) Parse(log string) ([]*Result, error) {
	return nil, nil
}

// func TestPantherLog(t *testing.T) {
// 	now := time.Now()
// 	strTime := (*timestamp.RFC3339)(&now).String()
// 	testEvent := TestEvent{
// 		Ts:       now,
// 		RemoteIP: "127.0.0.1",
// 	}

// 	r, _ := NewRegistry()
// 	r.Register(LogType{
// 		Name:        "test",
// 		Description: "foo",
// 		Schema: struct {
// 			TestEvent
// 			PantherLog
// 		}{},
// 		NewParser: func() Parser {
// 			return &testParser{}
// 		},
// 	})

// 	actual, err := r.RepackJSON(&testEvent)
// 	require.NoError(t, err)
// 	exp := &PantherLogJSON{
// 		JSON: []byte(`{
// 			"ts": "` + strTime + `",
// 			"remote_ip": "127.0.0.1",
// 			"p_any_ip_addresses": ["127.0.0.1"],
// 			"p_event_time": "` + strTime + `",
// 			"p_parse_time": "` + strTime + `",
// 			"p_log_type": "test"
// 		}`),
// 		LogType:   "test",
// 		EventTime: now,
// 	}
// 	// require.Equal(t, exp, actual)
// 	require.JSONEq(t, string(exp.JSON), string(actual.JSON))
// }

// func TestAnyStringUnmarshal(t *testing.T) {
// 	var any PantherAnyString

// 	// nil case
// 	jsonString := `[]`
// 	expectedAny := PantherAnyString{
// 		set: make(map[string]struct{}),
// 	}
// 	err := jsoniter.Unmarshal(([]byte)(jsonString), &any)
// 	require.NoError(t, err)
// 	require.Equal(t, expectedAny, any)

// 	// non-nil case
// 	jsonString = `["a"]`
// 	expectedAny = PantherAnyString{
// 		set: map[string]struct{}{
// 			"a": {},
// 		},
// 	}
// 	err = jsoniter.Unmarshal(([]byte)(jsonString), &any)
// 	require.NoError(t, err)
// 	require.Equal(t, expectedAny, any)
// }

// func TestAppendAnyString(t *testing.T) {
// 	value := "a"
// 	expectedAny := &PantherAnyString{
// 		set: map[string]struct{}{
// 			value: {},
// 		},
// 	}
// 	any := NewPantherAnyString()
// 	AppendAnyString(any, value)
// 	require.Equal(t, expectedAny, any)
// }

// func TestAppendAnyStringWithEmptyString(t *testing.T) {
// 	value := ""                                                  // should not be stored
// 	expectedAny := &PantherAnyString{set: map[string]struct{}{}} // empty map
// 	any := NewPantherAnyString()
// 	AppendAnyString(any, value)
// 	require.Equal(t, expectedAny, any)
// }

// func TestSetCoreFields(t *testing.T) {
// 	event := PantherLog{}
// 	logType := "Data.Source"
// 	eventTime := (timestamp.RFC3339)(time.Date(2020, 1, 2, 3, 0, 0, 0, time.UTC))
// 	expectedNow := timestamp.Now()
// 	expectedEvent := PantherLog{
// 		PantherLogType:   &logType,
// 		PantherEventTime: &eventTime,
// 		PantherParseTime: &expectedNow,
// 	}
// 	event.SetCoreFields(logType, &eventTime, nil)
// 	expectedEvent.PantherRowID = event.PantherRowID // set because it is random

// 	// PantherParseTime will be set to time.Now().UTC(), require it to be within one second of expectedNow
// 	delta := (*time.Time)(event.PantherParseTime).Sub(*(*time.Time)(expectedEvent.PantherParseTime)).Nanoseconds()
// 	require.Less(t, delta, 1*time.Second.Nanoseconds())
// 	require.Greater(t, delta, -1*time.Second.Nanoseconds())
// 	expectedEvent.PantherParseTime = event.PantherParseTime

// 	require.Equal(t, expectedEvent, event)
// }

// func TestSetCoreFieldsNilEventTime(t *testing.T) {
// 	event := PantherLog{}
// 	logType := "Data.Source"
// 	expectedNow := timestamp.Now()
// 	expectedEvent := PantherLog{
// 		PantherLogType:   &logType,
// 		PantherEventTime: &expectedNow,
// 		PantherParseTime: &expectedNow,
// 	}
// 	event.SetCoreFields(logType, nil, nil)
// 	expectedEvent.PantherRowID = event.PantherRowID // set because it is random

// 	// PantherEventTime will be set to time.Now().UTC(), require it to be within one second of expectedNow
// 	delta := (*time.Time)(event.PantherEventTime).Sub(*(*time.Time)(expectedEvent.PantherEventTime)).Nanoseconds()
// 	require.Less(t, delta, 1*time.Second.Nanoseconds())
// 	require.Greater(t, delta, -1*time.Second.Nanoseconds())
// 	// Require Panther set the EventTime to the ParseTime
// 	require.Equal(t, expectedEvent.PantherEventTime, expectedEvent.PantherParseTime)
// 	expectedEvent.PantherEventTime = event.PantherEventTime
// 	expectedEvent.PantherParseTime = event.PantherParseTime

// 	require.Equal(t, expectedEvent, event)
// }

// func TestAppendAnyIPsInField(t *testing.T) {
// 	event := PantherLog{}
// 	require.True(t, event.AppendAnyIPAddressInFieldPtr(aws.String("connection established from 192.168.1.1")))
// 	require.True(t, event.AppendAnyIPAddressInField("Accepted publickey for ubuntu from 192.168.1.2 port 54717 ssh2"))
// 	require.False(t, event.AppendAnyIPAddressInField("connection established"))

// 	expectedAny := &PantherAnyString{
// 		set: map[string]struct{}{
// 			"192.168.1.1": {},
// 			"192.168.1.2": {},
// 		},
// 	}
// 	require.Equal(t, expectedAny, event.PantherAnyIPAddresses)
// }

// func TestAppendAnyIPsInFieldMultiple(t *testing.T) {
// 	event := PantherLog{}
// 	require.True(t, event.AppendAnyIPAddressInFieldPtr(aws.String("connection established from 206.206.199.127 to 186.28.188.20")))
// 	require.True(t, event.AppendAnyIPAddressInField("Accepted publickey from 221.19.216.201 to 229.12.27.176 port 54717"))
// 	require.False(t, event.AppendAnyIPAddressInField("connection established"))

// 	expectedAny := &PantherAnyString{
// 		set: map[string]struct{}{
// 			"206.206.199.127": {},
// 			"186.28.188.20":   {},
// 			"221.19.216.201":  {},
// 			"229.12.27.176":   {},
// 		},
// 	}
// 	require.Equal(t, expectedAny, event.PantherAnyIPAddresses)
// }

// func TestAppendAnyIPV4(t *testing.T) {
// 	event := PantherLog{}
// 	require.True(t, event.AppendAnyIPAddressPtr(aws.String("192.168.1.1")))
// 	require.True(t, event.AppendAnyIPAddress("192.168.1.2"))
// 	require.False(t, event.AppendAnyIPAddress("not-an-ip"))

// 	expectedAny := &PantherAnyString{
// 		set: map[string]struct{}{
// 			"192.168.1.1": {},
// 			"192.168.1.2": {},
// 		},
// 	}
// 	require.Equal(t, expectedAny, event.PantherAnyIPAddresses)
// }

// func TestAppendAnyIPV6(t *testing.T) {
// 	event := PantherLog{}
// 	require.True(t, event.AppendAnyIPAddressPtr(aws.String("2001:db8:85a3:0:0:8a2e:370:7334")))
// 	require.True(t, event.AppendAnyIPAddress("::ffff:192.0.2.128"))
// 	require.False(t, event.AppendAnyIPAddress("not-an-ip"))

// 	expectedAny := &PantherAnyString{
// 		set: map[string]struct{}{
// 			"2001:db8:85a3:0:0:8a2e:370:7334": {},
// 			"::ffff:192.0.2.128":              {},
// 		},
// 	}
// 	require.Equal(t, expectedAny, event.PantherAnyIPAddresses)
// }

// func TestAppendAnyDomainNames(t *testing.T) {
// 	event := PantherLog{}
// 	value := "a"
// 	expectedAny := &PantherAnyString{
// 		set: map[string]struct{}{
// 			value: {},
// 		},
// 	}
// 	event.AppendAnyDomainNames(value)
// 	require.Equal(t, expectedAny, event.PantherAnyDomainNames)

// 	event = PantherLog{}
// 	event.AppendAnyDomainNamePtrs(&value)
// 	require.Equal(t, expectedAny, event.PantherAnyDomainNames)
// }

// func TestAppendAnySHA1Hashes(t *testing.T) {
// 	event := PantherLog{}
// 	value := "a"
// 	expectedAny := &PantherAnyString{
// 		set: map[string]struct{}{
// 			value: {},
// 		},
// 	}
// 	event.AppendAnySHA1Hashes(value)
// 	require.Equal(t, expectedAny, event.PantherAnySHA1Hashes)

// 	event = PantherLog{}
// 	event.AppendAnySHA1HashPtrs(&value)
// 	require.Equal(t, expectedAny, event.PantherAnySHA1Hashes)
// }

// func TestAppendAnyMD5Hashes(t *testing.T) {
// 	event := PantherLog{}
// 	value := "a"
// 	expectedAny := &PantherAnyString{
// 		set: map[string]struct{}{
// 			value: {},
// 		},
// 	}
// 	event.AppendAnyMD5Hashes(value)
// 	require.Equal(t, expectedAny, event.PantherAnyMD5Hashes)

// 	event = PantherLog{}
// 	event.AppendAnyMD5HashPtrs(&value)
// 	require.Equal(t, expectedAny, event.PantherAnyMD5Hashes)
// }
