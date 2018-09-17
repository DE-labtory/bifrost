/*
 * Copyright 2018 It-chain
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

package bifrost

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnectionStore(t *testing.T) {
	// when
	testConnStore := NewConnectionStore()

	// then
	assert.NotNil(t, testConnStore.RWMutex)
	assert.NotNil(t, testConnStore.connMap)
}
