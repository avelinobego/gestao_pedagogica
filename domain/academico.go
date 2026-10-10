package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode"

	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/google/uuid"
)

var (
	ErrEtapaPossuiTurmas       = errors.New("a etapa possui turmas vinculadas")
	ErrRelacionamentoAcademico = errors.New("os vínculos acadêmicos não correspondem")
)

func ValidarInstituicaoEnsino(instituicao model.InstituicaoEnsino) error {
	if err := validarTexto("nome da instituição", instituicao.Nome, 150); err != nil {
		return err
	}
	if len([]rune(instituicao.NomeFantasia)) > 150 {
		return errors.New("nome fantasia deve ter até 150 caracteres")
	}
	if len(instituicao.CNPJ) > 14 || !somenteDigitos(instituicao.CNPJ) {
		return errors.New("CNPJ deve conter até 14 dígitos")
	}
	if len([]rune(instituicao.CodigoINEP)) > 20 {
		return errors.New("código INEP deve ter até 20 caracteres")
	}
	if len([]rune(instituicao.Email)) > 150 || len([]rune(instituicao.Telefone)) > 30 ||
		len([]rune(instituicao.Endereco)) > 255 || len([]rune(instituicao.Cidade)) > 100 {
		return errors.New("os dados de contato ou endereço excedem o tamanho permitido")
	}
	if instituicao.Email != "" {
		endereco, err := mail.ParseAddress(instituicao.Email)
		if err != nil || endereco.Address != instituicao.Email {
			return errors.New("e-mail da instituição inválido")
		}
	}
	if len(instituicao.UF) != 0 && !ufValida(instituicao.UF) {
		return errors.New("UF deve conter duas letras maiúsculas")
	}
	return nil
}

func ValidarCurso(curso model.Curso) error {
	if curso.InstituicaoID == uuid.Nil {
		return errors.New("instituição do curso é obrigatória")
	}
	if err := validarTexto("nome do curso", curso.Nome, 150); err != nil {
		return err
	}
	if err := validarTexto("código do curso", curso.Codigo, 30); err != nil {
		return err
	}
	return nil
}

func ValidarEtapaEnsino(etapa model.EtapaEnsino) error {
	if etapa.CursoID == uuid.Nil {
		return errors.New("curso da etapa é obrigatório")
	}
	if err := validarTexto("nome da etapa", etapa.Nome, 100); err != nil {
		return err
	}
	if etapa.Ordem <= 0 {
		return errors.New("ordem da etapa deve ser maior que zero")
	}
	if etapa.CargaHoraria <= 0 {
		return errors.New("carga horária da etapa deve ser maior que zero")
	}
	return nil
}

func ValidarDisciplina(disciplina model.Disciplina) error {
	if disciplina.InstituicaoID == uuid.Nil {
		return errors.New("instituição da disciplina é obrigatória")
	}
	if err := validarTexto("nome da disciplina", disciplina.Nome, 150); err != nil {
		return err
	}
	if err := validarTexto("código da disciplina", disciplina.Codigo, 30); err != nil {
		return err
	}
	if disciplina.CargaHoraria <= 0 {
		return errors.New("carga horária da disciplina deve ser maior que zero")
	}
	return nil
}

func ValidarTurma(turma model.Turma) error {
	if turma.InstituicaoID == uuid.Nil || turma.CursoID == uuid.Nil || turma.EtapaID == uuid.Nil {
		return errors.New("instituição, curso e etapa da turma são obrigatórios")
	}
	if err := validarTexto("nome da turma", turma.Nome, 100); err != nil {
		return err
	}
	if err := validarTexto("código da turma", turma.Codigo, 30); err != nil {
		return err
	}
	if turma.AnoLetivo < 2000 {
		return errors.New("ano letivo deve ser igual ou posterior a 2000")
	}
	switch turma.Turno {
	case model.TurnoMatutino, model.TurnoVespertino, model.TurnoNoturno, model.TurnoIntegral:
	default:
		return errors.New("turno da turma inválido")
	}
	switch turma.Status {
	case model.StatusTurmaAtiva, model.StatusTurmaEncerrada, model.StatusTurmaPlanejada:
	default:
		return errors.New("status da turma inválido")
	}
	return nil
}

func validarTexto(campo, valor string, max int) error {
	valor = strings.TrimSpace(valor)
	if valor == "" {
		return fmt.Errorf("%s é obrigatório", campo)
	}
	if len([]rune(valor)) > max {
		return fmt.Errorf("%s deve ter até %d caracteres", campo, max)
	}
	return nil
}

func ufValida(uf string) bool {
	if len(uf) != 2 || uf != strings.ToUpper(uf) {
		return false
	}
	for _, caractere := range uf {
		if caractere < 'A' || caractere > 'Z' {
			return false
		}
	}
	return true
}

func somenteDigitos(valor string) bool {
	for _, caractere := range valor {
		if !unicode.IsDigit(caractere) || caractere > unicode.MaxASCII {
			return false
		}
	}
	return true
}
