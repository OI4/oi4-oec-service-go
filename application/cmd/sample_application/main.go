package main

import (
	"time"

	oi4 "github.com/mzeiher/oi4/api/v1"
	v1 "github.com/mzeiher/oi4/api/v1"
	"github.com/mzeiher/oi4/application"
	"github.com/mzeiher/oi4/application/pkg/mqtt"
)

func main() {

	oi4Application := application.CreateNewApplication(oi4.ServiceType_OTConnector, &oi4.MasterAssetModel{
		Manufacturer: oi4.LocalizedText{
			"en_us": "Bitschubser",
		},
		ManufacturerUri:    "bitschubser.de",
		Model:              "SampleApplication",
		ProductCode:        "08",
		HardwareRevision:   "0",
		SoftwareRevision:   "0",
		DeviceRevision:     "0",
		DeviceManual:       "",
		DeviceClass:        "PerpetuumMobile",
		SerialNumber:       "15",
		ProductInstanceUri: "aliexpress.com",
		RevisionCounter:    0,
		Description: oi4.LocalizedText{
			"en_us": "Cooles Teil",
		},
	})

	dataPublication := application.CreatePublication(v1.Resource_Data, nil).SetPublicationMode(v1.PublicationMode_APPLICATION_2)
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		counter := 0
		for {
			<-ticker.C
			dataPublication.SetData(counter)
			counter++
		}

	}()

	oi4Application.RegisterPublication(dataPublication)

	if err := oi4Application.Start(&mqtt.MQTTClientOptions{
		Host:     "localhost",
		Port:     1883,
		Tls:      false,
		Username: "test",
		Password: "test",
	}); err != nil {
		panic(err)
	}

	waiter := make(chan struct{})
	<-waiter
	// client.PublishResource(&oi4.Oi4Identifier{
	// 	ManufacturerUri: "bitschubser.dev",
	// 	Model:           "weather-test",
	// 	ProductCode:     "08",
	// 	SerialNumber:    "15",
	// }, oi4.ServiceType_Utility, oi4.Resource_Data, nil, &oi4.Oi4Data{
	// 	PrimaryValue: 9000,
	// })
}
