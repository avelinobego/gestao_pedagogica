-- Dados demonstrativos para desenvolvimento local.

INSERT INTO instituicao_ensino (
    id, nome, nome_fantasia, cnpj, codigo_inep, email, telefone,
    endereco, cidade, uf, ativa
)
VALUES (
    '90000000-0000-0000-0000-000000000001',
    'Escola Municipal Horizonte',
    'Escola Horizonte',
    '12345678000190',
    '12345678',
    'contato@escolahorizonte.edu.br',
    '(11) 3333-4444',
    'Rua das Palmeiras, 100',
    'Sao Paulo',
    'SP',
    TRUE
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO curso (id, instituicao_id, nome, codigo, descricao, ativo)
VALUES (
    '91000000-0000-0000-0000-000000000001',
    '90000000-0000-0000-0000-000000000001',
    'Ensino Fundamental',
    'EF',
    'Curso regular do ensino fundamental.',
    TRUE
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO etapa_ensino (id, curso_id, nome, ordem, carga_horaria)
VALUES (
    '92000000-0000-0000-0000-000000000001',
    '91000000-0000-0000-0000-000000000001',
    '6o Ano',
    1,
    800
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO disciplina (id, instituicao_id, nome, codigo, carga_horaria, ativa)
VALUES
    (
        '40000000-0000-0000-0000-000000000001',
        '90000000-0000-0000-0000-000000000001',
        'Matematica',
        'MAT',
        160,
        TRUE
    ),
    (
        '40000000-0000-0000-0000-000000000002',
        '90000000-0000-0000-0000-000000000001',
        'Ciencias',
        'CIE',
        120,
        TRUE
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO turma (
    id, instituicao_id, curso_id, etapa_id, nome, codigo,
    ano_letivo, turno, status
)
VALUES (
    '30000000-0000-0000-0000-000000000001',
    '90000000-0000-0000-0000-000000000001',
    '91000000-0000-0000-0000-000000000001',
    '92000000-0000-0000-0000-000000000001',
    '6o Ano A',
    '6A-2026',
    2026,
    'MATUTINO',
    'ATIVA'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO aluno (
    id, instituicao_id, nome, cpf, rg, data_nascimento,
    email, telefone, endereco, cidade, uf, status
)
VALUES
    (
        '60000000-0000-0000-0000-000000000001',
        '90000000-0000-0000-0000-000000000001',
        'Ana Beatriz Souza',
        '12345678901',
        '123456789',
        '2014-05-12',
        'ana.souza@example.com',
        '(11) 98888-0001',
        'Rua das Flores, 10',
        'Sao Paulo',
        'SP',
        'ATIVO'
    ),
    (
        '60000000-0000-0000-0000-000000000002',
        '90000000-0000-0000-0000-000000000001',
        'Bruno Henrique Lima',
        '23456789012',
        '234567890',
        '2014-09-21',
        'bruno.lima@example.com',
        '(11) 98888-0002',
        'Avenida Central, 25',
        'Sao Paulo',
        'SP',
        'ATIVO'
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO responsavel (
    id, nome, cpf, email, telefone, endereco
)
VALUES
    (
        '93000000-0000-0000-0000-000000000001',
        'Carla Souza',
        '34567890123',
        'carla.souza@example.com',
        '(11) 97777-0001',
        'Rua das Flores, 10'
    ),
    (
        '93000000-0000-0000-0000-000000000002',
        'Daniel Lima',
        '45678901234',
        'daniel.lima@example.com',
        '(11) 97777-0002',
        'Avenida Central, 25'
    )
ON CONFLICT (id) DO NOTHING;

INSERT INTO aluno_responsavel (aluno_id, responsavel_id, parentesco, principal)
VALUES
    (
        '60000000-0000-0000-0000-000000000001',
        '93000000-0000-0000-0000-000000000001',
        'Mae',
        TRUE
    ),
    (
        '60000000-0000-0000-0000-000000000002',
        '93000000-0000-0000-0000-000000000002',
        'Pai',
        TRUE
    )
ON CONFLICT (aluno_id, responsavel_id) DO NOTHING;

INSERT INTO matricula (
    id, aluno_id, turma_id, numero, data_matricula, status, observacao
)
VALUES
    (
        '94000000-0000-0000-0000-000000000001',
        '60000000-0000-0000-0000-000000000001',
        '30000000-0000-0000-0000-000000000001',
        '2026-0001',
        '2026-01-15',
        'ATIVA',
        ''
    ),
    (
        '94000000-0000-0000-0000-000000000002',
        '60000000-0000-0000-0000-000000000002',
        '30000000-0000-0000-0000-000000000001',
        '2026-0002',
        '2026-01-16',
        'ATIVA',
        ''
    )
ON CONFLICT (id) DO NOTHING;
