package micro

import (
	"os"
	"syscall"
	"os/signal"
	"sync"
	"errors"
	l "github.com/efureev/micro/log"
)

type Service interface {
	//Init(...Option)
	Options() Options

	setParent(p *Service)
	Parent() *Service
	Root() *Service
	Services() map[string]*Service
	AddService(*Service) *BaseService
	GetService(string) *Service
	GetChildByName(string) *Service
	IsRoot() bool
	GetLevel() int

	Run() error
	Start() error
	Stop() error
	GetStatus() ServiceStatus
	IsRunning() bool

	ToLogMsg(msg string) string
	String() string
}

type ServiceStatus int

const (
	ServiceStatusStop ServiceStatus = iota
	ServiceStatusRun
)

type BaseService struct {
	sync.RWMutex
	opts   Options
	status ServiceStatus

	parent   *Service
	services map[string]*Service
}

func (s *BaseService) Options() Options {
	return s.opts
}

func (s *BaseService) String() string {
	return s.opts.Name
}

func (s *BaseService) Parent() *Service {
	return s.parent
}

func (s *BaseService) Root() *Service {
	if s.parent == nil {
		S := Service(s)
		return &S
	}

	return (*s.parent).Root()
}

func (s *BaseService) setParent(p *Service) {
	s.parent = p
}

func (s *BaseService) Services() map[string]*Service {
	return s.services
}

func (s *BaseService) IsRoot() bool {
	return s.parent == nil
}

func (s *BaseService) GetLevel() int {
	if s.parent == nil {
		return 0
	}

	return (*s.parent).GetLevel() + 1
}

/*
func (s *BaseService) GetService(name string) *Service {
	return s.services[name]
}*/

func (s *BaseService) GetChildByName(name string) *Service {
	return s.services[name]
}

func (s *BaseService) GetService(path string) *Service {
	root := s.Root()

	//todo: сделать на рекурсивный путь с разделителем точка

	return (*root).GetChildByName(path)
}

func (s *BaseService) AddService(child *Service) *BaseService {
	S := Service(s)
	(*child).setParent(&S)

	s.services[(*child).Options().Name] = child

	return s
}

func (s *BaseService) GetStatus() ServiceStatus {
	return s.status
}

func (s *BaseService) IsRunning() bool {
	return s.status == ServiceStatusRun
}

func (s *BaseService) Start() error {
	S := Service(s)

	for _, fn := range s.opts.BeforeStart {
		if err := fn(&S); err != nil {
			return err
		}
	}

	s.status = ServiceStatusRun

	l.Log(s.ToLogMsg(`service started`))

	for _, childS := range s.Services() {
		if (*childS).IsRunning() {
			l.Fatal(errors.New(`service <` + s.opts.Name + `> is already running`))
			continue
		}

		go func(childS *Service) {
			if err := (*childS).Start(); err != nil {
				l.Fatal(err)
				return
			}
		}(childS)
	}

	for _, fn := range s.opts.AfterStart {
		if err := fn(&S); err != nil {
			return err
		}
	}

	return nil
}

func (s *BaseService) Stop() error {

	if !s.IsRunning() {
		return nil
	}

	var gErr error
	S := Service(s)

	for _, fn := range s.opts.BeforeStop {
		if err := fn(&S); err != nil {
			gErr = err
		}
	}

	for _, child := range s.Services() {

		if err := (*child).Stop(); err != nil {
			l.Log(err)
		}
	}

	s.status = ServiceStatusStop

	//l.Log(getOffset(s.GetLevel(), "---") + `[service <` + s.opts.Name + `>] stoped...`)
	//l.Log(GetLogPrefix() + `[service <` + s.opts.Name + `>] stoped...`)
	l.Log(s.ToLogMsg(`service stopped`))

	for _, fn := range s.opts.AfterStop {
		if err := fn(&S); err != nil {
			gErr = err
		}
	}

	return gErr
}

func (s *BaseService) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	// wait on kill signal
	case <-exit:
		// wait on context cancel
	case <-s.opts.Context.Done():
	}

	if err := s.Stop(); err != nil {
		return err
	}

	return nil
}

func (s *BaseService) ToLogMsg(msg string) string {
	return GetOffsetString(s.GetLevel(), "---", "\t", "> ") + `[` + s.opts.Name + `] ` + msg
}

/*

func (s *service) Init() *service {
	s.Lock()
	configPathPtr := flag.String("config", "config.toml", "путь к конфигу")
	s.cmd = *flag.String("cmd", "listing", "Команда для исполнения: test|listing")
	flag.Parse()

	if _, err := toml.DecodeFile(*configPathPtr, &s); err != nil {
		log.Fatalln(err)
	}

	s.Config.BasePath, _ = os.Getwd()

	s.Unlock()
	return s
}
*/

func NewService(opts ...Option) Service {
	options := newOptions(opts...)

	return &BaseService{
		opts:     options,
		status:   ServiceStatusStop,
		services: make(map[string]*Service),
	}
}

/*
func getOffset(offset int, s string) string {

	if offset == 0 {
		return ""
	}

	prefix := "\t"
	postfix := "> "
	var str string

	for i := 0; i < offset; i++ {
		str += s
	}

	return prefix + str + postfix
}
*/
