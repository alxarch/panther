package handlers

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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLowerSet(t *testing.T) {
	result := lowerSet([]string{"AWS", "aws", "CIS", "cis", "Panther"})
	sort.Strings(result)
	assert.Equal(t, []string{"aws", "cis", "panther"}, result)
}

func TestIntMin(t *testing.T) {
	assert.Equal(t, -2, intMin(-2, 0))
	assert.Equal(t, 0, intMin(0, 0))
	assert.Equal(t, 5, intMin(10, 5))
}

func TestSetDifference(t *testing.T) {
	assert.Empty(t, setDifference([]string{}, []string{}))
	assert.Empty(t, setDifference([]string{"a", "b", "c"}, []string{"c", "a", "b"}))
	assert.Equal(t, []string{"a", "b"}, setDifference([]string{"a", "b"}, nil))
	assert.Empty(t, setDifference(nil, []string{"a", "b"}))
	assert.Equal(t, []string{"panther", "labs"},
		setDifference([]string{"panther", "labs", "inc"}, []string{"inc", "runpanther.io"}))
}

func TestSetEquality(t *testing.T) {
	assert.True(t, setEquality(nil, []string{}))
	assert.True(t, setEquality([]string{"panther", "labs", "inc"}, []string{"inc", "labs", "panther"}))
	assert.False(t, setEquality([]string{"panther"}, []string{"panther", "labs"}))
	assert.False(t, setEquality([]string{"panther", "labs"}, []string{"panther", "inc"}))
}

func TestSortCaseInsensitive(t *testing.T) {
	input := []string{"AWS.EC2.VPC", "AWS.EC2.Volume"}
	sortCaseInsensitive(input)
	assert.Equal(t, []string{"AWS.EC2.Volume", "AWS.EC2.VPC"}, input)

	// Sort by case if lowercase versions are equal
	input = []string{"panther", "Panther", "Panna Cotta", "panna cotta"}
	sortCaseInsensitive(input)
	expected := []string{"Panna Cotta", "panna cotta", "Panther", "panther"}
	assert.Equal(t, expected, input)
}
