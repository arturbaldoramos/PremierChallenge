package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`        // "estados", "municipios", "medicos", etc.
	Data        json.RawMessage `json:"data"`        // Dados do batch
	Status      JobStatus       `json:"status"`
	Progress    int             `json:"progress"`    // 0-100
	TotalItems  int             `json:"total_items"`
	ProcessedItems int          `json:"processed_items"`
	CreatedAt   time.Time       `json:"created_at"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Error       string          `json:"error,omitempty"`
	SessionID   string          `json:"session_id"` // Para associar ao WebSocket
}

type BatchData struct {
	BatchNumber int             `json:"batch_number"`
	TotalBatches int            `json:"total_batches"`
	Data        json.RawMessage `json:"data"`
}

func NewRedisService(addr, password string, db int) *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// Testar conexão
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
	}

	return &RedisService{
		client: rdb,
		ctx:    ctx,
	}
}

// EnqueueJob adiciona um job à fila
func (r *RedisService) EnqueueJob(job *Job) error {
	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// Adicionar à fila de processamento
	queueKey := fmt.Sprintf("queue:%s", job.Type)
	err = r.client.LPush(r.ctx, queueKey, jobData).Err()
	if err != nil {
		return err
	}

	// Salvar status do job
	statusKey := fmt.Sprintf("job:%s", job.ID)
	err = r.client.Set(r.ctx, statusKey, jobData, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	// Adicionar à lista de jobs da sessão
	sessionKey := fmt.Sprintf("session:%s:jobs", job.SessionID)
	err = r.client.SAdd(r.ctx, sessionKey, job.ID).Err()
	if err != nil {
		return err
	}
	r.client.Expire(r.ctx, sessionKey, 24*time.Hour)

	return nil
}

// DequeueJob remove um job da fila para processamento
func (r *RedisService) DequeueJob(jobType string) (*Job, error) {
	queueKey := fmt.Sprintf("queue:%s", jobType)

	result, err := r.client.BRPop(r.ctx, 5*time.Second, queueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Nenhum job disponível
		}
		return nil, err
	}

	var job Job
	err = json.Unmarshal([]byte(result[1]), &job)
	if err != nil {
		return nil, err
	}

	// Marcar como processando
	now := time.Now()
	job.Status = JobStatusProcessing
	job.StartedAt = &now

	err = r.UpdateJobStatus(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// UpdateJobStatus atualiza o status de um job
func (r *RedisService) UpdateJobStatus(job *Job) error {
	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}

	statusKey := fmt.Sprintf("job:%s", job.ID)
	return r.client.Set(r.ctx, statusKey, jobData, 24*time.Hour).Err()
}

// GetJobStatus recupera o status de um job
func (r *RedisService) GetJobStatus(jobID string) (*Job, error) {
	statusKey := fmt.Sprintf("job:%s", jobID)

	result, err := r.client.Get(r.ctx, statusKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("job not found")
		}
		return nil, err
	}

	var job Job
	err = json.Unmarshal([]byte(result), &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// GetSessionJobs recupera todos os jobs de uma sessão
func (r *RedisService) GetSessionJobs(sessionID string) ([]*Job, error) {
	sessionKey := fmt.Sprintf("session:%s:jobs", sessionID)

	jobIDs, err := r.client.SMembers(r.ctx, sessionKey).Result()
	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, 0, len(jobIDs))
	for _, jobID := range jobIDs {
		job, err := r.GetJobStatus(jobID)
		if err != nil {
			continue // Pular jobs com erro
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// CleanupSession remove dados de uma sessão
func (r *RedisService) CleanupSession(sessionID string) error {
	sessionKey := fmt.Sprintf("session:%s:jobs", sessionID)

	jobIDs, err := r.client.SMembers(r.ctx, sessionKey).Result()
	if err != nil {
		return err
	}

	// Remover jobs individuais
	for _, jobID := range jobIDs {
		statusKey := fmt.Sprintf("job:%s", jobID)
		r.client.Del(r.ctx, statusKey)
	}

	// Remover sessão
	return r.client.Del(r.ctx, sessionKey).Err()
}

// GetQueueLength retorna o número de jobs pendentes na fila
func (r *RedisService) GetQueueLength(jobType string) (int64, error) {
	queueKey := fmt.Sprintf("queue:%s", jobType)
	return r.client.LLen(r.ctx, queueKey).Result()
}