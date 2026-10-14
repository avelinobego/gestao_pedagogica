package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/avelinobego/gestao_pedagogica/domain"
	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/google/uuid"
)

type fakeAcademicoUseCases struct {
	AcademicoUseCases
	instituicoes       []model.InstituicaoEnsino
	cursos             []model.Curso
	etapas             []model.EtapaEnsino
	createdInstitution *model.InstituicaoEnsino
	inactivationErr    error
	inactivationCalls  int
	deleteStageErr     error
}

func (f *fakeAcademicoUseCases) ListarTodasInstituicoes(context.Context) ([]model.InstituicaoEnsino, error) {
	return f.instituicoes, nil
}

func (f *fakeAcademicoUseCases) ListarTodosCursos(context.Context) ([]model.Curso, error) {
	return f.cursos, nil
}

func (f *fakeAcademicoUseCases) ListarTodasEtapasEnsino(context.Context) ([]model.EtapaEnsino, error) {
	return f.etapas, nil
}

func (f *fakeAcademicoUseCases) CriarInstituicaoEnsino(_ context.Context, instituicao *model.InstituicaoEnsino) error {
	copia := *instituicao
	f.createdInstitution = &copia
	return nil
}

func (f *fakeAcademicoUseCases) InativarInstituicaoEnsino(context.Context, uuid.UUID) error {
	f.inactivationCalls++
	return f.inactivationErr
}

func (f *fakeAcademicoUseCases) ExcluirEtapaEnsino(context.Context, uuid.UUID) error {
	return f.deleteStageErr
}

func novoHandlerAcademicoTeste(t *testing.T, academico AcademicoUseCases) http.Handler {
	t.Helper()
	handler, err := NewAcademicoHandler(academico)
	if err != nil {
		t.Fatalf("NewAcademicoHandler() retornou erro: %v", err)
	}
	return handler
}

func TestInicioAcademicoRedirecionaParaSecaoNaHome(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/academico", nil)
	res := httptest.NewRecorder()

	novoHandlerAcademicoTeste(t, &fakeAcademicoUseCases{}).ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusSeeOther)
	}
	if destino := res.Header().Get("Location"); destino != "/#academico" {
		t.Errorf("Location = %q; esperado %q", destino, "/#academico")
	}
}

func TestHandlerCompoePaginaInicialEModulosSemConflito(t *testing.T) {
	handler, err := NewHandler(&fakeMatriculaUseCases{}, &fakeAcademicoUseCases{})
	if err != nil {
		t.Fatalf("NewHandler() retornou erro: %v", err)
	}
	for _, rota := range []struct {
		caminho  string
		conteudo string
	}{
		{"/", "Organize a vida acadêmica em um só lugar"},
	} {
		req := httptest.NewRequest(http.MethodGet, rota.caminho, nil)
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		if res.Code != http.StatusOK {
			t.Errorf("%s: status = %d; esperado %d", rota.caminho, res.Code, http.StatusOK)
		}
		if !strings.Contains(res.Body.String(), rota.conteudo) {
			t.Errorf("%s: resposta não contém %q", rota.caminho, rota.conteudo)
		}
	}
	req := httptest.NewRequest(http.MethodGet, "/academico", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	if res.Code != http.StatusSeeOther || res.Header().Get("Location") != "/#academico" {
		t.Errorf("/academico: status = %d, Location = %q; esperado redirect para /#academico", res.Code, res.Header().Get("Location"))
	}
}

func TestListaInstituicoesRenderizaDadosENavegacao(t *testing.T) {
	useCases := &fakeAcademicoUseCases{instituicoes: []model.InstituicaoEnsino{
		{ID: uuid.New(), Nome: "Escola Horizonte", Cidade: "Campinas", UF: "SP", Ativa: true},
	}}
	req := httptest.NewRequest(http.MethodGet, "/academico/instituicoes", nil)
	res := httptest.NewRecorder()

	novoHandlerAcademicoTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusOK)
	}
	for _, texto := range []string{"Escola Horizonte", "Campinas/SP", "Ativa", "Editar", `href="/#academico"`} {
		if !strings.Contains(res.Body.String(), texto) {
			t.Errorf("resposta não contém %q", texto)
		}
	}
}

func TestNovoCadastroTurmaCarregaOpcoesDeReferencia(t *testing.T) {
	instituicaoID, cursoID, etapaID := uuid.New(), uuid.New(), uuid.New()
	useCases := &fakeAcademicoUseCases{
		instituicoes: []model.InstituicaoEnsino{{ID: instituicaoID, Nome: "Escola Horizonte", Ativa: true}},
		cursos:       []model.Curso{{ID: cursoID, Nome: "Ensino Fundamental", Ativo: true}},
		etapas:       []model.EtapaEnsino{{ID: etapaID, Nome: "6º ano"}},
	}
	req := httptest.NewRequest(http.MethodGet, "/academico/turmas/nova", nil)
	res := httptest.NewRecorder()

	novoHandlerAcademicoTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d; esperado %d: %s", res.Code, http.StatusOK, res.Body.String())
	}
	for _, texto := range []string{
		"Escola Horizonte", "Ensino Fundamental", "6º ano", `href="/#academico"`,
		instituicaoID.String(), cursoID.String(), etapaID.String(),
	} {
		if !strings.Contains(res.Body.String(), texto) {
			t.Errorf("formulário não contém %q", texto)
		}
	}
}

func TestCriarInstituicaoValidaFormularioERedireciona(t *testing.T) {
	useCases := &fakeAcademicoUseCases{}
	form := url.Values{
		"nome":   {"Escola Horizonte"},
		"cidade": {"Campinas"},
		"uf":     {"sp"},
		"ativa":  {"on"},
	}
	req := httptest.NewRequest(http.MethodPost, "/academico/instituicoes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()

	novoHandlerAcademicoTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther || res.Header().Get("Location") != "/academico/instituicoes" {
		t.Fatalf("status/location = %d/%q; esperado redirect para lista", res.Code, res.Header().Get("Location"))
	}
	if useCases.createdInstitution == nil || useCases.createdInstitution.UF != "SP" || !useCases.createdInstitution.Ativa {
		t.Fatalf("instituição criada não corresponde ao formulário: %#v", useCases.createdInstitution)
	}
}

func TestCriarInstituicaoRejeitaEntradaInvalidaSemPersistir(t *testing.T) {
	useCases := &fakeAcademicoUseCases{}
	form := url.Values{"nome": {"Escola Inválida"}, "uf": {"ABC"}}
	req := httptest.NewRequest(http.MethodPost, "/academico/instituicoes", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()

	novoHandlerAcademicoTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusBadRequest)
	}
	if useCases.createdInstitution != nil {
		t.Fatal("instituição inválida não deveria ser persistida")
	}
	if !strings.Contains(res.Body.String(), "UF deve conter duas letras maiúsculas") {
		t.Errorf("resposta não informa o erro de validação: %s", res.Body.String())
	}
}

func TestAcaoAcademicaCrossOriginEhBloqueada(t *testing.T) {
	useCases := &fakeAcademicoUseCases{}
	req := httptest.NewRequest(http.MethodPost, "/academico/instituicoes/"+uuid.NewString()+"/inativar", nil)
	req.Header.Set("Origin", "https://site-malicioso.example")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	res := httptest.NewRecorder()

	novoHandlerAcademicoTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusForbidden {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusForbidden)
	}
	if useCases.inactivationCalls != 0 {
		t.Fatal("requisição cross-origin não deveria alterar a instituição")
	}
}

func TestEtapaComTurmasRetornaConflito(t *testing.T) {
	useCases := &fakeAcademicoUseCases{deleteStageErr: domain.ErrEtapaPossuiTurmas}
	req := httptest.NewRequest(http.MethodPost, "/academico/etapas/"+uuid.NewString()+"/excluir", nil)
	res := httptest.NewRecorder()

	novoHandlerAcademicoTeste(t, useCases).ServeHTTP(res, req)

	if res.Code != http.StatusConflict {
		t.Fatalf("status = %d; esperado %d", res.Code, http.StatusConflict)
	}
	if !strings.Contains(res.Body.String(), "há turmas vinculadas") {
		t.Errorf("resposta não explica o conflito: %s", res.Body.String())
	}
}
