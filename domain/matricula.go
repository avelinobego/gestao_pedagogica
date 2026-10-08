package domain

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/google/uuid"
)

func ValidarMatricula(matricula model.Matricula) error {
	if matricula.AlunoID == uuid.Nil {
		return errors.New("aluno da matrícula é obrigatório")
	}
	if matricula.TurmaID == uuid.Nil {
		return errors.New("turma da matrícula é obrigatória")
	}
	if strings.TrimSpace(matricula.Numero) == "" {
		return errors.New("número da matrícula é obrigatório")
	}
	if utf8.RuneCountInString(matricula.Numero) > 30 {
		return errors.New("número da matrícula não pode exceder 30 caracteres")
	}
	if matricula.DataMatricula.IsZero() {
		return errors.New("data da matrícula é obrigatória")
	}
	if matricula.DataConclusao != nil &&
		dataSemHorario(*matricula.DataConclusao).Before(dataSemHorario(matricula.DataMatricula)) {
		return errors.New("data de conclusão não pode ser anterior à data da matrícula")
	}
	if !statusMatriculaValido(matricula.Status) {
		return errors.New("status da matrícula é inválido")
	}
	return nil
}

func statusMatriculaValido(status model.StatusMatricula) bool {
	switch status {
	case model.StatusMatriculaAtiva,
		model.StatusMatriculaConcluida,
		model.StatusMatriculaCancelada,
		model.StatusMatriculaTrancada:
		return true
	default:
		return false
	}
}

func dataSemHorario(data time.Time) time.Time {
	return time.Date(data.Year(), data.Month(), data.Day(), 0, 0, 0, 0, time.UTC)
}
