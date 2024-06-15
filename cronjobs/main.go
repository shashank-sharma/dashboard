package cronjobs

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/cron"
	"github.com/shashank-sharma/backend/logger"
)

type CronJob struct {
	Name     string
	Interval string
	JobFunc  func()
	IsActive bool
}

func Run(cronJobs []CronJob) error {
	scheduler := cron.New()

	for _, job := range cronJobs {
		if job.IsActive {
			logger.LogInfo("Running CRON:", job.Name)
			err := scheduler.Add(job.Name, job.Interval, job.JobFunc)
			if err != nil {
				logger.LogError("Failed to run CRON: ", job.Name)
			}
		} else {
			logger.LogInfo("Skipping CRON:", job.Name)
		}
	}
	return nil
}

func InitCronjobs(app *pocketbase.PocketBase) error {
	cronJobs := []CronJob{
		{
			Name:     "track-device",
			Interval: "*/6 * * * *",
			JobFunc: func() {
				TrackDevices(app)
			},
			IsActive: true,
		},
	}

	Run(cronJobs)
	return nil
}
