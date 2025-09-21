import { useState, useEffect } from 'react';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export interface TotalStats {
  total_hospitais: number;
  total_medicos: number;
  total_leitos: number;
  total_pacientes: number;
}

export interface Stats2 {
  total_estados: number;
  total_municipios: number;
}

export interface Cid10Item {
  cid10: string;
  doenca: string;
  count: number;
}

export interface HospitalAccess {
  hospital_uuid: string;
  hospital_nome: string;
  especialidades: string;
  count: number;
}

export function useTotalStats() {
  const [data, setData] = useState<TotalStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchTotalStats = async () => {
    try {
      setLoading(true);
      setError(null);
      
      console.log('Buscando dados de /stats/totals...');
      const response = await fetch(`${API_BASE_URL}/stats/totals`);
      
      if (!response.ok) {
        throw new Error(`Erro na API: ${response.status} - ${response.statusText}`);
      }
      
      const result = await response.json();
      console.log('Dados recebidos:', result);
      console.log('Dados processados:', result.data);
      setData(result.data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro desconhecido';
      setError(errorMessage);
      console.error('Erro ao buscar dados:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTotalStats();
  }, []);

  return { data, loading, error, refetch: fetchTotalStats };
}

// Hook para Stats2 (estados e municípios)
export function useStats2() {
  const [data, setData] = useState<Stats2 | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStats2 = async () => {
    try {
      setLoading(true);
      setError(null);
      
      console.log('Buscando dados de /stats2...');
      const response = await fetch(`${API_BASE_URL}/stats2`);
      
      if (!response.ok) {
        throw new Error(`Erro na API: ${response.status} - ${response.statusText}`);
      }
      
      const result = await response.json();
      console.log('Dados Stats2 recebidos:', result);
      setData(result.data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro desconhecido';
      setError(errorMessage);
      console.error('Erro ao buscar dados Stats2:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats2();
  }, []);

  return { data, loading, error, refetch: fetchStats2 };
}

// Hook para CID10 mais comuns
export function useCid10MaisComuns() {
  const [data, setData] = useState<Cid10Item[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchCid10 = async () => {
    try {
      setLoading(true);
      setError(null);
      
      console.log('Buscando dados de /stats/cid10-mais-comuns...');
      const response = await fetch(`${API_BASE_URL}/stats/cid10-mais-comuns`);
      
      if (!response.ok) {
        throw new Error(`Erro na API: ${response.status} - ${response.statusText}`);
      }
      
      const result = await response.json();
      console.log('Dados CID10 recebidos:', result);
      setData(result.data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro desconhecido';
      setError(errorMessage);
      console.error('Erro ao buscar dados CID10:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCid10();
  }, []);

  return { data, loading, error, refetch: fetchCid10 };
}

// Hook para hospitais mais acessados
export function useHospitaisMaisAcessados() {
  const [data, setData] = useState<HospitalAccess[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchHospitais = async () => {
    try {
      setLoading(true);
      setError(null);
      
      console.log('Buscando dados de /stats/hospitais-mais-acessados...');
      const response = await fetch(`${API_BASE_URL}/stats/hospitais-mais-acessados`);
      
      if (!response.ok) {
        throw new Error(`Erro na API: ${response.status} - ${response.statusText}`);
      }
      
      const result = await response.json();
      console.log('Dados Hospitais recebidos:', result);
      setData(result.data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro desconhecido';
      setError(errorMessage);
      console.error('Erro ao buscar dados Hospitais:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchHospitais();
  }, []);

  return { data, loading, error, refetch: fetchHospitais };
}
