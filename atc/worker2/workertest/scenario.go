package workertest

import (
	"code.cloudfoundry.org/lager/lagertest"
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/db"
	"github.com/concourse/concourse/atc/db/dbtest"
	"github.com/concourse/concourse/atc/runtime"
	"github.com/concourse/concourse/atc/worker2"
	. "github.com/onsi/gomega"
)

var dummyLogger = lagertest.NewTestLogger("dummy")

type Scenario struct {
	DB        *dbtest.Scenario
	DBBuilder dbtest.Builder

	Pool    worker2.Pool
	Factory *Factory

	TeamID   int
	JobID    int
	StepName string
}

type SetupFunc func(*Scenario)

func Setup(poolFactory PoolFactory, builder dbtest.Builder, setup ...SetupFunc) *Scenario {
	scenario := &Scenario{
		DB:        dbtest.Setup(),
		DBBuilder: builder,

		Factory: &Factory{},
	}
	scenario.Pool = poolFactory(scenario.Factory)
	scenario.Run(setup...)
	return scenario
}

func (s *Scenario) Run(setup ...SetupFunc) {
	for _, f := range setup {
		f(s)
	}
}

func WithBasicJob() SetupFunc {
	return func(s *Scenario) {
		s.DB.Run(s.DBBuilder.WithTeam("team"))
		s.DB.Run(s.DBBuilder.WithPipeline(atc.Config{
			Jobs: []atc.JobConfig{{Name: "job"}},
		}))
		job := s.DB.Job("job")
		s.TeamID = s.DB.Team.ID()
		s.JobID = job.ID()
		s.StepName = "some-step"
	}
}

func WithWorkers(workers ...Worker) SetupFunc {
	return func(s *Scenario) {
		for _, worker := range workers {
			_, _, found := s.Factory.FindWorker(worker.Name())
			Expect(found).To(BeFalse(), "cannot add worker twice: %s", worker.Name())
			worker.Setup(s)

			s.Factory.Workers = append(s.Factory.Workers, worker)
		}
	}
}

func (s *Scenario) Worker(name string) runtime.Worker {
	worker, found, err := s.Pool.FindWorker(dummyLogger, name)
	Expect(err).ToNot(HaveOccurred())
	Expect(found).To(BeTrue())
	return worker
}

func (s *Scenario) WorkerTaskCacheVolume(workerName string, path string) db.CreatedVolume {
	utc, found, err := s.DBBuilder.TaskCacheFactory.Find(s.JobID, s.StepName, path)
	Expect(err).ToNot(HaveOccurred())
	Expect(found).To(BeTrue())

	volume, found, err := s.DBBuilder.VolumeRepo.FindTaskCacheVolume(s.TeamID, workerName, utc)
	Expect(err).ToNot(HaveOccurred())
	Expect(found).To(BeTrue())

	return volume
}
