CREATE TABLE aluno (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    instituicao_id UUID NOT NULL REFERENCES instituicao_ensino (id),
    nome VARCHAR(150) NOT NULL,
    nome_social VARCHAR(150) NOT NULL DEFAULT '',
    cpf VARCHAR(11) UNIQUE,
    rg VARCHAR(30) NOT NULL DEFAULT '',
    data_nascimento DATE NOT NULL,
    email VARCHAR(150) NOT NULL DEFAULT '',
    telefone VARCHAR(30) NOT NULL DEFAULT '',
    endereco VARCHAR(255) NOT NULL DEFAULT '',
    cidade VARCHAR(100) NOT NULL DEFAULT '',
    uf CHAR(2) NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'ATIVO' CHECK (status IN ('ATIVO', 'INATIVO')),
    criado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (uf = '' OR uf ~ '^[A-Z]{2}$')
);

CREATE TABLE responsavel (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(150) NOT NULL,
    cpf VARCHAR(11) UNIQUE,
    email VARCHAR(150) NOT NULL DEFAULT '',
    telefone VARCHAR(30) NOT NULL DEFAULT '',
    endereco VARCHAR(255) NOT NULL DEFAULT '',
    criado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE aluno_responsavel (
    aluno_id UUID NOT NULL REFERENCES aluno (id) ON DELETE CASCADE,
    responsavel_id UUID NOT NULL REFERENCES responsavel (id) ON DELETE CASCADE,
    parentesco VARCHAR(50) NOT NULL,
    principal BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (aluno_id, responsavel_id)
);

CREATE TABLE matricula (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    aluno_id UUID NOT NULL REFERENCES aluno (id),
    turma_id UUID NOT NULL REFERENCES turma (id),
    numero VARCHAR(30) NOT NULL,
    data_matricula DATE NOT NULL,
    data_conclusao DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'ATIVA' CHECK (status IN ('ATIVA', 'CONCLUIDA', 'CANCELADA', 'TRANCADA')),
    observacao TEXT NOT NULL DEFAULT '',
    criado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (data_conclusao IS NULL OR data_conclusao >= data_matricula),
    UNIQUE (turma_id, numero),
    UNIQUE (aluno_id, turma_id)
);

CREATE INDEX idx_aluno_instituicao ON aluno (instituicao_id);
CREATE INDEX idx_aluno_nome ON aluno (nome);
CREATE INDEX idx_aluno_responsavel_responsavel ON aluno_responsavel (responsavel_id);
CREATE INDEX idx_matricula_aluno ON matricula (aluno_id);
CREATE INDEX idx_matricula_turma ON matricula (turma_id);