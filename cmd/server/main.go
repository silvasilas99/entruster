// Package main is the entry point of the Entruster API server.
//
// It connects to the Hyperledger Fabric network and starts the HTTP server.
// Business logic is kept in separate packages; this file only bootstraps the
// gateway connection and wires up the HTTP router.
package main

//	@title			Entruster API
//	@version		1.0
//	@description	REST API for registering and querying healthcare metadata on a Hyperledger Fabric blockchain.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Silas Silva
//	@contact.url	https://github.com/silvasilas99/entruster

//	@license.name	MIT

//	@host		localhost:8080
//	@BasePath	/api

//	@schemes	http

import (
	"fmt"

	_ "github.com/silvasilas99/entruster/docs" // generated Swagger docs
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
	fmt.Println("📄 Swagger UI at http://localhost:" + config.Port + "/swagger/index.html")
	r.Run(":" + config.Port)
}