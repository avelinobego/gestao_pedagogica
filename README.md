# Gestão Pedagógica

Projeto em Go para gestão pedagógica, com estruturas acadêmicas, avaliação e matrícula. O código está organizado em um único módulo Go:

```text
github.com/avelinobego/gestao_pedagogica
```

> **Escopo atual:** o projeto oferece páginas HTML server-rendered para matrículas e gestão acadêmica de instituições, cursos, etapas, disciplinas e turmas. A interface ainda não tem autenticação; mantenha-a em rede confiável.

## Requisitos

- Go na versão indicada em [`go.mod`](./go.mod) (`1.27.1`).
- PostgreSQL para executar as operações de persistência.

## Organização do código

| Diretório | Responsabilidade |
| --- | --- |
| `cmd/` | Inicialização e encerramento do servidor HTTP. |
| `model/` | Estruturas de dados e tipos de status do sistema. |
| `domain/` | Regras de negócio independentes de persistência, como validação de matrícula. |
| `service/` | Operações de leitura e persistência dos contextos implementados. |
| `web/` | Handlers HTTP, rotas e templates HTML. |
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
Os valores `DATABASE_DRIVER` e `DATABASE_DSN`, quando definidos no ambiente, substituem os campos correspondentes do arquivo.

O endereço HTTP pode ser configurado no mesmo arquivo:

```yaml
server:
  address: "127.0.0.1:8080"
```

O servidor carrega a configuração e conecta ao banco durante a inicialização. A conexão é fechada ao encerrar o processo. O driver PostgreSQL (`lib/pq`) é registrado pela camada de infraestrutura.

## Banco de dados e migrations

Os arquivos SQL estão em `shared/plataforma/migrations/`. Há migrations `up` para criar e popular as tabelas e migrations `down` correspondentes para revertê-las. A migration inicial habilita a extensão PostgreSQL `uuid-ossp`. Aplique os arquivos `up` em ordem numérica (`000001`, `000002` e assim por diante), usando a ferramenta de banco de dados adotada no seu ambiente.

O repositório não inclui um executor automático de migrations. Confirme que o banco e as extensões exigidas pelos scripts estão preparados antes de aplicá-los.

## Executar

Depois de configurar o PostgreSQL e aplicar as migrations:

Na raiz do repositório:

```sh
go run ./cmd
```

O endereço de escuta vem de `server.address` em `gestao.yaml`; se omitido, usa `127.0.0.1:8080`. A variável de ambiente `HTTP_ADDR`, quando definida, sobrescreve o valor do arquivo. O servidor fica acessível localmente em <http://127.0.0.1:8080/>. Por exemplo:

```sh
HTTP_ADDR=127.0.0.1:8080 go run ./cmd
```

O `Makefile` também define o alvo `run`:

```sh
make run
```

### Rotas

| Método | Rota | Ação |
| --- | --- | --- |
| `GET` | `/` | Exibe a página inicial com atalhos para matrículas. |
| `GET` | `/matriculas` | Lista matrículas cadastradas. |
| `GET` | `/matriculas/nova` | Exibe o formulário de cadastro. |
| `POST` | `/matriculas` | Valida e cria uma matrícula. |
| `POST` | `/matriculas/{id}/cancelar` | Cancela uma matrícula. |
| `POST` | `/matriculas/{id}/ativar` | Ativa uma matrícula cancelada ou trancada. |
| `GET` | `/academico` | Exibe a página inicial do módulo acadêmico. |
| `GET` | `/academico/{entidade}` | Lista instituições, cursos, etapas, disciplinas ou turmas. |
| `GET` | `/academico/{entidade}/nova` | Exibe o formulário para novo cadastro. |
| `POST` | `/academico/{entidade}` | Valida e cria um cadastro acadêmico. |
| `GET` | `/academico/{entidade}/{id}/editar` | Exibe o formulário de edição. |
| `POST` | `/academico/{entidade}/{id}` | Atualiza um cadastro acadêmico. |
| `POST` | `/academico/instituicoes/{id}/inativar` | Inativa uma instituição sem apagar o histórico. |
| `POST` | `/academico/cursos/{id}/inativar` | Inativa um curso sem apagar o histórico. |
| `POST` | `/academico/etapas/{id}/excluir` | Exclui uma etapa somente quando não houver turmas vinculadas. |
| `POST` | `/academico/disciplinas/{id}/inativar` | Inativa uma disciplina sem apagar o histórico. |
| `POST` | `/academico/turmas/{id}/encerrar` | Encerra uma turma sem apagar o histórico. |

Os valores de `{entidade}` são `instituicoes`, `cursos`, `etapas`, `disciplinas` e `turmas`. Matrículas canceladas e trancadas podem ser reativadas para o estado `ATIVA`; os demais estados não permitem essa transição. Instituições, cursos e disciplinas são inativados, turmas são encerradas, e etapas só podem ser excluídas se não houver turmas vinculadas. Os formulários usam `POST` para alterações e o servidor aplica proteção contra requisições cross-origin (CSRF). A autenticação ainda não está implementada: mantenha o serviço ligado apenas a `localhost` ou a uma rede interna confiável e não o exponha diretamente à internet.

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

Os testes de `web` usam dependências simuladas e não precisam do PostgreSQL:

```sh
go test ./web
```

Os testes de integração de `service` e de conexão usam Testcontainers-Go para criar um PostgreSQL efêmero, aplicar todas as migrations `up` e removê-lo ao fim da suíte. É necessário ter Docker (ou um runtime compatível com a API Docker) disponível; esses testes não usam o banco configurado para a aplicação e falham explicitamente se o container não puder ser iniciado.

Execute somente as suítes com integração de banco usando:

```sh
make test-db
```

Ou execute os testes sem integrações de banco com `go test -short ./...`.
