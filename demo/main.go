package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/OI4/oi4-oec-service-go/container"
	"math/rand/v2"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	v1 "github.com/OI4/oi4-oec-service-go/api/pkg/types"
	application "github.com/OI4/oi4-oec-service-go/service"
	"github.com/OI4/oi4-oec-service-go/service/pkg/oi4"
)

func main() {

	configuration, mam, err := getStorage()
	if err != nil {
		wd, _ := os.Getwd()
		fmt.Println("Working directory:", wd)
		fmt.Println("Failed to retrieve storage configuration:", err)
		panic(err)
	}

	applicationSource := oi4.NewApplicationSourceImpl(*mam)

	oi4Application := application.CreateNewApplication(v1.ServiceTypeOTConnector, applicationSource)

	dataApplicationPublication := application.CreatePublication[v1.Oi4Data](v1.ResourceData, false).SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5)

	applicationTicker := time.NewTicker(10 * time.Second)
	go func() {
		counter := 0
		for {
			<-applicationTicker.C

			var data any = v1.Oi4Data{PrimaryValue: counter}
			applicationSource.UpdateData(&data, "Oi4Data")

			counter++
		}

	}()

	err = oi4Application.RegisterPublication(dataApplicationPublication)
	if err != nil {
		fmt.Println("Failed to register publication:", err)
		panic(err)
	}

	metaDataApplicationPublication := application.CreatePublication[any](v1.ResourceMetadata, false).SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5)
	err = oi4Application.RegisterPublication(metaDataApplicationPublication)
	if err != nil {
		fmt.Println("Failed to register publication:", err)
		panic(err)
	}

	assetSource := oi4.NewSourceImpl(v1.MasterAssetModel{
		Manufacturer: v1.LocalizedText{
			Locale: "en-US",
			Text:   "ACME",
		},
		ManufacturerUri: "acme.com",
		Model: v1.LocalizedText{
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
		Description: v1.LocalizedText{
			Locale: "en-US",
			Text:   "Cool Asset",
		},
	})

	oi4Asset := application.CreateNewAsset(assetSource)

	dataAssetPublication := application.CreatePublication[v1.Oi4Data](v1.ResourceData, false).SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5)
	assetTicker := time.NewTicker(10 * time.Second)
	go func() {
		counter := 0
		for {
			<-assetTicker.C
			data := v1.Oi4Data{PrimaryValue: counter}
			dataAssetPublication.SetData(data)
			counter++
		}
	}()
	err = oi4Asset.RegisterPublication(dataAssetPublication)
	if err != nil {
		fmt.Println("Failed to register publication:", err)
		panic(err)
	}

	metaDataAssetPublication := application.CreatePublication[any](v1.ResourceMetadata, false).SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5)
	err = oi4Asset.RegisterPublication(metaDataAssetPublication)
	if err != nil {
		fmt.Println("Failed to register publication:", err)
		panic(err)
	}

	oi4Application.RegisterAsset(oi4Asset)

	if err := oi4Application.Start(*configuration); err != nil {
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
			applicationSource.UpdateHealth(v1.Health{Health: v1.Health_Normal, HealthScore: byte(score)})
			assetSource.UpdateHealth(v1.Health{Health: v1.Health_Normal, HealthScore: byte(aScore)})
		}

		applicationSource.UpdateHealth(v1.Health{Health: v1.Health_Normal, HealthScore: 100})
		assetSource.UpdateHealth(v1.Health{Health: v1.Health_Normal, HealthScore: 100})
		done <- true
	}()

	<-done

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	oi4Application.Stop()

	os.Exit(0)
}

func getStorage() (*container.Storage, *v1.MasterAssetModel, error) {
	baseDir := flag.String("base", "", "base dir of the configuration")
	runtime := flag.String("runtime", "container", "runtime environment (program, container)")
	flag.Parse()

	fmt.Printf("Using app as: %s with base dir: %s\n", *runtime, *baseDir)

	isContainer := *runtime == "container"

	var config container.StorageConfiguration
	var mam *v1.MasterAssetModel
	var err error
	if isContainer {
		config = *container.DefaultStorageConfiguration()
		mam, err = getMasterAssetModel(container.DefaultOi4Folder)
		if err != nil {
			return nil, nil, err
		}
	} else {
		mam, err = getMasterAssetModel(filepath.Join(*baseDir, container.DefaultOi4Folder))
		if err != nil {
			return nil, nil, err
		}
		config = container.StorageConfiguration{
			ContainerName:                        mam.SerialNumber,
			MessageBusStoragePath:                filepath.Join(*baseDir, container.DefaultMessageBusStorageSubFolder),
			Oi4CertificateStoragePath:            filepath.Join(*baseDir, container.DefaultOi4CertificateStorageSubFolder),
			SecretStoragePath:                    filepath.Join(*baseDir, container.DefaultSecretsFolder),
			ApplicationSpecificConfigurationPath: filepath.Join(*baseDir, container.DefaultApplicationSpecificConfigurationFolder),
			ApplicationSpecificDataPath:          filepath.Join(*baseDir, container.DefaultApplicationSpecificDataFolder),
		}
	}

	storage, err := container.NewContainerStorage(config)
	if err != nil {
		return nil, nil, err
	}

	return storage, mam, nil
}

func getMasterAssetModel(oi4Dir string) (*v1.MasterAssetModel, error) {
	mamFile := filepath.Join(oi4Dir, "config", "mam.json")
	fileBytes, err := os.ReadFile(mamFile)
	if err != nil {
		return nil, &v1.Error{
			Message: "Failed to read master asset model file from: " + mamFile,
			Err:     err,
		}
	}
	var mam v1.MasterAssetModel
	err = json.Unmarshal(fileBytes, &mam)
	if err != nil {
		return nil, &v1.Error{
			Message: "Failed to unmarshal master asset model file from: " + mamFile,
			Err:     err,
		}
	}

	return &mam, nil
}
