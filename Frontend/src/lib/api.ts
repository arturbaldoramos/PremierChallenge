// Configuração base da API
const API_BASE_URL = 'http://localhost:8080/api/v1/upload/'

// Tipos para as respostas da API
export interface ApiResponse<T> {
  data: T
  message?: string
  success: boolean
}

export interface ApiErrorResponse {
  message: string
  status: number
  error?: string
}

// Classe de erro personalizada
class ApiError extends Error {
  public status: number
  public error?: string

  constructor(message: string, status: number, error?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.error = error
  }
}

// Classe para gerenciar requisições HTTP
class ApiClient {
  private baseUrl: string

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl
  }

  // Método privado para fazer requisições HTTP
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseUrl}${endpoint.startsWith('/') ? endpoint.slice(1) : endpoint}`
    
    const defaultHeaders = {
      'Content-Type': 'application/json',
    }

    const config: RequestInit = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    }

    try {
      const response = await fetch(url, config)
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new ApiError(
          errorData.message || `HTTP Error: ${response.status}`,
          response.status,
          errorData.error
        )
      }

      const data = await response.json()
      return {
        data,
        success: true,
        message: data.message
      }
    } catch (error) {
      if (error instanceof ApiError) {
        throw error
      }
      
      // Erro de rede ou parsing
      throw new ApiError(
        error instanceof Error ? error.message : 'Erro desconhecido',
        0
      )
    }
  }

  // GET - Buscar dados
  async get<T>(entity: string, id?: string | number): Promise<ApiResponse<T>> {
    const endpoint = id ? `/${entity}/${id}` : `/${entity}`
    return this.request<T>(endpoint, {
      method: 'GET',
    })
  }

  // POST - Criar novo registro ou upload de arquivo
  async post<T>(entity: string, data: any): Promise<ApiResponse<T>> {
    // Se data é um File, criar FormData
    if (data instanceof File) {
      const formData = new FormData()
      formData.append('file', data)
      
      return this.request<T>(`/${entity}`, {
        method: 'POST',
        headers: {}, // Remove Content-Type para FormData
        body: formData,
      })
    }
    
    // Caso contrário, enviar como JSON
    return this.request<T>(`/${entity}`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // PUT - Atualizar registro existente
  async put<T>(entity: string, id: string | number, data: any): Promise<ApiResponse<T>> {
    return this.request<T>(`/${entity}/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  }

  // DELETE - Deletar registro
  async delete<T>(entity: string, id: string | number): Promise<ApiResponse<T>> {
    return this.request<T>(`/${entity}/${id}`, {
      method: 'DELETE',
    })
  }

  // PATCH - Atualização parcial
  async patch<T>(entity: string, id: string | number, data: any): Promise<ApiResponse<T>> {
    return this.request<T>(`/${entity}/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  }
}

// Instância singleton do cliente API
export const apiClient = new ApiClient()

// Funções utilitárias para entidades específicas
export const api = {
  // Hospitais
  hospitals: {
    getAll: () => apiClient.get('hospital'),
    getById: (id: string | number) => apiClient.get('hospital', id),
    create: (data: any) => apiClient.post('hospital', data),
    update: (id: string | number, data: any) => apiClient.put('hospital', id, data),
    delete: (id: string | number) => apiClient.delete('hospital', id),
    uploadFile: (file: File) => apiClient.post('hospital', file),
  },

  // Médicos
  doctors: {
    getAll: () => apiClient.get('medico'),
    getById: (id: string | number) => apiClient.get('medico', id),
    create: (data: any) => apiClient.post('medico', data),
    update: (id: string | number, data: any) => apiClient.put('medico', id, data),
    delete: (id: string | number) => apiClient.delete('medico', id),
    uploadFile: (file: File) => apiClient.post('medico', file),
  },

  // Estados
  states: {
    getAll: () => apiClient.get('estado'),
    getById: (id: string | number) => apiClient.get('estado', id),
    create: (data: any) => apiClient.post('estado', data),
    update: (id: string | number, data: any) => apiClient.put('estado', id, data),
    delete: (id: string | number) => apiClient.delete('estado', id),
    uploadFile: (file: File) => apiClient.post('estado', file),
  },

  // Municípios
  municipalities: {
    getAll: () => apiClient.get('municipio'),
    getById: (id: string | number) => apiClient.get('municipio', id),
    create: (data: any) => apiClient.post('municipio', data),
    update: (id: string | number, data: any) => apiClient.put('municipio', id, data),
    delete: (id: string | number) => apiClient.delete('municipio', id),
    uploadFile: (file: File) => apiClient.post('municipio', file),
  },

  // Pacientes
  patients: {
    getAll: () => apiClient.get('paciente'),
    getById: (id: string | number) => apiClient.get('paciente', id),
    create: (data: any) => apiClient.post('paciente', data),
    update: (id: string | number, data: any) => apiClient.put('paciente', id, data),
    delete: (id: string | number) => apiClient.delete('paciente', id),
    uploadFile: (file: File) => apiClient.post('paciente', file),
  },

  // CID (Classificação Internacional de Doenças)
  cid: {
    getAll: () => apiClient.get('cid'),
    getById: (id: string | number) => apiClient.get('cid', id),
    create: (data: any) => apiClient.post('cid', data),
    update: (id: string | number, data: any) => apiClient.put('cid', id, data),
    delete: (id: string | number) => apiClient.delete('cid', id),
    uploadFile: (file: File) => apiClient.post('cid', file),
  },
}

// Função genérica para criar endpoints para novas entidades
export const createEntityApi = (entityName: string) => ({
  getAll: () => apiClient.get(entityName),
  getById: (id: string | number) => apiClient.get(entityName, id),
  create: (data: any) => apiClient.post(entityName, data),
  update: (id: string | number, data: any) => apiClient.put(entityName, id, data),
  delete: (id: string | number) => apiClient.delete(entityName, id),
  uploadFile: (file: File) => apiClient.post(entityName, file),
})

// Exportar tipos e classes
export { ApiClient, ApiError }
export default apiClient
