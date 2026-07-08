package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/silvasilas99/entruster/app/core/audit"
	"github.com/silvasilas99/entruster/app/core/chaincode"
	"github.com/silvasilas99/entruster/app/core/elasticsearch"
	"github.com/silvasilas99/entruster/app/domain/metadata"
)

// Payload reflete os parâmetros que o smart contract espera em RegisterMetadataOnNetwork.
type Payload struct {
	PatientID     uint64
	AssetID       uint64
	ZKPProof      string
	Name          string
	Value         string
	Version       string
	Owner         string
	Rights        string
	TermsOfAccess string
	CreatedAt     string
	UpdatedAt     string
	CreatedBy     string
	UpdatedBy     string
}

// generateFHIRPatient cria um mock de um recurso FHIR Patient.
func generateFHIRPatient() string {
	data := map[string]interface{}{
		"resourceType": "Patient",
		"id":           fmt.Sprintf("pat-%d", rand.Intn(10000)),
		"active":       true,
		"name": []map[string]interface{}{
			{
				"use":    "official",
				"family": "Silva",
				"given":  []string{"João", "Carlos"},
			},
		},
		"gender":    "male",
		"birthDate": "1980-01-01",
	}
	b, _ := json.Marshal(data)
	return string(b)
}

// generateFHIRObservation cria um mock de um recurso FHIR Observation.
func generateFHIRObservation() string {
	data := map[string]interface{}{
		"resourceType": "Observation",
		"id":           fmt.Sprintf("obs-%d", rand.Intn(10000)),
		"status":       "final",
		"code": map[string]interface{}{
			"text": "Body Weight",
		},
		"subject": map[string]interface{}{
			"reference": fmt.Sprintf("Patient/pat-%d", rand.Intn(10000)),
		},
		"valueQuantity": map[string]interface{}{
			"value": rand.Float64()*50 + 50,
			"unit":  "kg",
		},
	}
	b, _ := json.Marshal(data)
	return string(b)
}

// generateFHIRDiagnosticReport cria um mock de um recurso FHIR DiagnosticReport.
func generateFHIRDiagnosticReport() string {
	data := map[string]interface{}{
		"resourceType": "DiagnosticReport",
		"id":           fmt.Sprintf("diag-%d", rand.Intn(10000)),
		"status":       "final",
		"code": map[string]interface{}{
			"text": "Complete Blood Count",
		},
		"subject": map[string]interface{}{
			"reference": fmt.Sprintf("Patient/pat-%d", rand.Intn(10000)),
		},
		"conclusion": "All results are within normal parameters.",
	}
	b, _ := json.Marshal(data)
	return string(b)
}

// generatePayloads gera os dados simulados combinando os perfis HL7 FHIR.
func generatePayloads(count int) []Payload {
	payloads := make([]Payload, count)
	now := time.Now().UTC().Format(time.RFC3339)
	for i := 0; i < count; i++ {
		t := rand.Intn(3)
		var fhirData, name string
		switch t {
		case 0:
			fhirData = generateFHIRPatient()
			name = "FHIR_Patient"
		case 1:
			fhirData = generateFHIRObservation()
			name = "FHIR_Observation"
		case 2:
			fhirData = generateFHIRDiagnosticReport()
			name = "FHIR_DiagnosticReport"
		}

		// Simulando a hash ZKP da informação
		hash := sha256.Sum256([]byte(fhirData))
		zkp := fmt.Sprintf("%x", hash)

		payloads[i] = Payload{
			PatientID:     uint64(rand.Intn(1000)),
			AssetID:       uint64(rand.Intn(1000000)),
			ZKPProof:      zkp,
			Name:          name,
			Value:         fhirData,
			Version:       "1.0",
			Owner:         "Hospital_Interop",
			Rights:        "Read",
			TermsOfAccess: "Consent Required",
			CreatedAt:     now,
			UpdatedAt:     now,
			CreatedBy:     "Seeder_Caliper_Sim",
			UpdatedBy:     "Seeder_Caliper_Sim",
		}
	}
	return payloads
}

// submitTx envia a transação para a rede Fabric e realiza o indexing no Elasticsearch.
func submitTx(svc *metadata.MetadataService, p Payload) error {
	dto := metadata.NewMetadataDTO(
		0,
		p.PatientID,
		p.AssetID,
		p.ZKPProof,
		p.Name,
		p.Value,
		p.Version,
		p.Owner,
		p.Rights,
		p.TermsOfAccess,
		p.CreatedAt,
		p.CreatedBy,
		p.UpdatedAt,
		p.UpdatedBy,
		"",
		"",
	)
	return svc.RegisterMetadata(dto)
}

// runSequential submete as transações de maneira síncrona/sequencial.
func runSequential(svc *metadata.MetadataService, payloads []Payload) time.Duration {
	start := time.Now()
	for i, p := range payloads {
		err := submitTx(svc, p)
		if err != nil {
			log.Printf("Erro Sequencial na idx %d: %v", i, err)
		}
	}
	return time.Since(start)
}

// runConcurrent submete as transações concorrentemente através de um Pool de Workers.
func runConcurrent(svc *metadata.MetadataService, payloads []Payload, workers int) time.Duration {
	start := time.Now()
	var wg sync.WaitGroup
	jobs := make(chan Payload, len(payloads))

	for w := 1; w <= workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range jobs {
				err := submitTx(svc, p)
				if err != nil {
					log.Printf("Erro Concorrente: %v", err)
				}
			}
		}()
	}

	for _, p := range payloads {
		jobs <- p
	}
	close(jobs)
	wg.Wait()
	return time.Since(start)
}

func main() {
	count := flag.Int("count", 20, "Número de registros a serem semeados por fase (sequencial / concorrente)")
	workers := flag.Int("workers", 10, "Número de workers paralelos para a fase concorrente")
	flag.Parse()

	log.Println("Iniciando Seeding de Metadados (Caliper-like Benchmark)...")
	log.Println("Conectando na rede Hyperledger Fabric...")
	
	// Utiliza o cliente gateway já implementado no projeto
	contract, gw, conn := chaincode.ConnectOnFabric()
	defer gw.Close()
	defer conn.Close()
	log.Println("✅ Conectado com sucesso no Fabric.")

	elasticSvc := elasticsearch.NewElasticService()
	chaincodeQuery := chaincode.NewChaincodeQuery(contract)
	auditSvc := audit.NewAuditService(elasticSvc, chaincodeQuery)
	metadataObserver := metadata.NewMetadataObserver(auditSvc, elasticSvc)
	metadataSvc := metadata.NewMetadataService(contract, metadataObserver)

	// 1. Gera os payloads HL7 FHIR
	log.Printf("Gerando %d registros no total...", *count*2)
	seqPayloads := generatePayloads(*count)
	concPayloads := generatePayloads(*count)

	// 2. Roda Benchmark Sequencial
	log.Printf("\n--- Executando Seeding Sequencial (%d registros) ---", *count)
	seqDuration := runSequential(metadataSvc, seqPayloads)
	seqTPS := float64(*count) / seqDuration.Seconds()
	log.Printf("Tempo Sequencial : %v", seqDuration)
	log.Printf("TPS Sequencial   : %.2f tx/sec", seqTPS)

	// 3. Roda Benchmark Concorrente
	log.Printf("\n--- Executando Seeding Concorrente (%d registros, %d workers) ---", *count, *workers)
	concDuration := runConcurrent(metadataSvc, concPayloads, *workers)
	concTPS := float64(*count) / concDuration.Seconds()
	log.Printf("Tempo Concorrente: %v", concDuration)
	log.Printf("TPS Concorrente  : %.2f tx/sec", concTPS)

	// 4. Calcula e Metrifica o Speedup e Resultados
	speedup := float64(seqDuration) / float64(concDuration)
	log.Printf("\n================ RESULTADOS DO BENCHMARK ================")
	log.Printf("Speedup: %.2fx (A execução concorrente foi %.2fx mais rápida)", speedup, speedup)
	
	fmt.Printf("\n📊 Resumo estilo Hyperledger Caliper:\n")
	fmt.Printf("| %-15s | %-10s | %-15s |\n", "Fase", "Workers", "Throughput (TPS)")
	fmt.Printf("| %-15s | %-10d | %-15.2f |\n", "Sequencial", 1, seqTPS)
	fmt.Printf("| %-15s | %-10d | %-15.2f |\n", "Concorrente", *workers, concTPS)
	fmt.Println("=========================================================")
}