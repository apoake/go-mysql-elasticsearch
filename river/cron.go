package river

import(
	"github.com/robfig/cron"
	"github.com/pkg/errors"
	"fmt"
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
	cronn := &crons{isSleep: true, cronCh: cronCh, c: r.c}
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
	cronn.Cron = *c
	return cronn, nil
}

func (c *crons) start() {
	if !c.c.IsCron {
		return
	}
	c.Cron.Start()
	fmt.Printf("cron start ....")
}

func (c *crons) status(method string) {
	fmt.Printf("method[%s] --> isSleep: %t", method, c.isSleep)
	fmt.Println("-----------------------------")
}

func (c *crons) suspend() {
	c.status("suspend")
	if !c.c.IsCron || c.isSleep {
		return
	}
	c.isSleep = true
	fmt.Println("do suspend")
}

func (c *crons) restart() {
	println("------------------5------------")
	c.status("restart")
	if !c.c.IsCron || !c.isSleep {
		return
	}
	c.isSleep = false
	c.cronCh <- false
	fmt.Println("do restart")
}

func (c *crons) canRun() {
	if !c.c.IsCron || !c.isSleep {
		return
	}
	c.isSleep = <- c.cronCh
	fmt.Println("do canRun")
}

func (c *crons) stop() {
	if !c.c.IsCron {
		return
	}
	c.Stop()
	fmt.Println("do stop")
}
