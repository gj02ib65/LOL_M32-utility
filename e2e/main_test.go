package main

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

var (
	targetIP   = os.Getenv("TARGET_IP")
	targetPort = 10023
)

func init() {
	if targetIP == "" {
		targetIP = "localhost"
	}
}

// startTestListener starts a UDP listener to act as the "Replica" console.
// It returns a channel that will receive incoming OSC messages.
func startTestListener(port int) (chan *osc.Message, func(), error) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, nil, err
	}

	msgChan := make(chan *osc.Message, 10)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				// Listener closed
				close(msgChan)
				return
			}
			
			packet, err := osc.ParsePacket(string(buf[:n]))
			if err != nil {
				continue
			}
			
			switch pkg := packet.(type) {
			case *osc.Message:
				if pkg.Address != "/xremote" {
					msgChan <- pkg
				}
			case *osc.Bundle:
				for _, p := range pkg.Messages {
					if p.Address != "/xremote" {
						msgChan <- p
					}
				}
			}
		}
	}()

	cleanup := func() {
		conn.Close()
	}

	return msgChan, cleanup, err
}

func sendOSCMessage(client *osc.Client, address string, value interface{}) error {
	msg := osc.NewMessage(address)
	msg.Append(value)
	
	// Start with default
	return client.Send(msg)
}

func sendSpoofedOSCMessage(targetAddr string, address string, value interface{}, sourcePort int) error {
	addr, err := net.ResolveUDPAddr("udp", targetAddr)
	if err != nil {
		return err
	}
	
	laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", sourcePort))
	if err != nil {
		return err
	}
	
	conn, err := net.DialUDP("udp", laddr, addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	
	msg := osc.NewMessage(address)
	msg.Append(value)
	
	p, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = conn.Write(p)
	return err
}

func TestMuteSync(t *testing.T) {
	msgChan, cleanup, err := startTestListener(10023)
	if err != nil {
		t.Fatalf("Failed to start listener: %v", err)
	}
	defer cleanup()

	time.Sleep(2 * time.Second)
	
	waitForMessage := func(expectedAddress string, expectedValue int32) {
		timeout := time.After(2 * time.Second)
		for {
			select {
			case msg := <-msgChan:
				fmt.Printf("Received message: %s with args: %v\n", msg.Address, msg.Arguments)
				if msg.Address == expectedAddress {
					if len(msg.Arguments) > 0 {
						if val, ok := msg.Arguments[0].(int32); ok && val == expectedValue {
							return // Success!
						}
					}
				}
			case <-timeout:
				t.Fatalf("Timeout waiting for %s with value %d", expectedAddress, expectedValue)
			}
		}
	}

	t.Run("Standard Channel Mute Sync", func(t *testing.T) {
		err := sendSpoofedOSCMessage(fmt.Sprintf("%s:%d", targetIP, targetPort), "/ch/05/mix/on", int32(0), 10024)
		if err != nil {
			t.Fatalf("Failed to send OSC: %v", err)
		}
		waitForMessage("/ch/05/mix/on", int32(0))
	})
	
	t.Run("Standard Channel Unmute Sync", func(t *testing.T) {
		err := sendSpoofedOSCMessage(fmt.Sprintf("%s:%d", targetIP, targetPort), "/ch/05/mix/on", int32(1), 10025)
		if err != nil {
			t.Fatalf("Failed to send OSC: %v", err)
		}
		waitForMessage("/ch/05/mix/on", int32(1))
	})
}
