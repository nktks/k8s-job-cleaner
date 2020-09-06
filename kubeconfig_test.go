package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetMasterURLFromKubeConfig(t *testing.T) {
	t.Run("could parse kubeconfig", func(t *testing.T) {
		kc := `apiVersion: v1
clusters:
- cluster:
    server: http://10.0.0.1
  name: a
- cluster:
    server: http://10.0.0.2
  name: b
`
		t.Run("could find name from kubeconfig", func(t *testing.T) {
			name := "a"
			ret, err := GetMasterURLFromKubeConfig([]byte(kc), name)
			require.NoError(t, err)
			require.Equal(t, "http://10.0.0.1", ret)

		})
		t.Run("could not find name from kubeconfig", func(t *testing.T) {
			name := "c"
			_, err := GetMasterURLFromKubeConfig([]byte(kc), name)
			require.Error(t, err)
		})
	})
	t.Run("could not parse kubeconfig yaml", func(t *testing.T) {
		kc := `hoge
`
		name := "a"
		_, err := GetMasterURLFromKubeConfig([]byte(kc), name)
		require.Error(t, err)
	})
}
