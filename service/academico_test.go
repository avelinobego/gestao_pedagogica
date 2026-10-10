package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/avelinobego/gestao_pedagogica/domain"
	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"github.com/google/uuid"
)

func TestServicosAcademicosCRUDELifecycle(t *testing.T) {
	requireDatabase(t)
	ctx := context.Background()

	instituicao := model.InstituicaoEnsino{
		ID: uuid.New(), Nome: "Instituição de teste " + uuid.NewString()[:8], UF: "SP", Ativa: true,
	}
	cleanupRegistro(t, "instituicao_ensino", instituicao.ID)
	if err := plataforma.Executar(ctx, CriarInstituicaoEnsino(&instituicao)); err != nil {
		t.Fatalf("CriarInstituicaoEnsino() retornou erro: %v", err)
	}
	if _, err := BuscarInstituicaoEnsino(ctx, instituicao.ID); err != nil {
		t.Fatalf("BuscarInstituicaoEnsino() retornou erro: %v", err)
	}
	instituicao.Nome = "Instituição atualizada"
	if err := plataforma.Executar(ctx, AtualizarInstituicaoEnsino(instituicao)); err != nil {
		t.Fatalf("AtualizarInstituicaoEnsino() retornou erro: %v", err)
	}

	curso := model.Curso{
		ID: uuid.New(), InstituicaoID: instituicao.ID, Nome: "Curso " + uuid.NewString()[:8],
		Codigo: uuid.NewString()[:8], Ativo: true,
	}
	cleanupRegistro(t, "curso", curso.ID)
	if err := plataforma.Executar(ctx, CriarCurso(&curso)); err != nil {
		t.Fatalf("CriarCurso() retornou erro: %v", err)
	}
	if _, err := BuscarCurso(ctx, curso.ID); err != nil {
		t.Fatalf("BuscarCurso() retornou erro: %v", err)
	}
	curso.Descricao = "Descrição atualizada"
	if err := plataforma.Executar(ctx, AtualizarCurso(curso)); err != nil {
		t.Fatalf("AtualizarCurso() retornou erro: %v", err)
	}

	etapa := model.EtapaEnsino{
		ID: uuid.New(), CursoID: curso.ID, Nome: "Etapa " + uuid.NewString()[:8],
		Ordem: 1, CargaHoraria: 800,
	}
	cleanupRegistro(t, "etapa_ensino", etapa.ID)
	if err := plataforma.Executar(ctx, CriarEtapaEnsino(&etapa)); err != nil {
		t.Fatalf("CriarEtapaEnsino() retornou erro: %v", err)
	}
	if _, err := BuscarEtapaEnsino(ctx, etapa.ID); err != nil {
		t.Fatalf("BuscarEtapaEnsino() retornou erro: %v", err)
	}
	etapa.Nome = "Etapa atualizada"
	if err := plataforma.Executar(ctx, AtualizarEtapaEnsino(etapa)); err != nil {
		t.Fatalf("AtualizarEtapaEnsino() retornou erro: %v", err)
	}

	disciplina := model.Disciplina{
		ID: uuid.New(), InstituicaoID: instituicao.ID, Nome: "Disciplina " + uuid.NewString()[:8],
		Codigo: uuid.NewString()[:8], CargaHoraria: 100, Ativa: true,
	}
	cleanupRegistro(t, "disciplina", disciplina.ID)
	if err := plataforma.Executar(ctx, CriarDisciplina(&disciplina)); err != nil {
		t.Fatalf("CriarDisciplina() retornou erro: %v", err)
	}
	if _, err := BuscarDisciplina(ctx, disciplina.ID); err != nil {
		t.Fatalf("BuscarDisciplina() retornou erro: %v", err)
	}
	disciplina.CargaHoraria = 120
	if err := plataforma.Executar(ctx, AtualizarDisciplina(disciplina)); err != nil {
		t.Fatalf("AtualizarDisciplina() retornou erro: %v", err)
	}

	turma := model.Turma{
		ID: uuid.New(), InstituicaoID: instituicao.ID, CursoID: curso.ID, EtapaID: etapa.ID,
		Nome: "Turma " + uuid.NewString()[:8], Codigo: uuid.NewString()[:8], AnoLetivo: 2026,
		Turno: model.TurnoMatutino, Status: model.StatusTurmaPlanejada,
	}
	cleanupRegistro(t, "turma", turma.ID)
	if err := plataforma.Executar(ctx, CriarTurmaAcademica(&turma)); err != nil {
		t.Fatalf("CriarTurmaAcademica() retornou erro: %v", err)
	}
	if _, err := BuscarTurmaAcademica(ctx, turma.ID); err != nil {
		t.Fatalf("BuscarTurmaAcademica() retornou erro: %v", err)
	}
	turma.Status = model.StatusTurmaAtiva
	if err := plataforma.Executar(ctx, AtualizarTurmaAcademica(turma)); err != nil {
		t.Fatalf("AtualizarTurmaAcademica() retornou erro: %v", err)
	}
	turmas, err := ListarTodasTurmasAcademicas(ctx)
	if err != nil {
		t.Fatalf("ListarTodasTurmasAcademicas() retornou erro: %v", err)
	}
	var encontrouTurma bool
	for _, item := range turmas {
		if item.ID == turma.ID {
			encontrouTurma = item.InstituicaoNome == instituicao.Nome &&
				item.CursoNome == curso.Nome && item.EtapaNome == etapa.Nome
		}
	}
	if !encontrouTurma {
		t.Fatal("lista de turmas não trouxe os nomes dos cadastros relacionados")
	}

	if err := plataforma.Executar(ctx, ExcluirEtapaEnsino(etapa.ID)); !errors.Is(err, domain.ErrEtapaPossuiTurmas) {
		t.Fatalf("ExcluirEtapaEnsino() erro = %v; esperado etapa com turmas vinculadas", err)
	}
	etapaSemTurmas := model.EtapaEnsino{
		ID: uuid.New(), CursoID: curso.ID, Nome: "Etapa removível", Ordem: 2, CargaHoraria: 400,
	}
	cleanupRegistro(t, "etapa_ensino", etapaSemTurmas.ID)
	if err := plataforma.Executar(ctx, CriarEtapaEnsino(&etapaSemTurmas)); err != nil {
		t.Fatalf("CriarEtapaEnsino() para exclusão retornou erro: %v", err)
	}
	if err := plataforma.Executar(ctx, ExcluirEtapaEnsino(etapaSemTurmas.ID)); err != nil {
		t.Fatalf("ExcluirEtapaEnsino() sem turmas retornou erro: %v", err)
	}

	for _, caso := range []struct {
		nome string
		fn   func() error
	}{
		{"inativar instituição", func() error { return plataforma.Executar(ctx, InativarInstituicaoEnsino(instituicao.ID)) }},
		{"inativar curso", func() error { return plataforma.Executar(ctx, InativarCurso(curso.ID)) }},
		{"inativar disciplina", func() error { return plataforma.Executar(ctx, InativarDisciplina(disciplina.ID)) }},
		{"encerrar turma", func() error { return plataforma.Executar(ctx, EncerrarTurma(turma.ID)) }},
	} {
		t.Run(caso.nome, func(t *testing.T) {
			if err := caso.fn(); err != nil {
				t.Fatalf("%s retornou erro: %v", caso.nome, err)
			}
		})
	}

	if _, err := ListarTodasInstituicoes(ctx); err != nil {
		t.Errorf("ListarTodasInstituicoes() retornou erro: %v", err)
	}
	if _, err := ListarTodosCursos(ctx); err != nil {
		t.Errorf("ListarTodosCursos() retornou erro: %v", err)
	}
	if _, err := ListarTodasEtapasEnsino(ctx); err != nil {
		t.Errorf("ListarTodasEtapasEnsino() retornou erro: %v", err)
	}
	if _, err := ListarTodasDisciplinas(ctx); err != nil {
		t.Errorf("ListarTodasDisciplinas() retornou erro: %v", err)
	}

	instituicaoFinal, err := BuscarInstituicaoEnsino(ctx, instituicao.ID)
	if err != nil || instituicaoFinal.Ativa {
		t.Errorf("instituição inativada = %#v, erro = %v", instituicaoFinal, err)
	}
	cursoFinal, err := BuscarCurso(ctx, curso.ID)
	if err != nil || cursoFinal.Ativo {
		t.Errorf("curso inativado = %#v, erro = %v", cursoFinal, err)
	}
	disciplinaFinal, err := BuscarDisciplina(ctx, disciplina.ID)
	if err != nil || disciplinaFinal.Ativa {
		t.Errorf("disciplina inativada = %#v, erro = %v", disciplinaFinal, err)
	}
	turmaFinal, err := BuscarTurmaAcademica(ctx, turma.ID)
	if err != nil || turmaFinal.Status != model.StatusTurmaEncerrada {
		t.Errorf("turma encerrada = %#v, erro = %v", turmaFinal, err)
	}
}

func TestValidacaoTurmaAcademicaRejeitaVinculosInconsistentes(t *testing.T) {
	requireDatabase(t)
	ctx := context.Background()
	err := plataforma.Executar(ctx, CriarTurmaAcademica(&model.Turma{
		InstituicaoID: instituicaoTesteID, CursoID: uuid.New(), EtapaID: uuid.New(),
		Nome: "Turma inválida", Codigo: fmt.Sprintf("INV-%s", uuid.NewString()[:8]),
		AnoLetivo: 2026, Turno: model.TurnoMatutino, Status: model.StatusTurmaPlanejada,
	}))
	if !errors.Is(err, domain.ErrRelacionamentoAcademico) {
		t.Fatalf("erro = %v; esperado vínculo acadêmico inconsistente", err)
	}
}
