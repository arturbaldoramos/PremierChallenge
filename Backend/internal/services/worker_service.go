package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"example.com/m/v2/internal/domain"
	"example.com/m/v2/internal/parsers"
)

type WorkerService struct {
	redisService *RedisService
	dataService  *DataService
	csvParser    *parsers.CSVParser
	workers      map[string]*Worker
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
}

type Worker struct {
	ID       string
	JobType  string
	Status   string
	LastJob  *Job
	StartedAt time.Time
}

type WorkerConfig struct {
	MaxWorkers map[string]int // Por tipo de job
	BatchSize  map[string]int // Tamanho do batch por tipo
}

func NewWorkerService(redisService *RedisService, dataService *DataService) *WorkerService {
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerService{
		redisService: redisService,
		dataService:  dataService,
		csvParser:    parsers.NewCSVParser(),
		workers:      make(map[string]*Worker),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// StartWorkers inicia workers para diferentes tipos de jobs
func (w *WorkerService) StartWorkers(config WorkerConfig) {
	jobTypes := []string{"estados", "municipios", "medicos", "hospitais", "pacientes", "cid10"}

	for _, jobType := range jobTypes {
		maxWorkers := config.MaxWorkers[jobType]
		if maxWorkers == 0 {
			maxWorkers = 2 // Default
		}

		for i := 0; i < maxWorkers; i++ {
			go w.startWorker(jobType, i)
		}
	}

	log.Printf("Started workers for job types: %v", jobTypes)
}

func (w *WorkerService) startWorker(jobType string, workerIndex int) {
	workerID := fmt.Sprintf("%s-worker-%d", jobType, workerIndex)

	worker := &Worker{
		ID:        workerID,
		JobType:   jobType,
		Status:    "idle",
		StartedAt: time.Now(),
	}

	w.mu.Lock()
	w.workers[workerID] = worker
	w.mu.Unlock()

	log.Printf("Worker %s started for job type: %s", workerID, jobType)

	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Worker %s stopping", workerID)
			return
		default:
			job, err := w.redisService.DequeueJob(jobType)
			if err != nil {
				log.Printf("Worker %s error dequeuing job: %v", workerID, err)
				time.Sleep(5 * time.Second)
				continue
			}

			if job == nil {
				// Nenhum job disponível, aguardar
				time.Sleep(2 * time.Second)
				continue
			}

			// Processar job
			w.processJob(worker, job)
		}
	}
}

func (w *WorkerService) processJob(worker *Worker, job *Job) {
	worker.Status = "processing"
	worker.LastJob = job

	log.Printf("Worker %s processing job %s (type: %s)", worker.ID, job.ID, job.Type)

	// Atualizar status para processando
	job.Status = JobStatusProcessing
	job.ProcessedItems = 0
	w.redisService.UpdateJobStatus(job)

	var err error
	switch job.Type {
	case "estados":
		err = w.processEstados(job)
	case "municipios":
		err = w.processMunicipios(job)
	case "medicos":
		err = w.processMedicos(job)
	case "hospitais":
		err = w.processHospitais(job)
	case "pacientes":
		err = w.processPacientes(job)
	case "cid10":
		err = w.processCID10(job)
	default:
		err = fmt.Errorf("unknown job type: %s", job.Type)
	}

	// Finalizar job
	now := time.Now()
	if err != nil {
		job.Status = JobStatusFailed
		job.Error = err.Error()
		log.Printf("Worker %s failed to process job %s: %v", worker.ID, job.ID, err)
	} else {
		job.Status = JobStatusCompleted
		job.Progress = 100
		job.ProcessedItems = job.TotalItems
		log.Printf("Worker %s completed job %s", worker.ID, job.ID)
	}

	job.CompletedAt = &now
	w.redisService.UpdateJobStatus(job)

	worker.Status = "idle"
	worker.LastJob = nil
}

func (w *WorkerService) processEstados(job *Job) error {
	var batchData BatchData
	if err := json.Unmarshal(job.Data, &batchData); err != nil {
		return err
	}

	var estados []domain.Estado
	if err := json.Unmarshal(batchData.Data, &estados); err != nil {
		return err
	}

	// Processar batch
	inserted, updated, err := w.dataService.UpsertEstados(estados)
	if err != nil {
		return err
	}

	log.Printf("Estados batch %d processed: %d inserted, %d updated",
		batchData.BatchNumber, inserted, updated)

	return nil
}

func (w *WorkerService) processMunicipios(job *Job) error {
	var batchData BatchData
	if err := json.Unmarshal(job.Data, &batchData); err != nil {
		return err
	}

	var municipios []domain.Municipio
	if err := json.Unmarshal(batchData.Data, &municipios); err != nil {
		return err
	}

	// Processar batch
	inserted, updated, err := w.dataService.UpsertMunicipios(municipios)
	if err != nil {
		return err
	}

	log.Printf("Municipios batch %d processed: %d inserted, %d updated",
		batchData.BatchNumber, inserted, updated)

	return nil
}

func (w *WorkerService) processMedicos(job *Job) error {
	var batchData BatchData
	if err := json.Unmarshal(job.Data, &batchData); err != nil {
		return err
	}

	var medicos []domain.Medico
	if err := json.Unmarshal(batchData.Data, &medicos); err != nil {
		return err
	}

	// Processar batch
	inserted, updated, err := w.dataService.UpsertMedicos(medicos)
	if err != nil {
		return err
	}

	log.Printf("Medicos batch %d processed: %d inserted, %d updated",
		batchData.BatchNumber, inserted, updated)

	return nil
}

func (w *WorkerService) processHospitais(job *Job) error {
	var batchData BatchData
	if err := json.Unmarshal(job.Data, &batchData); err != nil {
		return err
	}

	var hospitais []domain.Hospital
	if err := json.Unmarshal(batchData.Data, &hospitais); err != nil {
		return err
	}

	// Processar batch
	inserted, updated, err := w.dataService.UpsertHospitais(hospitais)
	if err != nil {
		return err
	}

	log.Printf("Hospitais batch %d processed: %d inserted, %d updated",
		batchData.BatchNumber, inserted, updated)

	return nil
}

func (w *WorkerService) processPacientes(job *Job) error {
	var batchData BatchData
	if err := json.Unmarshal(job.Data, &batchData); err != nil {
		return err
	}

	var pacientes []domain.Paciente
	if err := json.Unmarshal(batchData.Data, &pacientes); err != nil {
		return err
	}

	// Processar batch
	inserted, updated, err := w.dataService.UpsertPacientes(pacientes)
	if err != nil {
		return err
	}

	log.Printf("Pacientes batch %d processed: %d inserted, %d updated",
		batchData.BatchNumber, inserted, updated)

	return nil
}

func (w *WorkerService) processCID10(job *Job) error {
	var batchData BatchData
	if err := json.Unmarshal(job.Data, &batchData); err != nil {
		return err
	}

	var cid10s []domain.Cid10
	if err := json.Unmarshal(batchData.Data, &cid10s); err != nil {
		return err
	}

	// Processar batch
	inserted, updated, err := w.dataService.UpsertCID10(cid10s)
	if err != nil {
		return err
	}

	log.Printf("CID10 batch %d processed: %d inserted, %d updated",
		batchData.BatchNumber, inserted, updated)

	return nil
}

// Stop para todos os workers
func (w *WorkerService) Stop() {
	log.Println("Stopping all workers...")
	w.cancel()
}

// GetWorkerStats retorna estatísticas dos workers
func (w *WorkerService) GetWorkerStats() map[string]*Worker {
	w.mu.RLock()
	defer w.mu.RUnlock()

	stats := make(map[string]*Worker)
	for id, worker := range w.workers {
		stats[id] = &Worker{
			ID:        worker.ID,
			JobType:   worker.JobType,
			Status:    worker.Status,
			StartedAt: worker.StartedAt,
		}
		if worker.LastJob != nil {
			stats[id].LastJob = &Job{
				ID:     worker.LastJob.ID,
				Type:   worker.LastJob.Type,
				Status: worker.LastJob.Status,
			}
		}
	}

	return stats
}