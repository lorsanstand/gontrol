package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	gen "github.com/lorsanstand/gontrol/internal/hub/store/postgres/gen"
	"github.com/lorsanstand/gontrol/internal/models"
)

type CommandSpecsRepository struct {
	db *sql.DB
	q  gen.Querier
}

func (c *CommandSpecsRepository) GetCommandSpecs(ctx context.Context, agentID uuid.UUID) (models.CommandSpecs, error) {
	command, err := c.q.GetCommandsByAgentID(ctx, agentID)
	if err != nil {
		return models.CommandSpecs{}, fmt.Errorf("select command into db: %w", err)
	}

	return toCommands(command)
}

func (c *CommandSpecsRepository) SetCommandSpecs(ctx context.Context, agentID uuid.UUID) error {

}