package fabric

import (
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	networkClient "github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/silvasilas99/entruster/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TODO: 	Improve semanthics of this package. It should be a wrapper around the fabric-gateway, but it is currently
// 			doing more than that. It should not be responsible for registering users, for example.
// 			That should be done in a separate package, and this package should only be responsible for connecting
// 			to the gateway and providing a contract instance.

func Connect() (*networkClient.Contract, *networkClient.Gateway, *grpc.ClientConn) {
	if config.TestNetworkPath == "" {
		panic("chaincodeSdk.Connect: Internal error. TEST_NETWORK_PATH env var is not set")
	}
	conn := newGrpcConnection()
	gw, err := networkClient.Connect(
		newIdentity(),
		networkClient.WithSign(newSign()),
		networkClient.WithHash(hash.SHA256),
		networkClient.WithClientConnection(conn),
		networkClient.WithEvaluateTimeout(5*time.Second),
		networkClient.WithEndorseTimeout(15*time.Second),
		networkClient.WithSubmitTimeout(5*time.Second),
		networkClient.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.Connect: Internal error. Failed to Connect gateway: %w", err))
	}
	network := gw.GetNetwork(config.ChannelName)
	contract := network.GetContractWithName(config.ChaincodeName, config.ContractName)
	return contract, gw, conn
}

func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(config.TLSCertPath)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.newGrpcConnection: Internal error. Failed to read TLS certificate: %w", err))
	}
	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.newGrpcConnection: Internal error. Failed to create certificate: %w", err))
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, config.PeerHostOverride)
	connection, err := grpc.Dial(config.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.newGrpcConnection: Internal error. Failed to create gRPC connection: %w", err))
	}
	return connection
}

func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(config.CertPath)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.NewIdentity: Internal error. Failed to read certificate: %w", err))
	}
	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.NewIdentity: Internal error. Failed to create certificate: %w", err))
	}
	id, err := identity.NewX509Identity(config.MSPID, certificate)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.NewIdentity: Internal error. Failed to create identity: %w", err))
	}
	return id
}

func newSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(config.KeyPath)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.NewSign: Internal error. Failed to read private key: %w", err))
	}
	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.NewSign: Internal error. Failed to create private key: %w", err))
	}
	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.NewSign: Internal error. Failed to create sign: %w", err))
	}
	return sign
}

// open folder, read whatever file is inside
func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.ReadFirstFile: Internal error. Failed to open directory: %w", err))
	}
	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		panic(fmt.Errorf("chaincodeSdk.ReadFirstFile: Internal error. Failed to read directory contents: %w", err))
	}
	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

func RegisterMetadataOnNetwork(
	contract *networkClient.Contract,
	id string,
	name string,
	value string,
	version string,
	owner string,
	rights string,
	termsOfAccess string,
	createdAt string,
	updatedAt string,
	createdBy string,
	updatedBy string,
) error {
	fmt.Printf("--> Submit Transaction: RegisterMetadataOnNetwork | ID: %s\n", id)
	_, err := contract.SubmitTransaction(
		"RegisterMetadataOnNetwork",
		id,
		name,
		value,
		version,
		owner,
		rights,
		termsOfAccess,
		createdAt,
		updatedAt,
		createdBy,
		updatedBy,
	)
	if err != nil {
		return fmt.Errorf("chaincodeSdk.RegisterMetadataOnNetwork: Internal error. Failed to submit transaction: %w", err)
	}
	fmt.Println("*** Transaction committed successfully")
	return nil
}

func GetAllMetadataFromNetwork(
	contract *networkClient.Contract,
) ([]byte, error) {
	fmt.Printf("--> Evaluate Transaction: GetAllMetadataFromNetwork\n")
	result, err := contract.EvaluateTransaction("GetAllMetadataFromNetwork")
	if err != nil {
		return nil, fmt.Errorf("chaincodeSdk.GetAllMetadataFromNetwork: Internal error. Failed to evaluate transaction: %w", err)
	}
	return result, nil
}