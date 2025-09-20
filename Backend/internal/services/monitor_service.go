package services

import (
	"context"
	"log"
	"sync"
	"time"
)

type MonitorService struct {
	redisService *RedisService
	subscribers  map[string]chan ProgressUpdate
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
}

type ProgressUpdate struct {
	SessionID      string `json:"session_id"`
	JobID          string `json:"job_id"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	Progress       int    `json:"progress"`
	TotalItems     int    `json:"total_items"`
	ProcessedItems int    `json:"processed_items"`
	Message        string `json:"message,omitempty"`
	Error          string `json:"error,omitempty"`
}

type Subscriber interface {
	OnProgressUpdate(update ProgressUpdate)
}

func NewMonitorService(redisService *RedisService) *MonitorService {
	ctx, cancel := context.WithCancel(context.Background())

	return &MonitorService{
		redisService: redisService,
		subscribers:  make(map[string]chan ProgressUpdate),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start inicia o monitoramento
func (m *MonitorService) Start() {
	go m.monitorJobs()
	log.Println("Monitor service started")
}

// Stop para o monitoramento
func (m *MonitorService) Stop() {
	m.cancel()

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ch := range m.subscribers {
		close(ch)
	}

	log.Println("Monitor service stopped")
}

// Subscribe adiciona um subscriber para updates de progresso
func (m *MonitorService) Subscribe(sessionID string) <-chan ProgressUpdate {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch := make(chan ProgressUpdate, 100)
	m.subscribers[sessionID] = ch

	return ch
}

// Unsubscribe remove um subscriber
func (m *MonitorService) Unsubscribe(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if ch, exists := m.subscribers[sessionID]; exists {
		close(ch)
		delete(m.subscribers, sessionID)
	}
}

// monitorJobs monitora continuamente o status dos jobs
func (m *MonitorService) monitorJobs() {
	ticker := time.NewTicker(2 * time.Second) // Check every 2 seconds
	defer ticker.Stop()

	lastStates := make(map[string]*Job)

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.checkJobUpdates(lastStates)
		}
	}
}

func (m *MonitorService) checkJobUpdates(lastStates map[string]*Job) {
	m.mu.RLock()
	sessionIDs := make([]string, 0, len(m.subscribers))
	for sessionID := range m.subscribers {
		sessionIDs = append(sessionIDs, sessionID)
	}
	m.mu.RUnlock()

	for _, sessionID := range sessionIDs {
		jobs, err := m.redisService.GetSessionJobs(sessionID)
		if err != nil {
			log.Printf("Error getting session jobs for %s: %v", sessionID, err)
			continue
		}

		for _, job := range jobs {
			lastJob, exists := lastStates[job.ID]

			// Verificar se houve mudança significativa
			if m.shouldNotify(job, lastJob, exists) {
				update := ProgressUpdate{
					SessionID:      sessionID,
					JobID:          job.ID,
					Type:           job.Type,
					Status:         string(job.Status),
					Progress:       job.Progress,
					TotalItems:     job.TotalItems,
					ProcessedItems: job.ProcessedItems,
					Error:          job.Error,
				}

				m.notifySubscribers(sessionID, update)

				// Atualizar estado
				lastStates[job.ID] = &Job{
					ID:             job.ID,
					Status:         job.Status,
					Progress:       job.Progress,
					ProcessedItems: job.ProcessedItems,
					Error:          job.Error,
				}
			}
		}
	}
}

func (m *MonitorService) shouldNotify(current, last *Job, exists bool) bool {
	if !exists {
		return true // Novo job
	}

	// Notificar se status, progresso ou itens processados mudaram
	return current.Status != last.Status ||
		current.Progress != last.Progress ||
		current.ProcessedItems != last.ProcessedItems ||
		(current.Error != "" && current.Error != last.Error)
}

func (m *MonitorService) notifySubscribers(sessionID string, update ProgressUpdate) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if ch, exists := m.subscribers[sessionID]; exists {
		select {
		case ch <- update:
		default:
			// Channel full, skip this update
			log.Printf("Progress update channel full for session %s", sessionID)
		}
	}
}

// NotifyJobProgress força uma notificação de progresso para um job específico
func (m *MonitorService) NotifyJobProgress(job *Job) {
	update := ProgressUpdate{
		SessionID:      job.SessionID,
		JobID:          job.ID,
		Type:           job.Type,
		Status:         string(job.Status),
		Progress:       job.Progress,
		TotalItems:     job.TotalItems,
		ProcessedItems: job.ProcessedItems,
		Error:          job.Error,
	}

	m.notifySubscribers(job.SessionID, update)
}

// GetSessionProgress retorna o progresso atual de uma sessão
func (m *MonitorService) GetSessionProgress(sessionID string) ([]ProgressUpdate, error) {
	jobs, err := m.redisService.GetSessionJobs(sessionID)
	if err != nil {
		return nil, err
	}

	updates := make([]ProgressUpdate, len(jobs))
	for i, job := range jobs {
		updates[i] = ProgressUpdate{
			SessionID:      sessionID,
			JobID:          job.ID,
			Type:           job.Type,
			Status:         string(job.Status),
			Progress:       job.Progress,
			TotalItems:     job.TotalItems,
			ProcessedItems: job.ProcessedItems,
			Error:          job.Error,
		}
	}

	return updates, nil
}

// GetQueueStats retorna estatísticas das filas
func (m *MonitorService) GetQueueStats() (map[string]int64, error) {
	stats := make(map[string]int64)
	jobTypes := []string{"estados", "municipios", "medicos", "hospitais", "pacientes", "cid10"}

	for _, jobType := range jobTypes {
		length, err := m.redisService.GetQueueLength(jobType)
		if err != nil {
			return nil, err
		}
		stats[jobType] = length
	}

	return stats, nil
}