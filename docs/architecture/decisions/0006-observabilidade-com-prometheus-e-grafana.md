# 6. Observabilidade com Prometheus e Grafana

Data: 10 de dezembro de 2024

## Status

Aceito

## Contexto

Dado que o projeto está sendo desenvolvido, é desejável obtermos dados sobre a escala que ele deve alcançar, ter algum nível de __Observabilidade__ desde o início é muito útil. As métricas fornecidas por essa abordagem nos orientam sobre como nossos projetos crescem, tornando-os mais escaláveis e sustentáveis.

É crucial para obedecer ao `timeoutSLA` restrito de 100ms, entendermos em quais situações ocorrem os Erros, sua Duração e Taxas (`RED`), bem como os dados de conexão aos nossos `Bancos` e respectivas `fontes de verdade`, em conjunto com nossos testes de Performance/Desempenho/Carga.

Decisões de projeto sempre devem ser tomadas com base em __métricas__.

## Decisão

Como o [Prometheus](https://prometheus.io/) e o [Grafana](https://grafana.com/) são amplamente utilizados no mercado e a equipe possui experiência e conhecimento prévio com eles, além do fato de que é desejável a migração da ferramenta de `Cliente Sintético` para o `Grafana K6`, faz sentido a ampla adoção de ferramentas `Grafana`.

Faremos a implementação de um middleware usando o [Prometheus](https://github.com/prometheus/client_golang) e o [Gorm-Prometheus](https://github.com/go-gorm/prometheus) customizado no framework `Gin`, respeitando ao máximo a `Arquitetura Hexagonal` adotada no projeto, juntamente com suas devidas configurações em variáveis de ambiente. Também preparamos nosso `docker-compose.yml` para atender à estrutura dessas ferramentas, permitindo um desenvolvimento na máquina local que esteja alinhado com essa diretriz.

## Consequências

Obteremos um conjunto básico de métricas `RED (Rate, Errors e Duration)` e `API Basics`, amplamente utilizadas no mercado, que são indicadores do crescimento e escalabilidade de nossas APIs. Balizadores do desenvolvimento. Assim como a análise dos testes de Performance/Desempenho/Carga devem ser enriquecidas, nos dando um panorama geral mais abrangente da aplicação. Preparamos o terreno para adotar o `Grafana K6` como nosso `Cliente Sintético` oficial.
