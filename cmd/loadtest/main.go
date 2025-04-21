package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-faker/faker/v4"
	"math/rand/v2"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func getDynamicPayloadCreateContent() map[string]string {
	urls := []string{faker.URL() + "home/camisa/*", faker.URL() + "home/camisa", faker.URL() + "home"}
	url := urls[rand.IntN(4)]

	payload := map[string]string{
		"video_url":     faker.URL(),
		"thumbnail_url": faker.URL(),
		"endpoint":      url,
	}

	return payload
}

type TestScenario struct {
	Name        string
	Method      string
	URL         string
	PayloadFunc func() []byte
	Weight      int // Peso para distribuição de probabilidade
}

func main() {

	input := getDynamicPayloadCreateContent()
	url := input["endpoint"]
	urlEncoded := base64.StdEncoding.EncodeToString([]byte(url))

	scenarios := []TestScenario{
		{
			Name:   "Criar Conteúdo",
			Method: "POST",
			URL:    "http://localhost:8080/content",
			PayloadFunc: func() []byte {
				payload, _ := json.Marshal(input)
				return payload
			},
			Weight: 10, // 50% das requisições
		},
		{
			Name:   "Criar Produto",
			Method: "GET",
			URL:    fmt.Sprintf("http://localhost:8080/content/%s", urlEncoded),
			PayloadFunc: func() []byte {
				return nil
			},
			Weight: 90, // 30% das requisições
		},
	}

	targeter := func(tgt *vegeta.Target) error {
		// Escolhe um cenário baseado em pesos
		totalWeight := 0
		for _, s := range scenarios {
			totalWeight += s.Weight
		}

		r := rand.IntN(totalWeight)
		cumulativeWeight := 0

		var scenario TestScenario
		for _, s := range scenarios {
			cumulativeWeight += s.Weight
			if r < cumulativeWeight {
				scenario = s
				break
			}
		}

		// Configura o target com base no cenário selecionado
		tgt.Method = scenario.Method
		tgt.URL = scenario.URL
		tgt.Body = scenario.PayloadFunc()
		tgt.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}

		return nil
	}

	// Configuração do teste
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 30 * time.Second

	// Executa o ataque
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics

	// Usa o targeter dinâmico
	for res := range attacker.Attack(targeter, rate, duration, "Cenários Múltiplos") {
		metrics.Add(res)
	}
	metrics.Close()

	// Exiba os resultados
	fmt.Printf("Resultados:\n")
	fmt.Printf("Total de requisições: %d\n", metrics.Requests)
	fmt.Printf("Taxa de sucesso: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Latência média: %s\n", metrics.Latencies.Mean)
	fmt.Printf("Latência P99: %s\n", metrics.Latencies.P99)
	fmt.Printf("Latência máxima: %s\n", metrics.Latencies.Max)
	fmt.Printf("Throughput: %.2f req/s\n", metrics.Throughput)

	payloadMetrics, _ := json.Marshal(map[string]interface{}{
		"metrics": metrics,
		"payload": input,
	})

	os.WriteFile("payload_metrics.json", payloadMetrics, 0644)
}
