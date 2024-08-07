package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/OI4/oi4-oec-service-go/api/pkg/types"
	application "github.com/OI4/oi4-oec-service-go/service"
	"github.com/OI4/oi4-oec-service-go/service/pkg/mqtt"
)

func main() {

	oi4Application := application.CreateNewApplication(v1.ServiceTypeOTConnector, &v1.MasterAssetModel{
		Manufacturer: v1.LocalizedText{
			"en_us": "ACME",
		},
		ManufacturerUri:    "acme.com",
		Model:              "SampleApplication",
		ProductCode:        "FC#156",
		HardwareRevision:   "0",
		SoftwareRevision:   "0",
		DeviceRevision:     "0",
		DeviceManual:       "",
		DeviceClass:        `PerpetuumMobile`,
		SerialNumber:       "F87263976#4",
		ProductInstanceUri: "acme.com",
		RevisionCounter:    0,
		Description: v1.LocalizedText{
			"en_us": "Cool Application",
		},
	})

	dataApplicationPublication := application.CreatePublication[v1.Oi4Data](v1.ResourceData, false).SetPublicationMode(v1.PublicationMode_APPLICATION_2)
	applicationTicker := time.NewTicker(10 * time.Second)
	go func() {
		counter := 0
		for {
			<-applicationTicker.C
			dataApplicationPublication.SetData(v1.Oi4Data{PrimaryValue: counter})
			counter++
		}

	}()

	oi4Application.RegisterPublication(dataApplicationPublication)

	oi4Asset := application.CreateNewAsset(&v1.MasterAssetModel{
		Manufacturer: v1.LocalizedText{
			"en_us": "ACME",
		},
		ManufacturerUri:    "acme.com",
		Model:              "SampleAsset",
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
			"en_us": "Cool Asset",
		},
	})

	dataAssetPublication := application.CreatePublication[v1.Oi4Data](v1.ResourceData, false).SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5)
	assetTicker := time.NewTicker(10 * time.Second)
	go func() {
		counter := 0
		for {
			<-assetTicker.C
			dataAssetPublication.SetData(v1.Oi4Data{PrimaryValue: counter})
			counter++
		}

	}()
	oi4Asset.RegisterPublication(dataAssetPublication)

	oi4Application.RegisterAsset(oi4Asset)

	if err := oi4Application.Start(&mqtt.MQTTClientOptions{
		Host:     "127.0.0.1",
		Port:     8883,
		Tls:      true,
		Username: "oi4",
		Password: "oi4",
	}); err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	oi4Application.Stop()

	os.Exit(0)
}
