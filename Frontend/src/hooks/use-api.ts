import { useState, useEffect } from 'react';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export interface TotalStats {
  total_hospitais: number;
  total_medicos: number;
  total_leitos: number;
  total_pacientes: number;
}

export function useTotalStats() {
  const [data, setData] = useState<TotalStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTotalStats = async () => {
      try {
        setLoading(true);
        setError(null);
        
        const response = await fetch(`${API_BASE_URL}/stats/totals`);
        
        if (!response.ok) {
          throw new Error(`Erro na API: ${response.status}`);
        }
        
        const result = await response.json();
        setData(result);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Erro desconhecido');
        console.error('Erro ao buscar dados:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchTotalStats();
  }, []);

  return { data, loading, error };
}
