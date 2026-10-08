package plataforma_test

import (
	"context"
	"fmt"
	"github.com/avelinobego/gestao_pedagogica/shared/plataforma"
	"testing"

	"github.com/avelinobego/gestao_pedagogica/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestConexao(t *testing.T) {

	alunos := []model.Aluno{}

	ctx, cancel := plataforma.CincoSegudos()
	defer cancel()

	con, err := plataforma.GetConn(ctx)
	if err != nil {
		fmt.Println(err)
	}
	err = con.SelectContext(ctx, &alunos, "SELECT * FROM aluno ORDER BY nome ASC")
	if err != nil {
		t.Log(err)
	}

	plataforma.CloseConn()

	for _, a := range alunos {
		fmt.Printf("Aluno %v\n", a)
	}

	con, err = plataforma.GetConn(ctx)
	if err != nil {
		t.Log(err)
	}

	err = con.SelectContext(ctx, &alunos, "SELECT * FROM aluno ORDER BY nome ASC")
	if err != nil {
		t.Log(err)
	}

}

func TestTransacao(t *testing.T) {

	ctx, cancel := plataforma.UmMinuto()
	defer cancel()

	handler := func(tx *sqlx.Tx, ctx context.Context) (err error) {
		_, err = tx.ExecContext(ctx, "INSERT INTO aluno (nome) VALUES ($1)", "Aluno Teste")
		return
	}

	err := plataforma.Executar(ctx, handler)

	if err != nil {
		t.Log(err)
	}

}
