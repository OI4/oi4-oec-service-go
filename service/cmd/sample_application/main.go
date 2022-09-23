package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/mzeiher/oi4/api/pkg/types"
	application "github.com/mzeiher/oi4/service"
	"github.com/mzeiher/oi4/service/pkg/mqtt"
)

func main() {

	oi4Application := application.CreateNewApplication(v1.ServiceType_OTConnector, &v1.MasterAssetModel{
		Manufacturer: v1.LocalizedText{
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
		Description: v1.LocalizedText{
			"en_us": "Cooles Teil",
		},
	})

	dataPublication := application.CreatePublication[v1.Oi4Data](v1.Resource_Data, false).SetPublicationMode(v1.PublicationMode_APPLICATION_2)
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		counter := 0
		for {
			<-ticker.C
			dataPublication.SetData(v1.Oi4Data{PrimaryValue: counter})
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	oi4Application.Stop()

	os.Exit(0)
}
