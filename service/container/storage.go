// Package container provides the storage for the container according to the OI4 OEC Guideline.
package container

import (
	"bytes"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const DefaultOi4Folder = "/etc/oi4"
const DefaultSecretsFolder = "/run/secrets"
const DefaultOi4CertificateStorageSubFolder = DefaultOi4Folder + "/certs"
const DefaultMessageBusStorageSubFolder = DefaultOi4Folder + "/mqtt"

const DefaultApplicationSpecificConfigurationFolder = DefaultOi4Folder + "/app"
const DefaultApplicationSpecificDataFolder = "/opt/oi4/app"

type BaseStorage struct {
	FolderPath *string
}

func newStorage(folderPath string) (*BaseStorage, error) {
	_, err := os.Stat(folderPath)
	if err != nil {
		return nil, err
	}
	return &BaseStorage{FolderPath: &folderPath}, nil
}

// ****************************************************************
// ***                                  MessageBusStorage                                  ***
// ****************************************************************

type StorageConfiguration struct {
	ContainerName                        string
	MessageBusStoragePath                string
	Oi4CertificateStoragePath            string
	SecretStoragePath                    string
	ApplicationSpecificConfigurationPath string
	ApplicationSpecificDataPath          string
}

func DefaultStorageConfiguration() *StorageConfiguration {
	return &StorageConfiguration{
		ContainerName:                        GetContainerName(),
		MessageBusStoragePath:                DefaultMessageBusStorageSubFolder,
		Oi4CertificateStoragePath:            DefaultOi4CertificateStorageSubFolder,
		SecretStoragePath:                    DefaultSecretsFolder,
		ApplicationSpecificConfigurationPath: DefaultApplicationSpecificConfigurationFolder,
		ApplicationSpecificDataPath:          DefaultApplicationSpecificDataFolder,
	}
}

type Storage struct {
	ContainerName               string
	MessageBusStorage           *MessageBusStorage
	Oi4CertificateStorage       *Oi4CertificateStorage
	SecretStorage               *SecretStorage
	ApplicationSpecificStorages *ApplicationSpecificStorages
}

func NewContainerStorage(configuration StorageConfiguration, logger *zap.SugaredLogger) (*Storage, error) {
	messageBusStorage, err := NewMessageBusStorage(configuration.MessageBusStoragePath, logger)
	if err != nil {
		return nil, err
	}

	oi4CertificateStorage, err := NewOi4CertificateStorage(configuration.Oi4CertificateStoragePath, configuration.ContainerName, logger)
	if err != nil {
		return nil, err
	}

	secretStorage, err := NewSecretStorage(configuration.SecretStoragePath, logger)
	if err != nil {
		return nil, err
	}

	applicationSpecificStorages, err := NewApplicationSpecificStorages(configuration.ApplicationSpecificConfigurationPath, configuration.ApplicationSpecificDataPath)
	if err != nil {
		return nil, err
	}

	return &Storage{
		ContainerName:               configuration.ContainerName,
		MessageBusStorage:           messageBusStorage,
		Oi4CertificateStorage:       oi4CertificateStorage,
		SecretStorage:               secretStorage,
		ApplicationSpecificStorages: applicationSpecificStorages,
	}, nil
}

func GetContainerName() string {
	containerName, _ := os.Hostname()
	return containerName
}

// ****************************************************************
// ***                                  MessageBusStorage                                  ***
// ****************************************************************

type MessageBusStorage struct {
	BaseStorage
	// Public certificate of the broker.
	BrokerCertificate *x509.Certificate
	// Local CA used to sign the MQTT broker certificate.
	// If no local CA is present, the global CA of the secret storage is used.
	RootCaCertificate *x509.Certificate
	// Local Sub-CAs used to sign the MQTT broker certificate.
	//	If no local Sub-CAs are present, the global Sub-CAs of the secret storage is used.
	SubCaCertificates map[string]*x509.Certificate
	// Configuration containing the settings of the MQTT broker.
	BrokerConfiguration *BrokerConfiguration
}

type BrokerConfiguration struct {

	// Defines the IP address or host name of the MQTT broker used.
	Address string

	//Defines the port, where the MQTT broker is listening.
	SecurePort uint16

	// Defines the maximum size of a MQTT payload in KiB. A NetworkMessage, containing several DataSetMessages, shall not exceed this size.
	// The minimum value of MaxPacketSize is 262,144 bytes (256 kiB), the theoretical maximum is INT32max.
	// The intention of the MaxPacketSize is to protect the Message Bus broker and client.
	MaxPacketSize int32
}

func NewMessageBusStorage(folderPath string, logger *zap.SugaredLogger) (*MessageBusStorage, error) {
	storage, err := newStorage(folderPath)
	if err != nil {
		return nil, &Error{
			Message: "invalid message bus storage location",
			Err:     err,
		}
	}
	brokerPem := filepath.Join(*storage.FolderPath, "broker.pem")
	brokerCertificate, err := readPemFile(brokerPem, logger)
	if err != nil {
		// The guideline defines the broker certificate as mandatory.
		// Nevertheless, the broker certificate is not strictly required for connecting to the broker.
		logger.Error("invalid broker certificate. Certificate is going to be skipped", err)
	}

	rootCertificate, subCaCertificates, err := getCAs(folderPath, "broker_", logger)
	// CA and Sub-CA certificates are optional. But if present, they must be valid.
	if err != nil {
		return nil, &Error{
			Message: "invalid CA certificates",
			Err:     err,
		}
	}

	configuration, err := parseBrokerConfiguration(filepath.Join(folderPath, "broker.json"), logger)
	if err != nil {
		return nil, &Error{
			Message: "invalid broker configuration",
			Err:     err,
		}
	}

	return &MessageBusStorage{
		BaseStorage:         *storage,
		BrokerCertificate:   brokerCertificate,
		RootCaCertificate:   rootCertificate,
		SubCaCertificates:   subCaCertificates,
		BrokerConfiguration: configuration,
	}, nil
}

func parseBrokerConfiguration(configPath string, logger *zap.SugaredLogger) (*BrokerConfiguration, error) {
	configBytes, err := readFile(configPath, logger)
	if err != nil {
		return nil, errors.New("failed to read broker configuration")
	}

	var configuration BrokerConfiguration
	err = json.Unmarshal(configBytes, &configuration)
	if err != nil {
		return nil, &Error{
			Message: "invalid json format",
			Err:     err,
		}
	}

	return &configuration, nil
}

// ****************************************************************
// ***                                 Oi4CertificateStorage                                 ***
// ****************************************************************

type Oi4CertificateStorage struct {
	BaseStorage
	ClientCertificate *x509.Certificate
	// CA used for validation of certificates.
	RootCaCertificate *x509.Certificate
	// Sub-CAs used for validation of certificates.
	SubCaCertificates map[string]*x509.Certificate
}

func NewOi4CertificateStorage(folderPath string, containerName string, logger *zap.SugaredLogger) (*Oi4CertificateStorage, error) {
	storage, err := newStorage(folderPath)
	if err != nil {
		return nil, &Error{
			Message: "invalid OI4 certificate storage location",
			Err:     err,
		}
	}
	clientPem := filepath.Join(*storage.FolderPath, containerName+".pem")
	clientCertificate, err := readPemFile(clientPem, logger)
	if err != nil {
		return nil, &Error{
			Message: "invalid client certificate",
			Err:     err,
		}
	}

	rootCertificate, subCaCertificates, err := getCAs(folderPath, "", logger)
	if err != nil {
		return nil, &Error{
			Message: "invalid CA certificates",
			Err:     err,
		}
	}

	return &Oi4CertificateStorage{
		BaseStorage:       *storage,
		ClientCertificate: clientCertificate,
		RootCaCertificate: rootCertificate,
		SubCaCertificates: subCaCertificates,
	}, nil
}

// ****************************************************************
// ***                                       SecretStorage                                        ***
// ****************************************************************

type SecretStorage struct {
	BaseStorage
	// If present, it contains a username and password, which is used for MQTT authentication
	MqttCredentials *url.Userinfo
	// If present, it contains the private PEM encoded key to the according client certificate in the OI4 certificate store.
	MqttPrivateKey *pem.Block
	// If present, it contains the base64 encoded passphrase for the private key.
	MqttPassphrase *string
}

func NewSecretStorage(folderPath string, logger *zap.SugaredLogger) (*SecretStorage, error) {
	storage, err := newStorage(folderPath)
	if err != nil {
		return nil, &Error{
			Message: "invalid secret storage location",
			Err:     err,
		}
	}

	credentials, err := readCredentials(filepath.Join(folderPath, "mqtt_credentials"), logger)
	if err != nil {
		fmt.Println("no or invalid mqtt credentials. Credentials are going to be skipped")
		// TODO log error in trace level
	}

	privateKey, err := readPrivateKeyFile(filepath.Join(folderPath, "mqtt_private_key.pem"), logger)
	if err != nil {
		fmt.Println("no or invalid mqtt private key. Private key is going to be skipped")
		// TODO log error in trace level
	}

	passphrase := readPassphrase(filepath.Join(folderPath, "mqtt_passphrase"), logger)

	return &SecretStorage{
		BaseStorage:     *storage,
		MqttCredentials: credentials,
		MqttPrivateKey:  privateKey,
		MqttPassphrase:  passphrase,
	}, nil
}

func readCredentials(credentialFile string, logger *zap.SugaredLogger) (*url.Userinfo, error) {
	credentialBytes, err := readFile(credentialFile, logger)
	if err != nil {
		return nil, &Error{
			Message: "failed to read credentials",
			Err:     err,
		}
	}

	decoded, err := base64.StdEncoding.DecodeString(string(credentialBytes))
	if err != nil {
		return nil, &Error{
			Message: "invalid credential format",
			Err:     err,
		}
	}

	if !strings.Contains(string(decoded), ":") {
		return nil, errors.New("invalid credentials: no colon found")
	}

	credentials := strings.Split(string(decoded), ":")
	return url.UserPassword(credentials[0], credentials[1]), nil
}

func readPassphrase(passphraseFile string, logger *zap.SugaredLogger) *string {
	passphraseBytes, err := readFile(passphraseFile, logger)
	if err != nil {
		return nil
	}

	decoded, err := base64.StdEncoding.DecodeString(string(passphraseBytes))
	if err != nil {
		fmt.Println("invalid passphrase format", err)
		return nil
	}
	passphrase := string(decoded)
	return &passphrase
}

// ****************************************************************
// ***                            ApplicationSpecificStorages                            ***
// ****************************************************************

type ApplicationSpecificStorages struct {
	// Application specific configuration shall be stored here
	ConfigurationPath string

	// Application specific data shall be stored here
	DataPath string
}

func NewApplicationSpecificStorages(configurationPath string, dataPath string) (*ApplicationSpecificStorages, error) {
	_, err := os.Stat(configurationPath)
	if err != nil {
		errMsg, _ := strings.CutPrefix(err.Error(), "stat invalid: ")
		return nil, errors.New(fmt.Sprintf("invalid configuration path: %s", errMsg))
	}

	_, err = os.Stat(dataPath)
	if err != nil {
		errMsg, _ := strings.CutPrefix(err.Error(), "stat invalid: ")
		return nil, errors.New(fmt.Sprintf("invalid data path: %s", errMsg))
	}
	return &ApplicationSpecificStorages{ConfigurationPath: configurationPath, DataPath: dataPath}, nil
}

// ****************************************************************
// ***                                     Helper functions                                      ***
// ****************************************************************

func readPrivateKeyFile(pemFile string, logger *zap.SugaredLogger) (*pem.Block, error) {
	pemBytes, err := readFile(pemFile, logger)
	if err != nil {
		return nil, &Error{
			Message: "failed to read private key pem file",
			Err:     err,
		}
	}
	return readPrivateKey(pemBytes)
}

func readPrivateKey(pemBytes []byte) (*pem.Block, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid private key: not in PEM format")
	}
	if block.Type != "PRIVATE KEY" {
		return nil, errors.New(fmt.Sprintf("invalid private key: invalid type %s", block.Type))
	}
	return block, nil
}

func getCAs(folderPath string, caPrefix string, logger *zap.SugaredLogger) (*x509.Certificate, map[string]*x509.Certificate, error) {
	var rootCertificate *x509.Certificate
	subCaCertificates := make(map[string]*x509.Certificate)

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, nil, &Error{
			Message: "failed to read directory",
			Err:     err,
		}
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".pem" {
			continue
		}
		if file.Name() == caPrefix+"ca.pem" {
			rootCertificate, err = readPemFile(filepath.Join(folderPath, file.Name()), logger)
			if err != nil {
				return nil, nil, &Error{
					Message: "invalid root CA certificate",
					Err:     err,
				}
			}
		} else if strings.HasPrefix(filepath.Base(file.Name()), caPrefix+"ca.") {
			var subCaCertificate *x509.Certificate
			subCaCertificate, err = readPemFile(filepath.Join(folderPath, file.Name()), logger)
			if err != nil {
				return nil, nil, &Error{
					Message: fmt.Sprintf("invalid sub CA certificate %s", file.Name()),
					Err:     err,
				}
			}
			subCaCertificates[file.Name()] = subCaCertificate
		}

	}
	return rootCertificate, subCaCertificates, nil
}

func readPemFile(pemFile string, logger *zap.SugaredLogger) (*x509.Certificate, error) {
	pemBytes, err := readFile(pemFile, logger)
	if err != nil {
		return nil, errors.New("failed to read pem file")
	}
	return readPem(pemBytes)
}

func readPem(pemBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid PEM file")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse certificate")
	}

	return cert, nil
}

func readFile(filePath string, logger *zap.SugaredLogger) ([]byte, error) {
	// Open the file
	file, err := os.Open(filePath)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logger.Debug("Error closing file: "+filePath, err)
		}
	}(file)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	return buffer.Bytes(), err
}
