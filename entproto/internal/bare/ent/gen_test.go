// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ent

import (
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// Tests the basic example available at https://entgo.io/docs/grpc-generating-proto.
func TestDefaultEntproto(t *testing.T) {
	err := os.RemoveAll("proto")
	require.NoError(t, err)
	err = gen()
	require.NoError(t, err)
	_, err = os.Stat("proto")
	require.NoError(t, err)
	_, err = os.Stat("proto/entpb/generate.go")
	require.NoError(t, err)
}

func gen() error {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("go", "generate", "./...")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
