import React, { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { api } from '@/lib/api'

export function ApiDebug() {
  const [testResults, setTestResults] = useState<{
    totals: any
    hospitalsByState: any
    doctorsByState: any
    specialties: any
    doctorDistribution: any
    hospitalsBySpecialty: any
    hospitalsByMunicipality: any
  }>({
    totals: null,
    hospitalsByState: null,
    doctorsByState: null,
    specialties: null,
    doctorDistribution: null,
    hospitalsBySpecialty: null,
    hospitalsByMunicipality: null,
  })

  const [loading, setLoading] = useState<Record<string, boolean>>({})
  const [errors, setErrors] = useState<Record<string, string>>({})

  const testEndpoint = async (endpointName: string, testFunction: () => Promise<any>) => {
    setLoading(prev => ({ ...prev, [endpointName]: true }))
    setErrors(prev => ({ ...prev, [endpointName]: '' }))
    
    try {
      console.log(`Testando endpoint: ${endpointName}`)
      const result = await testFunction()
      console.log(`Sucesso no endpoint ${endpointName}:`, result)
      
      setTestResults(prev => ({ ...prev, [endpointName]: result }))
    } catch (error) {
      console.error(`Erro no endpoint ${endpointName}:`, error)
      const errorMessage = error instanceof Error ? error.message : 'Erro desconhecido'
      setErrors(prev => ({ ...prev, [endpointName]: errorMessage }))
    } finally {
      setLoading(prev => ({ ...prev, [endpointName]: false }))
    }
  }

  const testAllEndpoints = async () => {
    await Promise.all([
      testEndpoint('totals', () => api.stats.getTotals()),
      testEndpoint('hospitalsByState', () => api.stats.getHospitalsByState()),
      testEndpoint('doctorsByState', () => api.stats.getDoctorsByState()),
      testEndpoint('specialties', () => api.stats.getSpecialties()),
      testEndpoint('doctorDistribution', () => api.stats.getDoctorDistribution()),
      testEndpoint('hospitalsBySpecialty', () => api.stats.getHospitalsBySpecialty()),
      testEndpoint('hospitalsByMunicipality', () => api.stats.getHospitalsByMunicipality()),
    ])
  }

  const testDirectFetch = async () => {
    try {
      console.log('Testando fetch direto...')
      const response = await fetch('/api/v1/stats/totals')
      console.log('Response status:', response.status)
      console.log('Response headers:', Object.fromEntries(response.headers.entries()))
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }
      
      const data = await response.json()
      console.log('Dados recebidos:', data)
      
      setTestResults(prev => ({ ...prev, totals: data }))
    } catch (error) {
      console.error('Erro no fetch direto:', error)
      setErrors(prev => ({ ...prev, totals: error instanceof Error ? error.message : 'Erro desconhecido' }))
    }
  }

  return (
    <div className="space-y-6 p-6">
      <Card>
        <CardHeader>
          <CardTitle>Debug da API - Endpoints de Estatísticas</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-2">
            <Button onClick={testAllEndpoints} disabled={Object.values(loading).some(Boolean)}>
              Testar Todos os Endpoints
            </Button>
            <Button onClick={testDirectFetch} variant="outline">
              Testar Fetch Direto
            </Button>
          </div>

          {/* Teste de conectividade básica */}
          <div className="space-y-2">
            <h3 className="font-semibold">Teste de Conectividade</h3>
            <Button 
              onClick={() => testEndpoint('totals', () => api.stats.getTotals())}
              disabled={loading.totals}
              variant="outline"
              size="sm"
            >
              {loading.totals ? 'Testando...' : 'Testar /stats/totals'}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Resultados dos testes */}
      <div className="grid gap-4 md:grid-cols-2">
        {Object.entries(testResults).map(([endpoint, result]) => (
          <Card key={endpoint}>
            <CardHeader>
              <CardTitle className="flex items-center justify-between">
                {endpoint}
                {result ? (
                  <Badge variant="default">✓ Sucesso</Badge>
                ) : errors[endpoint] ? (
                  <Badge variant="destructive">✗ Erro</Badge>
                ) : (
                  <Badge variant="secondary">⏳ Não testado</Badge>
                )}
              </CardTitle>
            </CardHeader>
            <CardContent>
              {loading[endpoint] && <p className="text-blue-500">Carregando...</p>}
              {errors[endpoint] && (
                <Alert variant="destructive">
                  <AlertDescription>{errors[endpoint]}</AlertDescription>
                </Alert>
              )}
              {result && (
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">Dados recebidos:</p>
                  <pre className="text-xs bg-gray-100 p-2 rounded overflow-auto max-h-40">
                    {JSON.stringify(result, null, 2)}
                  </pre>
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Informações de debug */}
      <Card>
        <CardHeader>
          <CardTitle>Informações de Debug</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          <div>
            <strong>URL Base da API:</strong> <code>/api/v1/stats/</code>
          </div>
          <div>
            <strong>URL Completa:</strong> <code>{window.location.origin}/api/v1/stats/totals</code>
          </div>
          <div>
            <strong>Backend esperado:</strong> <code>http://localhost:8080</code>
          </div>
          <div>
            <strong>Proxy configurado:</strong> <code>/api → http://localhost:8080</code>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

export default ApiDebug