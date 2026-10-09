package model

import (
	"time"

	"github.com/google/uuid"
)

type StatusAluno string

const (
	StatusAlunoAtivo   StatusAluno = "ATIVO"
	StatusAlunoInativo StatusAluno = "INATIVO"
)

type StatusMatricula string

const (
	StatusMatriculaAtiva     StatusMatricula = "ATIVA"
	StatusMatriculaConcluida StatusMatricula = "CONCLUIDA"
	StatusMatriculaCancelada StatusMatricula = "CANCELADA"
	StatusMatriculaTrancada  StatusMatricula = "TRANCADA"
)

type Aluno struct {
	ID             uuid.UUID   `json:"id" db:"id"`
	InstituicaoID  uuid.UUID   `json:"instituicao_id" db:"instituicao_id"`
	Nome           string      `json:"nome" db:"nome"`
	NomeSocial     string      `json:"nome_social" db:"nome_social"`
	CPF            string      `json:"cpf" db:"cpf"`
	RG             string      `json:"rg" db:"rg"`
	DataNascimento time.Time   `json:"data_nascimento" db:"data_nascimento"`
	Email          string      `json:"email" db:"email"`
	Telefone       string      `json:"telefone" db:"telefone"`
	Endereco       string      `json:"endereco" db:"endereco"`
	Cidade         string      `json:"cidade" db:"cidade"`
	UF             string      `json:"uf" db:"uf"`
	Status         StatusAluno `json:"status" db:"status"`
	CriadoEm       time.Time   `json:"criado_em" db:"criado_em"`
	AtualizadoEm   time.Time   `json:"atualizado_em" db:"atualizado_em"`
}

type Responsavel struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Nome     string    `json:"nome" db:"nome"`
	CPF      string    `json:"cpf" db:"cpf"`
	Email    string    `json:"email" db:"email"`
	Telefone string    `json:"telefone" db:"telefone"`
	Endereco string    `json:"endereco" db:"endereco"`
	CriadoEm time.Time `json:"criado_em" db:"criado_em"`
}

type AlunoResponsavel struct {
	AlunoID       uuid.UUID `json:"aluno_id" db:"aluno_id"`
	ResponsavelID uuid.UUID `json:"responsavel_id" db:"responsavel_id"`
	Parentesco    string    `json:"parentesco" db:"parentesco"`
	Principal     bool      `json:"principal" db:"principal"`
}

type Matricula struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	AlunoID       uuid.UUID       `json:"aluno_id" db:"aluno_id"`
	AlunoNome     string          `json:"aluno_nome,omitempty" db:"aluno_nome"`
	TurmaID       uuid.UUID       `json:"turma_id" db:"turma_id"`
	TurmaNome     string          `json:"turma_nome,omitempty" db:"turma_nome"`
	Numero        string          `json:"numero" db:"numero"`
	DataMatricula time.Time       `json:"data_matricula" db:"data_matricula"`
	DataConclusao *time.Time      `json:"data_conclusao,omitempty" db:"data_conclusao"`
	Status        StatusMatricula `json:"status" db:"status"`
	Observacao    string          `json:"observacao" db:"observacao"`
	CriadoEm      time.Time       `json:"criado_em" db:"criado_em"`
	AtualizadoEm  time.Time       `json:"atualizado_em" db:"atualizado_em"`
}
