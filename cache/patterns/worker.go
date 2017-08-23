package patterns

import (
	"time"

	"gopkg.in/tomb.v2"

	"github.com/moira-alert/moira-alert"
	"github.com/moira-alert/moira-alert/cache"
	"github.com/moira-alert/moira-alert/metrics/graphite"
)

//RefreshPatternWorker realization
type RefreshPatternWorker struct {
	database       moira.Database
	logger         moira.Logger
	metrics        *graphite.CacheMetrics
	patternStorage *cache.PatternStorage
	tomb           tomb.Tomb
}

//NewRefreshPatternWorker creates new RefreshPatternWorker
func NewRefreshPatternWorker(database moira.Database, metrics *graphite.CacheMetrics, logger moira.Logger, patternStorage *cache.PatternStorage) *RefreshPatternWorker {
	return &RefreshPatternWorker{
		database:       database,
		metrics:        metrics,
		logger:         logger,
		patternStorage: patternStorage,
	}
}

//Run process to refresh pattern tree every second
func (worker *RefreshPatternWorker) Start() {
	worker.tomb.Go(func() error {
		for {
			checkTicker := time.NewTicker(time.Second)
			select {
			case <-worker.tomb.Dying():
				worker.logger.Infof("Moira Cache pattern updater stopped")
				return nil
			case <-checkTicker.C:
				timer := time.Now()
				err := worker.patternStorage.RefreshTree()
				if err != nil {
					worker.logger.Errorf("pattern refresh failed: %s", err.Error())
				}
				worker.metrics.BuildTreeTimer.UpdateSince(timer)
			}
		}
	})
	worker.logger.Infof("Moira Cache pattern updater started")
}

//Stop stops update pattern tree
func (worker *RefreshPatternWorker) Stop() error {
	worker.tomb.Kill(nil)
	return worker.tomb.Wait()
}
