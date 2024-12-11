# 2. Go, Gin, Gorm e PostgreSQL com Arquitetura Hexagonal e TDD

Data: 05 de outubro de 2024

## Status

Aceito

## Contexto

Precisamos definir o melhor fluxo de trabalho e testes para a go-payments-api.

Existem muitas opções para a arquitetura base, pois elementos cruciais para o projeto foram deixados em aberto (ainda que o desafiante tenha suas __preferências__), conforme descrito na seção [Sobre do arquivo README.md](../../../README.md).

O desafio cita implementação de testes, dos quais utilizarei os `unit tests` e `integration tests` em minhas `services`, `repositories` e `routes` para auxiliar no fluxo de desenvolvimento em `TDD`.

<div align="center">
<img src="../../assets/images/layout/graphics/test_pyramid.jpg">
<br/><i><a href="https://www.headspin.io/blog/the-testing-pyramid-simplified-for-one-and-all">*Imagem retirada do artigo: The Testing Pyramid: Simplified for One and All</a></i>
</div>


## Decisão

Este documento determina o fluxo de trabalho __Branch Based com Feature Branch__, design estrutural e a abordagem de testes para garantir um padrão para a aplicação.

O desafio sugere `Scala`, `Kotlin` e o `paradigma de programação funcional`. Porém, realizarei em Golang, em arquitetura [`hexagonal`](https://alistair.cockburn.us/hexagonal-architecture/), por ter maior familiaridade, além de considerá-las altamente performáticas.

O desafio deixa claro no trecho:

><br/>
> __Linguagem e bibliotecas:__
> 
> Na *********, usamos Scala e Kotlin no nosso dia a dia (e demonstrar experiência em alguma delas é um grande diferencial). No entanto, você pode implementar sua solução utilizando sua linguagem favorita, dando preferência ao paradigma de programação funcional.
> <br/>
> <br/>

Evidenciando preferências, mas aceitando submissões com outras linguagens e paradigmas.

### Justificativa

#### Arquitetura Hexagonal
A Arquitetura [Hexagonal](https://alistair.cockburn.us/hexagonal-architecture/) foi escolhida por sua capacidade de isolar o domínio do problema das implementações tecnológicas, permitindo que mudanças em frameworks ou bancos de dados não impactem o núcleo da aplicação. Esse design também facilita a testabilidade e a separação de responsabilidades, o que é crítico em um projeto que deve evoluir rapidamente sem comprometer a manutenibilidade.

<div align="center">
<img src="../../assets/images/layout/graphics/hexagonal_style-1.jpg">

<i><a href="https://elemarjr.com/arquivo/ensuring-the-quality-of-the-domain-model-through-the-hexagonal-architecture-pattern/">*Imagem retirada do artigo: Hexagonal Architecture Pattern</a></div>
</div>

#### GIN
Foi selecionado como o framework de API por sua alta performance e simplicidade. Ideal para aplicações de baixa latência como esta. Sua arquitetura minimalista também facilita a rápida implementação de novos endpoints, sem sobrecarregar o processo de desenvolvimento. 

#### GORM
Escolhemos GORM pela sua flexibilidade e capacidade de integração com os principais bancos de dados. A abstração oferecida pelo GORM simplifica a manipulação de dados, ao mesmo tempo que permite um controle refinado sobre queries mais complexas, quando necessário. Essa escolha também nos prepara para adoção futura de ferramentas de observabilidade e tracking, garantindo que a aplicação possa evoluir sem grandes refactors.

#### Postgres
Optamos pelo PostgreSQL devido à sua robustez e features modernas, como JSONB, que oferecem flexibilidade para lidar com dados estruturados e semiestruturados. Postgres também é conhecido pela sua confiabilidade em ambientes de alta carga, o que é essencial considerando o volume de transações esperadas.


#### TDD
A adoção de TDD garante que a aplicação seja desenvolvida com um foco claro na cobertura de testes, minimizando bugs e retrabalho ao longo do ciclo de vida do projeto. Isso também nos prepara para uma maior resiliência em produção, especialmente considerando o impacto de falhas em um sistema financeiro.

### Ferramentas Adicionais
Adotaremos também o **Swagger** para documentação da API, garantindo que o produto seja bem documentado desde o início, permitindo integração com outras equipes e sistemas. O uso do [**GitHub Projects**](https://github.com/users/jtonynet/projects/7/views/1) para um fluxo Kanban auxiliará no acompanhamento das entregas e na priorização adequada das tarefas.


## Consequências

A escolha dessas tecnologias, aliada a uma abordagem iterativa e incremental, permite que o projeto seja escalável e flexível. A documentação clara, através do Swagger, ADRs e diagramas, garantirá que as equipes futuras possam se integrar ao projeto com facilidade. Além disso, a arquitetura adotada nos prepara para futuras mudanças tecnológicas, sem comprometer o core business ou aumentar os custos operacionais.

