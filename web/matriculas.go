package web

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/avelinobego/gestao_pedagogica/domain"
	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/avelinobego/gestao_pedagogica/service"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"github.com/google/uuid"
)

//go:embed templates/*.html
var templateFiles embed.FS

type MatriculaUseCases interface {
	ListarMatriculas(context.Context) ([]model.Matricula, error)
	ListarAlunos(context.Context) ([]model.Aluno, error)
	ListarTurmas(context.Context) ([]model.Turma, error)
	CriarMatricula(context.Context, *model.Matricula) error
	CancelarMatricula(context.Context, uuid.UUID) error
	AtivarMatricula(context.Context, uuid.UUID) error
}

type MatriculaService struct{}

func (MatriculaService) ListarMatriculas(ctx context.Context) ([]model.Matricula, error) {
	return service.ListarTodasMatriculas(ctx)
}

func (MatriculaService) ListarAlunos(ctx context.Context) ([]model.Aluno, error) {
	return service.ListarTodosAlunos(ctx)
}

func (MatriculaService) ListarTurmas(ctx context.Context) ([]model.Turma, error) {
	return service.ListarTodasTurmas(ctx)
}

func (MatriculaService) CriarMatricula(ctx context.Context, matricula *model.Matricula) error {
	return plataforma.Executar(ctx, service.CriarMatricula(matricula))
}

func (MatriculaService) CancelarMatricula(ctx context.Context, id uuid.UUID) error {
	return plataforma.Executar(ctx, service.CancelarMatricula(id))
}

func (MatriculaService) AtivarMatricula(ctx context.Context, id uuid.UUID) error {
	return plataforma.Executar(ctx, service.AtivarMatricula(id))
}

type matriculasHandler struct {
	useCases  MatriculaUseCases
	templates *template.Template
}

type listPage struct {
	Matriculas []model.Matricula
}

type formValues struct {
	AlunoID       string
	TurmaID       string
	Numero        string
	DataMatricula string
	Observacao    string
}

type formPage struct {
	Alunos []model.Aluno
	Turmas []model.Turma
	Values formValues
	Erro   string
}

func NewMatriculasHandler(useCases MatriculaUseCases) (http.Handler, error) {
	if useCases == nil {
		return nil, errors.New("serviços de matrícula são obrigatórios")
	}
	templates, err := template.ParseFS(templateFiles,
		"templates/home.html",
		"templates/navigation.html",
		"templates/matriculas.html",
		"templates/matricula_form.html",
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar templates web: %w", err)
	}
	handler := &matriculasHandler{useCases: useCases, templates: templates}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handler.inicio)
	mux.HandleFunc("GET /matriculas", handler.listar)
	mux.HandleFunc("GET /matriculas/nova", handler.nova)
	mux.HandleFunc("POST /matriculas", handler.criar)
	mux.HandleFunc("POST /matriculas/{id}/cancelar", handler.cancelar)
	mux.HandleFunc("POST /matriculas/{id}/ativar", handler.ativar)
	return http.NewCrossOriginProtection().Handler(mux), nil
}

func (h *matriculasHandler) inicio(w http.ResponseWriter, r *http.Request) {
	h.render(w, http.StatusOK, "home.html", nil)
}

func (h *matriculasHandler) listar(w http.ResponseWriter, r *http.Request) {
	matriculas, err := h.useCases.ListarMatriculas(r.Context())
	if err != nil {
		slog.ErrorContext(r.Context(), "falha ao listar matrículas", "error", err)
		http.Error(w, "Não foi possível carregar as matrículas.", http.StatusInternalServerError)
		return
	}
	h.render(w, http.StatusOK, "matriculas.html", listPage{Matriculas: matriculas})
}

func (h *matriculasHandler) nova(w http.ResponseWriter, r *http.Request) {
	page, err := h.carregarFormulario(r.Context(), formValues{})
	if err != nil {
		slog.ErrorContext(r.Context(), "falha ao preparar formulário de matrícula", "error", err)
		http.Error(w, "Não foi possível carregar o formulário.", http.StatusInternalServerError)
		return
	}
	h.render(w, http.StatusOK, "matricula_form.html", page)
}

func (h *matriculasHandler) criar(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Formulário inválido.", http.StatusBadRequest)
		return
	}
	values := formValues{
		AlunoID:       r.FormValue("aluno_id"),
		TurmaID:       r.FormValue("turma_id"),
		Numero:        r.FormValue("numero"),
		DataMatricula: r.FormValue("data_matricula"),
		Observacao:    r.FormValue("observacao"),
	}
	matricula, err := parseMatriculaForm(values)
	if err != nil {
		h.renderFormularioComErro(w, r, values, err.Error())
		return
	}
	if err := domain.ValidarMatricula(matricula); err != nil {
		h.renderFormularioComErro(w, r, values, err.Error())
		return
	}
	if err := h.useCases.CriarMatricula(r.Context(), &matricula); err != nil {
		slog.ErrorContext(r.Context(), "falha ao criar matrícula", "error", err)
		http.Error(w, "Não foi possível criar a matrícula.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/matriculas", http.StatusSeeOther)
}

func (h *matriculasHandler) cancelar(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		http.Error(w, "Identificador de matrícula inválido.", http.StatusBadRequest)
		return
	}
	if err := h.useCases.CancelarMatricula(r.Context(), id); err != nil {
		slog.ErrorContext(r.Context(), "falha ao cancelar matrícula", "matricula_id", id, "error", err)
		http.Error(w, "Não foi possível cancelar a matrícula.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/matriculas", http.StatusSeeOther)
}

func (h *matriculasHandler) ativar(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		http.Error(w, "Identificador de matrícula inválido.", http.StatusBadRequest)
		return
	}
	if err := h.useCases.AtivarMatricula(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Matrícula não encontrada.", http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrMatriculaNaoAtivavel) {
			http.Error(w, "Somente matrículas canceladas ou trancadas podem ser ativadas.", http.StatusConflict)
			return
		}
		slog.ErrorContext(r.Context(), "falha ao ativar matrícula", "matricula_id", id, "error", err)
		http.Error(w, "Não foi possível ativar a matrícula.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/matriculas", http.StatusSeeOther)
}

func (h *matriculasHandler) carregarFormulario(ctx context.Context, values formValues) (formPage, error) {
	alunos, err := h.useCases.ListarAlunos(ctx)
	if err != nil {
		return formPage{}, fmt.Errorf("falha ao listar alunos: %w", err)
	}
	turmas, err := h.useCases.ListarTurmas(ctx)
	if err != nil {
		return formPage{}, fmt.Errorf("falha ao listar turmas: %w", err)
	}
	return formPage{Alunos: alunos, Turmas: turmas, Values: values}, nil
}

func (h *matriculasHandler) renderFormularioComErro(w http.ResponseWriter, r *http.Request, values formValues, mensagem string) {
	page, err := h.carregarFormulario(r.Context(), values)
	if err != nil {
		slog.ErrorContext(r.Context(), "falha ao recarregar formulário de matrícula", "error", err)
		http.Error(w, "Não foi possível carregar o formulário.", http.StatusInternalServerError)
		return
	}
	page.Erro = mensagem
	h.render(w, http.StatusBadRequest, "matricula_form.html", page)
}

func (h *matriculasHandler) render(w http.ResponseWriter, status int, nome string, dados any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := h.templates.ExecuteTemplate(w, nome, dados); err != nil {
		slog.Error("falha ao renderizar página", "template", nome, "error", err)
	}
}

func parseMatriculaForm(values formValues) (model.Matricula, error) {
	alunoID, err := uuid.Parse(values.AlunoID)
	if err != nil {
		return model.Matricula{}, errors.New("selecione um aluno válido")
	}
	turmaID, err := uuid.Parse(values.TurmaID)
	if err != nil {
		return model.Matricula{}, errors.New("selecione uma turma válida")
	}
	dataMatricula, err := time.Parse("2006-01-02", values.DataMatricula)
	if err != nil {
		return model.Matricula{}, errors.New("informe uma data de matrícula válida")
	}
	return model.Matricula{
		AlunoID:       alunoID,
		TurmaID:       turmaID,
		Numero:        values.Numero,
		DataMatricula: dataMatricula,
		Status:        model.StatusMatriculaAtiva,
		Observacao:    values.Observacao,
	}, nil
}
