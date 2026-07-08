package config

import (
    "os"
    "path"
)

func getEnvOrDefault(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

const (
    MSPID            = "Org1MSP"
    ChannelName      = "metadatachannel"
    ChaincodeName    = "basic"
    ContractName     = "MetadataContract"
    PeerHostOverride = "peer0.org1.example.com"
)

var (
	PeerEndpoint    = getEnvOrDefault("PEER_ENDPOINT", "dns:///localhost:7051")
	Port            = getEnvOrDefault("PORT", "8080")
	TestNetworkPath = os.Getenv("TEST_NETWORK_PATH")

	// JWT and User Mock API settings
	StaticToken    = getEnvOrDefault("STATIC_TOKEN", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c3JfNzc3IiwibmFtZSI6IkRyLiBTaWxhcyBTaWx2YSIsInJvbGUiOiJDaGllZiBDbGluaWNpYW4iLCJpYXQiOjE3ODIyOTc2MDB9.4M2R0_2nN0X4Y0tX6W7N8wYwS8uV2Fz3T9k_Y3Z1U_w") // Token estático pré-assinado
	JWTSecret      = getEnvOrDefault("JWT_SECRET", "super-secret-key-change-it-in-production")
	UserMockApiUrl = getEnvOrDefault("USER_MOCK_API_URL", "http://localhost:8080/api/mock/user")
)

var (
    TLSCertPath = path.Join(TestNetworkPath,
        "organizations/peerOrganizations/org1.example.com",
        "peers/peer0.org1.example.com/tls/ca.crt")
    // OrdererTLSCertPath is the TLS CA certificate of the orderer organization.
    // It must be added to the gRPC TLS cert pool so the fabric-gateway can
    // verify the orderer's certificate when submitting transactions.
    OrdererTLSCertPath = path.Join(TestNetworkPath,
        "organizations/ordererOrganizations/example.com",
        "orderers/orderer.example.com/tls/ca.crt")
    CertPath = path.Join(TestNetworkPath,
        "organizations/peerOrganizations/org1.example.com",
        "users/Admin@org1.example.com/msp/signcerts")
    KeyPath = path.Join(TestNetworkPath,
        "organizations/peerOrganizations/org1.example.com",
        "users/Admin@org1.example.com/msp/keystore")
)