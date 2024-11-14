//go:build demo

package main

import (
	"encoding/json"
	"github.com/OI4/oi4-oec-service-go/service/application"
	"github.com/OI4/oi4-oec-service-go/service/application/publication"
	"github.com/OI4/oi4-oec-service-go/service/application/source"
	"github.com/OI4/oi4-oec-service-go/service/container"
	"go.uber.org/zap"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/OI4/oi4-oec-service-go/service/api"
)

const baseDir = "./_demo/testdata/"

func main() {
	logger := getLogger()

	storage, mam, err := getStorage(logger)
	if err != nil {
		wd, _ := os.Getwd()
		logger.Info("Working directory:", wd)
		logger.Info("Failed to retrieve storage configuration:", err)
		panic(err)
	}

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

	dataApplicationPublication := publication.NewResourcePublication(oi4Application, applicationSource, api.ResourceData)

	applicationTicker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-applicationTicker.C

			data := api.NewOi4Data(rand.Float64())

			addValue := func(key string, value any) {
				dErr := data.AddSecondaryData(key, &value)

				if dErr != nil {
					logger.Error("Failed to add secondary data:", dErr)
				}
			}

			addValue("Sv1", rand.Float64())
			addValue("Sv2", rand.Float64())

			applicationSource.UpdateData(data, "Oi4Data")
		}

	}()

	err = oi4Application.RegisterPublication(dataApplicationPublication)
	if err != nil {
		logger.Fatal("Failed to register publication:", err)
		panic(err)
	}

	metaDataApplicationPublication := publication.NewResourcePublication(oi4Application, applicationSource, api.ResourceMetadata)

	err = oi4Application.RegisterPublication(metaDataApplicationPublication)
	if err != nil {
		logger.Fatal("Failed to register publication:", err)
		panic(err)
	}

	assetSource := source.NewAssetSourceImpl(api.MasterAssetModel{
		Manufacturer: api.LocalizedText{
			Locale: "en-US",
			Text:   "ACME",
		},
		ManufacturerUri: "acme.com",
		Model: api.LocalizedText{
			Locale: "en-US",
			Text:   "SampleAsset",
		},
		ProductCode:        "08",
		HardwareRevision:   "0",
		SoftwareRevision:   "0",
		DeviceRevision:     "0",
		DeviceManual:       "",
		DeviceClass:        "PerpetuumMobile",
		SerialNumber:       "15",
		ProductInstanceUri: "acme.com",
		RevisionCounter:    0,
		Description: api.LocalizedText{
			Locale: "en-US",
			Text:   "Cool Asset",
		},
	})

	oi4Asset := application.CreateNewAsset(assetSource, oi4Application)

	dataAssetPublication := publication.NewResourcePublicationWithFilter(oi4Application, assetSource, api.ResourceData, api.NewStringFilter("Oi4Data"))
	assetTicker := time.NewTicker(10 * time.Second)
	go func() {
		counter := 0
		for {
			<-assetTicker.C
			data := api.Oi4Data{PrimaryValue: counter}
			assetSource.UpdateData(&data, "Oi4Data")
			counter++
		}
	}()
	err = oi4Asset.RegisterPublication(dataAssetPublication)
	if err != nil {
		logger.Fatal("Failed to register publication:", err)
		panic(err)
	}

	metaDataAssetPublication := publication.NewResourcePublication(oi4Application, assetSource, api.ResourceMetadata)

	err = oi4Asset.RegisterPublication(metaDataAssetPublication)
	if err != nil {
		logger.Fatal("Failed to register publication:", err)
		panic(err)
	}

	oi4Application.RegisterAsset(oi4Asset)

	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(180 * time.Second)
		defer ticker.Stop()

		for i := 1; i <= 10; i++ {
			<-ticker.C

			maxH := 100
			minH := 0
			score := rand.IntN(maxH-minH) + minH
			aScore := score + 1
			if aScore > 100 {
				aScore = 100
			}
			applicationSource.UpdateHealth(api.Health{Health: api.Health_Normal, HealthScore: byte(score)})
			assetSource.UpdateHealth(api.Health{Health: api.Health_Normal, HealthScore: byte(aScore)})
		}

		applicationSource.UpdateHealth(api.Health{Health: api.Health_Normal, HealthScore: 100})
		assetSource.UpdateHealth(api.Health{Health: api.Health_Normal, HealthScore: 100})

		done <- true
	}()

	<-done

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	oi4Application.Stop()

	os.Exit(0)
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
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Print("Error flushing logger buffer: ", err)
		}
	}(logger) // flushes buffer, if any

	return logger.Sugar()
}
