package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	os.Setenv("PROCTOR_ENVIRONMENT", "development")

	viper.AutomaticEnv()

	assert.Equal(t, "development", Environment())
}

func TestLogLevel(t *testing.T) {
	os.Setenv("PROCTOR_LOG_LEVEL", "debug")

	viper.AutomaticEnv()

	assert.Equal(t, "debug", LogLevel())
}

func TestAppPort(t *testing.T) {
	os.Setenv("PROCTOR_APP_PORT", "3000")

	viper.AutomaticEnv()

	assert.Equal(t, "3000", AppPort())
}

func TestDefaultNamespace(t *testing.T) {
	os.Setenv("PROCTOR_DEFAULT_NAMESPACE", "default")

	viper.AutomaticEnv()

	assert.Equal(t, "default", DefaultNamespace())
}

func TestRedisAddress(t *testing.T) {
	os.Setenv("PROCTOR_REDIS_ADDRESS", "localhost:6379")

	viper.AutomaticEnv()

	assert.Equal(t, "localhost:6379", RedisAddress())
}

func TestKubeClusterHostName(t *testing.T) {
	os.Setenv("PROCTOR_KUBE_CLUSTER_HOST_NAME", "somekube.io")

	viper.AutomaticEnv()

	assert.Equal(t, "somekube.io", KubeClusterHostName())
}

func TestRedisMaxActiveConnections(t *testing.T) {
	os.Setenv("PROCTOR_REDIS_MAX_ACTIVE_CONNECTIONS", "50")

	viper.AutomaticEnv()

	assert.Equal(t, 50, RedisMaxActiveConnections())
}

func TestLogsStreamReadBufferSize(t *testing.T) {
	os.Setenv("PROCTOR_LOGS_STREAM_READ_BUFFER_SIZE", "140")

	viper.AutomaticEnv()

	assert.Equal(t, 140, LogsStreamReadBufferSize())
}

func TestLogsStreamWriteBufferSize(t *testing.T) {
	os.Setenv("PROCTOR_LOGS_STREAM_WRITE_BUFFER_SIZE", "4096")

	viper.AutomaticEnv()

	assert.Equal(t, 4096, LogsStreamWriteBufferSize())
}

func TestKubeJobActiveDeadlineSeconds(t *testing.T) {
	os.Setenv("PROCTOR_KUBE_JOB_ACTIVE_DEADLINE_SECONDS", "900")

	viper.AutomaticEnv()

	expectedValue := int64(900)
	assert.Equal(t, &expectedValue, KubeJobActiveDeadlineSeconds())
}
