package sync

import (
	"fmt"
	"github.com/lovelaze/nebula-sync/internal/config"
	"github.com/rs/zerolog/log"
)

func (target *target) ManualSync(syncConf *config.Sync) error {
	log.Info().Str("mode", "selective").Int("replicas", len(target.Replicas)).Msg("Running sync")

	if err := target.authenticate(); err != nil {
		return fmt.Errorf("authentication: %w", err)
	}

	if err := target.syncTeleporters(syncConf.GravitySettings); err != nil {
		return fmt.Errorf("sync teleporters: %w", err)
	}

	if err := target.syncConfigs(syncConf.ConfigSettings); err != nil {
		return fmt.Errorf("sync configs: %w", err)
	}

	if syncConf.RunGravity {
		if err := target.runGravity(); err != nil {
			return fmt.Errorf("run gravity: %w", err)
		}
	}

	if err := target.deleteSessions(); err != nil {
		return fmt.Errorf("delete sessions: %w", err)
	}

	return nil
}
