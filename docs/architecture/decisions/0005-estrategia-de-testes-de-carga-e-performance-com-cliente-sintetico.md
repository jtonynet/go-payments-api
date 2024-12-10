# 5. Estratégia de Testes de Carga e Performance com Cliente Sintético

**Data:** 10 de dezembro de 2024

## **Status**

Aceito

## **Contexto**

### Simular cargas aproximadas às de produção através de Cliente Sintético

_`Clientes Sintéticos` (como `Gatling`, `Jmeter`, `K6`...) utilizados em [`Testes de Carga`](https://pt.wikipedia.org/wiki/Teste_de_carga) são ferramentas ou serviços que simulam interações reais de usuários com uma aplicação ou API._ Não deve ser confundido com a técnica de [`Testes Sintéticos`](https://www.hipsters.tech/testes-sinteticos-no-c6-bank-hipsters-on-the-road-40/), que envolve monitoramento em produção.

O desenvolvimento da `payment API` com um `timeoutSLA` de `100ms`, como requisito de negócio, é desafiador. Embora em `Golang` existam recursos poderosos para controle de concorrência e cancelamento, como `context.timeout`, validar a concorrência com timeout em cenários próximos aos reais na máquina do desenvolvedor pode ser frustrante.

O uso de um `Cliente Sintético` é essencial. Este documento avalia abordagens e ferramentas para testes de `Performance/Desempenho`, executáveis localmente e em ambientes próximos à produção (como `pre-prod`, `homol`, `stg` etc.) Embora já tenha tido experiências com o [`Jmeter`](https://jmeter.apache.org/) e o [`Vegeta`](https://github.com/tsenart/vegeta), desconsiderei-os em favor de dois outros clientes mais modernos e em ascensão: o [`Gatling`](https://gatling.com/), pela fácil configuração, e o [`Grafana K6`](https://k6.io/) pela aderencia a `observabilidade`

### Referências e Opções de Clientes Sintéticos:

[`Grafana Load Testing`](https://grafana.com/load-testing/)
<br/>Embora da equipe `Grafana`, oferece overview abrangentes sobre estratégia, ferramentas e tipos de testes.


- [`Grafana K6`](https://k6.io/)
  - [PPT Slides 2023](https://pt.slideshare.net/slideshow/k6-teste-de-carga-e-desempenhopptx/257546892#2)
  - [Repositório](https://github.com/grafana/k6)
  - [Artigo do Blog Full Cycle](https://fullcycle.com.br/como-fazer-testes-de-carga-nas-suas-aplicacoes/)

- [`Gatling`](https://gatling.com/)
  - [PPT Slides TDC 2018](https://pt.slideshare.net/slideshow/tdc2018sp-trilha-testes-testes-de-carga-e-performance-com-gatlingio/108137696#2)
  - [Load Testing A Dockerized Application](https://gatling.io/blog/load-testing-a-dockerized-application)
  - [Step-by-Step: Gatling Load Tests with TestContainers & Docker](https://gatling.io/blog/step-by-step-gatling-load-tests-with-testcontainers-and-docker)

<br/>

## Decisão

Como o uso do script `.sh` do `Gatling` já é conhecido, utilizaremos para configurar inicialmente um teste de carga com esforço de desenvolvimento reduzido. Porém, a [configuração em novas versões](https://github.com/gatling/gatling/issues/4512) do mesmo foi [alterada](https://community.gatling.io/t/missing-command-line-options-in-gatling-3-11-bundles/9311), o que força a manter uma versão antiga (3.9.5). 

Esse cenário frustrante, embora incentive a pesquisa sobre outras maneiras de utilização do `Gatling`, nos levou a avaliar sua substituição pelo `K6` no futuro próximo. Além da modernidade da ferramenta com integrações a pipelines CI/CD, suas [`extensões escritas em Golang`](https://grafana.com/docs/k6/latest/extensions/) 🫶🏽 e ao fato de já existirem iniciativas (não documentadas) para que os [testes sejam escritos na mesma linguagem do projeto](https://github.com/szkiba/xk6-g0) (além do padrão em `TypeScript`).

Sendo assim, no momento, o projeto deve continuar com `Gatling` em versão antiga, mas tão logo a `Observabilidade` seja adicionada ao projeto, seu uso deve ser pivotado para o `K6`, o que deve servir também como estudo de sua integração com as ferramentas da família `Grafana` que fazem sentido nesse cenário.

<br/>

## Consequências

Inicialmente, teremos testes que nortearão o desenvolvimento e a implantação, mesmo com um `Cliente Sintético` desatualizado. À medida que os requisitos de `Observabilidade` forem atendidos, a migração para uma ferramenta atual, aderente a linguagem do projeto e de interesse do mercado (hype), como o `Grafana K6`, torna-se atrativa.

Cumprimos o `timeoutSLA` com testes simples, evoluindo para uma abordagem mais robusta conforme o projeto avança.

- Testes básicos serão realizados com `Gatling`, permitindo validar o `timeoutSLA`.
- A transição para o `Grafana K6` deve ocorrer com o amadurecimento da observabilidade, melhorando alinhamento com o mercado e modernizando a abordagem.

