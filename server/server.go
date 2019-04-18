package server

import (
	"errors"
	fmt "fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/molotovtv/tc/internal/config"
	"github.com/molotovtv/tc/tc"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type tcServer struct {
	lock   sync.RWMutex
	builds map[string]map[string]*tc.Build
}

func (s *tcServer) clear() {
	s.builds = make(map[string]map[string]*tc.Build)
}

func (s *tcServer) run() {
	s.clear()

	refresh := func() {
		c, err := config.Load()
		if err != nil {
			s.clear()
			log.Println(err)
			return
		}

		for project, envs := range c.BuildIDs {
			for env, buildID := range envs {
				build, err := tc.LastBuild(c, buildID)
				if err != nil {
					log.Println(err)
					s.clear()
					continue
				}
				if build == nil {
					continue
				}

				s.lock.Lock()
				_, ok := s.builds[project]
				if !ok {
					s.builds[project] = make(map[string]*tc.Build)
				}

				s.builds[project][env] = build
				s.lock.Unlock()
			}
		}
	}

	refresh()
	for range time.Tick(5 * time.Second) {
		refresh()
	}
}

func (s *tcServer) LastBuild(ctx context.Context, b *ProjectEnv) (*tc.Build, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	p, ok := s.builds[b.Project]
	if !ok {
		return nil, errors.New("unknow project / env")
	}

	build, ok := p[b.Env]
	if !ok {
		return nil, errors.New("unknow project / env")
	}

	return build, nil
}

func Run(listen string) error {
	t := &tcServer{}
	go t.run()
	grpcServer := grpc.NewServer()
	RegisterTCServiceServer(grpcServer, t)

	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Printf("listening on %s\n", listen)
	return grpcServer.Serve(lis)
}
