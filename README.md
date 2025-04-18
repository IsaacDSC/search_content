# Video Rendering Service

## Descrição do Desafio

Este projeto consiste em criar um serviço de alta performance para disponibilizar conteúdo de vídeo e thumbnails para frontend. O sistema precisa responder a diferentes endpoints e entregar o recurso apropriado (vídeo ou thumbnail) de acordo com a requisição recebida.

## Requisitos Técnicos

- **Performance Crítica**: O serviço precisa ser extremamente rápido no tempo de resposta
- **Alta Escala**: O sistema deve suportar milhões de requisições
- **Baixo Custo**: A implementação deve ser otimizada para minimizar custos operacionais
- **Sem Cache**: Devido ao volume de requisições, implementar cache seria economicamente inviável

## Funcionalidades

O serviço deve expor endpoints que:
- Entreguem vídeos otimizados quando solicitado
- Disponibilizem thumbnails correspondentes aos vídeos
- Adaptem-se às necessidades do frontend consumidor

## Por que não utilizar cache?

A implementação de cache para este volume de dados apresentaria desafios significativos:
- O custo de armazenamento para milhões de recursos diferentes seria proibitivo
- A taxa de hit do cache seria potencialmente baixa dependendo do padrão de acesso
- O gerenciamento e invalidação do cache adicionaria complexidade desnecessária

## Considerações de Arquitetura

Para atender aos requisitos de performance e custo, algumas estratégias podem incluir:
- Otimização no processamento de vídeo sob demanda
- Implementação de algoritmos eficientes para geração de thumbnails
- Distribuição geográfica estratégica dos recursos
- Escolha cuidadosa de formatos de compressão e codificação

## Métricas de Sucesso

- Tempo de resposta abaixo de X ms
- Custo por milhão de requisições dentro do orçamento estabelecido
- Estabilidade do serviço sob carga máxima


## Controle de Metrias
```shell
go test -bench=. ./pkg/filesystem >> filesystem_result.txt

```



### Atualizar mocks
```shell
mockgen -source=pkg/filesystem/adapter.go -destination=pkg/filesystem/driver_mock.go -package=filesystem

```