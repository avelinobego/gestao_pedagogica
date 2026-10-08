package model

import (
	"time"

	"github.com/google/uuid"
)

type StatusTurma string

const (
	StatusTurmaAtiva     StatusTurma = "ATIVA"
	StatusTurmaEncerrada StatusTurma = "ENCERRADA"
	StatusTurmaPlanejada StatusTurma = "PLANEJADA"
)

type Turno string

const (
	TurnoMatutino   Turno = "MATUTINO"
	TurnoVespertino Turno = "VESPERTINO"
	TurnoNoturno    Turno = "NOTURNO"
	TurnoIntegral   Turno = "INTEGRAL"
)

type InstituicaoEnsino struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Nome         string    `json:"nome" db:"nome"`
	NomeFantasia string    `json:"nome_fantasia" db:"nome_fantasia"`
	CNPJ         string    `json:"cnpj" db:"cnpj"`
	CodigoINEP   string    `json:"codigo_inep" db:"codigo_inep"`
	Email        string    `json:"email" db:"email"`
	Telefone     string    `json:"telefone" db:"telefone"`
	Endereco     string    `json:"endereco" db:"endereco"`
	Cidade       string    `json:"cidade" db:"cidade"`
	UF           string    `json:"uf" db:"uf"`
	Ativa        bool      `json:"ativa" db:"ativa"`
	CriadoEm     time.Time `json:"criado_em" db:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em" db:"atualizado_em"`
}

type Curso struct {
	ID            uuid.UUID `json:"id" db:"id"`
	InstituicaoID uuid.UUID `json:"instituicao_id" db:"instituicao_id"`
	Nome          string    `json:"nome" db:"nome"`
	Codigo        string    `json:"codigo" db:"codigo"`
	Descricao     string    `json:"descricao" db:"descricao"`
	Ativo         bool      `json:"ativo" db:"ativo"`
}

type EtapaEnsino struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CursoID      uuid.UUID `json:"curso_id" db:"curso_id"`
	Nome         string    `json:"nome" db:"nome"`
	Ordem        int       `json:"ordem" db:"ordem"`
	CargaHoraria int       `json:"carga_horaria" db:"carga_horaria"`
}

type Disciplina struct {
	ID            uuid.UUID `json:"id" db:"id"`
	InstituicaoID uuid.UUID `json:"instituicao_id" db:"instituicao_id"`
	Nome          string    `json:"nome" db:"nome"`
	Codigo        string    `json:"codigo" db:"codigo"`
	CargaHoraria  int       `json:"carga_horaria" db:"carga_horaria"`
	Ativa         bool      `json:"ativa" db:"ativa"`
}

type Turma struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	InstituicaoID uuid.UUID   `json:"instituicao_id" db:"instituicao_id"`
	CursoID       uuid.UUID   `json:"curso_id" db:"curso_id"`
	EtapaID       uuid.UUID   `json:"etapa_id" db:"etapa_id"`
	Nome          string      `json:"nome" db:"nome"`
	Codigo        string      `json:"codigo" db:"codigo"`
	AnoLetivo     int         `json:"ano_letivo" db:"ano_letivo"`
	Turno         Turno       `json:"turno" db:"turno"`
	Status        StatusTurma `json:"status" db:"status"`
	CriadoEm      time.Time   `json:"criado_em" db:"criado_em"`
}
