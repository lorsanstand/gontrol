package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lorsanstand/gontrol/internal/hub"
	"github.com/lorsanstand/gontrol/internal/hub/config"
	http2 "github.com/lorsanstand/gontrol/internal/hub/delivery/http"
	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/handler"
	"github.com/lorsanstand/gontrol/internal/hub/delivery/http/middlewares"
	"github.com/lorsanstand/gontrol/internal/hub/service"
	in_memory "github.com/lorsanstand/gontrol/internal/hub/store/in-memory"
	"github.com/lorsanstand/gontrol/internal/hub/store/postgres"
	"github.com/lorsanstand/gontrol/internal/utils/logger"
)

type DI struct {
	cfg    *config.Config
	logger *slog.Logger
	db     *sql.DB

	stateAgentRepo *in_memory.AgentStateStore
	agentRepo      *postgres.AgentRepository
	commandRepo    *postgres.CommandSpecsRepository
	taskRepo       *postgres.TaskRepository

	agentService     *service.AgentService
	taskService      *service.TaskService
	heartbeatService *service.HeartbeatService

	agentHandlers     *handler.AgentHandler
	taskHandlers      *handler.TaskHandler
	heartbeatHandlers *handler.HeartbeatHandler

	apiRouter *chi.Mux
	httpServer *http.Server
}

func (d *DI) GetCfg() *config.Config {
	if d.cfg == nil {
		cfg, err := config.Load()
		if err != nil {
			slog.Error("failed load config", slog.Any("error", err))
			os.Exit(1)
		}

		d.cfg = cfg
	}

	return d.cfg
}

func (d *DI) GetLogger() *slog.Logger {
	if d.logger == nil {
		cfg := d.GetCfg()

		d.logger = logger.Init(cfg.LogLevel, cfg.LogFormat)
	}

	return d.logger
}

func (d *DI) GetDB() *sql.DB {
	if d.db == nil {
		log := d.GetLogger()
		cfg := d.GetCfg()

		db, err := sql.Open("pgx", cfg.DatabaseURL())
		if err != nil {
			log.Error("unable connect to database", slog.Any("error", err))
			os.Exit(1)
		}

		d.db = db
	}

	return d.db
}

func (d *DI) GetAgentStateRepo() *in_memory.AgentStateStore {
	if d.stateAgentRepo == nil {
		d.stateAgentRepo = in_memory.NewAgentStateStore(hub.TimeTTLOnlineAgent)
	}

	return d.stateAgentRepo
}

func (d *DI) GetAgentRepo() *postgres.AgentRepository {
	if d.agentRepo == nil {
		d.agentRepo = postgres.NewAgentRepository(d.GetDB())
	}

	return d.agentRepo
}

func (d *DI) GetCommandRepo() *postgres.CommandSpecsRepository {
	if d.commandRepo == nil {
		d.commandRepo = postgres.NewCommandSpecsRepository(d.GetDB())
	}

	return d.commandRepo
}

func (d *DI) GetTaskRepo() *postgres.TaskRepository {
	if d.taskRepo == nil {
		d.taskRepo = postgres.NewTaskRepository(d.GetDB())
	}

	return d.taskRepo
}

func (d *DI) GetAgentService() *service.AgentService {
	if d.agentService == nil {
		d.agentService = service.NewAgentService(d.GetLogger(), d.GetAgentRepo(), d.GetAgentStateRepo(), d.GetCommandRepo())
	}

	return d.agentService
}

func (d *DI) GetTaskService() *service.TaskService {
	if d.taskService == nil {
		d.taskService = service.NewTasksService(d.GetLogger(), d.GetTaskRepo(), d.GetCommandRepo())
	}

	return d.taskService
}

func (d *DI) GetHeartbeatService() *service.HeartbeatService {
	if d.heartbeatService == nil {
		d.heartbeatService = service.NewHeartbeatService(d.GetAgentService(), d.GetTaskService())
	}

	return d.heartbeatService
}

func (d *DI) GetAgentHandlers() *handler.AgentHandler {
	if d.agentHandlers == nil {
		d.agentHandlers = handler.NewAgentHandler(d.GetLogger(), d.GetAgentService())
	}

	return d.agentHandlers
}

func (d *DI) GetTaskHandlers() *handler.TaskHandler {
	if d.taskHandlers == nil {
		d.taskHandlers = handler.NewTaskHandler(d.GetLogger(), d.GetTaskService())
	}

	return d.taskHandlers
}

func (d *DI) GetHeartbeatHandlers() *handler.HeartbeatHandler {
	if d.heartbeatHandlers == nil {
		d.heartbeatHandlers = handler.NewHeartbeatHandler(d.GetLogger(), d.GetHeartbeatService())
	}

	return d.heartbeatHandlers
}

func (d *DI) GetApiRouter() *chi.Mux {
	if d.apiRouter == nil {
		r := chi.NewRouter()
		middlewares.ConnectMiddlewares(r, d.GetLogger())
		apiRouter := http2.RegisterRoutes(d.GetTaskHandlers(), d.GetAgentHandlers(), d.GetHeartbeatHandlers())
		r.Mount("/api/v1", apiRouter)
		d.apiRouter = r
	}

	return d.apiRouter
}

func (d *DI) GetHttpServer() *http.Server {
	if d.httpServer == nil {
		cfg := d.GetCfg()
		r := d.GetApiRouter()

		d.httpServer = &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      r,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
	}

	return d.httpServer
}