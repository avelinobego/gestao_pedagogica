package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/avelinobego/gestao_pedagogica/domain"
	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const selectAluno = `
	SELECT id, instituicao_id, nome, nome_social, COALESCE(cpf, '') AS cpf,
		rg, data_nascimento, email, telefone, endereco, cidade, uf, status,
		criado_em, atualizado_em
	FROM aluno`

const selectResponsavel = `
	SELECT id, nome, COALESCE(cpf, '') AS cpf, email, telefone, endereco, criado_em
	FROM responsavel`

const selectMatricula = `
	SELECT id, aluno_id, turma_id, numero, data_matricula, data_conclusao, status,
		observacao, criado_em, atualizado_em
	FROM matricula`

func obterConn(ctx context.Context) (*sqlx.Conn, error) {
	conn, err := plataforma.GetConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter conexão com o banco: %w", err)
	}
	return conn, nil
}

func verificarLinhasAfetadasTx(result sql.Result, err error, entidade string) error {
	if err != nil {
		return fmt.Errorf("falha ao persistir %s: %w", entidade, err)
	}
	afetadas, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("falha ao verificar persistência de %s: %w", entidade, err)
	}
	if afetadas == 0 {
		return fmt.Errorf("%s: %w", entidade, sql.ErrNoRows)
	}
	return nil
}

func ListarTodosAlunos(ctx context.Context) ([]model.Aluno, error) {
	conn, err := obterConn(ctx)
	if err != nil {
		return nil, err
	}
	alunos := make([]model.Aluno, 0)
	if err := conn.SelectContext(ctx, &alunos, selectAluno+" ORDER BY nome ASC"); err != nil {
		return nil, fmt.Errorf("falha ao listar alunos: %w", err)
	}
	return alunos, nil
}

func BuscarAluno(ctx context.Context, id uuid.UUID) (model.Aluno, error) {
	conn, err := obterConn(ctx)
	if err != nil {
		return model.Aluno{}, err
	}
	var aluno model.Aluno
	if err := conn.GetContext(ctx, &aluno, selectAluno+" WHERE id = $1", id); err != nil {
		return model.Aluno{}, fmt.Errorf("falha ao buscar aluno %s: %w", id, err)
	}
	return aluno, nil
}

func CriarAluno(aluno *model.Aluno) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if aluno == nil {
			return errors.New("aluno não pode ser nil")
		}
		if aluno.Status == "" {
			aluno.Status = model.StatusAlunoAtivo
		}
		err := tx.QueryRowxContext(ctx, `
			INSERT INTO aluno (
				id, instituicao_id, nome, nome_social, cpf, rg, data_nascimento,
				email, telefone, endereco, cidade, uf, status
			) VALUES (
				COALESCE(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()),
				$2, $3, $4, NULLIF($5, ''), $6, $7, $8, $9, $10, $11, $12, $13
			)
			RETURNING id, criado_em, atualizado_em`,
			aluno.ID, aluno.InstituicaoID, aluno.Nome, aluno.NomeSocial, aluno.CPF,
			aluno.RG, aluno.DataNascimento, aluno.Email, aluno.Telefone, aluno.Endereco,
			aluno.Cidade, aluno.UF, aluno.Status,
		).Scan(&aluno.ID, &aluno.CriadoEm, &aluno.AtualizadoEm)
		if err != nil {
			return fmt.Errorf("falha ao criar aluno: %w", err)
		}
		return nil
	}
}

func AtualizarAluno(aluno model.Aluno) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `
			UPDATE aluno
			SET instituicao_id = $2, nome = $3, nome_social = $4, cpf = NULLIF($5, ''),
				rg = $6, data_nascimento = $7, email = $8, telefone = $9, endereco = $10,
				cidade = $11, uf = $12, status = $13, atualizado_em = CURRENT_TIMESTAMP
			WHERE id = $1`,
			aluno.ID, aluno.InstituicaoID, aluno.Nome, aluno.NomeSocial, aluno.CPF,
			aluno.RG, aluno.DataNascimento, aluno.Email, aluno.Telefone, aluno.Endereco,
			aluno.Cidade, aluno.UF, aluno.Status,
		)
		return verificarLinhasAfetadasTx(result, err, "aluno")
	}
}

func InativarAluno(id uuid.UUID) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `
			UPDATE aluno SET status = $2, atualizado_em = CURRENT_TIMESTAMP WHERE id = $1`,
			id, model.StatusAlunoInativo,
		)
		return verificarLinhasAfetadasTx(result, err, "aluno")
	}
}

func ListarTodosResponsaveis(ctx context.Context) ([]model.Responsavel, error) {
	conn, err := obterConn(ctx)
	if err != nil {
		return nil, err
	}
	responsaveis := make([]model.Responsavel, 0)
	if err := conn.SelectContext(ctx, &responsaveis, selectResponsavel+" ORDER BY nome ASC"); err != nil {
		return nil, fmt.Errorf("falha ao listar responsáveis: %w", err)
	}
	return responsaveis, nil
}

func BuscarResponsavel(ctx context.Context, id uuid.UUID) (model.Responsavel, error) {
	conn, err := obterConn(ctx)
	if err != nil {
		return model.Responsavel{}, err
	}
	var responsavel model.Responsavel
	if err := conn.GetContext(ctx, &responsavel, selectResponsavel+" WHERE id = $1", id); err != nil {
		return model.Responsavel{}, fmt.Errorf("falha ao buscar responsável %s: %w", id, err)
	}
	return responsavel, nil
}

func CriarResponsavel(responsavel *model.Responsavel) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if responsavel == nil {
			return errors.New("responsável não pode ser nil")
		}
		err := tx.QueryRowxContext(ctx, `
			INSERT INTO responsavel (id, nome, cpf, email, telefone, endereco)
			VALUES (
				COALESCE(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()),
				$2, NULLIF($3, ''), $4, $5, $6
			)
			RETURNING id, criado_em`,
			responsavel.ID, responsavel.Nome, responsavel.CPF, responsavel.Email,
			responsavel.Telefone, responsavel.Endereco,
		).Scan(&responsavel.ID, &responsavel.CriadoEm)
		if err != nil {
			return fmt.Errorf("falha ao criar responsável: %w", err)
		}
		return nil
	}
}

func AtualizarResponsavel(responsavel model.Responsavel) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `
			UPDATE responsavel
			SET nome = $2, cpf = NULLIF($3, ''), email = $4, telefone = $5, endereco = $6
			WHERE id = $1`,
			responsavel.ID, responsavel.Nome, responsavel.CPF, responsavel.Email,
			responsavel.Telefone, responsavel.Endereco,
		)
		return verificarLinhasAfetadasTx(result, err, "responsável")
	}
}

func ExcluirResponsavel(id uuid.UUID) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `
			DELETE FROM responsavel AS r
			WHERE r.id = $1
				AND NOT EXISTS (
					SELECT 1 FROM aluno_responsavel AS ar WHERE ar.responsavel_id = r.id
				)`, id)
		if err != nil {
			return fmt.Errorf("falha ao excluir responsável %s: %w", id, err)
		}
		afetadas, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("falha ao verificar exclusão do responsável %s: %w", id, err)
		}
		if afetadas > 0 {
			return nil
		}
		var existe bool
		if err := tx.GetContext(ctx, &existe, "SELECT EXISTS (SELECT 1 FROM responsavel WHERE id = $1)", id); err != nil {
			return fmt.Errorf("falha ao verificar responsável %s: %w", id, err)
		}
		if !existe {
			return fmt.Errorf("responsável %s: %w", id, sql.ErrNoRows)
		}
		return fmt.Errorf("responsável %s possui vínculos com alunos e não pode ser excluído", id)
	}
}

func ListarVinculosPorAluno(ctx context.Context, alunoID uuid.UUID) ([]model.AlunoResponsavel, error) {
	conn, err := obterConn(ctx)
	if err != nil {
		return nil, err
	}
	vinculos := make([]model.AlunoResponsavel, 0)
	if err := conn.SelectContext(ctx, &vinculos, `
		SELECT aluno_id, responsavel_id, parentesco, principal
		FROM aluno_responsavel WHERE aluno_id = $1
		ORDER BY principal DESC, parentesco ASC`, alunoID); err != nil {
		return nil, fmt.Errorf("falha ao listar vínculos do aluno %s: %w", alunoID, err)
	}
	return vinculos, nil
}

func BuscarVinculoAlunoResponsavel(ctx context.Context, alunoID, responsavelID uuid.UUID) (model.AlunoResponsavel, error) {
	conn, err := obterConn(ctx)
	if err != nil {
		return model.AlunoResponsavel{}, err
	}
	var vinculo model.AlunoResponsavel
	if err := conn.GetContext(ctx, &vinculo, `
		SELECT aluno_id, responsavel_id, parentesco, principal
		FROM aluno_responsavel
		WHERE aluno_id = $1 AND responsavel_id = $2`,
		alunoID, responsavelID,
	); err != nil {
		return model.AlunoResponsavel{}, fmt.Errorf(
			"falha ao buscar vínculo do aluno %s com responsável %s: %w",
			alunoID, responsavelID, err,
		)
	}
	return vinculo, nil
}

func CriarVinculoAlunoResponsavel(vinculo model.AlunoResponsavel) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO aluno_responsavel (aluno_id, responsavel_id, parentesco, principal)
			VALUES ($1, $2, $3, $4)`,
			vinculo.AlunoID, vinculo.ResponsavelID, vinculo.Parentesco, vinculo.Principal,
		); err != nil {
			return fmt.Errorf("falha ao criar vínculo aluno-responsável: %w", err)
		}
		return nil
	}
}

func AtualizarVinculoAlunoResponsavel(vinculo model.AlunoResponsavel) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `
			UPDATE aluno_responsavel SET parentesco = $3, principal = $4
			WHERE aluno_id = $1 AND responsavel_id = $2`,
			vinculo.AlunoID, vinculo.ResponsavelID, vinculo.Parentesco, vinculo.Principal,
		)
		return verificarLinhasAfetadasTx(result, err, "vínculo aluno-responsável")
	}
}

func ExcluirVinculoAlunoResponsavel(alunoID, responsavelID uuid.UUID) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `
			DELETE FROM aluno_responsavel WHERE aluno_id = $1 AND responsavel_id = $2`,
			alunoID, responsavelID,
		)
		return verificarLinhasAfetadasTx(result, err, "vínculo aluno-responsável")
	}
}

func ListarTodasMatriculas(ctx context.Context) ([]model.Matricula, error) {
	conn, err := obterConn(ctx)
	if err != nil {
		return nil, err
	}
	matriculas := make([]model.Matricula, 0)
	if err := conn.SelectContext(ctx, &matriculas, selectMatricula+" ORDER BY data_matricula DESC, numero ASC"); err != nil {
		return nil, fmt.Errorf("falha ao listar matrículas: %w", err)
	}
	return matriculas, nil
}

func ListarMatriculasPorAluno(ctx context.Context, alunoID uuid.UUID) ([]model.Matricula, error) {
	conn, err := obterConn(ctx)
	if err != nil {
		return nil, err
	}
	matriculas := make([]model.Matricula, 0)
	if err := conn.SelectContext(ctx, &matriculas, selectMatricula+" WHERE aluno_id = $1 ORDER BY data_matricula DESC, numero ASC", alunoID); err != nil {
		return nil, fmt.Errorf("falha ao listar matrículas do aluno %s: %w", alunoID, err)
	}
	return matriculas, nil
}

func BuscarMatricula(ctx context.Context, id uuid.UUID) (model.Matricula, error) {
	conn, err := obterConn(ctx)
	if err != nil {
		return model.Matricula{}, err
	}
	var matricula model.Matricula
	if err := conn.GetContext(ctx, &matricula, selectMatricula+" WHERE id = $1", id); err != nil {
		return model.Matricula{}, fmt.Errorf("falha ao buscar matrícula %s: %w", id, err)
	}
	return matricula, nil
}

func CriarMatricula(matricula *model.Matricula) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if matricula == nil {
			return errors.New("matrícula não pode ser nil")
		}
		if matricula.Status == "" {
			matricula.Status = model.StatusMatriculaAtiva
		}
		if err := domain.ValidarMatricula(*matricula); err != nil {
			return fmt.Errorf("matrícula inválida: %w", err)
		}
		err := tx.QueryRowxContext(ctx, `
			INSERT INTO matricula (
				id, aluno_id, turma_id, numero, data_matricula, data_conclusao, status, observacao
			) VALUES (
				COALESCE(NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), uuid_generate_v4()),
				$2, $3, $4, $5, $6, $7, $8
			)
			RETURNING id, criado_em, atualizado_em`,
			matricula.ID, matricula.AlunoID, matricula.TurmaID, matricula.Numero,
			matricula.DataMatricula, matricula.DataConclusao, matricula.Status,
			matricula.Observacao,
		).Scan(&matricula.ID, &matricula.CriadoEm, &matricula.AtualizadoEm)
		if err != nil {
			return fmt.Errorf("falha ao criar matrícula: %w", err)
		}
		return nil
	}
}

func AtualizarMatricula(matricula model.Matricula) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		if err := domain.ValidarMatricula(matricula); err != nil {
			return fmt.Errorf("matrícula inválida: %w", err)
		}
		result, err := tx.ExecContext(ctx, `
			UPDATE matricula
			SET aluno_id = $2, turma_id = $3, numero = $4, data_matricula = $5,
				data_conclusao = $6, status = $7, observacao = $8,
				atualizado_em = CURRENT_TIMESTAMP
			WHERE id = $1`,
			matricula.ID, matricula.AlunoID, matricula.TurmaID, matricula.Numero,
			matricula.DataMatricula, matricula.DataConclusao, matricula.Status,
			matricula.Observacao,
		)
		return verificarLinhasAfetadasTx(result, err, "matrícula")
	}
}

func CancelarMatricula(id uuid.UUID) plataforma.Handler {
	return func(tx *sqlx.Tx, ctx context.Context) error {
		result, err := tx.ExecContext(ctx, `
			UPDATE matricula
			SET status = $2, atualizado_em = CURRENT_TIMESTAMP
			WHERE id = $1`,
			id, model.StatusMatriculaCancelada,
		)
		return verificarLinhasAfetadasTx(result, err, "matrícula")
	}
}
