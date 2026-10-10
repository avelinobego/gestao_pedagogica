package web

import (
	"fmt"
	"net/http"
)

func NewHandler(matriculas MatriculaUseCases, academico AcademicoUseCases) (http.Handler, error) {
	matriculasHandler, err := NewMatriculasHandler(matriculas)
	if err != nil {
		return nil, err
	}
	academicoHandler, err := NewAcademicoHandler(academico)
	if err != nil {
		return nil, fmt.Errorf("falha ao configurar módulo acadêmico: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/academico", academicoHandler)
	mux.Handle("/academico/", academicoHandler)
	mux.Handle("/", matriculasHandler)
	return mux, nil
}
