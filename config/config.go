package config

import (
    "os"
    "path"
)
const (
    MSPID            = "Org1MSP"
    ChannelName      = "metadatachannel"
    ChaincodeName    = "basic"
    ContractName     = "MetadataContract"
    PeerEndpoint     = "dns:///localhost:7051"
    PeerHostOverride = "peer0.org1.example.com"
    Port             = "8080"
)
var TestNetworkPath = os.Getenv("TEST_NETWORK_PATH")
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