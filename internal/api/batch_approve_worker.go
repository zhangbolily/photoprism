package api

import "github.com/photoprism/photoprism/internal/entity"

type ApproveJob struct {
	Photo entity.Photo
}

type ApproveResult struct {
	Photos entity.Photos
}

func ApproveWorker(jobs <-chan ApproveJob) (result ApproveResult) {
	var err error

	for job := range jobs {
		p := job.Photo

		if err = p.Approve(); err != nil {
			log.Errorf("approve: %s", err)
		} else {
			result.Photos = append(result.Photos, p)
			SavePhotoAsYaml(p)
		}
	}

	return
}
