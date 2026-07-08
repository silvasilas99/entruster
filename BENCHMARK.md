# Benchmark de Seeding

Este documento contém os resultados do benchmark de seeding na rede blockchain Hyperledger Fabric utilizando metadados HL7 FHIR.
A avaliação consistiu na submissão de transações em modo **Sequencial** (síncrono) e **Concorrente** (assíncrono utilizando goroutines/workers), para fins de verificação do ganho de velocidade (speedup) e do Throughput em Transações por Segundo (TPS).

O ambiente e as simulações baseiam-se em uma lógica inspirada no Hyperledger Caliper, calculando as métricas e tempos de resposta.

## Resultados da Execução

A execução foi feita submetendo 5 registros de forma sequencial e 5 registros de forma concorrente utilizando 2 workers. O comando executado foi:
`TEST_NETWORK_PATH=./fabric-samples/test-network go run ./app/cmd/seeding/MetadataSeeding.go -count 5 -workers 2`

### Métricas Gerais

- **Tempo Sequencial**: 11.05s
- **TPS Sequencial**: 0.45 tx/sec
- **Tempo Concorrente**: 6.16s
- **TPS Concorrente**: 0.81 tx/sec

### Speedup Calculado

A execução concorrente obteve um **Speedup de 1.79x**, sendo aproximadamente 1.79 vezes mais rápida que a execução puramente sequencial para a quantidade de registros avaliada.

### Resumo em Formato de Tabela

| Fase            | Workers | Throughput (TPS) |
|-----------------|---------|------------------|
| Sequencial      | 1       | 0.45             |
| Concorrente     | 2       | 0.81             |

---
**Observação**: Durante o modo concorrente, é possível que surjam exceções do tipo `MVCC_READ_CONFLICT`. Isso ocorre como comportamento padrão e esperado no Hyperledger Fabric quando múltiplas transações tentam interagir com a mesma chave ou alteram estados concorrentemente, causando conflito de versão de leitura antes do commit dos blocos na rede.
