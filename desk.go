package main

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"tinygo.org/x/bluetooth"
)

const (
	MIN_DESK_HEIGHT = 65  // cm
	MAX_DESK_HEIGHT = 130 // cm

	// Values sent by 'STANDUP! pls' appliaction
	UP_COMMAND   int = 0xd9ff01633c
	DOWN_COMMAND int = 0xd9ff02603a

	UP_DIRECTION   = "up"
	DOWN_DIRECTION = "down"
)

type desk struct {
	// UUID address of the bluetooth device
	address string
	// Current height and speed of the desk
	height uint8
	speed  uint8
	// BLE device characteristics used for communication
	rxChar, txChar *bluetooth.DeviceCharacteristic

	mux *sync.RWMutex
}

func NewDesk(address string) *desk {
	return &desk{
		address: address,
		mux:     &sync.RWMutex{},
	}
}

func (d *desk) GetHeight() uint8 {
	d.mux.RLock()
	height := d.height
	d.mux.RUnlock()

	return height
}

func (d *desk) SetHeight(height uint8) {
	d.mux.Lock()
	d.height = height
	d.mux.Unlock()
}

func (d *desk) GetSpeed() uint8 {
	d.mux.RLock()
	speed := d.speed
	d.mux.RUnlock()

	return speed
}

func (d *desk) SetSpeed(speed uint8) {
	d.mux.Lock()
	d.speed = speed
	d.mux.Unlock()
}

func (d *desk) Connect() error {
	// Enable BLE interface.
	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return fmt.Errorf("could not enable BLE stack: %w", err)
	}

	scanResults := make(chan bluetooth.ScanResult, 1)

	// Start scanning.
	slog.Info("Scanning BLE devices")
	err := adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		slog.Debug("found device", "address", result.Address.String(), "rssi", result.RSSI, "local_name", result.LocalName())
		if result.Address.String() == d.address {
			adapter.StopScan()
			scanResults <- result
		}
	})
	if err != nil {
		return fmt.Errorf("could not initiate bluetooth scan: %w", err)
	}

	var device bluetooth.Device
	select {
	case result := <-scanResults:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		if err != nil {
			return fmt.Errorf("could not connect to the desk: %w", err)
		}
		slog.Info("Connected to the device", "address", result.Address.String(), "local_name", result.LocalName())
	}

	// Get UART service for the device
	slog.Debug("discovering services/characteristics")
	services, err := device.DiscoverServices([]bluetooth.UUID{bluetooth.ServiceUUIDNordicUART})
	if err != nil {
		return fmt.Errorf("discovering services: %w", err)
	}
	if len(services) == 0 {
		return fmt.Errorf("could not find UART service: %w", err)
	}

	service := services[0]
	slog.Debug("found service", "uuid", service.UUID().String())

	chars, err := service.DiscoverCharacteristics(
		[]bluetooth.UUID{
			// https://docs.nordicsemi.com/bundle/ncs-latest/page/nrf/libraries/bluetooth_services/services/nus.html
			bluetooth.CharacteristicUUIDUARTRX, // Write
			bluetooth.CharacteristicUUIDUARTTX, // Notifications
		},
	)
	if err != nil {
		return fmt.Errorf("could not get device characteristics: %w", err)
	}

	d.rxChar, d.txChar = &(chars[0]), &(chars[1])
	slog.Debug("found RX characteristic", "uuid", d.rxChar.UUID().String())
	slog.Debug("found TX characteristic", "uuid", d.txChar.UUID().String())

	// Set up updating desk height
	d.txChar.EnableNotifications(func(buf []byte) {
		// println("data:", buf)
		slog.Debug("Notification Received", "buffer", buf)
		if len(buf) > 3 { // Ignore shorter messages
			d.SetHeight(uint8(buf[3]))
			d.SetSpeed(uint8(buf[1]))
			slog.Debug("Parameters updated", "height", d.GetHeight(), "speed", d.GetSpeed())
		}
	})

	// Initalise desk state by sending stop command to trigger notifications
	d.rxChar.Write([]byte{0x0, 0x0, 0x0, 0x0, 0x0})
	time.Sleep(time.Second)
	slog.Info("Device successfully initialised", "height", d.GetHeight())

	return nil
}

func (d *desk) reachedTargetHeight(direction string, targetHeight uint8) bool {
	if direction == UP_DIRECTION {
		return targetHeight-5 <= d.GetHeight()
	} else if direction == DOWN_DIRECTION {
		return targetHeight+5 >= d.GetHeight()
	}

	panic("Unknown direction!")
}

func (d *desk) MoveTo(targetHeight uint8) {
	if d.GetHeight() == targetHeight {
		return
	}

	direction, command := "up", UP_COMMAND
	if targetHeight < d.GetHeight() {
		direction, command = "down", DOWN_COMMAND
	}
	slog.Info("Staring to move", "direction", direction, "height", d.GetHeight(), "target_height", targetHeight, "speed", d.GetSpeed())

	// encode command
	data := make([]byte, 5)
	for i := 4; i >= 0; i-- {
		data[4-i] = byte((command >> (8 * i)) & 0xFF)
	}

	for !d.reachedTargetHeight(direction, targetHeight) {
		time.Sleep(200 * time.Millisecond)
		_, err := d.rxChar.Write(data)
		if err != nil {
			fmt.Println("Cannot write data: ", err)
			os.Exit(1)
		}
		slog.Debug("Sent command", "direction", direction, "target_height", targetHeight, "current_height", d.GetHeight())
	}
	slog.Info("Finished moving", "direction", direction, "height", d.GetHeight(), "target_height", targetHeight, "speed", d.GetSpeed())
	d.rxChar.Write([]byte{0x0, 0x0, 0x0, 0x0, 0x0})
}
