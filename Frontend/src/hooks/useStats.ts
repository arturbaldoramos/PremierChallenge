import { useState, useEffect } from 'react'
import { api } from '@/lib/api'
import type { 
  StatsTotals,
  StatsByState,
  HospitalBySpecialty,
  HospitalByMunicipality,
  DoctorDistribution,
  Hospital
} from '@/types/stats'

// Hook para estatísticas gerais (totais)
export interface UseStatsReturn {
  data: StatsTotals | null
  isLoading: boolean
  error: string | null
  refetch: () => Promise<void>
}

export function useStats(): UseStatsReturn {
  const [data, setData] = useState<StatsTotals | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchStats = async () => {
    try {
      setIsLoading(true)
      setError(null)
      
      // Timeout para evitar loading infinito
      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), 10000) // 10 segundos
      
      const response = await api.stats.getTotals()
      clearTimeout(timeoutId)
      
      setData(response.data)
    } catch (err) {
      console.error('Erro ao carregar estatísticas:', err)
      
      // Fallback com dados mockados se a API falhar
      setData({
        total_hospitais: 6844,
        total_medicos: 280000,
        total_estados: 27,
        total_municipios: 5570
      })
      
      if (err instanceof Error && err.name === 'AbortError') {
        setError('Timeout: API não respondeu a tempo - usando dados de exemplo')
      } else if (err instanceof Error) {
        setError(`Erro ao carregar dados reais: ${err.message} - usando dados de exemplo`)
      } else {
        setError('Erro desconhecido ao carregar dados reais - usando dados de exemplo')
      }
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchStats()
  }, [])

  return {
    data,
    isLoading,
    error,
    refetch: fetchStats,
  }
}

// Hook para hospitais por estado
export interface UseHospitalsByStateReturn {
  data: StatsByState[] | null
  isLoading: boolean
  error: string | null
  refetch: () => Promise<void>
}

export function useHospitalsByState(): UseHospitalsByStateReturn {
  const [data, setData] = useState<StatsByState[] | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchData = async () => {
    try {
      setIsLoading(true)
      setError(null)
      
      const response = await api.stats.getHospitalsByState()
      setData(response.data)
    } catch (err) {
      console.error('Erro ao carregar hospitais por estado:', err)
      setError('Erro ao carregar dados de hospitais por estado')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  }
}

// Hook para médicos por estado
export interface UseDoctorsByStateReturn {
  data: StatsByState[] | null
  isLoading: boolean
  error: string | null
  refetch: () => Promise<void>
}

export function useDoctorsByState(): UseDoctorsByStateReturn {
  const [data, setData] = useState<StatsByState[] | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchData = async () => {
    try {
      setIsLoading(true)
      setError(null)
      
      const response = await api.stats.getDoctorsByState()
      setData(response.data)
    } catch (err) {
      console.error('Erro ao carregar médicos por estado:', err)
      setError('Erro ao carregar dados de médicos por estado')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  }
}

// Hook para hospitais por especialidade
export interface UseHospitalsBySpecialtyReturn {
  data: HospitalBySpecialty[] | Hospital[] | null
  isLoading: boolean
  error: string | null
  refetch: (specialty?: string) => Promise<void>
}

export function useHospitalsBySpecialty(): UseHospitalsBySpecialtyReturn {
  const [data, setData] = useState<HospitalBySpecialty[] | Hospital[] | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchData = async (specialty?: string) => {
    try {
      setIsLoading(true)
      setError(null)
      
      const response = await api.stats.getHospitalsBySpecialty(specialty)
      setData(response.data)
    } catch (err) {
      console.error('Erro ao carregar hospitais por especialidade:', err)
      setError('Erro ao carregar dados de hospitais por especialidade')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  }
}

// Hook para hospitais por município
export interface UseHospitalsByMunicipalityReturn {
  data: HospitalByMunicipality[] | Hospital[] | null
  isLoading: boolean
  error: string | null
  refetch: (municipality?: string) => Promise<void>
}

export function useHospitalsByMunicipality(): UseHospitalsByMunicipalityReturn {
  const [data, setData] = useState<HospitalByMunicipality[] | Hospital[] | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchData = async (municipality?: string) => {
    try {
      setIsLoading(true)
      setError(null)
      
      const response = await api.stats.getHospitalsByMunicipality(municipality)
      setData(response.data)
    } catch (err) {
      console.error('Erro ao carregar hospitais por município:', err)
      setError('Erro ao carregar dados de hospitais por município')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  }
}

// Hook para distribuição de médicos
export interface UseDoctorDistributionReturn {
  data: DoctorDistribution | null
  isLoading: boolean
  error: string | null
  refetch: () => Promise<void>
}

export function useDoctorDistribution(): UseDoctorDistributionReturn {
  const [data, setData] = useState<DoctorDistribution | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchData = async () => {
    try {
      setIsLoading(true)
      setError(null)
      
      const response = await api.stats.getDoctorDistribution()
      setData(response.data)
    } catch (err) {
      console.error('Erro ao carregar distribuição de médicos:', err)
      setError('Erro ao carregar dados de distribuição de médicos')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  }
}

// Hook para especialidades
export interface UseSpecialtiesReturn {
  data: string[] | null
  isLoading: boolean
  error: string | null
  refetch: () => Promise<void>
}

export function useSpecialties(): UseSpecialtiesReturn {
  const [data, setData] = useState<string[] | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchData = async () => {
    try {
      setIsLoading(true)
      setError(null)
      
      const response = await api.stats.getSpecialties()
      setData(response.data)
    } catch (err) {
      console.error('Erro ao carregar especialidades:', err)
      setError('Erro ao carregar lista de especialidades')
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  return {
    data,
    isLoading,
    error,
    refetch: fetchData,
  }
}