//go:build demo

package main

import (
	"encoding/json"
	"fmt"
	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/OI4/oi4-oec-service-go/service/application"
	"github.com/OI4/oi4-oec-service-go/service/application/source"
	"github.com/OI4/oi4-oec-service-go/service/application/subscription"
	"github.com/OI4/oi4-oec-service-go/service/container"
	tp "github.com/OI4/oi4-oec-service-go/service/topic"
	"go.uber.org/zap"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const baseDir = "./_demo/subscriber/testdata/"

func main() {
	logger := getLogger()

	storage, mam, err := getStorage(logger)
	if err != nil {
		wd, _ := os.Getwd()
		logger.Info("Working directory:", wd)
		logger.Info("Failed to retrieve storage configuration:", err)
		panic(err)
	}

	/**********************************/
	/* Create an application and source */
	/**********************************/
	option := source.WithHealthFn(
		func(_ api.BaseSource) api.Health {
			maxH := 100
			minH := 0
			score := rand.IntN(maxH-minH) + minH
			aScore := score + 1
			if aScore > 100 {
				aScore = 100
			}

			return api.Health{
				Health:      api.Health_Normal,
				HealthScore: byte(score),
			}
		})

	ltOption := source.WithLicenseText(map[string]api.LicenseText{
		"MIT": {
			"Lorem Ipsum",
		},
	})

	applicationSource := source.NewApplicationSourceImpl(*mam, option, ltOption)

	oi4Application, err := application.CreateNewApplication(api.ServiceTypeOTConnector, applicationSource, logger)
	if err != nil {
		logger.Fatal("Failed to create application:", err)
		panic(err)
	}

	if err = oi4Application.Start(*storage); err != nil {
		logger.Fatal("Failed to start application:", err)
		panic(err)
	}

	subscriptionTopic := fmt.Sprintf("%s/+/+/+/+/+/%s/%s/#", tp.Oi4Namespace, api.MethodPub, api.ResourceHealth)
	dataApplicationSubscription := subscription.NewTopicSubscription(subscriptionTopic, handler(oi4Application))
	err = oi4Application.RegisterSubscription(dataApplicationSubscription)
	if err != nil {
		logger.Fatal("Failed to register publication:", err)
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	oi4Application.Stop()

	os.Exit(0)
}

func handler(app *application.Oi4ApplicationImpl) api.MessageHandler {
	return subscription.NewMessageHandler(app, func(resource api.ResourceType, source *api.Oi4Identifier, networkMessage api.NetworkMessage, topic *tp.Topic) {
		logger := app.GetLogger()
		logger.Infof("Received message on topic: %s", topic.ToString())
		logger.Infof("Message ID: %s", networkMessage.MessageId)
		logger.Infof("Resource Type: %s", resource)
		logger.Infof("For OI4 ID: %s", source.ToString())
	})
}

func getStorage(logger *zap.SugaredLogger) (*container.Storage, *api.MasterAssetModel, error) {
	var config container.StorageConfiguration

	mam, err := getMasterAssetModel(filepath.Join(baseDir, container.DefaultOi4Folder))
	if err != nil {
		return nil, nil, err
	}
	config = container.StorageConfiguration{
		ContainerName:                        mam.SerialNumber,
		MessageBusStoragePath:                filepath.Join(baseDir, container.DefaultMessageBusStorageSubFolder),
		Oi4CertificateStoragePath:            filepath.Join(baseDir, container.DefaultOi4CertificateStorageSubFolder),
		SecretStoragePath:                    filepath.Join(baseDir, container.DefaultSecretsFolder),
		ApplicationSpecificConfigurationPath: filepath.Join(baseDir, container.DefaultApplicationSpecificConfigurationFolder),
		ApplicationSpecificDataPath:          filepath.Join(baseDir, container.DefaultApplicationSpecificDataFolder),
	}

	storage, err := container.NewContainerStorage(config, logger)
	if err != nil {
		return nil, nil, err
	}

	return storage, mam, nil
}

func getMasterAssetModel(oi4Dir string) (*api.MasterAssetModel, error) {
	mamFile := filepath.Join(oi4Dir, "config", "mam.json")
	fileBytes, err := os.ReadFile(mamFile)
	if err != nil {
		return nil, &api.Error{
			Message: "Failed to read master asset model file from: " + mamFile,
			Err:     err,
		}
	}
	var mam api.MasterAssetModel
	err = json.Unmarshal(fileBytes, &mam)
	if err != nil {
		return nil, &api.Error{
			Message: "Failed to unmarshal master asset model file from: " + mamFile,
			Err:     err,
		}
	}

	return &mam, nil
}

func getLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, _ := config.Build()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Print("Error flushing logger buffer: ", err)
		}
	}(logger) // flushes buffer, if any

	return logger.Sugar()
}
