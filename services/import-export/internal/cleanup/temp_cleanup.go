package cleanup

import (
	"context"
	"log"
	"time"

	"crm-system/services/import-export/internal/repo"
)

type TempCleanup struct {
	repo     *repo.RunRepo
	interval time.Duration
}

func NewTempCleanup(repo *repo.RunRepo, interval time.Duration) TempCleanup {
	if interval <= 0 {
		interval = time.Hour
	}
	return TempCleanup{repo: repo, interval: interval}
}

func (c TempCleanup) RunOnce(ctx context.Context, now time.Time) error {
	count, err := c.repo.MarkExpiredRunsDeleted(ctx, now)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("import-export cleanup marked %d expired runs deleted", count)
	}
	return nil
}

func (c TempCleanup) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				if err := c.RunOnce(ctx, now.UTC()); err != nil {
					log.Printf("import-export cleanup failed: %v", err)
				}
			}
		}
	}()
}
