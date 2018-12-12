//shows how to watch for new devices and list them
package legacy

import (
	"os"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/emitter"
	log "github.com/sirupsen/logrus"
)

const adapterID = "hci0"

// func main() {
// 	manager, err := api.NewManager()
// 	if err != nil {
// 		log.Error(err)
// 		os.Exit(1)
// 	}
//
// 	err = manager.RefreshState()
// 	if err != nil {
// 		log.Error(err)
// 		os.Exit(1)
// 	}
//
//   var address = "C4:7C:8D:66:D5:27"
//
//   dev, err := api.GetDeviceByAddress(address)
//   if err != nil {
//     log.Error(err)
//     os.Exit(1)
//   }
//
//   if dev == nil {
//     log.Infof("No device found!")
//     os.Exit(1)
//   }
//
//   dev.GetAllServicesAndUUID()
//
//   props, err := dev.GetProperties()
//   if err != nil {
//     log.Error(err)
//     os.Exit(1)
//   }
//
//   log.Infof("name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)
//
//   err = dev.Connect()
//   if err != nil {
//     log.Error(err)
//     os.Exit(1)
//   }
//
//   if !dev.IsConnected() {
//     log.Infof("Device not connected!")
//     os.Exit(1)
//   }
//
//   y, err := dev.GetCharByUUID("00002a00-0000-1000-8000-00805f9b34fb")
//   if err != nil {
//     log.Error(err)
//     os.Exit(1)
//   }
//
//   log.Infof(y.Path)
//
//   err = dev.Disconnect()
//   if err != nil {
//     log.Error(err)
//     os.Exit(1)
//   }
// }

// func main() {
//   log.SetLevel(log.DebugLevel)
//
//   // clean up connection on exit
//   defer api.Exit()
//
//   manager, err := api.NewManager()
//   if err != nil {
//     log.Error(err)
//     os.Exit(1)
//   }
//
//   err = manager.RefreshState()
//   if err != nil {
//     log.Error(err)
//     os.Exit(1)
//   }
//
//   boo, err := api.AdapterExists(adapterID)
//   if err != nil {
//     log.Error(err)
//     os.Exit(1)
//   }
//   log.Debugf("AdapterExists: %b", boo)
//
//   err = api.StartDiscoveryOn(adapterID)
//   if err != nil {
//     log.Error(err)
//     os.Exit(1)
//   }
//   log.Debugf("Started discovery")
//
//   err = api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {
//     discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
//     dev := discoveryEvent.Device
//     handleDevice(dev)
//   }))
//
//   select {}
// }
//
// func handleDevice(dev *api.Device) {
//   if dev == nil {
//     return
//   }
//
//   props, err := dev.GetProperties()
//   if err != nil {
//     log.Error(err)
//     return
//   }
//
//   log.Infof("name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)
// }
