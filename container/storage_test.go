package container

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

const validPath = "testdata/valid"

func TestNewMessageBusStorage_ValidPath(t *testing.T) {
	storage, err := NewMessageBusStorage(filepath.Join(validPath, DefaultMessageBusStorageSubFolder))
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestNewMessageBusStorage_InvalidPath(t *testing.T) {
	storage, err := NewMessageBusStorage("invalid/path")
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid message bus storage location: stat invalid/path: no such file or directory", err.Error())
}

func TestNewMessageBusStorage_InvalidCa(t *testing.T) {
	storage, err := NewMessageBusStorage(filepath.Join("testdata", "mqtt_invalid_ca"))
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid CA certificates: invalid root CA certificate: invalid PEM file", err.Error())
}

func TestNewMessageBusStorage_InvalidBrokerConfiguration(t *testing.T) {
	storage, err := NewMessageBusStorage(filepath.Join("testdata", "mqtt_invalid_broker_configuration"))
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid broker configuration: invalid json format: unexpected end of JSON input", err.Error())
}

func TestNewOi4CertificateStorage_ValidPath(t *testing.T) {
	storage, err := NewOi4CertificateStorage(filepath.Join(validPath, DefaultOi4CertificateStorageSubFolder), "F12AB35")
	require.NoError(t, err)
	require.NotNil(t, storage)
}

func TestNewOi4CertificateStorage_InvalidPath(t *testing.T) {
	storage, err := NewOi4CertificateStorage("invalid/path", "containerName")
	require.Error(t, err)
	require.Nil(t, storage)
}

func TestNewOi4CertificateStorage_InvalidClientPem(t *testing.T) {
	storage, err := NewOi4CertificateStorage(filepath.Join(validPath, DefaultOi4CertificateStorageSubFolder), "NoId")
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid client certificate: failed to read pem file", err.Error())
}

func TestNewOi4CertificateStorage_InvalidCA(t *testing.T) {
	storage, err := NewOi4CertificateStorage(filepath.Join("testdata", "certificate_invalid_ca"), "F12AB35")
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid CA certificates: invalid root CA certificate: invalid PEM file", err.Error())
}

func TestNewSecretStorage_ValidPath(t *testing.T) {
	storage, err := NewSecretStorage(filepath.Join(validPath, DefaultSecretsFolder))
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
	storage, err := NewSecretStorage("invalid/path")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid secret storage location: stat invalid/path: no such file or directory", err.Error())
}

func TestNewSecretStorage_InvalidCredentials(t *testing.T) {
	storage, err := NewSecretStorage(filepath.Join("testdata", "empty"))
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
	storage, err := NewApplicationSpecificStorages("invalid", filepath.Join("testdata", "empty"))
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid configuration path: no such file or directory", err.Error())
}

func TestNewApplicationSpecificStorages_InvalidDataPath(t *testing.T) {
	storage, err := NewApplicationSpecificStorages(filepath.Join("testdata", "empty"), "invalid")
	require.Error(t, err)
	require.Nil(t, storage)
	assert.Equal(t, "invalid data path: no such file or directory", err.Error())
}

func TestReadCredentials_ValidFile(t *testing.T) {
	credentials, err := readCredentials(filepath.Join(validPath, DefaultSecretsFolder, "mqtt_credentials"))
	require.NoError(t, err)
	require.NotNil(t, credentials)
	assert.Equal(t, "oi4User", credentials.Username())
	pwd, set := credentials.Password()
	assert.True(t, set)
	assert.Equal(t, "oi4Password", pwd)
}

func TestReadCredentials_InvalidFile(t *testing.T) {
	credentials, err := readCredentials("invalid")
	require.Error(t, err)
	require.Nil(t, credentials)
	assert.Equal(t, "failed to read credentials: open invalid: no such file or directory", err.Error())
}

func TestReadCredentials_InvalidBase64(t *testing.T) {
	credentials, err := readCredentials(filepath.Join("testdata", "credentials", "invalid_base64"))
	require.Error(t, err)
	require.Nil(t, credentials)
	assert.Equal(t, "invalid credential format: illegal base64 data at input byte 12", err.Error())
}

func TestReadCredentials_InvalidNoColon(t *testing.T) {
	credentials, err := readCredentials(filepath.Join("testdata", "credentials", "no_colon"))
	require.Error(t, err)
	require.Nil(t, credentials)
	assert.Equal(t, "invalid credentials: no colon found", err.Error())
}

func TestReadPassphrase_ValidFile(t *testing.T) {
	passphrase := readPassphrase(filepath.Join(validPath, DefaultSecretsFolder, "mqtt_passphrase"))
	require.NotNil(t, passphrase)
	assert.Equal(t, "secretPassphrase", *passphrase)
}

func TestReadPassphrase_InvalidFile(t *testing.T) {
	passphrase := readPassphrase("invalid")
	require.Nil(t, passphrase)
}

func TestReadPassphrase_InvalidBase64(t *testing.T) {
	passphrase := readPassphrase(filepath.Join("testdata", "credentials", "invalid_base64"))
	require.Nil(t, passphrase)
}

func TestParseBrokerConfiguration_InvalidFile(t *testing.T) {
	config, err := parseBrokerConfiguration("invalid/broker/configuration")
	require.Error(t, err)
	require.Nil(t, config)
}

func TestReadPrivateKeyFile_Valid(t *testing.T) {
	privateKey, err := readPrivateKeyFile(filepath.Join("testdata", "private_key", "valid.pem"))
	require.NoError(t, err)
	require.NotNil(t, privateKey)
}

func TestReadPrivateKeyFile_InvalidPemType(t *testing.T) {
	privateKey, err := readPrivateKeyFile(filepath.Join("testdata", "private_key", "invalid_pem_type.pem"))
	require.Error(t, err)
	require.Nil(t, privateKey)
	assert.Equal(t, "invalid private key: invalid type CERTIFICATE", err.Error())

	privateKey, err = readPrivateKeyFile(filepath.Join("testdata", "private_key", "invalid_pem_file.pem"))
	require.Error(t, err)
	require.Nil(t, privateKey)
	assert.Equal(t, "invalid private key: not in PEM format", err.Error())
}

func TestGetCAs_InvalidPath(t *testing.T) {
	root, cas, err := getCAs("invalid/path", "")
	require.Error(t, err)
	require.Nil(t, root)
	require.Nil(t, cas)
	assert.Equal(t, "failed to read directory: open invalid/path: no such file or directory", err.Error())
}

func TestGetCAs_InvalidPem(t *testing.T) {
	root, cas, err := getCAs(filepath.Join("testdata", "ca_invalid_pem"), "")
	require.Error(t, err)
	require.Nil(t, root)
	require.Nil(t, cas)
	assert.Equal(t, "invalid sub CA certificate ca.invalid_pem_file.pem: invalid PEM file", err.Error())
}

func TestReadPem_InvalidFile(t *testing.T) {
	pemBytes, err := os.ReadFile(filepath.Join("testdata", "private_key", "valid.pem"))
	require.NoError(t, err)
	pem, err := readPem(pemBytes)
	require.Error(t, err)
	require.Nil(t, pem)
	assert.Equal(t, "failed to parse certificate", err.Error())
}
