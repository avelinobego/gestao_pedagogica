# Gestão Pedagógica

Projeto em Go para gestão pedagógica, com estruturas acadêmicas, avaliação e matrícula. O código está organizado em um único módulo Go:

```text
github.com/avelinobego/gestao_pedagogica
```

> **Estado atual:** o executável em `cmd/` é apenas demonstrativo e imprime o nome do sistema. Os pacotes `model`, `domain` e `service` contêm estruturas e funcionalidades de matrícula, mas ainda não há uma API ou interface de usuário integrada.

## Requisitos

- Go na versão indicada em [`go.mod`](./go.mod) (`1.27.1`).
- PostgreSQL para executar as operações de persistência.

## Organização do código

| Diretório | Responsabilidade |
| --- | --- |
| `cmd/` | Ponto de entrada do executável. |
| `model/` | Estruturas de dados e tipos de status do sistema. |
| `domain/` | Regras de negócio independentes de persistência, como validação de matrícula. |
| `service/` | Operações de leitura e persistência dos contextos implementados. |
| `shared/` | Configuração e infraestrutura compartilhada. |
| `shared/plataforma/` | Acesso ao banco, composição de handlers transacionais e migrations SQL. |

## Configuração do banco

A configuração usa Viper com o nome `gestao`. O código adiciona `$HOME/.config` e `.` aos diretórios de busca. Como `$HOME` é passado literalmente ao Viper e não é expandido pelo shell, use `gestao.yaml` no diretório de execução para a busca funcionar como esperado. Exemplo:

```yaml
database:
  driver: postgres
  dsn: "host=localhost port=5432 user=postgres password=SUA_SENHA dbname=gestao_pedagogica sslmode=disable"
```

Substitua usuário, senha e demais parâmetros pelos valores do seu ambiente. Não versione credenciais reais.

O pacote `shared/plataforma` tenta conectar ao banco durante sua inicialização. A configuração precisa estar disponível antes de executar código que importe esse pacote. O executável `cmd` atual não importa essa camada nem acessa o banco.

> **Nota de implementação:** embora `lib/pq` esteja listado como dependência, o pacote do driver ainda precisa ser registrado na aplicação que usa `shared/plataforma`. Sem esse registro, a conexão pode falhar com `sql: unknown driver "postgres"`.

## Banco de dados e migrations

Os arquivos SQL estão em `shared/plataforma/migrations/`. Há migrations `up` para criar e popular as tabelas e migrations `down` correspondentes para revertê-las. A migration inicial habilita a extensão PostgreSQL `uuid-ossp`. Aplique os arquivos `up` em ordem numérica (`000001`, `000002` e assim por diante), usando a ferramenta de banco de dados adotada no seu ambiente.

O repositório não inclui um executor automático de migrations. Confirme que o banco e as extensões exigidas pelos scripts estão preparados antes de aplicá-los.

## Executar

Na raiz do repositório:

```sh
go run ./cmd
```

O comando atual imprime:

```text
Gestão Pedagógica - Sistema de Avaliação
```

O `Makefile` também define o alvo `run`:

```sh
make run
```

## Validar

Compilar todos os pacotes:

```sh
go build ./...
```

Executar o analisador estático:

```sh
go vet ./...
```

Executar os testes de domínio, que não precisam de conexão com o banco:

```sh
go test ./domain
```

Os testes de `service` dependem do PostgreSQL e da configuração disponível. Eles também requerem que o driver PostgreSQL esteja registrado; veja a nota na seção de configuração.
