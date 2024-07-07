package main

import (
    "context"
    "log"
    "os"

    "golang.org/x/sync/errgroup"
    "github.com/xtls/xray-core/core"
    "github.com/xtls/xray-core/core/common/log"
    "github.com/xtls/xray-core/infra/conf/serial"
)

func main() {
    configStr := `
{
  "log": {
    "loglevel": "none"
  },
  "inbounds": [
    {
      "port": 8080,
      "protocol": "dokodemo-door",
      "settings": {
        "address": "amazonfront.egy.pp.ua",
        "port": 8080,
        "network": "tcp,udp"
      }
    }
  ],
  "outbounds": [
    {
      "protocol": "freedom",
      "settings": {}
    }
  ]
}`

    config, err := serial.LoadJSONConfig(configStr)
    if err != nil {
        log.Fatalf("Failed to load Xray config: %v", err)
    }

    instance, err := core.New(config)
    if err != nil {
        log.Fatalf("Failed to create Xray instance: %v", err)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    g, ctx := errgroup.WithContext(ctx)
    g.Go(func() error {
        return instance.Start()
    })

    go func() {
        <-ctx.Done()
        instance.Close()
    }()

    if err := g.Wait(); err != nil {
        log.Fatalf("Xray instance stopped with error: %v", err)
    } else {
        log.Println("Xray instance stopped gracefully")
    }
}
