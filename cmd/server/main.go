package main

// This file is the entry point of the application. It is responsible for connecting to the Fabric network and starting the HTTP server.
// The main function should be kept as simple as possible, and should not contain any business logic.
// That should be done in a separate package, and this package should only be responsible for connecting
// to the gateway and providing a contract instance. It should not be responsible for registering users, for example.

import (
    "fmt"
    "github.com/silvasilas99/entruster/config"
    "github.com/silvasilas99/entruster/fabric"
    "github.com/silvasilas99/entruster/routes"
)
func main() {
    contract, gw, conn := fabric.Connect()

    defer gw.Close()
    defer conn.Close()

    fmt.Println("✅ Connected to Fabric —", config.ChannelName)
    r := routes.SetupRoutes(contract)

    fmt.Println("🚀 Server running on http://localhost:" + config.Port)
    r.Run(":" + config.Port)
}