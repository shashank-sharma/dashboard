package cronjobs

import (
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/shashank-sharma/backend/internal/logger"
)

type CronJob struct {
	Name     string
	Interval string
	JobFunc  func()
	IsActive bool
	LastRun  time.Time
}

// GetStatusString returns a user-friendly status of the cron job
func (c *CronJob) GetStatusString() string {
	if c.LastRun.IsZero() {
		return "Pending"
	}
	return "Active"
}

// Store active jobs for status reporting
var (
	activeJobs []CronJob
	jobsMutex  sync.RWMutex
)

// GetActiveJobs returns a copy of the active cron jobs
func GetActiveJobs() []CronJob {
	jobsMutex.RLock()
	defer jobsMutex.RUnlock()
	
	// Create a copy to avoid race conditions
	result := make([]CronJob, len(activeJobs))
	copy(result, activeJobs)
	return result
}

func Run(cronJobs []CronJob) error {
	scheduler := cron.New()
	
	// Store active jobs for status reporting
	jobsMutex.Lock()
	activeJobs = make([]CronJob, 0, len(cronJobs))
	jobsMutex.Unlock()

	for _, job := range cronJobs {
		if job.IsActive {
			logger.LogInfo("Running CRON", "job", job.Name)
			
			// Create a wrapper function to track execution time
			wrapperFunc := func(j CronJob) func() {
				return func() {
					// Update last run time
					jobsMutex.Lock()
					for i := range activeJobs {
						if activeJobs[i].Name == j.Name {
							activeJobs[i].LastRun = time.Now()
							break
						}
					}
					jobsMutex.Unlock()
					
					// Execute the actual job
					j.JobFunc()
				}
			}
			
			err := scheduler.Add(job.Name, job.Interval, wrapperFunc(job))
			if err != nil {
				logger.LogError("Failed to run CRON: ", job.Name)
			} else {
				// Add to active jobs list
				jobsMutex.Lock()
				activeJobs = append(activeJobs, job)
				jobsMutex.Unlock()
			}
		} else {
			logger.LogInfo("Skipping CRON:", job.Name)
		}
	}

	scheduler.Start()
	return nil
}
