//go:build linux && android
// +build linux,android

package anywherelan

import (
	"context"
	"os"
	"strconv"

	"github.com/anywherelan/awl"
	"github.com/anywherelan/awl/config"
	"github.com/anywherelan/awl/vpn"
	"github.com/ipfs/go-log/v2"
)

var (
	app           *awl.Application
	logger        *log.ZapEventLogger
	globalDataDir string
)

// All public functions are part of a library

func InitServer(dataDir string, tunFD int32) error {
	globalDataDir = dataDir
	_ = os.Setenv(config.AppDataDirEnvKey, dataDir)
	_ = os.Setenv(vpn.TunFDEnvKey, strconv.Itoa(int(tunFD)))

	app = awl.New()
	logger = app.SetupLoggerAndConfig()
	//ctx, ctxCancel := context.WithCancel(context.Background())
	ctx := context.Background()

	err := app.Init(ctx, nil)
	if err != nil {
		app.Close()
		app = nil
	}
	return err
}

func StopServer() {
	if app != nil {
		app.Close()
		app = nil
	}
}

func ImportConfig(data string) error {
	if app != nil || globalDataDir == "" {
		panic("call to ImportConfig before server shutdown")
	}

	return config.ImportConfig([]byte(data), globalDataDir)
}

func GetApiAddress() string {
	if app != nil && app.Api != nil {
		return app.Api.Address()
	}
	return ""
}
