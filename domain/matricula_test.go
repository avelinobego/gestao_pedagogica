package domain

import (
	"testing"
	"time"

	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/google/uuid"
)

func TestValidarMatricula(t *testing.T) {
	dataMatricula := time.Date(2026, time.March, 10, 0, 0, 0, 0, time.UTC)
	dataConclusao := time.Date(2026, time.March, 20, 0, 0, 0, 0, time.UTC)
	matriculaValida := model.Matricula{
		AlunoID:       uuid.New(),
		TurmaID:       uuid.New(),
		Numero:        "MAT-001",
		DataMatricula: dataMatricula,
		DataConclusao: &dataConclusao,
		Status:        model.StatusMatriculaAtiva,
	}

	testes := []struct {
		nome      string
		modificar func(*model.Matricula)
		erro      string
	}{
		{
			nome: "válida",
		},
		{
			nome: "aluno obrigatório",
			modificar: func(matricula *model.Matricula) {
				matricula.AlunoID = uuid.Nil
			},
			erro: "aluno da matrícula é obrigatório",
		},
		{
			nome: "turma obrigatória",
			modificar: func(matricula *model.Matricula) {
				matricula.TurmaID = uuid.Nil
			},
			erro: "turma da matrícula é obrigatória",
		},
		{
			nome: "número obrigatório",
			modificar: func(matricula *model.Matricula) {
				matricula.Numero = " \t "
			},
			erro: "número da matrícula é obrigatório",
		},
		{
			nome: "número limitado a 30 caracteres",
			modificar: func(matricula *model.Matricula) {
				matricula.Numero = "1234567890123456789012345678901"
			},
			erro: "número da matrícula não pode exceder 30 caracteres",
		},
		{
			nome: "data da matrícula obrigatória",
			modificar: func(matricula *model.Matricula) {
				matricula.DataMatricula = time.Time{}
			},
			erro: "data da matrícula é obrigatória",
		},
		{
			nome: "conclusão não pode anteceder matrícula",
			modificar: func(matricula *model.Matricula) {
				data := time.Date(2026, time.March, 9, 23, 59, 0, 0, time.UTC)
				matricula.DataConclusao = &data
			},
			erro: "data de conclusão não pode ser anterior à data da matrícula",
		},
		{
			nome: "status precisa ser permitido",
			modificar: func(matricula *model.Matricula) {
				matricula.Status = "DESCONHECIDO"
			},
			erro: "status da matrícula é inválido",
		},
	}

	for _, teste := range testes {
		t.Run(teste.nome, func(t *testing.T) {
			matricula := matriculaValida
			if teste.modificar != nil {
				teste.modificar(&matricula)
			}

			err := ValidarMatricula(matricula)
			if teste.erro == "" && err != nil {
				t.Fatalf("esperava matrícula válida, recebeu erro: %v", err)
			}
			if teste.erro != "" && (err == nil || err.Error() != teste.erro) {
				t.Fatalf("esperava erro %q, recebeu %v", teste.erro, err)
			}
		})
	}
}

func TestValidarMatriculaPermiteDatasDeConclusaoNoMesmoDia(t *testing.T) {
	dataMatricula := time.Date(2026, time.March, 10, 18, 0, 0, 0, time.FixedZone("UTC-3", -3*60*60))
	dataConclusao := time.Date(2026, time.March, 10, 8, 0, 0, 0, time.UTC)
	matricula := model.Matricula{
		AlunoID:       uuid.New(),
		TurmaID:       uuid.New(),
		Numero:        "MAT-002",
		DataMatricula: dataMatricula,
		DataConclusao: &dataConclusao,
		Status:        model.StatusMatriculaConcluida,
	}

	if err := ValidarMatricula(matricula); err != nil {
		t.Fatalf("esperava datas no mesmo dia como válidas, recebeu erro: %v", err)
	}
}

func TestValidarMatriculaPermiteStatusConhecidos(t *testing.T) {
	statuses := []model.StatusMatricula{
		model.StatusMatriculaAtiva,
		model.StatusMatriculaConcluida,
		model.StatusMatriculaCancelada,
		model.StatusMatriculaTrancada,
	}
	for _, status := range statuses {
		t.Run(string(status), func(t *testing.T) {
			matricula := model.Matricula{
				AlunoID:       uuid.New(),
				TurmaID:       uuid.New(),
				Numero:        "MAT-003",
				DataMatricula: time.Date(2026, time.March, 10, 0, 0, 0, 0, time.UTC),
				Status:        status,
			}
			if err := ValidarMatricula(matricula); err != nil {
				t.Fatalf("esperava status %q como válido, recebeu erro: %v", status, err)
			}
		})
	}
}
