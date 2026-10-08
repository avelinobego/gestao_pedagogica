package service

import (
	"context"
	"testing"
	"time"

	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"github.com/google/uuid"
)

// TestComposicaoHandlersAluno demonstrates atomic composition of multiple handlers
// creating an aluno in a single transaction via plataforma.Executar
func TestComposicaoHandlersAluno(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create an aluno with all fields
	aluno := &model.Aluno{
		ID:             uuid.New(),
		InstituicaoID:  uuid.New(),
		Nome:           "João Silva",
		NomeSocial:     "João",
		CPF:            "12345678900",
		RG:             "123456789",
		DataNascimento: time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC),
		Email:          "joao@example.com",
		Telefone:       "11999999999",
		Endereco:       "Rua A, 123",
		Cidade:         "São Paulo",
		UF:             "SP",
		Status:         model.StatusAlunoAtivo,
	}

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
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create a responsável
	responsavel := &model.Responsavel{
		ID:       uuid.New(),
		Nome:     "Maria Silva",
		CPF:      "98765432100",
		Email:    "maria@example.com",
		Telefone: "11988888888",
		Endereco: "Rua B, 456",
	}

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
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create an aluno with initial data
	aluno := &model.Aluno{
		ID:             uuid.New(),
		InstituicaoID:  uuid.New(),
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
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	alunoID := uuid.New()
	responsavelID := uuid.New()
	turmaID := uuid.New()
	instituicaoID := uuid.New()

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

	// Create responsável
	responsavel := &model.Responsavel{
		ID:       responsavelID,
		Nome:     "Ana Oliveira",
		CPF:      "11122233344",
		Email:    "ana@example.com",
		Telefone: "11984444444",
		Endereco: "Rua D, 321",
	}

	// Create relationship
	vinculo := model.AlunoResponsavel{
		AlunoID:       alunoID,
		ResponsavelID: responsavelID,
		Parentesco:    "MÃE",
		Principal:     true,
	}

	// Create matrícula
	matricula := &model.Matricula{
		ID:            uuid.New(),
		AlunoID:       alunoID,
		TurmaID:       turmaID,
		Numero:        "001",
		DataMatricula: time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		DataConclusao: nil,
		Status:        model.StatusMatriculaAtiva,
		Observacao:    "Matrícula automática",
	}

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
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	alunoID := uuid.New()

	// Create an aluno
	aluno := &model.Aluno{
		ID:             alunoID,
		InstituicaoID:  uuid.New(),
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
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	alunoID := uuid.New()
	turmaID := uuid.New()
	matriculaID := uuid.New()

	// Create aluno and matrícula, then cancel
	aluno := &model.Aluno{
		ID:             alunoID,
		InstituicaoID:  uuid.New(),
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

	matricula := &model.Matricula{
		ID:            matriculaID,
		AlunoID:       alunoID,
		TurmaID:       turmaID,
		Numero:        "002",
		DataMatricula: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
		DataConclusao: nil,
		Status:        model.StatusMatriculaAtiva,
		Observacao:    "",
	}

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
