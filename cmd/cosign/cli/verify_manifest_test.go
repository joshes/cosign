// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"reflect"
	"testing"
)

const SingleContainerManifest = `
apiVersion: v1
kind: Pod
metadata:
  name: single-pod
spec:
  restartPolicy: Never
  containers:
    - name: nginx-container
      image: nginx:1.21.1
`

const MultiContainerManifest = `
apiVersion: v1
kind: Pod
metadata:
  name: multi-pod
spec:
  restartPolicy: Never
  volumes:
    - name: shared-data
      emptyDir: {}
  containers:
    - name: nginx-container
      image: nginx:1.21.1
      volumeMounts:
        - name: shared-data
          mountPath: /usr/share/nginx/html
    - name: ubuntu-container
      image: ubuntu:21.10
      volumeMounts:
        - name: shared-data
          mountPath: /pod-data
      command: ["/bin/sh"]
      args: ["-c", "echo Hello, World > /pod-data/index.html"]
`

func TestGetImagesFromYamlManifest(t *testing.T) {
	testCases := []struct {
		name         string
		fileContents string
		expected     []string
	}{
		{
			name:         "single image",
			fileContents: SingleContainerManifest,
			expected:     []string{"nginx:1.21.1"},
		},
		{
			name:         "multi image",
			fileContents: MultiContainerManifest,
			expected:     []string{"nginx:1.21.1", "ubuntu:21.10"},
		},
		{
			name:         "no images found",
			fileContents: ``,
			expected:     nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := getImagesFromManifest(tc.fileContents)
			if err != nil {
				t.Fatalf("getImagesFromManifest returned error: %v", err)
			}
			if !reflect.DeepEqual(tc.expected, got) {
				t.Errorf("getImagesFromManifest returned %v, wanted %v", got, tc.expected)
			}
		})
	}
}