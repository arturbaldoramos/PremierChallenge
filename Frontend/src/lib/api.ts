// Configuração base da API
const API_BASE_URL = '/api/v1/upload/'

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

  // POST - Criar novo registro
  async post<T>(entity: string, data: any): Promise<ApiResponse<T>> {
    return this.request<T>(`/${entity}`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
  }

  // POST - Upload de arquivo
  async uploadFile<T>(entity: string, file: File): Promise<ApiResponse<T>> {
    const formData = new FormData()
    formData.append('file', file)

    return this.request<T>(`/${entity}`, {
      method: 'POST',
      headers: {}, // Remove Content-Type para FormData
      body: formData,
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
    getAll: () => apiClient.get('hospitals'),
    getById: (id: string | number) => apiClient.get('hospitals', id),
    create: (data: any) => apiClient.post('hospitals', data),
    update: (id: string | number, data: any) => apiClient.put('hospitals', id, data),
    delete: (id: string | number) => apiClient.delete('hospitals', id),
    uploadFile: (file: File) => apiClient.uploadFile('hospitals', file),
  },

  // Médicos
  doctors: {
    getAll: () => apiClient.get('doctors'),
    getById: (id: string | number) => apiClient.get('doctors', id),
    create: (data: any) => apiClient.post('doctors', data),
    update: (id: string | number, data: any) => apiClient.put('doctors', id, data),
    delete: (id: string | number) => apiClient.delete('doctors', id),
    uploadFile: (file: File) => apiClient.uploadFile('doctors', file),
  },

  // Estados
  states: {
    getAll: () => apiClient.get('states'),
    getById: (id: string | number) => apiClient.get('states', id),
    create: (data: any) => apiClient.post('states', data),
    update: (id: string | number, data: any) => apiClient.put('states', id, data),
    delete: (id: string | number) => apiClient.delete('states', id),
    uploadFile: (file: File) => apiClient.uploadFile('states', file),
  },

  // Cidades
  cities: {
    getAll: () => apiClient.get('cities'),
    getById: (id: string | number) => apiClient.get('cities', id),
    create: (data: any) => apiClient.post('cities', data),
    update: (id: string | number, data: any) => apiClient.put('cities', id, data),
    delete: (id: string | number) => apiClient.delete('cities', id),
    uploadFile: (file: File) => apiClient.uploadFile('cities', file),
  },

  // Pacientes
  patients: {
    getAll: () => apiClient.get('patients'),
    getById: (id: string | number) => apiClient.get('patients', id),
    create: (data: any) => apiClient.post('patients', data),
    update: (id: string | number, data: any) => apiClient.put('patients', id, data),
    delete: (id: string | number) => apiClient.delete('patients', id),
    uploadFile: (file: File) => apiClient.uploadFile('patients', file),
  },

  // Equipamentos Médicos
  medicalEquipment: {
    getAll: () => apiClient.get('medical-equipment'),
    getById: (id: string | number) => apiClient.get('medical-equipment', id),
    create: (data: any) => apiClient.post('medical-equipment', data),
    update: (id: string | number, data: any) => apiClient.put('medical-equipment', id, data),
    delete: (id: string | number) => apiClient.delete('medical-equipment', id),
    uploadFile: (file: File) => apiClient.uploadFile('medical-equipment', file),
  },
}

// Função genérica para criar endpoints para novas entidades
export const createEntityApi = (entityName: string) => ({
  getAll: () => apiClient.get(entityName),
  getById: (id: string | number) => apiClient.get(entityName, id),
  create: (data: any) => apiClient.post(entityName, data),
  update: (id: string | number, data: any) => apiClient.put(entityName, id, data),
  delete: (id: string | number) => apiClient.delete(entityName, id),
  uploadFile: (file: File) => apiClient.uploadFile(entityName, file),
})

// Exportar tipos e classes
export { ApiClient, ApiError }
export default apiClient
