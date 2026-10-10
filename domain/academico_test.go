package domain

import (
	"testing"

	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/google/uuid"
)

func TestValidarEntidadesAcademicas(t *testing.T) {
	instituicaoID, cursoID := uuid.New(), uuid.New()
	testes := []struct {
		nome    string
		validar func() error
	}{
		{
			nome: "instituição válida",
			validar: func() error {
				return ValidarInstituicaoEnsino(model.InstituicaoEnsino{Nome: "Escola", UF: "SP"})
			},
		},
		{
			nome: "curso válido",
			validar: func() error {
				return ValidarCurso(model.Curso{InstituicaoID: instituicaoID, Nome: "Ensino médio", Codigo: "EM"})
			},
		},
		{
			nome: "etapa válida",
			validar: func() error {
				return ValidarEtapaEnsino(model.EtapaEnsino{CursoID: cursoID, Nome: "1º ano", Ordem: 1, CargaHoraria: 800})
			},
		},
		{
			nome: "disciplina válida",
			validar: func() error {
				return ValidarDisciplina(model.Disciplina{InstituicaoID: instituicaoID, Nome: "Matemática", Codigo: "MAT", CargaHoraria: 160})
			},
		},
		{
			nome: "turma válida",
			validar: func() error {
				return ValidarTurma(model.Turma{
					InstituicaoID: instituicaoID, CursoID: cursoID, EtapaID: uuid.New(),
					Nome: "1º ano A", Codigo: "1A", AnoLetivo: 2026,
					Turno: model.TurnoMatutino, Status: model.StatusTurmaPlanejada,
				})
			},
		},
	}
	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			if err := tt.validar(); err != nil {
				t.Fatalf("validação inesperada: %v", err)
			}
		})
	}
}

func TestValidarEntidadesAcademicasRejeitaDadosInvalidos(t *testing.T) {
	testes := []struct {
		nome    string
		validar func() error
	}{
		{
			nome: "instituição sem nome",
			validar: func() error {
				return ValidarInstituicaoEnsino(model.InstituicaoEnsino{})
			},
		},
		{
			nome: "CNPJ com letras",
			validar: func() error {
				return ValidarInstituicaoEnsino(model.InstituicaoEnsino{Nome: "Escola", CNPJ: "12ABC"})
			},
		},
		{
			nome: "UF minúscula",
			validar: func() error {
				return ValidarInstituicaoEnsino(model.InstituicaoEnsino{Nome: "Escola", UF: "sp"})
			},
		},
		{
			nome: "UF não alfabética",
			validar: func() error {
				return ValidarInstituicaoEnsino(model.InstituicaoEnsino{Nome: "Escola", UF: "12"})
			},
		},
		{
			nome: "e-mail inválido",
			validar: func() error {
				return ValidarInstituicaoEnsino(model.InstituicaoEnsino{Nome: "Escola", Email: "inválido"})
			},
		},
		{
			nome: "curso sem instituição",
			validar: func() error {
				return ValidarCurso(model.Curso{Nome: "Curso", Codigo: "CUR"})
			},
		},
		{
			nome: "etapa com ordem inválida",
			validar: func() error {
				return ValidarEtapaEnsino(model.EtapaEnsino{CursoID: uuid.New(), Nome: "Etapa", Ordem: 0, CargaHoraria: 100})
			},
		},
		{
			nome: "disciplina sem carga horária",
			validar: func() error {
				return ValidarDisciplina(model.Disciplina{InstituicaoID: uuid.New(), Nome: "Matemática", Codigo: "MAT"})
			},
		},
		{
			nome: "turma com turno inválido",
			validar: func() error {
				return ValidarTurma(model.Turma{
					InstituicaoID: uuid.New(), CursoID: uuid.New(), EtapaID: uuid.New(),
					Nome: "Turma", Codigo: "T1", AnoLetivo: 2026,
					Turno: "MADRUGADA", Status: model.StatusTurmaPlanejada,
				})
			},
		},
	}
	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			if err := tt.validar(); err == nil {
				t.Fatal("esperava erro de validação")
			}
		})
	}
}
