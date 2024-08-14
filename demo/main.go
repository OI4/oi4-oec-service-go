package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/OI4/oi4-oec-service-go/api/pkg/types"
	application "github.com/OI4/oi4-oec-service-go/service"
	"github.com/OI4/oi4-oec-service-go/service/pkg/mqtt"
	"github.com/OI4/oi4-oec-service-go/service/pkg/oi4"
)

func main() {

	applicationSource := oi4.NewApplicationSourceImpl(v1.MasterAssetModel{
		Manufacturer: v1.LocalizedText{
			Locale: "en-US",
			Text:   "ACME",
		},
		ManufacturerUri: "acme.com",
		Model: v1.LocalizedText{
			Locale: "en-US",
			Text:   "SampleApplication",
		},
		ProductCode:        "FC#156",
		HardwareRevision:   "",
		SoftwareRevision:   "0",
		DeviceRevision:     "0",
		DeviceManual:       "",
		DeviceClass:        fmt.Sprintf("Oi4.%s", string(v1.ServiceTypeOTConnector)),
		SerialNumber:       "F87263976#4",
		ProductInstanceUri: "acme.com",
		RevisionCounter:    0,
		Description: v1.LocalizedText{
			Locale: "en-US",
			Text:   "Cool Application",
		},
	})

	oi4Application := application.CreateNewApplication(v1.ServiceTypeOTConnector, applicationSource)

	dataApplicationPublication := application.CreatePublication[v1.Oi4Data](v1.ResourceData, false).SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5)
	applicationTicker := time.NewTicker(10 * time.Second)
	go func() {
		counter := 0
		for {
			<-applicationTicker.C
			dataApplicationPublication.SetData(v1.Oi4Data{PrimaryValue: counter})

			data := v1.Oi4Data{PrimaryValue: counter}
			//applicationSource.UpdateData(*data)
			dataApplicationPublication.SetData(data)

			counter++
		}

	}()

	oi4Application.RegisterPublication(dataApplicationPublication)

	metaDataApplicationPublication := application.CreatePublication[any](v1.ResourceMetadata, false).SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5)
	oi4Application.RegisterPublication(metaDataApplicationPublication) //.SetData().SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5))

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
	oi4Asset.RegisterPublication(dataAssetPublication)

	metaDataAssetPublication := application.CreatePublication[any](v1.ResourceMetadata, false).SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5)
	oi4Asset.RegisterPublication(metaDataAssetPublication) //.SetData().SetPublicationMode(v1.PublicationMode_APPLICATION_SOURCE_5))

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
