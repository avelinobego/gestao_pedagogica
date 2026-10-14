package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/avelinobego/gestao_pedagogica/domain"
	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"github.com/google/uuid"
)

var instituicaoTesteID = uuid.MustParse("90000000-0000-0000-0000-000000000001")
var turmaTesteID = uuid.MustParse("30000000-0000-0000-0000-000000000001")

func cpfTeste() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")[:11]
}

func numeroMatriculaTeste() string {
	return uuid.NewString()[:30]
}

func cleanupRegistro(t *testing.T, tabela string, id uuid.UUID) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db, err := plataforma.GetDB(ctx)
		if err != nil {
			t.Errorf("failed to get test database for cleanup: %v", err)
			return
		}
		if _, err := db.ExecContext(ctx, "DELETE FROM "+tabela+" WHERE id = $1", id); err != nil {
			t.Errorf("failed to clean up test %s record: %v", tabela, err)
		}
	})
}

func cleanupVinculoTeste(t *testing.T, alunoID, responsavelID uuid.UUID) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db, err := plataforma.GetDB(ctx)
		if err != nil {
			t.Errorf("failed to get test database for cleanup: %v", err)
			return
		}
		if _, err := db.ExecContext(ctx, `
			DELETE FROM aluno_responsavel WHERE aluno_id = $1 AND responsavel_id = $2`,
			alunoID, responsavelID,
		); err != nil {
			t.Errorf("failed to clean up test student/guardian relationship: %v", err)
		}
	})
}

func requireDatabase(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := plataforma.GetDB(ctx)
	if err != nil {
		t.Fatalf("test database was not initialized: %v", err)
	}
	var possuiDadosBase bool
	err = db.GetContext(ctx, &possuiDadosBase, `
		SELECT EXISTS (
			SELECT 1 FROM instituicao_ensino WHERE id = $1
		) AND EXISTS (
			SELECT 1 FROM turma WHERE id = $2
		)`, instituicaoTesteID, turmaTesteID)
	if err != nil {
		t.Fatalf("failed to verify test database seed data: %v", err)
	}
	if !possuiDadosBase {
		t.Fatal("test database is missing the seeded institution or turma")
	}
}

// TestComposicaoHandlersAluno demonstrates atomic composition of multiple handlers
// creating an aluno in a single transaction via plataforma.Executar
func TestComposicaoHandlersAluno(t *testing.T) {
	requireDatabase(t)

	ctx := context.Background()

	// Create an aluno with all fields
	aluno := &model.Aluno{
		ID:             uuid.New(),
		InstituicaoID:  uuid.MustParse("90000000-0000-0000-0000-000000000001"),
		Nome:           "João Silva",
		NomeSocial:     "João",
		CPF:            cpfTeste(),
		RG:             "123456789",
		DataNascimento: time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC),
		Email:          "joao@example.com",
		Telefone:       "11999999999",
		Endereco:       "Rua A, 123",
		Cidade:         "São Paulo",
		UF:             "SP",
		Status:         model.StatusAlunoAtivo,
	}
	cleanupRegistro(t, "aluno", aluno.ID)

	// Execute handler in transaction
	if err := plataforma.Executar(ctx, CriarAluno(aluno)); err != nil {
		t.Fatalf("failed to execute CriarAluno handler: %v", err)
	}

	// Verify the aluno was created by reading it back
	retrieved, err := BuscarAluno(ctx, aluno.ID)
	if err != nil {
		t.Fatalf("failed to retrieve created aluno: %v", err)
	}

	if retrieved.Nome != aluno.Nome {
		t.Errorf("expected name %q, got %q", aluno.Nome, retrieved.Nome)
	}
	if retrieved.Status != model.StatusAlunoAtivo {
		t.Errorf("expected status %q, got %q", model.StatusAlunoAtivo, retrieved.Status)
	}
}

// TestComposicaoHandlersResponsavel demonstrates atomic composition for responsável operations
func TestComposicaoHandlersResponsavel(t *testing.T) {
	requireDatabase(t)

	ctx := context.Background()

	// Create a responsável
	responsavel := &model.Responsavel{
		ID:       uuid.New(),
		Nome:     "Maria Silva",
		CPF:      cpfTeste(),
		Email:    "maria@example.com",
		Telefone: "11988888888",
		Endereco: "Rua B, 456",
	}
	cleanupRegistro(t, "responsavel", responsavel.ID)

	// Execute handler in transaction
	if err := plataforma.Executar(ctx, CriarResponsavel(responsavel)); err != nil {
		t.Fatalf("failed to execute CriarResponsavel handler: %v", err)
	}

	// Verify the responsável was created
	retrieved, err := BuscarResponsavel(ctx, responsavel.ID)
	if err != nil {
		t.Fatalf("failed to retrieve created responsável: %v", err)
	}

	if retrieved.Nome != responsavel.Nome {
		t.Errorf("expected name %q, got %q", responsavel.Nome, retrieved.Nome)
	}
}

// TestComposicaoMultiplosHandlers demonstrates composing multiple handlers in a single transaction.
// This creates an aluno, then updates it, ensuring atomicity.
func TestComposicaoMultiplosHandlers(t *testing.T) {
	requireDatabase(t)

	ctx := context.Background()

	// Create an aluno with initial data
	aluno := &model.Aluno{
		ID:             uuid.New(),
		InstituicaoID:  instituicaoTesteID,
		Nome:           "Pedro Santos",
		NomeSocial:     "",
		CPF:            "",
		RG:             "",
		DataNascimento: time.Date(2011, 8, 20, 0, 0, 0, 0, time.UTC),
		Email:          "pedro@example.com",
		Telefone:       "11987654321",
		Endereco:       "Rua C, 789",
		Cidade:         "São Paulo",
		UF:             "SP",
		Status:         model.StatusAlunoAtivo,
	}
	cleanupRegistro(t, "aluno", aluno.ID)

	// Update data for the update handler
	alunoAtualizado := *aluno
	alunoAtualizado.NomeSocial = "Pedrinho"
	alunoAtualizado.Telefone = "11987654320"

	// Compose both create and update in a single transaction.
	// The transaction will commit both operations atomically or rollback on any error.
	err := plataforma.Executar(ctx,
		CriarAluno(aluno),
		AtualizarAluno(alunoAtualizado),
	)
	if err != nil {
		t.Fatalf("failed to execute composite handlers: %v", err)
	}

	// Verify the final state includes the update
	retrieved, err := BuscarAluno(ctx, aluno.ID)
	if err != nil {
		t.Fatalf("failed to retrieve aluno after update: %v", err)
	}

	if retrieved.NomeSocial != "Pedrinho" {
		t.Errorf("expected nome social %q, got %q", "Pedrinho", retrieved.NomeSocial)
	}
	if retrieved.Telefone != "11987654320" {
		t.Errorf("expected telefone %q, got %q", "11987654320", retrieved.Telefone)
	}
}

// TestComposicaoAlunoResponsavelMatricula demonstrates composing handlers for multiple entities:
// creating an aluno, responsável, their relationship, and a matrícula all atomically.
func TestComposicaoAlunoResponsavelMatricula(t *testing.T) {
	requireDatabase(t)

	ctx := context.Background()

	alunoID := uuid.New()
	responsavelID := uuid.New()
	turmaID := turmaTesteID
	instituicaoID := instituicaoTesteID

	// Create aluno
	aluno := &model.Aluno{
		ID:             alunoID,
		InstituicaoID:  instituicaoID,
		Nome:           "Lucas Oliveira",
		NomeSocial:     "",
		CPF:            "",
		RG:             "",
		DataNascimento: time.Date(2012, 3, 10, 0, 0, 0, 0, time.UTC),
		Email:          "lucas@example.com",
		Telefone:       "11985555555",
		Endereco:       "Rua D, 321",
		Cidade:         "São Paulo",
		UF:             "SP",
		Status:         model.StatusAlunoAtivo,
	}
	cleanupRegistro(t, "aluno", aluno.ID)

	// Create responsável
	responsavel := &model.Responsavel{
		ID:       responsavelID,
		Nome:     "Ana Oliveira",
		CPF:      cpfTeste(),
		Email:    "ana@example.com",
		Telefone: "11984444444",
		Endereco: "Rua D, 321",
	}
	cleanupRegistro(t, "responsavel", responsavel.ID)

	// Create relationship
	vinculo := model.AlunoResponsavel{
		AlunoID:       alunoID,
		ResponsavelID: responsavelID,
		Parentesco:    "MÃE",
		Principal:     true,
	}
	cleanupVinculoTeste(t, alunoID, responsavelID)

	// Create matrícula
	matricula := &model.Matricula{
		ID:            uuid.New(),
		AlunoID:       alunoID,
		TurmaID:       turmaID,
		Numero:        numeroMatriculaTeste(),
		DataMatricula: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		DataConclusao: nil,
		Status:        model.StatusMatriculaAtiva,
		Observacao:    "Matrícula automática",
	}
	cleanupRegistro(t, "matricula", matricula.ID)

	// Execute all handlers in a single atomic transaction
	err := plataforma.Executar(ctx,
		CriarAluno(aluno),
		CriarResponsavel(responsavel),
		CriarVinculoAlunoResponsavel(vinculo),
		CriarMatricula(matricula),
	)
	if err != nil {
		t.Fatalf("failed to execute composite handlers: %v", err)
	}

	// Verify aluno was created
	retrievedAluno, err := BuscarAluno(ctx, alunoID)
	if err != nil {
		t.Fatalf("failed to retrieve created aluno: %v", err)
	}
	if retrievedAluno.Nome != aluno.Nome {
		t.Errorf("expected aluno name %q, got %q", aluno.Nome, retrievedAluno.Nome)
	}

	// Verify responsável was created
	retrievedResponsavel, err := BuscarResponsavel(ctx, responsavelID)
	if err != nil {
		t.Fatalf("failed to retrieve created responsável: %v", err)
	}
	if retrievedResponsavel.Nome != responsavel.Nome {
		t.Errorf("expected responsável name %q, got %q", responsavel.Nome, retrievedResponsavel.Nome)
	}

	// Verify relationship was created
	retrievedVinculo, err := BuscarVinculoAlunoResponsavel(ctx, alunoID, responsavelID)
	if err != nil {
		t.Fatalf("failed to retrieve created relationship: %v", err)
	}
	if !retrievedVinculo.Principal {
		t.Errorf("expected principal to be true, got false")
	}

	// Verify matrícula was created
	retrievedMatricula, err := BuscarMatricula(ctx, matricula.ID)
	if err != nil {
		t.Fatalf("failed to retrieve created matrícula: %v", err)
	}
	if retrievedMatricula.Status != model.StatusMatriculaAtiva {
		t.Errorf("expected status %q, got %q", model.StatusMatriculaAtiva, retrievedMatricula.Status)
	}
}

// TestInativarAlunoHandler demonstrates the inactivation handler in a transaction
func TestInativarAlunoHandler(t *testing.T) {
	requireDatabase(t)

	ctx := context.Background()

	alunoID := uuid.New()

	// Create an aluno
	aluno := &model.Aluno{
		ID:             alunoID,
		InstituicaoID:  instituicaoTesteID,
		Nome:           "Carlos Mendes",
		NomeSocial:     "",
		CPF:            "",
		RG:             "",
		DataNascimento: time.Date(2010, 12, 25, 0, 0, 0, 0, time.UTC),
		Email:          "carlos@example.com",
		Telefone:       "11986666666",
		Endereco:       "Rua E, 654",
		Cidade:         "São Paulo",
		UF:             "SP",
		Status:         model.StatusAlunoAtivo,
	}
	cleanupRegistro(t, "aluno", alunoID)

	// Create and then inactivate in a single transaction
	err := plataforma.Executar(ctx,
		CriarAluno(aluno),
		InativarAluno(alunoID),
	)
	if err != nil {
		t.Fatalf("failed to execute handlers: %v", err)
	}

	// Verify the aluno is now inactive
	retrieved, err := BuscarAluno(ctx, alunoID)
	if err != nil {
		t.Fatalf("failed to retrieve aluno: %v", err)
	}

	if retrieved.Status != model.StatusAlunoInativo {
		t.Errorf("expected status %q, got %q", model.StatusAlunoInativo, retrieved.Status)
	}
}

// TestCancelarMatriculaHandler demonstrates the cancellation handler in a transaction
func TestCancelarMatriculaHandler(t *testing.T) {
	requireDatabase(t)

	ctx := context.Background()

	alunoID := uuid.New()
	turmaID := turmaTesteID
	matriculaID := uuid.New()

	// Create aluno and matrícula, then cancel
	aluno := &model.Aluno{
		ID:             alunoID,
		InstituicaoID:  instituicaoTesteID,
		Nome:           "Fernanda Costa",
		NomeSocial:     "",
		CPF:            "",
		RG:             "",
		DataNascimento: time.Date(2011, 6, 30, 0, 0, 0, 0, time.UTC),
		Email:          "fernanda@example.com",
		Telefone:       "11983333333",
		Endereco:       "Rua F, 987",
		Cidade:         "São Paulo",
		UF:             "SP",
		Status:         model.StatusAlunoAtivo,
	}
	cleanupRegistro(t, "aluno", alunoID)

	matricula := &model.Matricula{
		ID:            matriculaID,
		AlunoID:       alunoID,
		TurmaID:       turmaID,
		Numero:        numeroMatriculaTeste(),
		DataMatricula: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
		DataConclusao: nil,
		Status:        model.StatusMatriculaAtiva,
		Observacao:    "",
	}
	cleanupRegistro(t, "matricula", matriculaID)

	// Create aluno and matrícula, then cancel matrícula in a single transaction
	err := plataforma.Executar(ctx,
		CriarAluno(aluno),
		CriarMatricula(matricula),
		CancelarMatricula(matriculaID),
	)
	if err != nil {
		t.Fatalf("failed to execute handlers: %v", err)
	}

	// Verify matrícula is now cancelled
	retrieved, err := BuscarMatricula(ctx, matriculaID)
	if err != nil {
		t.Fatalf("failed to retrieve matrícula: %v", err)
	}

	if retrieved.Status != model.StatusMatriculaCancelada {
		t.Errorf("expected status %q, got %q", model.StatusMatriculaCancelada, retrieved.Status)
	}
}

func TestAtivarMatriculaHandler(t *testing.T) {
	requireDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	aluno := &model.Aluno{
		ID:             uuid.New(),
		InstituicaoID:  instituicaoTesteID,
		Nome:           "Aluno teste de ativação",
		CPF:            cpfTeste(),
		DataNascimento: time.Date(2012, 4, 5, 0, 0, 0, 0, time.UTC),
		Status:         model.StatusAlunoAtivo,
	}
	cleanupRegistro(t, "aluno", aluno.ID)
	matricula := &model.Matricula{
		ID:            uuid.New(),
		AlunoID:       aluno.ID,
		TurmaID:       turmaTesteID,
		Numero:        numeroMatriculaTeste(),
		DataMatricula: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Status:        model.StatusMatriculaTrancada,
	}
	cleanupRegistro(t, "matricula", matricula.ID)

	if err := plataforma.Executar(ctx, CriarAluno(aluno), CriarMatricula(matricula)); err != nil {
		t.Fatalf("failed to create test student and matrícula: %v", err)
	}
	if err := plataforma.Executar(ctx, AtivarMatricula(matricula.ID)); err != nil {
		t.Fatalf("failed to activate matrícula: %v", err)
	}

	retrieved, err := BuscarMatricula(ctx, matricula.ID)
	if err != nil {
		t.Fatalf("failed to retrieve activated matrícula: %v", err)
	}
	if retrieved.Status != model.StatusMatriculaAtiva {
		t.Errorf("status = %q; expected %q", retrieved.Status, model.StatusMatriculaAtiva)
	}
	if err := plataforma.Executar(ctx, AtivarMatricula(matricula.ID)); !errors.Is(err, domain.ErrMatriculaNaoAtivavel) {
		t.Errorf("second activation error = %v; expected ErrMatriculaNaoAtivavel", err)
	}
}

func TestAtivarMatriculaInexistenteRetornaNoRows(t *testing.T) {
	requireDatabase(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := plataforma.Executar(ctx, AtivarMatricula(uuid.New()))
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("AtivarMatricula() error = %v; expected sql.ErrNoRows", err)
	}
}
