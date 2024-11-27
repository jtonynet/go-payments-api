# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]
### Added

## [0.1.6] - 2024-25-11
### Adicionado

- [Issue 33](https://github.com/jtonynet/go-payments-api/issues/33)
  - Segregadas as APIs `REST` e `Processor` (anteriormente chamada de `Worker`) por questões de segurança
  - Projeto dividido em dois diretórios `cmd` e, consequentemente, dois `binários` 
  - Adicionada comunicação via `gRPC` entre esses `serviços`
  - Ajustes em `testes` e `CI`
  - Ajustes no `logger`

## [0.1.5] - 2024-21-11
### Adicionado

  - [Issue 36](https://github.com/jtonynet/go-payments-api/issues/35)
    - Removido `MVP` da solução L4 com `Exponential Backoff Retry` como `Growth Hack`
    - Adicionado `Redis Keyspace Notification` como melhoria `L4`
    - Unlock expira o `ttl` da chave de memory lock
    - Timeout de 100ms para memory lock

## [0.1.4] - 2024-20-11
### Adicionado

- [Issue 38](https://github.com/jtonynet/go-payments-api/issues/38)
  - Alteração visando aumentar a a resiliência a consistência eventual do banco.

## [0.1.3] - 2024-15-11
### Adicionado

- Melhorias na `repository` aderente ao padrão Hexagonal
- Implementando `MVP` da solução L4 com `Exponential Backoff Retry` como `Growth Hack`
- Melhorias na documentação

## [0.1.2] - 2024-10-11
### Adicionado

- Adaptador de Roteador aderente ao padrão Hexagonal
- Correções e pequenas refatorações no projeto
- Estratégia de cache com Redis adicionada para o merchant
- Versão mais madura do diagrama da solução L4
- Melhorias na documentação

## [0.1.1] - 2024-23-10
### Added

- Aumento da cobertura de testes
- Categorias e MCCs extraidos para a database
- Diagrama descritivo L4 com uso de `fila` e `memory lock`
- Outras pequenas melhorias


### Fixed
## [0.1.0] - 2024-23-10
### Fixed

  - Bump version
  - Status no README movido para AGUARDANDO
  - Removido arquivo antigo desnecessário

### Fixed
## [0.0.7] - 2024-23-10
### Fixed

  - Movendo Custom Error para domain

### Fixed
## [0.0.6] - 2024-23-10
### Added

- [Issue 6](https://github.com/jtonynet/go-payments-api/issues/12)
  - L4 finalizada
  - melhorias no README


## [0.0.5] - 2024-23-10
### Added

- [Issue 10](https://github.com/jtonynet/go-payments-api/issues/10)
  - L3 finalizada
  - Acerto na L2, debita os fundos da categoria principal disponivel e o restante de `fallback` CASH
  - melhorias nos testes
  - melhorias na arquitetura

## [0.0.4] - 2024-20-10
### Added

- [Issue 8](https://github.com/jtonynet/go-payments-api/issues/8)
  - L2 finalizada
  - Adiciona a lógica para debitar o saldo de CASH se o saldo da categoria principal for insuficiente.
  - melhorias nos testes
  - melhorias na arquitetura

## [0.0.3] - 2024-20-10
### Added

- [Issue 6](https://github.com/jtonynet/go-payments-api/issues/6)
  - L1 finalizada
  - Testes de integração na routes, unitários de repository e routes adicionados
  -  Acertos no Github Actions
  - Adicionado Swagger
  - Adicionado hot reload
  


## [0.0.2] - 2024-09-10
### Added

- [Issue 4](https://github.com/jtonynet/go-payments-api/issues/4)
  - Criar estrutura base para a implementação das tarefas/cards L1 a L3


## [0.0.1] - 2024-05-10
### Added

- [Issue 2](https://github.com/jtonynet/go-payments-api/issues/2)
   - Banco Postgres
   - API
   - docker compose
   - Dockerfile da API


## [0.0.0] - 2024-05-10
### Added

- [Issue 1](https://github.com/users/jtonynet/projects/7/views/1?pane=issue&itemId=82288146)
  - [Kanban Project View Iniciado](https://github.com/users/jtonynet/projects/7) com o commit inicial. 
  - Documentação base: Readme Rico, [Diagramas Mermaid](https://github.com/jtonynet/go-products-api/tree/main#diagrams), ADRs: [0001: Registro de Decisões de Arquitetura (ADR)](./docs/architecture/decisions/registro-de-decisoes-de-arquitetura.md) e [0002: Go, Gin, Gorm e PostgreSQL com Arquitetura Hexagonal e TDD](./docs/architecture/decisions/0002-go-gin-gorm-e-postgres-com-arquitetura-hexagonal-tdd.md).
  - Sabemos o que fazer, graças às definições do arquivo __README.md__. Sabemos como fazer graças aos __ADRs__ e documentações vinculadas. Devemos nos organizar em estrutura __Kanban__, guiados pelo modelo Agile, em nosso __Github Project__, e dar o devido prosseguimento às tarefas.

[0.1.6]: https://github.com/jtonynet/go-payments-api/compare/v0.1.5...v0.1.6
[0.1.5]: https://github.com/jtonynet/go-payments-api/compare/v0.1.4...v0.1.5
[0.1.4]: https://github.com/jtonynet/go-payments-api/compare/v0.1.3...v0.1.4
[0.1.3]: https://github.com/jtonynet/go-payments-api/compare/v0.1.2...v0.1.3
[0.1.2]: https://github.com/jtonynet/go-payments-api/compare/v0.1.1...v0.1.2
[0.1.1]: https://github.com/jtonynet/go-payments-api/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/jtonynet/go-payments-api/compare/v0.0.7...v0.1.0
[0.0.7]: https://github.com/jtonynet/go-payments-api/compare/v0.0.6...v0.0.7
[0.0.6]: https://github.com/jtonynet/go-payments-api/compare/v0.0.5...v0.0.6
[0.0.5]: https://github.com/jtonynet/go-payments-api/compare/v0.0.4...v0.0.5
[0.0.4]: https://github.com/jtonynet/go-payments-api/compare/v0.0.3...v0.0.4
[0.0.3]: https://github.com/jtonynet/go-payments-api/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/jtonynet/go-payments-api/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/jtonynet/go-payments-api/compare/v0.0.0...v0.0.1
[0.0.0]: https://github.com/jtonynet/go-payments-api/releases/tag/v0.0.0
