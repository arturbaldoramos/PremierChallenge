-- Migração para adicionar campos necessários para o dashboard
-- Executar este script após as migrações básicas

-- Atualizar tabela pacientes
ALTER TABLE pacientes 
ADD COLUMN IF NOT EXISTS codigo UUID DEFAULT gen_random_uuid(),
ADD COLUMN IF NOT EXISTS cod_municipio VARCHAR(10),
ADD COLUMN IF NOT EXISTS bairro VARCHAR(100),
ADD COLUMN IF NOT EXISTS cid_codigo VARCHAR(10);

-- Criar índices para melhor performance
CREATE INDEX IF NOT EXISTS idx_pacientes_cod_municipio ON pacientes(cod_municipio);
CREATE INDEX IF NOT EXISTS idx_pacientes_cid_codigo ON pacientes(cid_codigo);
CREATE INDEX IF NOT EXISTS idx_pacientes_convenio ON pacientes(convenio);

-- Atualizar tabela hospitais
ALTER TABLE hospitais 
ADD COLUMN IF NOT EXISTS especialidade_primaria VARCHAR(100),
ADD COLUMN IF NOT EXISTS especialidade_secundaria VARCHAR(100),
ADD COLUMN IF NOT EXISTS especialidade_terciaria VARCHAR(100),
ADD COLUMN IF NOT EXISTS especialidade_quaternaria VARCHAR(100);

-- Atualizar tabela cid10s
ALTER TABLE cid10s 
ADD COLUMN IF NOT EXISTS especialidade VARCHAR(100);

-- Criar índices para melhor performance
CREATE INDEX IF NOT EXISTS idx_cid10s_especialidade ON cid10s(especialidade);

-- Atualizar tabela municipios para incluir população se não existir
ALTER TABLE municipios 
ADD COLUMN IF NOT EXISTS populacao INTEGER DEFAULT 0;

-- Criar índices para coordenadas geográficas
CREATE INDEX IF NOT EXISTS idx_municipios_coordenadas ON municipios USING GIST (ST_Point(longitude::float, latitude::float));
CREATE INDEX IF NOT EXISTS idx_estados_coordenadas ON estados USING GIST (ST_Point(longitude::float, latitude::float));

-- Atualizar encoding para UTF-8 (se necessário)
-- ALTER DATABASE premiere_challenge SET client_encoding = 'UTF8';

-- Criar view para estatísticas de doenças por região
CREATE OR REPLACE VIEW v_doencas_por_regiao AS
SELECT 
    e.regiao,
    c.codigo as cid_codigo,
    c.descricao as cid_descricao,
    c.especialidade,
    COUNT(p.id) as total_pacientes,
    SUM(m.populacao) as total_populacao_regiao
FROM estados e
INNER JOIN municipios m ON m.codigo_uf = e.codigo
INNER JOIN pacientes p ON p.cod_municipio = m.codigo
INNER JOIN cid10s c ON c.codigo = p.cid_codigo
GROUP BY e.regiao, c.codigo, c.descricao, c.especialidade;

-- Criar view para médicos próximos a hospitais
CREATE OR REPLACE VIEW v_medicos_proximos AS
SELECT 
    h.uuid as hospital_uuid,
    h.nome as hospital_nome,
    m.uuid as medico_uuid,
    m.nome as medico_nome,
    m.especialidade as medico_especialidade,
    ST_Distance(
        ST_Point(mun_med.longitude::float, mun_med.latitude::float)::geography,
        ST_Point(mun_hosp.longitude::float, mun_hosp.latitude::float)::geography
    ) / 1000 as distancia_km
FROM hospitais h
INNER JOIN municipios mun_hosp ON mun_hosp.codigo = h.cod_municipio
CROSS JOIN medicos m
INNER JOIN municipios mun_med ON mun_med.codigo = m.cod_municipio
WHERE ST_DWithin(
    ST_Point(mun_med.longitude::float, mun_med.latitude::float)::geography,
    ST_Point(mun_hosp.longitude::float, mun_hosp.latitude::float)::geography,
    30000 -- 30km em metros
)
AND (
    m.especialidade = ANY(string_to_array(h.especialidades::text, ','))
    OR m.especialidade = 'Clínica Geral'
);

-- Criar função para calcular distância Haversine
CREATE OR REPLACE FUNCTION calcular_distancia_haversine(
    lat1 FLOAT,
    lon1 FLOAT,
    lat2 FLOAT,
    lon2 FLOAT
) RETURNS FLOAT AS $$
DECLARE
    raio_terra FLOAT := 6371; -- Raio da Terra em km
    dlat FLOAT;
    dlon FLOAT;
    a FLOAT;
    c FLOAT;
BEGIN
    -- Converter graus para radianos
    lat1 := lat1 * PI() / 180;
    lat2 := lat2 * PI() / 180;
    dlat := (lat2 - lat1) * PI() / 180;
    dlon := (lon2 - lon1) * PI() / 180;
    
    -- Fórmula de Haversine
    a := SIN(dlat/2) * SIN(dlat/2) + 
         COS(lat1) * COS(lat2) * 
         SIN(dlon/2) * SIN(dlon/2);
    c := 2 * ATAN2(SQRT(a), SQRT(1-a));
    
    RETURN raio_terra * c;
END;
$$ LANGUAGE plpgsql;

-- Criar função para buscar médicos dentro de um raio
CREATE OR REPLACE FUNCTION buscar_medicos_proximos(
    hospital_uuid_param UUID,
    raio_km FLOAT DEFAULT 30
) RETURNS TABLE (
    medico_uuid UUID,
    medico_nome VARCHAR,
    medico_especialidade VARCHAR,
    distancia_km FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        m.uuid,
        m.nome,
        m.especialidade,
        calcular_distancia_haversine(
            mun_hosp.latitude::float,
            mun_hosp.longitude::float,
            mun_med.latitude::float,
            mun_med.longitude::float
        ) as distancia
    FROM hospitais h
    INNER JOIN municipios mun_hosp ON mun_hosp.codigo = h.cod_municipio
    CROSS JOIN medicos m
    INNER JOIN municipios mun_med ON mun_med.codigo = m.cod_municipio
    WHERE h.uuid = hospital_uuid_param
    AND calcular_distancia_haversine(
        mun_hosp.latitude::float,
        mun_hosp.longitude::float,
        mun_med.latitude::float,
        mun_med.longitude::float
    ) <= raio_km
    AND (
        m.especialidade = ANY(string_to_array(h.especialidades::text, ','))
        OR m.especialidade = 'Clínica Geral'
    )
    ORDER BY distancia;
END;
$$ LANGUAGE plpgsql;
