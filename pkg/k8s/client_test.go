package k8s

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientset(t *testing.T) {
	t.Run("creates clientset with available config", func(t *testing.T) {
		clientset, err := NewClientset()
		if err != nil {
			t.Skipf("unable to create clientset (expected when no kubeconfig available): %v", err)
			return
		}
		assert.NotNil(t, clientset)
	})
}

func TestGetCurrentNamespace(t *testing.T) {
	t.Run("returns POD_NAMESPACE if set", func(t *testing.T) {
		os.Setenv("POD_NAMESPACE", "test-namespace")
		defer os.Unsetenv("POD_NAMESPACE")

		ns := GetCurrentNamespace()
		assert.Equal(t, "test-namespace", ns)
	})

	t.Run("returns default when no namespace available", func(t *testing.T) {
		os.Unsetenv("POD_NAMESPACE")

		ns := GetCurrentNamespace()
		assert.Equal(t, "default", ns)
	})
}

func TestGetPodName(t *testing.T) {
	t.Run("returns HOSTNAME env if set", func(t *testing.T) {
		originalHostname := os.Getenv("HOSTNAME")
		os.Setenv("HOSTNAME", "test-pod-12345")
		defer func() {
			if originalHostname != "" {
				os.Setenv("HOSTNAME", originalHostname)
			} else {
				os.Unsetenv("HOSTNAME")
			}
		}()

		podName := GetPodName()
		assert.Equal(t, "test-pod-12345", podName)
	})

	t.Run("returns os.Hostname when HOSTNAME not set", func(t *testing.T) {
		originalHostname := os.Getenv("HOSTNAME")
		os.Unsetenv("HOSTNAME")
		defer func() {
			if originalHostname != "" {
				os.Setenv("HOSTNAME", originalHostname)
			}
		}()

		podName := GetPodName()
		expectedHostname, _ := os.Hostname()
		assert.Equal(t, expectedHostname, podName)
	})
}
