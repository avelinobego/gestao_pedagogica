CREATE TABLE instituicao_ensino (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(150) NOT NULL,
    nome_fantasia VARCHAR(150) NOT NULL DEFAULT '',
    cnpj VARCHAR(14) UNIQUE,
    codigo_inep VARCHAR(20) UNIQUE,
    email VARCHAR(150) NOT NULL DEFAULT '',
    telefone VARCHAR(30) NOT NULL DEFAULT '',
    endereco VARCHAR(255) NOT NULL DEFAULT '',
    cidade VARCHAR(100) NOT NULL DEFAULT '',
    uf CHAR(2) NOT NULL DEFAULT '',
    ativa BOOLEAN NOT NULL DEFAULT TRUE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (uf = '' OR uf ~ '^[A-Z]{2}$')
);

CREATE TABLE curso (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    instituicao_id UUID NOT NULL REFERENCES instituicao_ensino (id),
    nome VARCHAR(150) NOT NULL,
    codigo VARCHAR(30) NOT NULL,
    descricao TEXT NOT NULL DEFAULT '',
    ativo BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (instituicao_id, codigo)
);

CREATE TABLE etapa_ensino (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    curso_id UUID NOT NULL REFERENCES curso (id) ON DELETE CASCADE,
    nome VARCHAR(100) NOT NULL,
    ordem INTEGER NOT NULL CHECK (ordem > 0),
    carga_horaria INTEGER NOT NULL CHECK (carga_horaria > 0),
    UNIQUE (curso_id, ordem)
);

CREATE TABLE disciplina (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    instituicao_id UUID NOT NULL REFERENCES instituicao_ensino (id),
    nome VARCHAR(150) NOT NULL,
    codigo VARCHAR(30) NOT NULL,
    carga_horaria INTEGER NOT NULL CHECK (carga_horaria > 0),
    ativa BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (instituicao_id, codigo)
);

CREATE TABLE turma (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    instituicao_id UUID NOT NULL REFERENCES instituicao_ensino (id),
    curso_id UUID NOT NULL REFERENCES curso (id),
    etapa_id UUID NOT NULL REFERENCES etapa_ensino (id),
    nome VARCHAR(100) NOT NULL,
    codigo VARCHAR(30) NOT NULL,
    ano_letivo INTEGER NOT NULL CHECK (ano_letivo >= 2000),
    turno VARCHAR(20) NOT NULL CHECK (turno IN ('MATUTINO', 'VESPERTINO', 'NOTURNO', 'INTEGRAL')),
    status VARCHAR(20) NOT NULL DEFAULT 'PLANEJADA' CHECK (status IN ('ATIVA', 'ENCERRADA', 'PLANEJADA')),
    criado_em TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (instituicao_id, codigo, ano_letivo),
    UNIQUE (curso_id, etapa_id, ano_letivo, nome)
);

CREATE INDEX idx_curso_instituicao ON curso (instituicao_id);
CREATE INDEX idx_disciplina_instituicao ON disciplina (instituicao_id);
CREATE INDEX idx_turma_instituicao_ano ON turma (instituicao_id, ano_letivo);
CREATE INDEX idx_turma_curso_etapa ON turma (curso_id, etapa_id);