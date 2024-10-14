package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/OI4/oi4-oec-service-go/service/application"
	"github.com/OI4/oi4-oec-service-go/service/container"
	"math/rand/v2"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/OI4/oi4-oec-service-go/service/api"
)

func main() {

	storage, mam, err := getStorage()
	if err != nil {
		wd, _ := os.Getwd()
		fmt.Println("Working directory:", wd)
		fmt.Println("Failed to retrieve storage configuration:", err)
		panic(err)
	}

	applicationSource := application.NewApplicationSourceImpl(*mam)

	oi4Application := application.CreateNewApplication(api.ServiceTypeOTConnector, applicationSource)

	source := api.Oi4Source(applicationSource)
	dataApplicationPublication := application.CreatePublication(api.ResourceData, &source).SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5)

	applicationTicker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			<-applicationTicker.C

			data := api.NewOi4Data(rand.Float64())

			addValue := func(key string, value any) {
				dErr := data.AddSecondaryData(key, &value)

				if dErr != nil {
					fmt.Println("Failed to add secondary data:", dErr)
				}
			}

			addValue("Sv1", rand.Float64())
			addValue("Sv2", rand.Float64())

			applicationSource.UpdateData(data, "Oi4Data")
		}

	}()

	err = oi4Application.RegisterPublication(dataApplicationPublication)
	if err != nil {
		fmt.Println("Failed to register publication:", err)
		panic(err)
	}

	metaDataApplicationPublication := application.CreatePublication(api.ResourceMetadata, &source).SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5)
	err = oi4Application.RegisterPublication(metaDataApplicationPublication)
	if err != nil {
		fmt.Println("Failed to register publication:", err)
		panic(err)
	}

	assetSource := application.NewSourceImpl(api.MasterAssetModel{
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

	oi4Asset := application.CreateNewAsset(assetSource)

	source = api.Oi4Source(assetSource)
	dataAssetPublication := application.CreatePublication(api.ResourceData, &source).SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5)
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
		fmt.Println("Failed to register publication:", err)
		panic(err)
	}

	metaDataAssetPublication := application.CreatePublication(api.ResourceMetadata, &source).SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5)
	err = oi4Asset.RegisterPublication(metaDataAssetPublication)
	if err != nil {
		fmt.Println("Failed to register publication:", err)
		panic(err)
	}

	oi4Application.RegisterAsset(oi4Asset)

	if err := oi4Application.Start(*storage); err != nil {
		panic(err)
	}

	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
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

func getStorage() (*container.Storage, *api.MasterAssetModel, error) {
	baseDir, runtime := getEnvironment()

	fmt.Printf("Using app as: %s with base dir: %s\n", runtime, baseDir)

	isContainer := runtime == "container"

	var config container.StorageConfiguration
	var mam *api.MasterAssetModel
	var err error
	if isContainer {
		config = *container.DefaultStorageConfiguration()
		mam, err = getMasterAssetModel(container.DefaultOi4Folder)
		if err != nil {
			return nil, nil, err
		}
	} else {
		mam, err = getMasterAssetModel(filepath.Join(baseDir, container.DefaultOi4Folder))
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
	}

	storage, err := container.NewContainerStorage(config)
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

func getEnvironment() (string, string) {
	baseDir, hasBaseDirEnv := os.LookupEnv("BASE_DIR")
	if !hasBaseDirEnv {
		baseDir = *flag.String("base", "", "base dir of the configuration")
	}

	runtime, hasRuntimeEnv := os.LookupEnv("RUNTIME")
	if !hasRuntimeEnv {
		runtime = *flag.String("runtime", "container", "runtime environment (program, container)")
	}

	if !hasBaseDirEnv || !hasRuntimeEnv {
		flag.Parse()
	}

	return baseDir, runtime
}
