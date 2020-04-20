// make sure to execute `go install` before tests
package main

import (
	"testing"

	"github.com/Divya063/pingApp/cmd"
)

func TestIPV4(t *testing.T) {
	pinger, initErr := cmd.NewPinger("8.8.8.8")
	_, _, _, err := cmd.SendAndReceiveRequests(pinger)
	if initErr != nil {
		t.Errorf("Expected to succeed, but failed: %v", initErr)
	}

	if err != nil {
		t.Errorf("%v", err)
	}

}

func TestIPV6(t *testing.T) {
	pinger, initErr := cmd.NewPinger("0:0:0:0:0:ffff:7f00:1")
	_, _, _, err := cmd.SendAndReceiveRequests(pinger)
	if initErr != nil {
		t.Errorf("Expected to succeed, but failed: %v", initErr)
	}

	if err != nil {
		t.Errorf("%v", err)
	}

}

func TestAddr(t *testing.T) {
	pinger, initErr := cmd.NewPinger("google.com")
	_, _, _, err := cmd.SendAndReceiveRequests(pinger)
	if initErr != nil {
		t.Errorf("Expected to succeed, but failed: %v", initErr)
	}

	if err != nil {
		t.Errorf("%v", err)
	}

}
