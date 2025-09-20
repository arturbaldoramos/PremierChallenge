import { useCallback, useEffect, useRef, useState } from 'react';

export interface WebSocketMessage {
  type: string;
  session_id?: string;
  data?: any;
  error?: string;
}

export interface UploadProgress {
  job_id: string;
  type: string;
  status: string;
  progress: number;
  total_items: number;
  processed_items: number;
  message: string;
}

export interface QueueStatus {
  [key: string]: number;
}

export interface UseWebSocketOptions {
  url: string;
  autoConnect?: boolean;
  reconnectAttempts?: number;
  reconnectInterval?: number;
}

export interface UseWebSocketReturn {
  isConnected: boolean;
  isConnecting: boolean;
  sessionId: string | null;
  connect: () => void;
  disconnect: () => void;
  sendMessage: (message: WebSocketMessage) => void;
  uploadFile: (file: File, fileType: string) => Promise<void>;
  getProgress: () => void;
  getQueueStatus: () => void;
  lastMessage: WebSocketMessage | null;
  progressData: UploadProgress[];
  queueStatus: QueueStatus | null;
  error: string | null;
}

export const useWebSocket = ({
  url,
  autoConnect = false,
  reconnectAttempts = 3,
  reconnectInterval = 5000
}: UseWebSocketOptions): UseWebSocketReturn => {
  console.log('useWebSocket inicializado com:', { url, autoConnect, reconnectAttempts, reconnectInterval });

  const [isConnected, setIsConnected] = useState(false);
  const [isConnecting, setIsConnecting] = useState(false);
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null);
  const [progressData, setProgressData] = useState<UploadProgress[]>([]);
  const [queueStatus, setQueueStatus] = useState<QueueStatus | null>(null);
  const [error, setError] = useState<string | null>(null);

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  const connect = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log('WebSocket já está conectado');
      return;
    }

    console.log('Tentando conectar ao WebSocket:', url);
    setIsConnecting(true);
    setError(null);

    try {
      wsRef.current = new WebSocket(url);

      wsRef.current.onopen = () => {
        console.log('WebSocket conectado com sucesso');
        setIsConnected(true);
        setIsConnecting(false);
        reconnectAttemptsRef.current = 0;

        // Iniciar sessão de upload
        const startMessage: WebSocketMessage = {
          type: 'upload_start'
        };
        console.log('Enviando mensagem de início:', startMessage);
        wsRef.current?.send(JSON.stringify(startMessage));
      };

      wsRef.current.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          console.log('Mensagem recebida:', message);
          setLastMessage(message);
          handleMessage(message);
        } catch (err) {
          console.error('Erro ao parsear mensagem WebSocket:', err);
        }
      };

      wsRef.current.onclose = () => {
        console.log('WebSocket desconectado');
        setIsConnected(false);
        setIsConnecting(false);
        setSessionId(null);

        if (reconnectAttemptsRef.current < reconnectAttempts) {
          reconnectAttemptsRef.current++;
          console.log(`Tentativa de reconexão ${reconnectAttemptsRef.current}/${reconnectAttempts}`);
          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, reconnectInterval);
        }
      };

      wsRef.current.onerror = (error) => {
        console.error('Erro no WebSocket:', error);
        setError('Erro na conexão WebSocket');
        setIsConnecting(false);
      };

    } catch (err) {
      setError('Falha ao conectar ao WebSocket');
      setIsConnecting(false);
    }
  }, [url, reconnectAttempts, reconnectInterval]);

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }

    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }

    setIsConnected(false);
    setIsConnecting(false);
    setSessionId(null);
    reconnectAttemptsRef.current = 0;
  }, []);

  const sendMessage = useCallback((message: WebSocketMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    } else {
      setError('WebSocket não está conectado');
    }
  }, []);

  const handleMessage = useCallback((message: WebSocketMessage) => {
    switch (message.type) {
      case 'upload_ready':
        setSessionId(message.session_id || null);
        break;

      case 'progress_update':
        if (message.data) {
          try {
            const progressArray: UploadProgress[] = Array.isArray(message.data)
              ? message.data
              : JSON.parse(message.data);
            setProgressData(progressArray);
          } catch (err) {
            console.error('Erro ao parsear dados de progresso:', err);
          }
        }
        break;

      case 'queue_status':
        if (message.data) {
          try {
            const queueData: QueueStatus = typeof message.data === 'object'
              ? message.data
              : JSON.parse(message.data);
            setQueueStatus(queueData);
          } catch (err) {
            console.error('Erro ao parsear status da fila:', err);
          }
        }
        break;

      case 'error':
        setError(message.error || 'Erro desconhecido');
        break;

      case 'upload_processing':
        // Arquivo foi aceito e está sendo processado
        break;

      case 'chunk_received':
        // Chunk foi recebido com sucesso
        break;
    }
  }, []);

  const fileToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => {
        if (typeof reader.result === 'string') {
          // Remover o prefixo "data:...;base64,"
          const base64 = reader.result.split(',')[1];
          resolve(base64);
        } else {
          reject(new Error('Falha ao converter arquivo para base64'));
        }
      };
      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  };

  const createChunks = (base64Data: string, chunkSize: number = 64 * 1024) => {
    const chunks = [];
    const totalChunks = Math.ceil(base64Data.length / chunkSize);

    for (let i = 0; i < totalChunks; i++) {
      const start = i * chunkSize;
      const end = Math.min(start + chunkSize, base64Data.length);
      const chunkData = base64Data.slice(start, end);

      chunks.push({
        index: i,
        data: chunkData,
        total: totalChunks
      });
    }

    return chunks;
  };

  const uploadFile = useCallback(async (file: File, fileType: string) => {
    console.log('uploadFile chamado com:', { fileName: file.name, fileType, sessionId });

    if (!sessionId) {
      throw new Error('Sessão WebSocket não iniciada');
    }

    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket não está conectado');
    }

    try {
      // Converter arquivo para base64
      const base64Data = await fileToBase64(file);
      console.log('Arquivo convertido para base64, tamanho:', base64Data.length);

      // Dividir em chunks
      const chunks = createChunks(base64Data);
      console.log('Arquivo dividido em', chunks.length, 'chunks');

      // Enviar dados do upload
      const uploadMessage: WebSocketMessage = {
        type: 'upload_complete',
        session_id: sessionId,
        data: {
          file_type: fileType,
          file_name: file.name,
          chunks: chunks
        }
      };

      console.log('Enviando mensagem de upload:', uploadMessage);
      sendMessage(uploadMessage);
    } catch (err) {
      throw new Error(`Falha no upload: ${err instanceof Error ? err.message : 'Erro desconhecido'}`);
    }
  }, [sessionId, sendMessage]);

  const getProgress = useCallback(() => {
    if (!sessionId) {
      setError('Sessão não iniciada');
      return;
    }

    const message: WebSocketMessage = {
      type: 'get_progress',
      session_id: sessionId
    };
    sendMessage(message);
  }, [sessionId, sendMessage]);

  const getQueueStatus = useCallback(() => {
    const message: WebSocketMessage = {
      type: 'get_queue_status'
    };
    sendMessage(message);
  }, [sendMessage]);

  useEffect(() => {
    console.log('useWebSocket useEffect executado - autoConnect:', autoConnect, 'URL:', url);
    if (autoConnect) {
      connect();
    }

    return () => {
      disconnect();
    };
  }, [autoConnect, connect, disconnect]);

  return {
    isConnected,
    isConnecting,
    sessionId,
    connect,
    disconnect,
    sendMessage,
    uploadFile,
    getProgress,
    getQueueStatus,
    lastMessage,
    progressData,
    queueStatus,
    error
  };
};