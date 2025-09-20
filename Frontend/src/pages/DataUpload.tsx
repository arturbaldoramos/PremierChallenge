import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { FileUploader } from '@/components/FileUploader';
import {
  Database,
  Users,
  Building2,
  MapPin,
  Stethoscope,
  FileSpreadsheet,
  Activity,
  CheckCircle2,
  AlertCircle,
  Clock
} from 'lucide-react';

interface UploadStats {
  [key: string]: {
    total: number;
    completed: number;
    errors: number;
  };
}

const DATA_TYPES = [
  {
    id: 'estados',
    label: 'Estados',
    description: 'Dados dos estados brasileiros',
    icon: MapPin,
    color: 'text-blue-500',
    fields: ['codigo/cod/uf', 'nome/name', 'unidade_federativa/uf/sigla', 'regiao/region', 'latitude/lat', 'longitude/lng']
  },
  {
    id: 'municipios',
    label: 'Municípios',
    description: 'Dados dos municípios brasileiros',
    icon: Building2,
    color: 'text-green-500',
    fields: ['codigo_ibge', 'nome', 'latitude', 'longitude', 'capital', 'codigo_uf', 'siafi_id', 'ddd', 'fuso_horario', 'populacao']
  },
  {
    id: 'medicos',
    label: 'Médicos',
    description: 'Cadastro de médicos',
    icon: Stethoscope,
    color: 'text-purple-500',
    fields: ['codigo', 'nome_completo', 'especialidade', 'cidade']
  },
  {
    id: 'hospitais',
    label: 'Hospitais',
    description: 'Cadastro de hospitais',
    icon: Building2,
    color: 'text-red-500',
    fields: ['nome', 'endereco', 'cidade', 'estado', 'tipo']
  },
  {
    id: 'pacientes',
    label: 'Pacientes',
    description: 'Cadastro de pacientes',
    icon: Users,
    color: 'text-orange-500',
    fields: ['nome', 'idade', 'sexo', 'endereco', 'cidade']
  },
  {
    id: 'cid10',
    label: 'CID-10',
    description: 'Códigos de classificação de doenças',
    icon: FileSpreadsheet,
    color: 'text-indigo-500',
    fields: ['codigo', 'descricao', 'categoria']
  }
];

export default function DataUpload() {
  const [activeTab, setActiveTab] = useState('estados');
  const [uploadStats, setUploadStats] = useState<UploadStats>({});
  const [recentUploads, setRecentUploads] = useState<Array<{
    id: string;
    type: string;
    fileName: string;
    timestamp: Date;
    status: 'success' | 'error';
  }>>([]);

  const handleUploadComplete = (fileType: string, fileName: string) => {
    // Atualizar estatísticas
    setUploadStats(prev => ({
      ...prev,
      [fileType]: {
        total: (prev[fileType]?.total || 0) + 1,
        completed: (prev[fileType]?.completed || 0) + 1,
        errors: prev[fileType]?.errors || 0
      }
    }));

    // Adicionar aos uploads recentes
    setRecentUploads(prev => [{
      id: `${Date.now()}-${fileType}`,
      type: fileType,
      fileName,
      timestamp: new Date(),
      status: 'success'
    }, ...prev.slice(0, 9)]); // Manter apenas os 10 mais recentes
  };

  const handleUploadError = (error: string) => {
    console.error('Upload error:', error);

    // Adicionar erro aos uploads recentes
    setRecentUploads(prev => [{
      id: `${Date.now()}-error`,
      type: activeTab,
      fileName: 'Upload falhou',
      timestamp: new Date(),
      status: 'error'
    }, ...prev.slice(0, 9)]);

    // Atualizar estatísticas de erro
    setUploadStats(prev => ({
      ...prev,
      [activeTab]: {
        total: (prev[activeTab]?.total || 0) + 1,
        completed: prev[activeTab]?.completed || 0,
        errors: (prev[activeTab]?.errors || 0) + 1
      }
    }));
  };

  const getCurrentDataType = () => {
    return DATA_TYPES.find(dt => dt.id === activeTab);
  };

  const formatTimestamp = (date: Date) => {
    return date.toLocaleString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="flex-1 space-y-6 p-4 md:p-8 pt-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Upload de Dados</h2>
          <p className="text-muted-foreground">
            Faça upload de arquivos CSV via WebSocket com processamento em lotes
          </p>
        </div>
        <Badge variant="outline" className="flex items-center gap-2">
          <Activity className="w-4 h-4" />
          WebSocket Upload
        </Badge>
      </div>

      {/* Estatísticas Gerais */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total de Uploads</CardTitle>
            <Database className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {Object.values(uploadStats).reduce((sum, stat) => sum + stat.total, 0)}
            </div>
            <p className="text-xs text-muted-foreground">
              Arquivos enviados
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Concluídos</CardTitle>
            <CheckCircle2 className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {Object.values(uploadStats).reduce((sum, stat) => sum + stat.completed, 0)}
            </div>
            <p className="text-xs text-muted-foreground">
              Processados com sucesso
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Erros</CardTitle>
            <AlertCircle className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {Object.values(uploadStats).reduce((sum, stat) => sum + stat.errors, 0)}
            </div>
            <p className="text-xs text-muted-foreground">
              Falhas no processamento
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Tipos de Dados</CardTitle>
            <FileSpreadsheet className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {Object.keys(uploadStats).length}
            </div>
            <p className="text-xs text-muted-foreground">
              Categorias utilizadas
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Upload Interface */}
        <div className="lg:col-span-2">
          <Tabs value={activeTab} onValueChange={setActiveTab}>
            <TabsList className="grid w-full grid-cols-3 lg:grid-cols-6">
              {DATA_TYPES.map((dataType) => (
                <TabsTrigger
                  key={dataType.id}
                  value={dataType.id}
                  className="flex items-center gap-2 text-xs"
                >
                  <dataType.icon className={`w-4 h-4 ${dataType.color}`} />
                  <span className="hidden sm:inline">{dataType.label}</span>
                </TabsTrigger>
              ))}
            </TabsList>

            {DATA_TYPES.map((dataType) => (
              <TabsContent key={dataType.id} value={dataType.id} className="space-y-4">
                <Card>
                  <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                      <dataType.icon className={`w-5 h-5 ${dataType.color}`} />
                      {dataType.label}
                    </CardTitle>
                    <p className="text-sm text-muted-foreground">
                      {dataType.description}
                    </p>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      <div>
                        <h4 className="text-sm font-medium mb-2">Campos esperados no CSV:</h4>
                        <div className="flex flex-wrap gap-2">
                          {dataType.fields.map((field) => (
                            <Badge key={field} variant="outline" className="text-xs">
                              {field}
                            </Badge>
                          ))}
                        </div>
                      </div>

                      {uploadStats[dataType.id] && (
                        <div className="grid grid-cols-3 gap-4 text-center">
                          <div>
                            <div className="text-2xl font-bold">{uploadStats[dataType.id].total}</div>
                            <div className="text-xs text-muted-foreground">Total</div>
                          </div>
                          <div>
                            <div className="text-2xl font-bold text-green-600">{uploadStats[dataType.id].completed}</div>
                            <div className="text-xs text-muted-foreground">Sucesso</div>
                          </div>
                          <div>
                            <div className="text-2xl font-bold text-red-600">{uploadStats[dataType.id].errors}</div>
                            <div className="text-xs text-muted-foreground">Erros</div>
                          </div>
                        </div>
                      )}
                    </div>
                  </CardContent>
                </Card>

                <FileUploader
                  onUploadComplete={handleUploadComplete}
                  onError={handleUploadError}
                />
              </TabsContent>
            ))}
          </Tabs>
        </div>

        {/* Sidebar com uploads recentes */}
        <div className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Clock className="w-5 h-5" />
                Uploads Recentes
              </CardTitle>
            </CardHeader>
            <CardContent>
              {recentUploads.length === 0 ? (
                <p className="text-sm text-muted-foreground text-center py-4">
                  Nenhum upload realizado ainda
                </p>
              ) : (
                <div className="space-y-3">
                  {recentUploads.map((upload) => {
                    const dataType = DATA_TYPES.find(dt => dt.id === upload.type);
                    return (
                      <div key={upload.id} className="flex items-center gap-3 p-2 rounded-lg bg-muted/50">
                        {dataType && (
                          <dataType.icon className={`w-4 h-4 ${dataType.color}`} />
                        )}
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium truncate">
                            {upload.fileName}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {dataType?.label} • {formatTimestamp(upload.timestamp)}
                          </p>
                        </div>
                        {upload.status === 'success' ? (
                          <CheckCircle2 className="w-4 h-4 text-green-500" />
                        ) : (
                          <AlertCircle className="w-4 h-4 text-red-500" />
                        )}
                      </div>
                    );
                  })}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}