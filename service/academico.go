package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/avelinobego/gestao_pedagogica/domain"
	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const selectInstituicaoEnsino = `
	SELECT id, nome, nome_fantasia, COALESCE(cnpj, '') AS cnpj,
		COALESCE(codigo_inep, '') AS codigo_inep, email, telefone, endereco,
		cidade, uf, ativa, criado_em, atualizado_em
	FROM instituicao_ensino`

const selectCurso = `
	SELECT c.id, c.instituicao_id, i.nome AS instituicao_nome, c.nome, c.codigo,
		c.descricao, c.ativo
	FROM curso AS c
	JOIN instituicao_ensino AS i ON i.id = c.instituicao_id`

const selectEtapaEnsino = `
	SELECT e.id, e.curso_id, c.nome AS curso_nome, e.nome, e.ordem, e.carga_horaria
	FROM etapa_ensino AS e
	JOIN curso AS c ON c.id = e.curso_id`

const selectDisciplina = `
	SELECT d.id, d.instituicao_id, i.nome AS instituicao_nome, d.nome, d.codigo,
		d.carga_horaria, d.ativa
	FROM disciplina AS d
	JOIN instituicao_ensino AS i ON i.id = d.instituicao_id`

const selectTurmaAcademica = `
	SELECT t.id, t.instituicao_id, i.nome AS instituicao_nome, t.curso_id,
		c.nome AS curso_nome, t.etapa_id, e.nome AS etapa_nome, t.nome, t.codigo,
		t.ano_letivo, t.turno, t.status, t.criado_em
	FROM turma AS t
	JOIN instituicao_ensino AS i ON i.id = t.instituicao_id
	JOIN curso AS c ON c.id = t.curso_id
	JOIN etapa_ensino AS e ON e.id = t.etapa_id`

func ListarTodasInstituicoes(ctx context.Context) ([]model.InstituicaoEnsino, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return nil, err
	}
	instituicoes := make([]model.InstituicaoEnsino, 0)
	if err := db.SelectContext(ctx, &instituicoes, selectInstituicaoEnsino+" ORDER BY nome ASC"); err != nil {
		return nil, fmt.Errorf("falha ao listar instituições: %w", err)
	}
	return instituicoes, nil
}

func BuscarInstituicaoEnsino(ctx context.Context, id uuid.UUID) (model.InstituicaoEnsino, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return model.InstituicaoEnsino{}, err
	}
	var instituicao model.InstituicaoEnsino
	if err := db.GetContext(ctx, &instituicao, selectInstituicaoEnsino+" WHERE id = $1", id); err != nil {
		return model.InstituicaoEnsino{}, fmt.Errorf("falha ao buscar instituição %s: %w", id, err)
	}
	return instituicao, nil
}

func CriarInstituicaoEnsino(instituicao *model.InstituicaoEnsino) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if instituicao == nil {
			return errors.New("instituição não pode ser nil")
		}
		if err := domain.ValidarInstituicaoEnsino(*instituicao); err != nil {
			return fmt.Errorf("instituição inválida: %w", err)
		}
		err := tx.QueryRowxContext(ctx, `
			INSERT INTO instituicao_ensino (
				id, nome, nome_fantasia, cnpj, codigo_inep, email, telefone,
				endereco, cidade, uf, ativa
			) VALUES (
				COALESCE(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()),
				$2, $3, NULLIF($4, ''), NULLIF($5, ''), $6, $7, $8, $9, $10, $11
			)
			RETURNING id, criado_em, atualizado_em`,
			instituicao.ID, instituicao.Nome, instituicao.NomeFantasia, instituicao.CNPJ,
			instituicao.CodigoINEP, instituicao.Email, instituicao.Telefone,
			instituicao.Endereco, instituicao.Cidade, instituicao.UF, instituicao.Ativa,
		).Scan(&instituicao.ID, &instituicao.CriadoEm, &instituicao.AtualizadoEm)
		if err != nil {
			return fmt.Errorf("falha ao criar instituição: %w", err)
		}
		return nil
	}
}

func AtualizarInstituicaoEnsino(instituicao model.InstituicaoEnsino) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if instituicao.ID == uuid.Nil {
			return errors.New("identificador da instituição é obrigatório")
		}
		if err := domain.ValidarInstituicaoEnsino(instituicao); err != nil {
			return fmt.Errorf("instituição inválida: %w", err)
		}
		result, err := tx.ExecContext(ctx, `
			UPDATE instituicao_ensino
			SET nome = $2, nome_fantasia = $3, cnpj = NULLIF($4, ''),
				codigo_inep = NULLIF($5, ''), email = $6, telefone = $7,
				endereco = $8, cidade = $9, uf = $10, ativa = $11,
				atualizado_em = CURRENT_TIMESTAMP
			WHERE id = $1`,
			instituicao.ID, instituicao.Nome, instituicao.NomeFantasia, instituicao.CNPJ,
			instituicao.CodigoINEP, instituicao.Email, instituicao.Telefone,
			instituicao.Endereco, instituicao.Cidade, instituicao.UF, instituicao.Ativa,
		)
		return verificarLinhasAfetadasTx(result, err, "instituição")
	}
}

func InativarInstituicaoEnsino(id uuid.UUID) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `UPDATE instituicao_ensino SET ativa = FALSE, atualizado_em = CURRENT_TIMESTAMP WHERE id = $1`, id)
		return verificarLinhasAfetadasTx(result, err, "instituição")
	}
}

func ListarTodosCursos(ctx context.Context) ([]model.Curso, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return nil, err
	}
	cursos := make([]model.Curso, 0)
	if err := db.SelectContext(ctx, &cursos, selectCurso+" ORDER BY i.nome, c.nome"); err != nil {
		return nil, fmt.Errorf("falha ao listar cursos: %w", err)
	}
	return cursos, nil
}

func BuscarCurso(ctx context.Context, id uuid.UUID) (model.Curso, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return model.Curso{}, err
	}
	var curso model.Curso
	if err := db.GetContext(ctx, &curso, selectCurso+" WHERE c.id = $1", id); err != nil {
		return model.Curso{}, fmt.Errorf("falha ao buscar curso %s: %w", id, err)
	}
	return curso, nil
}

func CriarCurso(curso *model.Curso) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if curso == nil {
			return errors.New("curso não pode ser nil")
		}
		if err := domain.ValidarCurso(*curso); err != nil {
			return fmt.Errorf("curso inválido: %w", err)
		}
		err := tx.QueryRowxContext(ctx, `
			INSERT INTO curso (id, instituicao_id, nome, codigo, descricao, ativo)
			VALUES (
				COALESCE(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()),
				$2, $3, $4, $5, $6
			)
			RETURNING id`,
			curso.ID, curso.InstituicaoID, curso.Nome, curso.Codigo, curso.Descricao, curso.Ativo,
		).Scan(&curso.ID)
		if err != nil {
			return fmt.Errorf("falha ao criar curso: %w", err)
		}
		return nil
	}
}

func AtualizarCurso(curso model.Curso) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if curso.ID == uuid.Nil {
			return errors.New("identificador do curso é obrigatório")
		}
		if err := domain.ValidarCurso(curso); err != nil {
			return fmt.Errorf("curso inválido: %w", err)
		}
		var possuiTurmasDeOutraInstituicao bool
		if err := tx.GetContext(ctx, &possuiTurmasDeOutraInstituicao, `
			SELECT EXISTS (
				SELECT 1 FROM turma
				WHERE curso_id = $1 AND instituicao_id <> $2
			)`, curso.ID, curso.InstituicaoID,
		); err != nil {
			return fmt.Errorf("falha ao verificar turmas do curso %s: %w", curso.ID, err)
		}
		if possuiTurmasDeOutraInstituicao {
			return fmt.Errorf("curso possui turmas de outra instituição: %w", domain.ErrRelacionamentoAcademico)
		}
		result, err := tx.ExecContext(ctx, `
			UPDATE curso SET instituicao_id = $2, nome = $3, codigo = $4,
				descricao = $5, ativo = $6
			WHERE id = $1`,
			curso.ID, curso.InstituicaoID, curso.Nome, curso.Codigo, curso.Descricao, curso.Ativo,
		)
		return verificarLinhasAfetadasTx(result, err, "curso")
	}
}

func InativarCurso(id uuid.UUID) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `UPDATE curso SET ativo = FALSE WHERE id = $1`, id)
		return verificarLinhasAfetadasTx(result, err, "curso")
	}
}

func ListarTodasEtapasEnsino(ctx context.Context) ([]model.EtapaEnsino, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return nil, err
	}
	etapas := make([]model.EtapaEnsino, 0)
	if err := db.SelectContext(ctx, &etapas, selectEtapaEnsino+" ORDER BY c.nome, e.ordem"); err != nil {
		return nil, fmt.Errorf("falha ao listar etapas: %w", err)
	}
	return etapas, nil
}

func BuscarEtapaEnsino(ctx context.Context, id uuid.UUID) (model.EtapaEnsino, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return model.EtapaEnsino{}, err
	}
	var etapa model.EtapaEnsino
	if err := db.GetContext(ctx, &etapa, selectEtapaEnsino+" WHERE e.id = $1", id); err != nil {
		return model.EtapaEnsino{}, fmt.Errorf("falha ao buscar etapa %s: %w", id, err)
	}
	return etapa, nil
}

func CriarEtapaEnsino(etapa *model.EtapaEnsino) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if etapa == nil {
			return errors.New("etapa não pode ser nil")
		}
		if err := domain.ValidarEtapaEnsino(*etapa); err != nil {
			return fmt.Errorf("etapa inválida: %w", err)
		}
		err := tx.QueryRowxContext(ctx, `
			INSERT INTO etapa_ensino (id, curso_id, nome, ordem, carga_horaria)
			VALUES (
				COALESCE(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()),
				$2, $3, $4, $5
			)
			RETURNING id`,
			etapa.ID, etapa.CursoID, etapa.Nome, etapa.Ordem, etapa.CargaHoraria,
		).Scan(&etapa.ID)
		if err != nil {
			return fmt.Errorf("falha ao criar etapa: %w", err)
		}
		return nil
	}
}

func AtualizarEtapaEnsino(etapa model.EtapaEnsino) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if etapa.ID == uuid.Nil {
			return errors.New("identificador da etapa é obrigatório")
		}
		if err := domain.ValidarEtapaEnsino(etapa); err != nil {
			return fmt.Errorf("etapa inválida: %w", err)
		}
		var possuiTurmasDeOutroCurso bool
		if err := tx.GetContext(ctx, &possuiTurmasDeOutroCurso, `
			SELECT EXISTS (
				SELECT 1 FROM turma
				WHERE etapa_id = $1 AND curso_id <> $2
			)`, etapa.ID, etapa.CursoID,
		); err != nil {
			return fmt.Errorf("falha ao verificar turmas da etapa %s: %w", etapa.ID, err)
		}
		if possuiTurmasDeOutroCurso {
			return fmt.Errorf("etapa possui turmas vinculadas a outro curso: %w", domain.ErrRelacionamentoAcademico)
		}
		result, err := tx.ExecContext(ctx, `
			UPDATE etapa_ensino SET curso_id = $2, nome = $3, ordem = $4, carga_horaria = $5
			WHERE id = $1`,
			etapa.ID, etapa.CursoID, etapa.Nome, etapa.Ordem, etapa.CargaHoraria,
		)
		return verificarLinhasAfetadasTx(result, err, "etapa")
	}
}

func ExcluirEtapaEnsino(id uuid.UUID) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if id == uuid.Nil {
			return errors.New("identificador da etapa é obrigatório")
		}
		var possuiTurmas bool
		if err := tx.GetContext(ctx, &possuiTurmas,
			"SELECT EXISTS (SELECT 1 FROM turma WHERE etapa_id = $1)", id,
		); err != nil {
			return fmt.Errorf("falha ao verificar turmas da etapa %s: %w", id, err)
		}
		if possuiTurmas {
			return fmt.Errorf("etapa %s: %w", id, domain.ErrEtapaPossuiTurmas)
		}
		result, err := tx.ExecContext(ctx, "DELETE FROM etapa_ensino WHERE id = $1", id)
		return verificarLinhasAfetadasTx(result, err, "etapa")
	}
}

func ListarTodasDisciplinas(ctx context.Context) ([]model.Disciplina, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return nil, err
	}
	disciplinas := make([]model.Disciplina, 0)
	if err := db.SelectContext(ctx, &disciplinas, selectDisciplina+" ORDER BY i.nome, d.nome"); err != nil {
		return nil, fmt.Errorf("falha ao listar disciplinas: %w", err)
	}
	return disciplinas, nil
}

func BuscarDisciplina(ctx context.Context, id uuid.UUID) (model.Disciplina, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return model.Disciplina{}, err
	}
	var disciplina model.Disciplina
	if err := db.GetContext(ctx, &disciplina, selectDisciplina+" WHERE d.id = $1", id); err != nil {
		return model.Disciplina{}, fmt.Errorf("falha ao buscar disciplina %s: %w", id, err)
	}
	return disciplina, nil
}

func CriarDisciplina(disciplina *model.Disciplina) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if disciplina == nil {
			return errors.New("disciplina não pode ser nil")
		}
		if err := domain.ValidarDisciplina(*disciplina); err != nil {
			return fmt.Errorf("disciplina inválida: %w", err)
		}
		err := tx.QueryRowxContext(ctx, `
			INSERT INTO disciplina (id, instituicao_id, nome, codigo, carga_horaria, ativa)
			VALUES (
				COALESCE(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()),
				$2, $3, $4, $5, $6
			)
			RETURNING id`,
			disciplina.ID, disciplina.InstituicaoID, disciplina.Nome, disciplina.Codigo,
			disciplina.CargaHoraria, disciplina.Ativa,
		).Scan(&disciplina.ID)
		if err != nil {
			return fmt.Errorf("falha ao criar disciplina: %w", err)
		}
		return nil
	}
}

func AtualizarDisciplina(disciplina model.Disciplina) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if disciplina.ID == uuid.Nil {
			return errors.New("identificador da disciplina é obrigatório")
		}
		if err := domain.ValidarDisciplina(disciplina); err != nil {
			return fmt.Errorf("disciplina inválida: %w", err)
		}
		result, err := tx.ExecContext(ctx, `
			UPDATE disciplina SET instituicao_id = $2, nome = $3, codigo = $4,
				carga_horaria = $5, ativa = $6
			WHERE id = $1`,
			disciplina.ID, disciplina.InstituicaoID, disciplina.Nome, disciplina.Codigo,
			disciplina.CargaHoraria, disciplina.Ativa,
		)
		return verificarLinhasAfetadasTx(result, err, "disciplina")
	}
}

func InativarDisciplina(id uuid.UUID) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, "UPDATE disciplina SET ativa = FALSE WHERE id = $1", id)
		return verificarLinhasAfetadasTx(result, err, "disciplina")
	}
}

func ListarTodasTurmasAcademicas(ctx context.Context) ([]model.Turma, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return nil, err
	}
	turmas := make([]model.Turma, 0)
	if err := db.SelectContext(ctx, &turmas, selectTurmaAcademica+" ORDER BY t.ano_letivo DESC, t.nome"); err != nil {
		return nil, fmt.Errorf("falha ao listar turmas acadêmicas: %w", err)
	}
	return turmas, nil
}

func BuscarTurmaAcademica(ctx context.Context, id uuid.UUID) (model.Turma, error) {
	db, err := obterDB(ctx)
	if err != nil {
		return model.Turma{}, err
	}
	var turma model.Turma
	if err := db.GetContext(ctx, &turma, selectTurmaAcademica+" WHERE t.id = $1", id); err != nil {
		return model.Turma{}, fmt.Errorf("falha ao buscar turma %s: %w", id, err)
	}
	return turma, nil
}

func CriarTurmaAcademica(turma *model.Turma) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if turma == nil {
			return errors.New("turma não pode ser nil")
		}
		if turma.Status == "" {
			turma.Status = model.StatusTurmaPlanejada
		}
		if err := domain.ValidarTurma(*turma); err != nil {
			return fmt.Errorf("turma inválida: %w", err)
		}
		if err := validarRelacionamentosTurma(ctx, tx, *turma); err != nil {
			return err
		}
		err := tx.QueryRowxContext(ctx, `
			INSERT INTO turma (
				id, instituicao_id, curso_id, etapa_id, nome, codigo,
				ano_letivo, turno, status
			) VALUES (
				COALESCE(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()),
				$2, $3, $4, $5, $6, $7, $8, $9
			)
			RETURNING id, criado_em`,
			turma.ID, turma.InstituicaoID, turma.CursoID, turma.EtapaID,
			turma.Nome, turma.Codigo, turma.AnoLetivo, turma.Turno, turma.Status,
		).Scan(&turma.ID, &turma.CriadoEm)
		if err != nil {
			return fmt.Errorf("falha ao criar turma: %w", err)
		}
		return nil
	}
}

func AtualizarTurmaAcademica(turma model.Turma) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if turma.ID == uuid.Nil {
			return errors.New("identificador da turma é obrigatório")
		}
		if err := domain.ValidarTurma(turma); err != nil {
			return fmt.Errorf("turma inválida: %w", err)
		}
		if err := validarRelacionamentosTurma(ctx, tx, turma); err != nil {
			return err
		}
		result, err := tx.ExecContext(ctx, `
			UPDATE turma SET instituicao_id = $2, curso_id = $3, etapa_id = $4,
				nome = $5, codigo = $6, ano_letivo = $7, turno = $8, status = $9
			WHERE id = $1`,
			turma.ID, turma.InstituicaoID, turma.CursoID, turma.EtapaID,
			turma.Nome, turma.Codigo, turma.AnoLetivo, turma.Turno, turma.Status,
		)
		return verificarLinhasAfetadasTx(result, err, "turma")
	}
}

func validarRelacionamentosTurma(ctx context.Context, tx *sqlx.Tx, turma model.Turma) error {
	var relacionamentosValidos bool
	if err := tx.GetContext(ctx, &relacionamentosValidos, `
		SELECT EXISTS (
			SELECT 1
			FROM curso AS c
			JOIN etapa_ensino AS e ON e.curso_id = c.id
			WHERE c.id = $1 AND c.instituicao_id = $2 AND e.id = $3
		)`, turma.CursoID, turma.InstituicaoID, turma.EtapaID,
	); err != nil {
		return fmt.Errorf("falha ao validar vínculos da turma %s: %w", turma.ID, err)
	}
	if !relacionamentosValidos {
		return fmt.Errorf("instituição, curso e etapa da turma não correspondem: %w", domain.ErrRelacionamentoAcademico)
	}
	return nil
}

func EncerrarTurma(id uuid.UUID) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `
			UPDATE turma SET status = $2 WHERE id = $1`,
			id, model.StatusTurmaEncerrada,
		)
		return verificarLinhasAfetadasTx(result, err, "turma")
	}
}
