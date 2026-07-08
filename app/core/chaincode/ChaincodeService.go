package chaincode

import (
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/silvasilas99/entruster/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ChaincodeService struct {
	Contract *client.Contract
	Gateway  *client.Gateway
	Connection *grpc.ClientConn
}

func NewChaincodeService() *ChaincodeService {
	contract, gateway, connection := ConnectOnFabric();

	return &ChaincodeService{
		Contract: contract,
		Gateway:  gateway,
		Connection: connection,
	};
}

func ConnectOnFabric() (*client.Contract, *client.Gateway, *grpc.ClientConn) {
	if config.TestNetworkPath == "" {
		panic("[ChaincodeService][ConnectOnFabric]: Internal error. TEST_NETWORK_PATH env var is not set");
	}

	grpcConnection := connectOnGrpc();
	fabricGateway, err := client.Connect(
		generateIdentity(),
		client.WithSign(generateSign()),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(grpcConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	);

	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][ConnectOnFabric]: Internal error. Failed to connect on gateway: %w", err));
	}

	network := fabricGateway.GetNetwork(config.ChannelName);
	contract := network.GetContractWithName(config.ChaincodeName, config.ContractName);

	return contract, fabricGateway, grpcConnection;
}

func connectOnGrpc() *grpc.ClientConn {
	// Build a cert pool with BOTH the peer TLS CA and the orderer TLS CA.
	//
	// The fabric-gateway SubmitTransaction flow is:
	//   1. Endorse  -> peer   (uses this gRPC connection)
	//   2. Submit   -> orderer (gateway discovers orderer address from channel config
	//                           and opens a NEW TLS connection to it)
	//
	// The orderer is issued by a DIFFERENT CA (ca.example.com) than the peer
	// (ca.org1.example.com). If the orderer's CA is missing from the cert pool,
	// Go's TLS stack raises:
	//   "x509: certificate signed by unknown authority (possibly because of
	//    x509: ECDSA verification failure while trying to verify candidate
	//    authority certificate ca.org1.example.com)"
	//
	// Fix: add both CAs so every participant's TLS cert can be verified.
	certPool := x509.NewCertPool();

	// Add peer TLS CA
	peerCertPEM, err := os.ReadFile(config.TLSCertPath);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][connectOnGrpc]: Internal error. Failed to read peer TLS certificate: %w", err));
	}
	peerCert, err := identity.CertificateFromPEM(peerCertPEM);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][connectOnGrpc]: Internal error. Failed to parse peer TLS certificate: %w", err));
	}
	certPool.AddCert(peerCert);

	// Add orderer TLS CA — required for SubmitTransaction to reach the orderer
	ordererCertPEM, err := os.ReadFile(config.OrdererTLSCertPath);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][connectOnGrpc]: Internal error. Failed to read orderer TLS certificate: %w", err));
	}
	ordererCert, err := identity.CertificateFromPEM(ordererCertPEM);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][connectOnGrpc]: Internal error. Failed to parse orderer TLS certificate: %w", err));
	}
	certPool.AddCert(ordererCert);

	transportCredentials := credentials.NewClientTLSFromCert(certPool, config.PeerHostOverride);
	connection, err := grpc.Dial(config.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials));
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][connectOnGrpc]: Internal error. Failed to create gRPC connection: %w", err));
	}

	return connection;
}

func generateIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(config.CertPath);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][generateIdentity]: Internal error. Failed to read certificate: %w", err));
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][generateIdentity]: Internal error. Failed to create certificate: %w", err));
	}

	generatedIdentity, err := identity.NewX509Identity(config.MSPID, certificate);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][generateIdentity]: Internal error. Failed to create identity: %w", err));
	}

	return generatedIdentity;
}

func generateSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(config.KeyPath);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][generateSign]: Internal error. Failed to read private key: %w", err));
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][generateSign]: Internal error. Failed to create private key: %w", err));
	}

	generatedSign, err := identity.NewPrivateKeySign(privateKey);
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][generateSign]: Internal error. Failed to create sign: %w", err));
	}
	return generatedSign;
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][readFirstFile]: Internal error. Failed to open directory: %w", err))
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		panic(fmt.Errorf("[ChaincodeService][readFirstFile]: Internal error. Failed to read directory contents: %w", err))
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}
