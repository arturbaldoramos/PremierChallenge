export interface StatsTotals {
  total_hospitais: number
  total_medicos: number
  total_estados: number
  total_municipios: number
}

export interface StatsResponse {
  data: StatsTotals
}

// Tipos para estatísticas por estado
export interface StatsByState {
  estado: string
  nome_estado: string
  count: number
}

export interface StatsByStateResponse {
  success: boolean
  data: StatsByState[]
}

// Tipos para hospitais por especialidade
export interface HospitalBySpecialty {
  especialidade: string
  count: number
}

export interface HospitalBySpecialtyResponse {
  success: boolean
  data: HospitalBySpecialty[]
}

// Tipos para hospitais por município
export interface HospitalByMunicipality {
  municipio: string
  count: number
}

export interface HospitalByMunicipalityResponse {
  success: boolean
  data: HospitalByMunicipality[]
}

// Tipos para distribuição de médicos
export interface DoctorDistribution {
  com_um_hospital: number
  com_dois_hospitais: number
  com_tres_hospitais: number
  sem_hospitais: number
}

export interface DoctorDistributionResponse {
  success: boolean
  data: DoctorDistribution
  message?: string
}

// Tipos para lista de especialidades
export interface SpecialtiesResponse {
  success: boolean
  data: string[]
}

// Tipos para hospitais filtrados por especialidade
export interface Hospital {
  id: number
  uuid: string
  nome: string
  cep: string
  especialidades: string
  leitos_totais: number
  cod_municipio: string
  bairro: string
}

export interface HospitalBySpecialtyFilterResponse {
  success: boolean
  data: Hospital[]
  message?: string
}

// Tipos para hospitais filtrados por município
export interface HospitalByMunicipalityFilterResponse {
  success: boolean
  data: Hospital[]
  message?: string
}