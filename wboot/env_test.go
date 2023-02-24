package wboot

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestEnvStr(t *testing.T) {
	os.Setenv("KEY1", "VAL1  ")
	require.Equal(t, "VAL1", EnvStr("KEY1"))
}

func TestEnvStrOr(t *testing.T) {
	os.Setenv("KEY1", "VAL1  ")
	require.Equal(t, "VAL1", EnvStrOr("KEY1", "B"))
	require.Equal(t, "B", EnvStrOr("KEY2", "B"))
}

func TestEnvBool(t *testing.T) {
	os.Setenv("KEY1", "false")
	os.Setenv("KEY2", "true")
	require.True(t, EnvBool("KEY2"))
	require.False(t, EnvBool("KEY1"))
}

func TestEnvBoolOr(t *testing.T) {
	os.Setenv("KEY1", "false")
	os.Setenv("KEY2", "rue")
	require.False(t, EnvBoolOr("KEY1", false))
	require.False(t, EnvBoolOr("KEY2", false))

	os.Setenv("KEY1", "true")
	os.Setenv("KEY2", "rue")
	require.True(t, EnvBoolOr("KEY1", true))
	require.True(t, EnvBoolOr("KEY2", true))
}
