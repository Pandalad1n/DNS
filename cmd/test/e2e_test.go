package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

const webAddr = "localhost"
const webPort = 8090

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	server := exec.CommandContext(
		ctx,
		"go", "run", "../dns/.",
	)
	server.Env = append(os.Environ(), fmt.Sprintf("LISTEN_WEB_ADDRESS=:%d", webPort), "SECTOR_ID=1", "LOG_LVL=debug")
	err := server.Start()
	if err != nil {
		os.Exit(1)
	}
	var tries int
	for {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(fmt.Sprintf("http://%s:%d/health", webAddr, webPort))
		if err == nil {
			_ = resp.Body.Close()
			break
		}
		tries++
		if tries > 100 {
			os.Exit(1)
		}
	}
	code := m.Run()
	cancel()
	os.Exit(code)
}

func TestEmpty(t *testing.T) {

}

func TestHealth(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/health", webAddr, webPort))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLocation(t *testing.T) {
	droneData := []byte(`{
		"x": "123.12",
		"y": "456.56",
		"z": "789.89",
		"vel": "20.0"
	}`)
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/v1/locate", webAddr, webPort), "application/json", bytes.NewBuffer(droneData))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	expectedResp := map[string]interface{}{"loc": 1389.57}
	var actualResp interface{}
	err = json.NewDecoder(resp.Body).Decode(&actualResp)
	require.NoError(t, err)
	assert.EqualValues(t, expectedResp, actualResp)
}

func TestMetrics(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/metrics", webAddr, webPort))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
