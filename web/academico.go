package web

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/avelinobego/gestao_pedagogica/domain"
	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/avelinobego/gestao_pedagogica/service"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AcademicoUseCases interface {
	ListarTodasInstituicoes(context.Context) ([]model.InstituicaoEnsino, error)
	BuscarInstituicaoEnsino(context.Context, uuid.UUID) (model.InstituicaoEnsino, error)
	CriarInstituicaoEnsino(context.Context, *model.InstituicaoEnsino) error
	AtualizarInstituicaoEnsino(context.Context, model.InstituicaoEnsino) error
	InativarInstituicaoEnsino(context.Context, uuid.UUID) error
	ListarTodosCursos(context.Context) ([]model.Curso, error)
	BuscarCurso(context.Context, uuid.UUID) (model.Curso, error)
	CriarCurso(context.Context, *model.Curso) error
	AtualizarCurso(context.Context, model.Curso) error
	InativarCurso(context.Context, uuid.UUID) error
	ListarTodasEtapasEnsino(context.Context) ([]model.EtapaEnsino, error)
	BuscarEtapaEnsino(context.Context, uuid.UUID) (model.EtapaEnsino, error)
	CriarEtapaEnsino(context.Context, *model.EtapaEnsino) error
	AtualizarEtapaEnsino(context.Context, model.EtapaEnsino) error
	ExcluirEtapaEnsino(context.Context, uuid.UUID) error
	ListarTodasDisciplinas(context.Context) ([]model.Disciplina, error)
	BuscarDisciplina(context.Context, uuid.UUID) (model.Disciplina, error)
	CriarDisciplina(context.Context, *model.Disciplina) error
	AtualizarDisciplina(context.Context, model.Disciplina) error
	InativarDisciplina(context.Context, uuid.UUID) error
	ListarTodasTurmasAcademicas(context.Context) ([]model.Turma, error)
	BuscarTurmaAcademica(context.Context, uuid.UUID) (model.Turma, error)
	CriarTurmaAcademica(context.Context, *model.Turma) error
	AtualizarTurmaAcademica(context.Context, model.Turma) error
	EncerrarTurma(context.Context, uuid.UUID) error
}

type AcademicoService struct{}

func (AcademicoService) ListarTodasInstituicoes(ctx context.Context) ([]model.InstituicaoEnsino, error) {
	return service.ListarTodasInstituicoes(ctx)
}

func (AcademicoService) BuscarInstituicaoEnsino(ctx context.Context, id uuid.UUID) (model.InstituicaoEnsino, error) {
	return service.BuscarInstituicaoEnsino(ctx, id)
}

func (AcademicoService) CriarInstituicaoEnsino(ctx context.Context, instituicao *model.InstituicaoEnsino) error {
	return plataforma.Executar(ctx, service.CriarInstituicaoEnsino(instituicao))
}

func (AcademicoService) AtualizarInstituicaoEnsino(ctx context.Context, instituicao model.InstituicaoEnsino) error {
	return plataforma.Executar(ctx, service.AtualizarInstituicaoEnsino(instituicao))
}

func (AcademicoService) InativarInstituicaoEnsino(ctx context.Context, id uuid.UUID) error {
	return plataforma.Executar(ctx, service.InativarInstituicaoEnsino(id))
}

func (AcademicoService) ListarTodosCursos(ctx context.Context) ([]model.Curso, error) {
	return service.ListarTodosCursos(ctx)
}

func (AcademicoService) BuscarCurso(ctx context.Context, id uuid.UUID) (model.Curso, error) {
	return service.BuscarCurso(ctx, id)
}

func (AcademicoService) CriarCurso(ctx context.Context, curso *model.Curso) error {
	return plataforma.Executar(ctx, service.CriarCurso(curso))
}

func (AcademicoService) AtualizarCurso(ctx context.Context, curso model.Curso) error {
	return plataforma.Executar(ctx, service.AtualizarCurso(curso))
}

func (AcademicoService) InativarCurso(ctx context.Context, id uuid.UUID) error {
	return plataforma.Executar(ctx, service.InativarCurso(id))
}

func (AcademicoService) ListarTodasEtapasEnsino(ctx context.Context) ([]model.EtapaEnsino, error) {
	return service.ListarTodasEtapasEnsino(ctx)
}

func (AcademicoService) BuscarEtapaEnsino(ctx context.Context, id uuid.UUID) (model.EtapaEnsino, error) {
	return service.BuscarEtapaEnsino(ctx, id)
}

func (AcademicoService) CriarEtapaEnsino(ctx context.Context, etapa *model.EtapaEnsino) error {
	return plataforma.Executar(ctx, service.CriarEtapaEnsino(etapa))
}

func (AcademicoService) AtualizarEtapaEnsino(ctx context.Context, etapa model.EtapaEnsino) error {
	return plataforma.Executar(ctx, service.AtualizarEtapaEnsino(etapa))
}

func (AcademicoService) ExcluirEtapaEnsino(ctx context.Context, id uuid.UUID) error {
	return plataforma.Executar(ctx, service.ExcluirEtapaEnsino(id))
}

func (AcademicoService) ListarTodasDisciplinas(ctx context.Context) ([]model.Disciplina, error) {
	return service.ListarTodasDisciplinas(ctx)
}

func (AcademicoService) BuscarDisciplina(ctx context.Context, id uuid.UUID) (model.Disciplina, error) {
	return service.BuscarDisciplina(ctx, id)
}

func (AcademicoService) CriarDisciplina(ctx context.Context, disciplina *model.Disciplina) error {
	return plataforma.Executar(ctx, service.CriarDisciplina(disciplina))
}

func (AcademicoService) AtualizarDisciplina(ctx context.Context, disciplina model.Disciplina) error {
	return plataforma.Executar(ctx, service.AtualizarDisciplina(disciplina))
}

func (AcademicoService) InativarDisciplina(ctx context.Context, id uuid.UUID) error {
	return plataforma.Executar(ctx, service.InativarDisciplina(id))
}

func (AcademicoService) ListarTodasTurmasAcademicas(ctx context.Context) ([]model.Turma, error) {
	return service.ListarTodasTurmasAcademicas(ctx)
}

func (AcademicoService) BuscarTurmaAcademica(ctx context.Context, id uuid.UUID) (model.Turma, error) {
	return service.BuscarTurmaAcademica(ctx, id)
}

func (AcademicoService) CriarTurmaAcademica(ctx context.Context, turma *model.Turma) error {
	return plataforma.Executar(ctx, service.CriarTurmaAcademica(turma))
}

func (AcademicoService) AtualizarTurmaAcademica(ctx context.Context, turma model.Turma) error {
	return plataforma.Executar(ctx, service.AtualizarTurmaAcademica(turma))
}

func (AcademicoService) EncerrarTurma(ctx context.Context, id uuid.UUID) error {
	return plataforma.Executar(ctx, service.EncerrarTurma(id))
}

type academicoHandler struct {
	useCases  AcademicoUseCases
	templates *template.Template
}

type academicOption struct {
	Value string
	Label string
}

type academicField struct {
	Name     string
	Label    string
	Type     string
	Value    string
	Required bool
	Options  []academicOption
}

type academicFormPage struct {
	Title  string
	Action string
	Fields []academicField
	Erro   string
}

type academicRow struct {
	Cells       []string
	EditURL     string
	ActionURL   string
	ActionLabel string
}

type academicListPage struct {
	Title   string
	NewURL  string
	Headers []string
	Rows    []academicRow
}

func NewAcademicoHandler(useCases AcademicoUseCases) (http.Handler, error) {
	if useCases == nil {
		return nil, errors.New("serviços acadêmicos são obrigatórios")
	}
	templates, err := template.ParseFS(templateFiles,
		"templates/navigation.html",
		"templates/academico_list.html",
		"templates/academico_form.html",
	)
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar templates do módulo acadêmico: %w", err)
	}
	h := &academicoHandler{useCases: useCases, templates: templates}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /academico", h.inicio)
	mux.HandleFunc("GET /academico/{$}", h.inicio)
	h.registrarRotas(mux)
	return http.NewCrossOriginProtection().Handler(mux), nil
}

func (h *academicoHandler) registrarRotas(mux *http.ServeMux) {
	for _, rota := range []struct {
		caminho string
		tipo    string
	}{
		{"/academico/instituicoes", "instituicoes"},
		{"/academico/cursos", "cursos"},
		{"/academico/etapas", "etapas"},
		{"/academico/disciplinas", "disciplinas"},
		{"/academico/turmas", "turmas"},
	} {
		mux.HandleFunc("GET "+rota.caminho, h.listar(rota.tipo))
		mux.HandleFunc("GET "+rota.caminho+"/nova", h.novo(rota.tipo))
		mux.HandleFunc("POST "+rota.caminho, h.salvar(rota.tipo, false))
		mux.HandleFunc("GET "+rota.caminho+"/{id}/editar", h.editar(rota.tipo))
		mux.HandleFunc("POST "+rota.caminho+"/{id}", h.salvar(rota.tipo, true))
	}
	mux.HandleFunc("POST /academico/instituicoes/{id}/inativar", h.alterarEstado("instituicoes"))
	mux.HandleFunc("POST /academico/cursos/{id}/inativar", h.alterarEstado("cursos"))
	mux.HandleFunc("POST /academico/etapas/{id}/excluir", h.alterarEstado("etapas"))
	mux.HandleFunc("POST /academico/disciplinas/{id}/inativar", h.alterarEstado("disciplinas"))
	mux.HandleFunc("POST /academico/turmas/{id}/encerrar", h.alterarEstado("turmas"))
}

func (h *academicoHandler) inicio(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/#academico", http.StatusSeeOther)
}

func (h *academicoHandler) listar(tipo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagina, err := h.montarLista(r, tipo)
		if err != nil {
			slog.ErrorContext(r.Context(), "falha ao listar cadastros acadêmicos", "tipo", tipo, "error", err)
			http.Error(w, "Não foi possível carregar os cadastros acadêmicos.", http.StatusInternalServerError)
			return
		}
		h.render(w, http.StatusOK, "academico_list.html", pagina)
	}
}

func (h *academicoHandler) montarLista(r *http.Request, tipo string) (academicListPage, error) {
	base := "/academico/" + tipo
	pagina := academicListPage{NewURL: base + "/nova", Rows: []academicRow{}}
	switch tipo {
	case "instituicoes":
		pagina.Title, pagina.Headers = "Instituições de ensino", []string{"Nome", "Nome fantasia", "Cidade/UF", "Status"}
		items, err := h.useCases.ListarTodasInstituicoes(r.Context())
		if err != nil {
			return pagina, err
		}
		for _, item := range items {
			estado, acao := "Inativa", ""
			if item.Ativa {
				estado, acao = "Ativa", "inativar"
			}
			pagina.Rows = append(pagina.Rows, academicRow{
				Cells:     []string{item.Nome, item.NomeFantasia, localidade(item.Cidade, item.UF), estado},
				EditURL:   fmt.Sprintf("%s/%s/editar", base, item.ID),
				ActionURL: actionURL(base, item.ID, acao), ActionLabel: "Inativar",
			})
		}
	case "cursos":
		pagina.Title, pagina.Headers = "Cursos", []string{"Nome", "Código", "Instituição", "Status"}
		items, err := h.useCases.ListarTodosCursos(r.Context())
		if err != nil {
			return pagina, err
		}
		for _, item := range items {
			estado, acao := "Inativo", ""
			if item.Ativo {
				estado, acao = "Ativo", "inativar"
			}
			pagina.Rows = append(pagina.Rows, academicRow{
				Cells:     []string{item.Nome, item.Codigo, item.InstituicaoNome, estado},
				EditURL:   fmt.Sprintf("%s/%s/editar", base, item.ID),
				ActionURL: actionURL(base, item.ID, acao), ActionLabel: "Inativar",
			})
		}
	case "etapas":
		pagina.Title, pagina.Headers = "Etapas de ensino", []string{"Nome", "Curso", "Ordem", "Carga horária"}
		items, err := h.useCases.ListarTodasEtapasEnsino(r.Context())
		if err != nil {
			return pagina, err
		}
		for _, item := range items {
			pagina.Rows = append(pagina.Rows, academicRow{
				Cells:     []string{item.Nome, item.CursoNome, strconv.Itoa(item.Ordem), strconv.Itoa(item.CargaHoraria)},
				EditURL:   fmt.Sprintf("%s/%s/editar", base, item.ID),
				ActionURL: actionURL(base, item.ID, "excluir"), ActionLabel: "Excluir",
			})
		}
	case "disciplinas":
		pagina.Title, pagina.Headers = "Disciplinas", []string{"Nome", "Código", "Instituição", "Carga horária", "Status"}
		items, err := h.useCases.ListarTodasDisciplinas(r.Context())
		if err != nil {
			return pagina, err
		}
		for _, item := range items {
			estado, acao := "Inativa", ""
			if item.Ativa {
				estado, acao = "Ativa", "inativar"
			}
			pagina.Rows = append(pagina.Rows, academicRow{
				Cells:     []string{item.Nome, item.Codigo, item.InstituicaoNome, strconv.Itoa(item.CargaHoraria), estado},
				EditURL:   fmt.Sprintf("%s/%s/editar", base, item.ID),
				ActionURL: actionURL(base, item.ID, acao), ActionLabel: "Inativar",
			})
		}
	case "turmas":
		pagina.Title, pagina.Headers = "Turmas", []string{"Nome", "Código", "Curso / etapa", "Ano letivo", "Turno", "Status"}
		items, err := h.useCases.ListarTodasTurmasAcademicas(r.Context())
		if err != nil {
			return pagina, err
		}
		for _, item := range items {
			acao := ""
			if item.Status != model.StatusTurmaEncerrada {
				acao = "encerrar"
			}
			pagina.Rows = append(pagina.Rows, academicRow{
				Cells: []string{item.Nome, item.Codigo, item.CursoNome + " / " + item.EtapaNome,
					strconv.Itoa(item.AnoLetivo), string(item.Turno), string(item.Status)},
				EditURL:   fmt.Sprintf("%s/%s/editar", base, item.ID),
				ActionURL: actionURL(base, item.ID, acao), ActionLabel: "Encerrar",
			})
		}
	default:
		return pagina, errors.New("tipo de cadastro acadêmico inválido")
	}
	return pagina, nil
}

func (h *academicoHandler) novo(tipo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagina, err := h.montarFormulario(r, tipo, uuid.Nil, nil)
		if err != nil {
			h.erroFormulario(r, w, err)
			return
		}
		h.render(w, http.StatusOK, "academico_form.html", pagina)
	}
}

func (h *academicoHandler) editar(tipo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil || id == uuid.Nil {
			http.Error(w, "Identificador inválido.", http.StatusBadRequest)
			return
		}
		valores, err := h.buscarRegistro(r, tipo, id)
		if err != nil {
			h.responderErroRegistro(r, w, err)
			return
		}
		pagina, err := h.montarFormulario(r, tipo, id, valores)
		if err != nil {
			h.erroFormulario(r, w, err)
			return
		}
		h.render(w, http.StatusOK, "academico_form.html", pagina)
	}
}

func (h *academicoHandler) salvar(tipo string, atualizacao bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Formulário inválido.", http.StatusBadRequest)
			return
		}
		id := uuid.Nil
		if atualizacao {
			var err error
			id, err = uuid.Parse(r.PathValue("id"))
			if err != nil || id == uuid.Nil {
				http.Error(w, "Identificador inválido.", http.StatusBadRequest)
				return
			}
		}
		registro, err := h.parseRegistro(tipo, id, r)
		if err != nil {
			h.alterarFormulario(w, r, tipo, id, http.StatusBadRequest, err.Error())
			return
		}
		err = h.persistir(r, tipo, atualizacao, registro)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Cadastro não encontrado.", http.StatusNotFound)
				return
			}
			if errors.Is(err, domain.ErrRelacionamentoAcademico) {
				h.alterarFormulario(w, r, tipo, id, http.StatusConflict,
					"Não foi possível salvar: verifique se a instituição, o curso e a etapa pertencem aos mesmos cadastros.")
				return
			}
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && (pqErr.Code == "23505" || pqErr.Code == "23503") {
				pagina, formErr := h.montarFormulario(r, tipo, id, valuesFromRequest(r))
				if formErr != nil {
					h.erroFormulario(r, w, formErr)
					return
				}
				pagina.Erro = "Não foi possível salvar: verifique os vínculos e se já existe um cadastro com esse código."
				h.render(w, http.StatusConflict, "academico_form.html", pagina)
				return
			}
			slog.ErrorContext(r.Context(), "falha ao salvar cadastro acadêmico", "tipo", tipo, "error", err)
			http.Error(w, "Não foi possível salvar o cadastro acadêmico.", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/academico/"+tipo, http.StatusSeeOther)
	}
}

func (h *academicoHandler) parseRegistro(tipo string, id uuid.UUID, r *http.Request) (any, error) {
	get := func(nome string) string { return strings.TrimSpace(r.FormValue(nome)) }
	switch tipo {
	case "instituicoes":
		item := model.InstituicaoEnsino{
			ID: id, Nome: get("nome"), NomeFantasia: get("nome_fantasia"), CNPJ: get("cnpj"),
			CodigoINEP: get("codigo_inep"), Email: get("email"), Telefone: get("telefone"),
			Endereco: get("endereco"), Cidade: get("cidade"), UF: strings.ToUpper(get("uf")),
			Ativa: r.FormValue("ativa") == "on",
		}
		return item, domain.ValidarInstituicaoEnsino(item)
	case "cursos":
		instituicaoID, err := parseAcademicUUID("instituicao_id", get("instituicao_id"))
		if err != nil {
			return nil, err
		}
		item := model.Curso{ID: id, InstituicaoID: instituicaoID, Nome: get("nome"),
			Codigo: get("codigo"), Descricao: get("descricao"), Ativo: r.FormValue("ativo") == "on"}
		return item, domain.ValidarCurso(item)
	case "etapas":
		cursoID, err := parseAcademicUUID("curso_id", get("curso_id"))
		if err != nil {
			return nil, err
		}
		ordem, err := parseAcademicInt("ordem", get("ordem"))
		if err != nil {
			return nil, err
		}
		carga, err := parseAcademicInt("carga_horaria", get("carga_horaria"))
		if err != nil {
			return nil, err
		}
		item := model.EtapaEnsino{ID: id, CursoID: cursoID, Nome: get("nome"), Ordem: ordem, CargaHoraria: carga}
		return item, domain.ValidarEtapaEnsino(item)
	case "disciplinas":
		instituicaoID, err := parseAcademicUUID("instituicao_id", get("instituicao_id"))
		if err != nil {
			return nil, err
		}
		carga, err := parseAcademicInt("carga_horaria", get("carga_horaria"))
		if err != nil {
			return nil, err
		}
		item := model.Disciplina{ID: id, InstituicaoID: instituicaoID, Nome: get("nome"),
			Codigo: get("codigo"), CargaHoraria: carga, Ativa: r.FormValue("ativa") == "on"}
		return item, domain.ValidarDisciplina(item)
	case "turmas":
		instituicaoID, err := parseAcademicUUID("instituicao_id", get("instituicao_id"))
		if err != nil {
			return nil, err
		}
		cursoID, err := parseAcademicUUID("curso_id", get("curso_id"))
		if err != nil {
			return nil, err
		}
		etapaID, err := parseAcademicUUID("etapa_id", get("etapa_id"))
		if err != nil {
			return nil, err
		}
		ano, err := parseAcademicInt("ano_letivo", get("ano_letivo"))
		if err != nil {
			return nil, err
		}
		item := model.Turma{
			ID: id, InstituicaoID: instituicaoID, CursoID: cursoID, EtapaID: etapaID,
			Nome: get("nome"), Codigo: get("codigo"), AnoLetivo: ano,
			Turno: model.Turno(get("turno")), Status: model.StatusTurma(get("status")),
		}
		if item.Status == "" && id == uuid.Nil {
			item.Status = model.StatusTurmaPlanejada
		}
		return item, domain.ValidarTurma(item)
	default:
		return nil, errors.New("tipo de cadastro acadêmico inválido")
	}
}

func (h *academicoHandler) persistir(r *http.Request, tipo string, atualizacao bool, registro any) error {
	ctx := r.Context()
	switch item := registro.(type) {
	case model.InstituicaoEnsino:
		if atualizacao {
			return h.useCases.AtualizarInstituicaoEnsino(ctx, item)
		}
		return h.useCases.CriarInstituicaoEnsino(ctx, &item)
	case model.Curso:
		if atualizacao {
			return h.useCases.AtualizarCurso(ctx, item)
		}
		return h.useCases.CriarCurso(ctx, &item)
	case model.EtapaEnsino:
		if atualizacao {
			return h.useCases.AtualizarEtapaEnsino(ctx, item)
		}
		return h.useCases.CriarEtapaEnsino(ctx, &item)
	case model.Disciplina:
		if atualizacao {
			return h.useCases.AtualizarDisciplina(ctx, item)
		}
		return h.useCases.CriarDisciplina(ctx, &item)
	case model.Turma:
		if atualizacao {
			return h.useCases.AtualizarTurmaAcademica(ctx, item)
		}
		return h.useCases.CriarTurmaAcademica(ctx, &item)
	default:
		return fmt.Errorf("tipo de cadastro acadêmico inválido: %s", tipo)
	}
}

func (h *academicoHandler) alterarEstado(tipo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil || id == uuid.Nil {
			http.Error(w, "Identificador inválido.", http.StatusBadRequest)
			return
		}
		switch tipo {
		case "instituicoes":
			err = h.useCases.InativarInstituicaoEnsino(r.Context(), id)
		case "cursos":
			err = h.useCases.InativarCurso(r.Context(), id)
		case "etapas":
			err = h.useCases.ExcluirEtapaEnsino(r.Context(), id)
		case "disciplinas":
			err = h.useCases.InativarDisciplina(r.Context(), id)
		case "turmas":
			err = h.useCases.EncerrarTurma(r.Context(), id)
		default:
			http.Error(w, "Operação acadêmica inválida.", http.StatusBadRequest)
			return
		}
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Cadastro não encontrado.", http.StatusNotFound)
				return
			}
			if errors.Is(err, domain.ErrEtapaPossuiTurmas) {
				http.Error(w, "A etapa não pode ser excluída porque há turmas vinculadas.", http.StatusConflict)
				return
			}
			var pqErr *pq.Error
			if tipo == "etapas" && errors.As(err, &pqErr) && pqErr.Code == "23503" {
				http.Error(w, "A etapa não pode ser excluída porque há turmas vinculadas.", http.StatusConflict)
				return
			}
			slog.ErrorContext(r.Context(), "falha ao alterar estado de cadastro acadêmico",
				"tipo", tipo, "id", id, "error", err)
			http.Error(w, "Não foi possível alterar o cadastro acadêmico.", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/academico/"+tipo, http.StatusSeeOther)
	}
}

func (h *academicoHandler) buscarRegistro(r *http.Request, tipo string, id uuid.UUID) (map[string]string, error) {
	switch tipo {
	case "instituicoes":
		item, err := h.useCases.BuscarInstituicaoEnsino(r.Context(), id)
		return map[string]string{
			"nome": item.Nome, "nome_fantasia": item.NomeFantasia, "cnpj": item.CNPJ,
			"codigo_inep": item.CodigoINEP, "email": item.Email, "telefone": item.Telefone,
			"endereco": item.Endereco, "cidade": item.Cidade, "uf": item.UF, "ativa": boolString(item.Ativa),
		}, err
	case "cursos":
		item, err := h.useCases.BuscarCurso(r.Context(), id)
		return map[string]string{"instituicao_id": item.InstituicaoID.String(), "nome": item.Nome,
			"codigo": item.Codigo, "descricao": item.Descricao, "ativo": boolString(item.Ativo)}, err
	case "etapas":
		item, err := h.useCases.BuscarEtapaEnsino(r.Context(), id)
		return map[string]string{"curso_id": item.CursoID.String(), "nome": item.Nome,
			"ordem": strconv.Itoa(item.Ordem), "carga_horaria": strconv.Itoa(item.CargaHoraria)}, err
	case "disciplinas":
		item, err := h.useCases.BuscarDisciplina(r.Context(), id)
		return map[string]string{"instituicao_id": item.InstituicaoID.String(), "nome": item.Nome,
			"codigo": item.Codigo, "carga_horaria": strconv.Itoa(item.CargaHoraria),
			"ativa": boolString(item.Ativa)}, err
	case "turmas":
		item, err := h.useCases.BuscarTurmaAcademica(r.Context(), id)
		return map[string]string{"instituicao_id": item.InstituicaoID.String(), "curso_id": item.CursoID.String(),
			"etapa_id": item.EtapaID.String(), "nome": item.Nome, "codigo": item.Codigo,
			"ano_letivo": strconv.Itoa(item.AnoLetivo), "turno": string(item.Turno), "status": string(item.Status)}, err
	default:
		return nil, errors.New("tipo de cadastro acadêmico inválido")
	}
}

func (h *academicoHandler) montarFormulario(r *http.Request, tipo string, id uuid.UUID, valores map[string]string) (academicFormPage, error) {
	base := "/academico/" + tipo
	pagina := academicFormPage{Action: base, Fields: []academicField{}}
	if id != uuid.Nil {
		pagina.Action = base + "/" + id.String()
	}
	get := func(nome, padrao string) string {
		if valor, existe := valores[nome]; existe {
			return valor
		}
		return padrao
	}
	optionList := func(tipo string) ([]academicOption, error) {
		switch tipo {
		case "instituicoes":
			itens, err := h.useCases.ListarTodasInstituicoes(r.Context())
			if err != nil {
				return nil, err
			}
			opcoes := make([]academicOption, 0, len(itens))
			for _, item := range itens {
				opcoes = append(opcoes, academicOption{item.ID.String(), item.Nome})
			}
			return opcoes, nil
		case "cursos":
			itens, err := h.useCases.ListarTodosCursos(r.Context())
			if err != nil {
				return nil, err
			}
			opcoes := make([]academicOption, 0, len(itens))
			for _, item := range itens {
				opcoes = append(opcoes, academicOption{item.ID.String(), item.Nome})
			}
			return opcoes, nil
		case "etapas":
			itens, err := h.useCases.ListarTodasEtapasEnsino(r.Context())
			if err != nil {
				return nil, err
			}
			opcoes := make([]academicOption, 0, len(itens))
			for _, item := range itens {
				opcoes = append(opcoes, academicOption{item.ID.String(), item.Nome})
			}
			return opcoes, nil
		default:
			return nil, nil
		}
	}
	addSelect := func(name, label, entity, current string) error {
		opcoes, err := optionList(entity)
		if err != nil {
			return err
		}
		pagina.Fields = append(pagina.Fields, academicField{
			Name: name, Label: label, Type: "select", Value: current, Required: true, Options: opcoes,
		})
		return nil
	}
	add := func(name, label, kind, def string, required bool) {
		pagina.Fields = append(pagina.Fields, academicField{
			Name: name, Label: label, Type: kind, Value: get(name, def), Required: required,
		})
	}
	switch tipo {
	case "instituicoes":
		pagina.Title = "Instituição de ensino"
		add("nome", "Nome", "text", "", true)
		add("nome_fantasia", "Nome fantasia", "text", "", false)
		add("cnpj", "CNPJ", "text", "", false)
		add("codigo_inep", "Código INEP", "text", "", false)
		add("email", "E-mail", "email", "", false)
		add("telefone", "Telefone", "text", "", false)
		add("endereco", "Endereço", "text", "", false)
		add("cidade", "Cidade", "text", "", false)
		add("uf", "UF", "text", "", false)
		add("ativa", "Instituição ativa", "checkbox", "true", false)
	case "cursos":
		pagina.Title = "Curso"
		if err := addSelect("instituicao_id", "Instituição", "instituicoes", get("instituicao_id", "")); err != nil {
			return pagina, err
		}
		add("nome", "Nome", "text", "", true)
		add("codigo", "Código", "text", "", true)
		add("descricao", "Descrição", "textarea", "", false)
		add("ativo", "Curso ativo", "checkbox", "true", false)
	case "etapas":
		pagina.Title = "Etapa de ensino"
		if err := addSelect("curso_id", "Curso", "cursos", get("curso_id", "")); err != nil {
			return pagina, err
		}
		add("nome", "Nome", "text", "", true)
		add("ordem", "Ordem", "number", "", true)
		add("carga_horaria", "Carga horária", "number", "", true)
	case "disciplinas":
		pagina.Title = "Disciplina"
		if err := addSelect("instituicao_id", "Instituição", "instituicoes", get("instituicao_id", "")); err != nil {
			return pagina, err
		}
		add("nome", "Nome", "text", "", true)
		add("codigo", "Código", "text", "", true)
		add("carga_horaria", "Carga horária", "number", "", true)
		add("ativa", "Disciplina ativa", "checkbox", "true", false)
	case "turmas":
		pagina.Title = "Turma"
		if err := addSelect("instituicao_id", "Instituição", "instituicoes", get("instituicao_id", "")); err != nil {
			return pagina, err
		}
		if err := addSelect("curso_id", "Curso", "cursos", get("curso_id", "")); err != nil {
			return pagina, err
		}
		if err := addSelect("etapa_id", "Etapa", "etapas", get("etapa_id", "")); err != nil {
			return pagina, err
		}
		add("nome", "Nome", "text", "", true)
		add("codigo", "Código", "text", "", true)
		add("ano_letivo", "Ano letivo", "number", "", true)
		pagina.Fields = append(pagina.Fields, academicField{Name: "turno", Label: "Turno", Type: "select", Value: get("turno", ""),
			Required: true, Options: enumOptions(string(model.TurnoMatutino), string(model.TurnoVespertino), string(model.TurnoNoturno), string(model.TurnoIntegral))})
		pagina.Fields = append(pagina.Fields, academicField{Name: "status", Label: "Status", Type: "select", Value: get("status", string(model.StatusTurmaPlanejada)),
			Required: true, Options: enumOptions(string(model.StatusTurmaPlanejada), string(model.StatusTurmaAtiva), string(model.StatusTurmaEncerrada))})
	default:
		return pagina, errors.New("tipo de cadastro acadêmico inválido")
	}
	return pagina, nil
}

func (h *academicoHandler) alterarFormulario(w http.ResponseWriter, r *http.Request, tipo string, id uuid.UUID, status int, mensagem string) {
	pagina, err := h.montarFormulario(r, tipo, id, valuesFromRequest(r))
	if err != nil {
		h.erroFormulario(r, w, err)
		return
	}
	pagina.Erro = mensagem
	h.render(w, status, "academico_form.html", pagina)
}

func (h *academicoHandler) erroFormulario(r *http.Request, w http.ResponseWriter, err error) {
	slog.ErrorContext(r.Context(), "falha ao preparar formulário acadêmico", "error", err)
	http.Error(w, "Não foi possível carregar o formulário acadêmico.", http.StatusInternalServerError)
}

func (h *academicoHandler) responderErroRegistro(r *http.Request, w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Cadastro não encontrado.", http.StatusNotFound)
		return
	}
	slog.ErrorContext(r.Context(), "falha ao buscar cadastro acadêmico", "error", err)
	http.Error(w, "Não foi possível carregar o cadastro acadêmico.", http.StatusInternalServerError)
}

func (h *academicoHandler) render(w http.ResponseWriter, status int, nome string, dados any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := h.templates.ExecuteTemplate(w, nome, dados); err != nil {
		slog.Error("falha ao renderizar página acadêmica", "template", nome, "error", err)
	}
}

func parseAcademicUUID(campo, valor string) (uuid.UUID, error) {
	id, err := uuid.Parse(valor)
	if err != nil || id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("selecione %s válido(a)", campo)
	}
	return id, nil
}

func parseAcademicInt(campo, valor string) (int, error) {
	numero, err := strconv.Atoi(valor)
	if err != nil {
		return 0, fmt.Errorf("%s deve ser um número inteiro", campo)
	}
	return numero, nil
}

func valuesFromRequest(r *http.Request) map[string]string {
	valores := make(map[string]string)
	for nome, values := range r.PostForm {
		if len(values) > 0 {
			valores[nome] = values[0]
			if values[0] == "on" {
				valores[nome] = "true"
			}
		}
	}
	return valores
}

func enumOptions(valores ...string) []academicOption {
	opcoes := make([]academicOption, 0, len(valores))
	for _, valor := range valores {
		opcoes = append(opcoes, academicOption{Value: valor, Label: valor})
	}
	return opcoes
}

func actionURL(base string, id uuid.UUID, action string) string {
	if action == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s", base, id, action)
}

func localidade(cidade, uf string) string {
	if cidade == "" {
		return uf
	}
	if uf == "" {
		return cidade
	}
	return cidade + "/" + uf
}

func boolString(valor bool) string {
	if valor {
		return "true"
	}
	return "false"
}
