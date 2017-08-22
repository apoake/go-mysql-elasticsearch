package river

import(
	"github.com/robfig/cron"
	"github.com/pkg/errors"
)

type crons struct {
	cron.Cron
	c *Config
	cronCh chan bool
	isSleep bool
}

func NewCron(r *River) (*crons, error) {
	c := cron.New()
	// add
	cronCh := make(chan bool)
	cronn := &crons{*c, r.c, cronCh, true}
	if r.c.IsCron {
		if r.c.StartJob == "" {
			return nil, errors.New("not config start_job")
		}
		if r.c.StopJob == "" {
			return nil, errors.New("not config stop_job")
		}
		c.AddFunc(r.c.StartJob, cronn.restart)
		c.AddFunc(r.c.StopJob, cronn.suspend)
	}
	return cronn, nil
}

func (c *crons) start() {
	if c.c.IsCron {
		return
	}
	c.isSleep = false
	c.Start()
}

func (c *crons) suspend() {
	if c.c.IsCron || c.isSleep {
		return
	}
	c.isSleep = true
}

func (c *crons) restart() {
	if c.c.IsCron || !c.isSleep {
		return
	}
	c.isSleep = false
	c.cronCh <- false
}

func (c *crons) canRun() {
	if c.c.IsCron || !c.isSleep {
		return
	}
	c.isSleep = <- c.cronCh
}

func (c *crons) stop() {
	if c.c.IsCron {
		return
	}
	c.Stop()
}
