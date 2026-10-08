CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE periodo_letivo (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(100) NOT NULL,
    ordem INTEGER NOT NULL CHECK (ordem > 0),
    data_inicio TIMESTAMPTZ NOT NULL,
    data_fim TIMESTAMPTZ NOT NULL,
    CHECK (data_fim >= data_inicio),
    UNIQUE (ordem)
);

CREATE TABLE consolidado_periodo (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    aluno_id UUID NOT NULL,
    disciplina_id UUID NOT NULL,
    periodo_letivo_id UUID NOT NULL REFERENCES periodo_letivo (id),
    media_calculada INTEGER NOT NULL CHECK (media_calculada >= 0),
    total_ausencias INTEGER NOT NULL DEFAULT 0 CHECK (total_ausencias >= 0),
    calculado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (aluno_id, disciplina_id, periodo_letivo_id)
);


CREATE TABLE avaliacao (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    turma_id UUID NOT NULL,
    disciplina_id UUID NOT NULL,
    periodo_letivo_id UUID NOT NULL REFERENCES periodo_letivo (id),
    titulo VARCHAR(255) NOT NULL,
    descricao TEXT NOT NULL DEFAULT '',
    tipo VARCHAR(30) NOT NULL CHECK (tipo IN ('PROVA', 'TRABALHO', 'ATIVIDADE', 'EXAME_FINAL')),
    peso INTEGER NOT NULL CHECK (peso > 0),
    nota_maxima INTEGER NOT NULL CHECK (nota_maxima > 0),
    data_aplicacao TIMESTAMPTZ NOT NULL,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_avaliacao_turma_disciplina
    ON avaliacao (turma_id, disciplina_id);

CREATE INDEX idx_avaliacao_periodo_letivo
    ON avaliacao (periodo_letivo_id);

CREATE TABLE nota_registro (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    avaliacao_id UUID NOT NULL REFERENCES avaliacao (id) ON DELETE CASCADE,
    aluno_id UUID NOT NULL,
    valor INTEGER NOT NULL CHECK (valor >= 0),
    observacao TEXT NOT NULL DEFAULT '',
    lancado_por UUID NOT NULL,
    lancado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (avaliacao_id, aluno_id)
);

CREATE INDEX idx_nota_registro_aluno
    ON nota_registro (aluno_id);

CREATE INDEX idx_consolidado_periodo_aluno
    ON consolidado_periodo (aluno_id);