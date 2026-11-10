package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lorsanstand/gontrol/internal/hub/store"
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


func (c *CommandSpecsRepository) SetCommandSpecs(ctx context.Context, agentID uuid.UUID, command models.CommandSpecs) error {
	data, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("commands marshal in json: %w", err)
	}

	param := gen.SetCommandsByAgentIDParams{AgentID: agentID, Commands: data}

	result, err := c.q.SetCommandsByAgentID(ctx, param)
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			if pgErr.Code == pgerrcode.ForeignKeyViolation {
				return store.ErrNotFound
			}
		}

		return fmt.Errorf("insert command into db: %w", err)
	}

	return ifChangesInDB(result)
}