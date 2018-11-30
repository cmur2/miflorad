package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "strings"
  "time"

  "github.com/currantlabs/gatt"
  "github.com/currantlabs/gatt/examples/option"
)

var done = make(chan struct{})

func onStateChanged(d gatt.Device, s gatt.State) {
  fmt.Println("State:", s)
  switch s {
  case gatt.StatePoweredOn:
    fmt.Println("Scanning...")
    d.Scan([]gatt.UUID{}, false)
    return
  default:
    d.StopScanning()
  }
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
  id := strings.ToUpper(flag.Args()[0])
  if strings.ToUpper(p.ID()) != id {
    return
  }

  // Stop scanning once we've got the peripheral we're looking for.
  p.Device().StopScanning()

  fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
  fmt.Println("  Local Name        =", a.LocalName)
  fmt.Println("  TX Power Level    =", a.TxPowerLevel)
  fmt.Println("  Manufacturer Data =", a.ManufacturerData)
  fmt.Println("  Service Data      =", a.ServiceData)
  fmt.Println("")

  p.Device().Connect(p)
}

func onPeriphConnected(p gatt.Peripheral, err error) {
  fmt.Println("Connected")
  defer p.Device().CancelConnection(p)

  if err := p.SetMTU(500); err != nil {
    fmt.Printf("Failed to set MTU, err: %s\n", err)
  }

  // Discovery services
  ss, err := p.DiscoverServices(nil)
  if err != nil {
    fmt.Printf("Failed to discover services, err: %s\n", err)
    return
  }

  for _, s := range ss {
    msg := "Service: " + s.UUID().String()
    if len(s.Name()) > 0 {
      msg += " (" + s.Name() + ")"
    }
    fmt.Println(msg)

    if s.UUID().String() != "0000120400001000800000805f9b34fb" {
      fmt.Printf("Skipping uninteresting service\n")
      continue
    }

    // Discovery characteristics
    cs, err := p.DiscoverCharacteristics(nil, s)
    if err != nil {
      fmt.Printf("Failed to discover characteristics, err: %s\n", err)
      continue
    }

    for _, c := range cs {
      msg := "  Characteristic  " + c.UUID().String()
      if len(c.Name()) > 0 {
        msg += " (" + c.Name() + ")"
      }
      fmt.Println(msg)

      if c.UUID().String() == "00001a0000001000800000805f9b34fb" {
        err := p.WriteCharacteristic(c, []byte{0xa0, 0x1f}, false)
        if err != nil {
          fmt.Printf("Failed to write characteristic, err: %s\n", err)
          continue
        }
      }

      if c.UUID().String() == "00001a0200001000800000805f9b34fb" {
        b, err := p.ReadCharacteristic(c)
        if err != nil {
          fmt.Printf("Failed to read characteristic, err: %s\n", err)
          continue
        }
        // fmt.Printf("    value         %x | %q\n", b, b)
        fmt.Printf("    battery level: %d%%\n", b[0])
        fmt.Printf("    firmware version: %s\n", string(b[2:]))
      }

      if c.UUID().String() == "00001a0100001000800000805f9b34fb" {
        b, err := p.ReadCharacteristic(c)
        if err != nil {
          fmt.Printf("Failed to read characteristic, err: %s\n", err)
          continue
        }
        fmt.Printf("    temparature: %f °C\n", float64((int32(b[1]) << 8) + int32(b[0])) / 10.0)
        fmt.Printf("    brightness:  %d lux\n", (int32(b[6]) << 24) + (int32(b[5]) << 16) + (int32(b[4]) << 8) + int32(b[3]))
        fmt.Printf("    moisture:    %d %%\n", int32(b[7]))
        fmt.Printf("    conductivity: %d µS/cm\n", (int32(b[9]) << 8) + int32(b[8]))
      }

      fmt.Printf("Skipping uninteresting char\n")
      continue

      // // Discovery descriptors
      // ds, err := p.DiscoverDescriptors(nil, c)
      // if err != nil {
      //   fmt.Printf("Failed to discover descriptors, err: %s\n", err)
      //   continue
      // }
      //
      // for _, d := range ds {
      //   msg := "  Descriptor      " + d.UUID().String()
      //   if len(d.Name()) > 0 {
      //     msg += " (" + d.Name() + ")"
      //   }
      //   fmt.Println(msg)
      //
      //   // Read descriptor (could fail, if it's not readable)
      //   b, err := p.ReadDescriptor(d)
      //   if err != nil {
      //     fmt.Printf("Failed to read descriptor, err: %s\n", err)
      //     continue
      //   }
      //   fmt.Printf("    value         %x | %q\n", b, b)
      // }
      //
      // // Subscribe the characteristic, if possible.
      // if (c.Properties() & (gatt.CharNotify | gatt.CharIndicate)) != 0 {
      //   f := func(c *gatt.Characteristic, b []byte, err error) {
      //     fmt.Printf("notified: % X | %q\n", b, b)
      //   }
      //   if err := p.SetNotifyValue(c, f); err != nil {
      //     fmt.Printf("Failed to subscribe characteristic, err: %s\n", err)
      //     continue
      //   }
      // }

    }
    fmt.Println()
  }

  fmt.Printf("Waiting for 5 seconds to get some notifiations, if any.\n")
  time.Sleep(5 * time.Second)
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
  fmt.Println("Disconnected")
  close(done)
}

func main() {
  flag.Parse()
  if len(flag.Args()) != 1 {
    log.Fatalf("usage: %s [options] peripheral-id\n", os.Args[0])
  }

  d, err := gatt.NewDevice(option.DefaultClientOptions...)
  if err != nil {
    log.Fatalf("Failed to open device, err: %s\n", err)
    return
  }

  // Register handlers.
  d.Handle(
    gatt.PeripheralDiscovered(onPeriphDiscovered),
    gatt.PeripheralConnected(onPeriphConnected),
    gatt.PeripheralDisconnected(onPeriphDisconnected),
  )

  d.Init(onStateChanged)
  <-done
  fmt.Println("Done")
}
