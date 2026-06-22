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
    PeerEndpoint = getEnvOrDefault("PEER_ENDPOINT", "dns:///localhost:7051")
    Port         = getEnvOrDefault("PORT", "8080")
    TestNetworkPath = os.Getenv("TEST_NETWORK_PATH")
)

var (
    TLSCertPath = path.Join(TestNetworkPath,
        "organizations/peerOrganizations/org1.example.com",
        "peers/peer0.org1.example.com/tls/ca.crt")
    CertPath = path.Join(TestNetworkPath,
        "organizations/peerOrganizations/org1.example.com",
        "users/Admin@org1.example.com/msp/signcerts")
    KeyPath = path.Join(TestNetworkPath,
        "organizations/peerOrganizations/org1.example.com",
        "users/Admin@org1.example.com/msp/keystore")
)