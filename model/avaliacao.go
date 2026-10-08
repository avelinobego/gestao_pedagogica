package model

import (
	"time"

	"github.com/google/uuid"
)

// StatusAvaliacao indica a situação do aluno na disciplina.
type StatusAvaliacao string

const (
	StatusEmAndamento StatusAvaliacao = "EM_ANDAMENTO"
	StatusAprovado    StatusAvaliacao = "APROVADO"
	StatusReprovado   StatusAvaliacao = "REPROVADO"
	StatusRecuperacao StatusAvaliacao = "RECUPERACAO"
)

// TipoAvaliacao define a natureza do instrumento de avaliação.
type TipoAvaliacao string

const (
	TipoProva      TipoAvaliacao = "PROVA"
	TipoTrabalho   TipoAvaliacao = "TRABALHO"
	TipoAtividade  TipoAvaliacao = "ATIVIDADE"
	TipoExameFinal TipoAvaliacao = "EXAME_FINAL"
)

// PeriodoLetivo representa o bimestre/trimestre letivo.
type PeriodoLetivo struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Nome       string    `json:"nome" db:"nome"`   // Ex: "1º Bimestre"
	Ordem      int       `json:"ordem" db:"ordem"` // Ex: 1
	DataInicio time.Time `json:"data_inicio" db:"data_inicio"`
	DataFim    time.Time `json:"data_fim" db:"data_fim"`
}

// Avaliacao representa a atividade/prova criada pelo professor.
type Avaliacao struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	TurmaID         uuid.UUID     `json:"turma_id" db:"turma_id"`
	DisciplinaID    uuid.UUID     `json:"disciplina_id" db:"disciplina_id"`
	PeriodoLetivoID uuid.UUID     `json:"periodo_letivo_id" db:"periodo_letivo_id"`
	Titulo          string        `json:"titulo" db:"titulo"` // Ex: "Prova de Álgebra"
	Descricao       string        `json:"descricao" db:"descricao"`
	Tipo            TipoAvaliacao `json:"tipo" db:"tipo"`
	Peso            int           `json:"peso" db:"peso"`               // Ex: 2 para peso 2
	NotaMaxima      int           `json:"nota_maxima" db:"nota_maxima"` // Armazenado como inteiro (ex: 1000 = 10.00)
	DataAplicacao   time.Time     `json:"data_aplicacao" db:"data_aplicacao"`
	CriadoEm        time.Time     `json:"criado_em" db:"criado_em"`
}

// NotaRegistro representa o lançamento efetivo da nota de um aluno em uma avaliação.
type NotaRegistro struct {
	ID           uuid.UUID `json:"id" db:"id"`
	AvaliacaoID  uuid.UUID `json:"avaliacao_id" db:"avaliacao_id"`
	AlunoID      uuid.UUID `json:"aluno_id" db:"aluno_id"`
	Valor        int       `json:"valor" db:"valor"`             // Armazenado multiplicado por 100 (ex: 850 = 8.50)
	Observacao   string    `json:"observacao" db:"observacao"`   // Parecer descritivo ou justificativa
	LancadoPor   uuid.UUID `json:"lancado_por" db:"lancado_por"` // ID do usuário/professor
	LancadoEm    time.Time `json:"lancado_em" db:"lancado_em"`
	AtualizadoEm time.Time `json:"atualizado_em" db:"atualizado_em"`
}

// ConsolidadoPeriodo representa a média calculada do aluno em um período específico.
type ConsolidadoPeriodo struct {
	ID              uuid.UUID `json:"id" db:"id"`
	AlunoID         uuid.UUID `json:"aluno_id" db:"aluno_id"`
	DisciplinaID    uuid.UUID `json:"disciplina_id" db:"disciplina_id"`
	PeriodoLetivoID uuid.UUID `json:"periodo_letivo_id" db:"periodo_letivo_id"`
	MediaCalculada  int       `json:"media_calculada" db:"media_calculada"` // Armazenado multiplicado por 100 (ex: 750 = 7.50)
	TotalAusencias  int       `json:"total_ausencias" db:"total_ausencias"`
	CalculadoEm     time.Time `json:"calculado_em" db:"calculado_em"`
}

// BoletimAluno representa a situação final/geral do aluno no ano letivo.
type BoletimAluno struct {
	AlunoID     uuid.UUID          `json:"aluno_id" db:"aluno_id"`
	TurmaID     uuid.UUID          `json:"turma_id" db:"turma_id"`
	AnoLetivo   int                `json:"ano_letivo" db:"ano_letivo"`
	Disciplinas []DisciplinaResumo `json:"disciplinas" db:"disciplinas"`
	StatusGeral StatusAvaliacao    `json:"status_geral" db:"status_geral"`
}

// DisciplinaResumo consolida as notas bimestrais e nota final de uma disciplina.
type DisciplinaResumo struct {
	DisciplinaID uuid.UUID            `json:"disciplina_id" db:"disciplina_id"`
	Nome         string               `json:"nome" db:"nome"`
	NotasPeriodo []ConsolidadoPeriodo `json:"notas_periodo" db:"notas_periodo"`
	MediaFinal   int                  `json:"media_final" db:"media_final"`
	Status       StatusAvaliacao      `json:"status" db:"status"`
}
