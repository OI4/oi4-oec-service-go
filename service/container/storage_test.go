package container

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const validPath = "docker_configs/valid"

func TestNewMessageBusStorage_ValidPath(t *testing.T) {
	storage, err := NewMessageBusStorage(filepath.Join(validPath, DefaultMessageBusStorageSubFolder), getLogger())
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestNewMessageBusStorage_InvalidPath(t *testing.T) {
	storage, err := NewMessageBusStorage("invalid/path", getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid message bus storage location: stat invalid/path: no such file or directory", err.Error())
}

func TestNewMessageBusStorage_InvalidCa(t *testing.T) {
	storage, err := NewMessageBusStorage(filepath.Join("docker_configs", "mqtt_invalid_ca"), getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid CA certificates: invalid root CA certificate: invalid PEM file", err.Error())
}

func TestNewMessageBusStorage_InvalidBrokerConfiguration(t *testing.T) {
	storage, err := NewMessageBusStorage(filepath.Join("docker_configs", "mqtt_invalid_broker_configuration"), getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid broker configuration: invalid json format: unexpected end of JSON input", err.Error())
}

func TestNewOi4CertificateStorage_ValidPath(t *testing.T) {
	storage, err := NewOi4CertificateStorage(filepath.Join(validPath, DefaultOi4CertificateStorageSubFolder), "F12AB35", getLogger())
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestNewOi4CertificateStorage_InvalidPath(t *testing.T) {
	storage, err := NewOi4CertificateStorage("invalid/path", "containerName", getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
}

func TestNewOi4CertificateStorage_InvalidClientPem(t *testing.T) {
	storage, err := NewOi4CertificateStorage(filepath.Join(validPath, DefaultOi4CertificateStorageSubFolder), "NoId", getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid client certificate: failed to read pem file", err.Error())
}

func TestNewOi4CertificateStorage_InvalidCA(t *testing.T) {
	storage, err := NewOi4CertificateStorage(filepath.Join("docker_configs", "certificate_invalid_ca"), "F12AB35", getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid CA certificates: invalid root CA certificate: invalid PEM file", err.Error())
}

func TestNewSecretStorage_ValidPath(t *testing.T) {
	storage, err := NewSecretStorage(filepath.Join(validPath, DefaultSecretsFolder), getLogger())
	require.NoError(t, err)
	require.NotNil(t, storage)
	credentials := storage.MqttCredentials
	require.NotNil(t, credentials)
	assert.Equal(t, "oi4User", credentials.Username())
	pwd, set := credentials.Password()
	assert.True(t, set)
	assert.Equal(t, "oi4Password", pwd)
}

func TestNewSecretStorage_InvalidPath(t *testing.T) {
	storage, err := NewSecretStorage("invalid/path", getLogger())
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid secret storage location: stat invalid/path: no such file or directory", err.Error())
}

func TestNewSecretStorage_InvalidCredentials(t *testing.T) {
	storage, err := NewSecretStorage(filepath.Join("docker_configs", "empty"), getLogger())
	require.NoError(t, err)
	require.NotNil(t, storage)
	require.Nil(t, storage.MqttCredentials)
}

func TestNewApplicationSpecificStorages_ValidPaths(t *testing.T) {
	storage, err := NewApplicationSpecificStorages(filepath.Join(validPath, DefaultApplicationSpecificConfigurationFolder), filepath.Join(validPath, DefaultApplicationSpecificDataFolder))
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestNewApplicationSpecificStorages_InvalidConfigPath(t *testing.T) {
	storage, err := NewApplicationSpecificStorages("invalid", filepath.Join("docker_configs", "empty"))
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid configuration path: no such file or directory", err.Error())
}

func TestNewApplicationSpecificStorages_InvalidDataPath(t *testing.T) {
	storage, err := NewApplicationSpecificStorages(filepath.Join("docker_configs", "empty"), "invalid")
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid data path: no such file or directory", err.Error())
}

func TestReadCredentials_ValidFile(t *testing.T) {
	credentials, err := readCredentials(filepath.Join(validPath, DefaultSecretsFolder, "mqtt_credentials"), getLogger())
	require.NoError(t, err)
	require.NotNil(t, credentials)
	assert.Equal(t, "oi4User", credentials.Username())
	pwd, set := credentials.Password()
	assert.True(t, set)
	assert.Equal(t, "oi4Password", pwd)
}

func TestReadCredentials_InvalidFile(t *testing.T) {
	credentials, err := readCredentials("invalid", getLogger())
	require.Error(t, err)
	require.Nil(t, credentials)
	assert.Equal(t, "failed to read credentials: open invalid: no such file or directory", err.Error())
}

func TestReadCredentials_InvalidBase64(t *testing.T) {
	credentials, err := readCredentials(filepath.Join("docker_configs", "credentials", "invalid_base64"), getLogger())
	require.Error(t, err)
	require.Nil(t, credentials)
	assert.Equal(t, "invalid credential format: illegal base64 data at input byte 12", err.Error())
}

func TestReadCredentials_InvalidNoColon(t *testing.T) {
	credentials, err := readCredentials(filepath.Join("docker_configs", "credentials", "no_colon"), getLogger())
	require.Error(t, err)
	require.Nil(t, credentials)
	assert.Equal(t, "invalid credentials: no colon found", err.Error())
}

func TestReadPassphrase_ValidFile(t *testing.T) {
	passphrase := readPassphrase(filepath.Join(validPath, DefaultSecretsFolder, "mqtt_passphrase"), getLogger())
	require.NotNil(t, passphrase)
	assert.Equal(t, "secretPassphrase", *passphrase)
}

func TestReadPassphrase_InvalidFile(t *testing.T) {
	passphrase := readPassphrase("invalid", getLogger())
	require.Nil(t, passphrase)
}

func TestReadPassphrase_InvalidBase64(t *testing.T) {
	passphrase := readPassphrase(filepath.Join("docker_configs", "credentials", "invalid_base64"), getLogger())
	require.Nil(t, passphrase)
}

func TestParseBrokerConfiguration_InvalidFile(t *testing.T) {
	config, err := parseBrokerConfiguration("invalid/broker/configuration", getLogger())
	require.Error(t, err)
	require.Nil(t, config)
}

func TestReadPrivateKeyFile_Valid(t *testing.T) {
	privateKey, err := readPrivateKeyFile(filepath.Join("docker_configs", "private_key", "valid.pem"), getLogger())
	require.NoError(t, err)
	require.NotNil(t, privateKey)
}

func TestReadPrivateKeyFile_InvalidPemType(t *testing.T) {
	privateKey, err := readPrivateKeyFile(filepath.Join("docker_configs", "private_key", "invalid_pem_type.pem"), getLogger())
	require.Error(t, err)
	require.Nil(t, privateKey)
	assert.Equal(t, "invalid private key: invalid type CERTIFICATE", err.Error())

	privateKey, err = readPrivateKeyFile(filepath.Join("docker_configs", "private_key", "invalid_pem_file.pem"), getLogger())
	require.Error(t, err)
	require.Nil(t, privateKey)
	assert.Equal(t, "invalid private key: not in PEM format", err.Error())
}

func TestGetCAs_InvalidPath(t *testing.T) {
	root, cas, err := getCAs("invalid/path", "", getLogger())
	require.Error(t, err)
	require.Nil(t, root)
	require.Nil(t, cas)
	assert.Equal(t, "failed to read directory: open invalid/path: no such file or directory", err.Error())
}

func TestGetCAs_InvalidPem(t *testing.T) {
	root, cas, err := getCAs(filepath.Join("docker_configs", "ca_invalid_pem"), "", getLogger())
	require.Error(t, err)
	require.Nil(t, root)
	require.Nil(t, cas)
	assert.Equal(t, "invalid sub CA certificate ca.invalid_pem_file.pem: invalid PEM file", err.Error())
}

func TestReadPem_InvalidFile(t *testing.T) {
	pemBytes, err := os.ReadFile(filepath.Join("docker_configs", "private_key", "valid.pem"))
	require.NoError(t, err)
	pem, err := readPem(pemBytes)
	require.Error(t, err)
	require.Nil(t, pem)
	assert.Equal(t, "failed to parse certificate", err.Error())
}

func TestGetContainerName(t *testing.T) {
	hostName, err := os.Hostname()
	require.NoError(t, err)
	assert.Equal(t, hostName, GetContainerName())
}

func TestNewContainerStorage_valid(t *testing.T) {
	configuration := StorageConfiguration{
		ContainerName:                        "F12AB35",
		MessageBusStoragePath:                filepath.Join(validPath, DefaultMessageBusStorageSubFolder),
		Oi4CertificateStoragePath:            filepath.Join(validPath, DefaultOi4CertificateStorageSubFolder),
		SecretStoragePath:                    filepath.Join(validPath, DefaultSecretsFolder),
		ApplicationSpecificConfigurationPath: validPath,
		ApplicationSpecificDataPath:          validPath,
	}

	storage, err := NewContainerStorage(configuration, getLogger())
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestNewContainerStorage_invalidMessageBus(t *testing.T) {
	configuration := StorageConfiguration{
		ContainerName:                        "F12AB35",
		MessageBusStoragePath:                filepath.Join("invalid", "path"),
		Oi4CertificateStoragePath:            "",
		SecretStoragePath:                    "",
		ApplicationSpecificConfigurationPath: "",
		ApplicationSpecificDataPath:          "",
	}

	storage, err := NewContainerStorage(configuration, getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid message bus storage location: stat invalid/path: no such file or directory", err.Error())
}

func TestNewContainerStorage_invalidCertificateStorage(t *testing.T) {
	configuration := StorageConfiguration{
		ContainerName:                        "F12AB35",
		MessageBusStoragePath:                filepath.Join(validPath, DefaultMessageBusStorageSubFolder),
		Oi4CertificateStoragePath:            filepath.Join("invalid", "path"),
		SecretStoragePath:                    "",
		ApplicationSpecificConfigurationPath: "",
		ApplicationSpecificDataPath:          "",
	}

	storage, err := NewContainerStorage(configuration, getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid OI4 certificate storage location: stat invalid/path: no such file or directory", err.Error())
}

func TestNewContainerStorage_invalidSecretsStorage(t *testing.T) {
	configuration := StorageConfiguration{
		ContainerName:                        "F12AB35",
		MessageBusStoragePath:                filepath.Join(validPath, DefaultMessageBusStorageSubFolder),
		Oi4CertificateStoragePath:            filepath.Join(validPath, DefaultOi4CertificateStorageSubFolder),
		SecretStoragePath:                    filepath.Join("invalid", "path"),
		ApplicationSpecificConfigurationPath: "",
		ApplicationSpecificDataPath:          "",
	}

	storage, err := NewContainerStorage(configuration, getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid secret storage location: stat invalid/path: no such file or directory", err.Error())
}

func TestNewContainerStorage_invalidConfigurationStorage(t *testing.T) {
	configuration := StorageConfiguration{
		ContainerName:                        "F12AB35",
		MessageBusStoragePath:                filepath.Join(validPath, DefaultMessageBusStorageSubFolder),
		Oi4CertificateStoragePath:            filepath.Join(validPath, DefaultOi4CertificateStorageSubFolder),
		SecretStoragePath:                    filepath.Join(validPath, DefaultSecretsFolder),
		ApplicationSpecificConfigurationPath: "invalid",
		ApplicationSpecificDataPath:          "",
	}

	storage, err := NewContainerStorage(configuration, getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid configuration path: no such file or directory", err.Error())
}

func TestNewContainerStorage_invalidDataStorage(t *testing.T) {
	configuration := StorageConfiguration{
		ContainerName:                        "F12AB35",
		MessageBusStoragePath:                filepath.Join(validPath, DefaultMessageBusStorageSubFolder),
		Oi4CertificateStoragePath:            filepath.Join(validPath, DefaultOi4CertificateStorageSubFolder),
		SecretStoragePath:                    filepath.Join(validPath, DefaultSecretsFolder),
		ApplicationSpecificConfigurationPath: validPath,
		ApplicationSpecificDataPath:          "invalid",
	}

	storage, err := NewContainerStorage(configuration, getLogger())
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid data path: no such file or directory", err.Error())
}

func TestDefaultStoragePaths(t *testing.T) {
	storage := DefaultStorageConfiguration()
	hostName, err := os.Hostname()
	require.NoError(t, err)
	require.NotNil(t, storage)
	assert.Equal(t, hostName, storage.ContainerName)
	assert.Equal(t, DefaultMessageBusStorageSubFolder, storage.MessageBusStoragePath)
	assert.Equal(t, DefaultOi4CertificateStorageSubFolder, storage.Oi4CertificateStoragePath)
	assert.Equal(t, DefaultSecretsFolder, storage.SecretStoragePath)
	assert.Equal(t, DefaultApplicationSpecificConfigurationFolder, storage.ApplicationSpecificConfigurationPath)
	assert.Equal(t, DefaultApplicationSpecificDataFolder, storage.ApplicationSpecificDataPath)
}

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Print("Error flushing logger buffer: ", err)
		}
	}(logger) // flushes buffer, if any

	return logger.Sugar()
}
