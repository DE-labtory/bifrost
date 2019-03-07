/*
 * Copyright 2018 DE-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package bifrost_test

import (
	"testing"

	"github.com/DE-labtory/bifrost"
	"github.com/stretchr/testify/assert"
)

func TestToAddress(t *testing.T) {
	// given
	tests := map[string]struct {
		input  string
		output bifrost.Address
		err    error
	}{
		"correct IP1": {
			input:  "127.0.0.1:1321",
			output: bifrost.Address{IP: "127.0.0.1:1321"},
			err:    nil,
		},
		"correct IP2": {
			input:  "127.0.0.1:41234",
			output: bifrost.Address{IP: "127.0.0.1:41234"},
			err:    nil,
		},
		"wrong IP1": {
			input:  "1271.0.0.1:1234",
			output: bifrost.Address{},
			err:    bifrost.ErrIPv4Format,
		},
		"wrong IP2": {
			input:  "127.0.0.1:123455",
			output: bifrost.Address{},
			err:    bifrost.ErrIPv4Format,
		},
		"wrong IP3": {
			input:  "127.1344.0.1:1235",
			output: bifrost.Address{},
			err:    bifrost.ErrIPv4Format,
		},
		"wrong IP4": {
			input:  "127.0.325125.1:1233",
			output: bifrost.Address{},
			err:    bifrost.ErrIPv4Format,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case [%s]", testName)

		// when
		address, err := bifrost.ToAddress(test.input)

		// then
		assert.Equal(t, address, test.output)
		assert.Equal(t, err, test.err)
	}

}
