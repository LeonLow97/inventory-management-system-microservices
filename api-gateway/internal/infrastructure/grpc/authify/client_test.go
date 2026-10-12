package authify

import (
	"testing"

	"github.com/LeonLow97/internal/infrastructure/config"
)

func TestNewClientCreatesClosableLazyClient(t *testing.T) {
	client, err := NewClient(config.Config{AuthifyService: config.AuthifyServiceConfig{BaseUrl: "127.0.0.1", Port: 50051}})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client.AuthifyService() == nil {
		t.Fatal("AuthifyService() returned nil")
	}
	if err := client.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestCloseWithNilConnection(t *testing.T) {
	client := &client{}
	if err := client.Close(); err != nil {
		t.Errorf("Close() with nil connection error = %v, want nil", err)
	}
}
