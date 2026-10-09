package web

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/avelinobego/gestao_pedagogica/domain"
	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/google/uuid"
)

type fakeMatriculaUseCases struct {
	matriculas   []model.Matricula
	alunos       []model.Aluno
	turmas       []model.Turma
	listErr      error
	referenceErr error
	createErr    error
	cancelErr    error
	activateErr  error
	criada       *model.Matricula
	canceladaID  uuid.UUID
	ativadaID    uuid.UUID
}

func (f *fakeMatriculaUseCases) ListarMatriculas(context.Context) ([]model.Matricula, error) {
	return f.matriculas, f.listErr
}

func (f *fakeMatriculaUseCases) ListarAlunos(context.Context) ([]model.Aluno, error) {
	return f.alunos, f.referenceErr
}

func (f *fakeMatriculaUseCases) ListarTurmas(context.Context) ([]model.Turma, error) {
	return f.turmas, f.referenceErr
}

func (f *fakeMatriculaUseCases) CriarMatricula(_ context.Context, matricula *model.Matricula) error {
	if f.createErr != nil {
		return f.createErr
	}
	copia := *matricula
	f.criada = &copia
	return nil
}

func (f *fakeMatriculaUseCases) CancelarMatricula(_ context.Context, id uuid.UUID) error {
	f.canceladaID = id
	return f.cancelErr
}

func (f *fakeMatriculaUseCases) AtivarMatricula(_ context.Context, id uuid.UUID) error {
	f.ativadaID = id
	return f.activateErr
}

func novoHandlerTeste(t *testing.T, useCases *fakeMatriculaUseCases) http.Handler {
	t.Helper()
	handler, err := NewMatriculasHandler(useCases)
	if err != nil {
		t.Fatalf("NewMatriculasHandler() retornou erro: %v", err)
	}
	return handler
}

func TestListarMatriculasRenderizaHTML(t *testing.T) {
	id := uuid.New()
	useCases := &fakeMatriculaUseCases{
		matriculas: []model.Matricula{
			{
				ID:            id,
				AlunoNome:     "Aluno Exemplo",
				TurmaNome:     "Turma Exemplo",
				Numero:        "MAT-42",
				DataMatricula: time.Date(2026, time.January, 5, 0, 0, 0, 0, time.UTC),
				Status:        model.StatusMatriculaCancelada,
			},
			{
				ID:            uuid.New(),
				Numero:        "MAT-43",
				DataMatricula: time.Date(2026, time.January, 6, 0, 0, 0, 0, time.UTC),
				Status:        model.StatusMatriculaTrancada,
			},
			{
				ID:            uuid.New(),
				Numero:        "MAT-44",
				DataMatricula: time.Date(2026, time.January, 7, 0, 0, 0, 0, time.UTC),
				Status:        model.StatusMatriculaAtiva,
			},
			{
				ID:            uuid.New(),
				Numero:        "MAT-45",
				DataMatricula: time.Date(2026, time.January, 8, 0, 0, 0, 0, time.UTC),
				Status:        model.StatusMatriculaConcluida,
			},
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/matriculas", nil)
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusOK)
	}
	if contentType := res.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("Content-Type = %q; esperado HTML", contentType)
	}
	if !strings.Contains(res.Body.String(), "MAT-42") ||
		!strings.Contains(res.Body.String(), "Aluno Exemplo") ||
		!strings.Contains(res.Body.String(), "Turma Exemplo") ||
		!strings.Contains(res.Body.String(), id.String()) ||
		strings.Count(res.Body.String(), "Ativar Matrícula") != 2 ||
		strings.Count(res.Body.String(), "Cancelar</button>") != 1 {
		t.Errorf("resposta não contém dados da matrícula: %s", res.Body.String())
	}
}

func TestFormularioDeMatriculaCarregaAlunoETurma(t *testing.T) {
	alunoID := uuid.New()
	turmaID := uuid.New()
	useCases := &fakeMatriculaUseCases{
		alunos: []model.Aluno{{ID: alunoID, Nome: "Aluna Teste"}},
		turmas: []model.Turma{{ID: turmaID, Nome: "Turma A", AnoLetivo: 2026}},
	}
	req := httptest.NewRequest(http.MethodGet, "/matriculas/nova", nil)
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusOK)
	}
	for _, texto := range []string{alunoID.String(), "Aluna Teste", turmaID.String(), "Turma A"} {
		if !strings.Contains(res.Body.String(), texto) {
			t.Errorf("resposta não contém %q", texto)
		}
	}
}

func TestCriarMatriculaValidaFormularioERedireciona(t *testing.T) {
	alunoID := uuid.New()
	turmaID := uuid.New()
	useCases := &fakeMatriculaUseCases{}
	form := url.Values{
		"aluno_id":       {alunoID.String()},
		"turma_id":       {turmaID.String()},
		"numero":         {"MAT-43"},
		"data_matricula": {"2026-02-01"},
		"observacao":     {"Cadastro web"},
	}
	req := httptest.NewRequest(http.MethodPost, "/matriculas", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther || res.Header().Get("Location") != "/matriculas" {
		t.Fatalf("resposta = %d, Location %q; esperado redirect para /matriculas", res.Code, res.Header().Get("Location"))
	}
	if useCases.criada == nil {
		t.Fatal("o caso de uso de criação não foi chamado")
	}
	if useCases.criada.AlunoID != alunoID || useCases.criada.TurmaID != turmaID ||
		useCases.criada.Numero != "MAT-43" || useCases.criada.Status != model.StatusMatriculaAtiva {
		t.Errorf("matrícula criada com dados inesperados: %+v", useCases.criada)
	}
}

func TestCriarMatriculaInvalidaRetornaErroDeFormulario(t *testing.T) {
	useCases := &fakeMatriculaUseCases{}
	form := url.Values{
		"aluno_id":       {uuid.New().String()},
		"turma_id":       {uuid.New().String()},
		"numero":         {"MAT-44"},
		"data_matricula": {"data-invalida"},
	}
	req := httptest.NewRequest(http.MethodPost, "/matriculas", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusBadRequest)
	}
	if !strings.Contains(res.Body.String(), "informe uma data de matrícula válida") {
		t.Errorf("resposta não informa o erro de data: %s", res.Body.String())
	}
	if useCases.criada != nil {
		t.Error("o caso de uso não deveria ser chamado com data inválida")
	}
}

func TestCancelarMatriculaUsaPOSTERedireciona(t *testing.T) {
	id := uuid.New()
	useCases := &fakeMatriculaUseCases{}
	req := httptest.NewRequest(http.MethodPost, "/matriculas/"+id.String()+"/cancelar", nil)
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther || useCases.canceladaID != id {
		t.Fatalf("status = %d, matrícula cancelada = %s", res.Code, useCases.canceladaID)
	}
}

func TestAtivarMatriculaUsaPOSTERedireciona(t *testing.T) {
	id := uuid.New()
	useCases := &fakeMatriculaUseCases{}
	req := httptest.NewRequest(http.MethodPost, "/matriculas/"+id.String()+"/ativar", nil)
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther || res.Header().Get("Location") != "/matriculas" {
		t.Fatalf("status = %d, Location = %q; esperado redirect para /matriculas",
			res.Code, res.Header().Get("Location"))
	}
	if useCases.ativadaID != id {
		t.Errorf("matrícula ativada = %s; esperado %s", useCases.ativadaID, id)
	}
}

func TestAtivarMatriculaComIDInvalidoRetornaBadRequest(t *testing.T) {
	useCases := &fakeMatriculaUseCases{}
	req := httptest.NewRequest(http.MethodPost, "/matriculas/invalido/ativar", nil)
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusBadRequest)
	}
	if useCases.ativadaID != uuid.Nil {
		t.Error("o caso de uso não deveria ser chamado com ID inválido")
	}
}

func TestAtivarMatriculaInexistenteRetornaNotFound(t *testing.T) {
	useCases := &fakeMatriculaUseCases{activateErr: sql.ErrNoRows}
	req := httptest.NewRequest(http.MethodPost, "/matriculas/"+uuid.NewString()+"/ativar", nil)
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusNotFound)
	}
}

func TestAtivarMatriculaEstadoNaoPermitidoRetornaConflict(t *testing.T) {
	useCases := &fakeMatriculaUseCases{activateErr: domain.ErrMatriculaNaoAtivavel}
	req := httptest.NewRequest(http.MethodPost, "/matriculas/"+uuid.NewString()+"/ativar", nil)
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusConflict {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusConflict)
	}
}

func TestProtecaoCSRFRejeitaOrigemExterna(t *testing.T) {
	useCases := &fakeMatriculaUseCases{}
	form := url.Values{
		"aluno_id":       {uuid.New().String()},
		"turma_id":       {uuid.New().String()},
		"numero":         {"MAT-45"},
		"data_matricula": {"2026-02-01"},
	}
	req := httptest.NewRequest(http.MethodPost, "/matriculas", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://site-malicioso.example")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusForbidden)
	}
	if useCases.criada != nil {
		t.Error("requisição cross-site não deveria criar matrícula")
	}
}

func TestErroInternoNaoVazaDetalhesNaResposta(t *testing.T) {
	useCases := &fakeMatriculaUseCases{listErr: errors.New("erro interno detalhado")}
	req := httptest.NewRequest(http.MethodGet, "/matriculas", nil)
	res := httptest.NewRecorder()

	novoHandlerTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusInternalServerError)
	}
	if strings.Contains(res.Body.String(), "erro interno detalhado") {
		t.Error("resposta HTTP expôs detalhe interno")
	}
}
