package micro

import (
	"context"
	"testing"
	"sync"
	"fmt"
	"log"
	"time"
)

func TestService(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	serviceChild := NewService(
		Name("children"),
		AfterStart(func(s *Service) error {
			fmt.Println("AfterStart #1 children: ", time.Now())
			return nil
		}),
		AfterStart(func(s *Service) error {
			fmt.Println("AfterStart #2 children: ", time.Now())
			fmt.Println("parent is: ", *(*s).Parent())
			wg.Done()
			return nil
		}),
		BeforeStart(func(s *Service) error {
			fmt.Println("BeforeStart children: ", time.Now())
			return nil
		}),
	)

	service := NewService(
		Name("main"),
		Context(ctx),
		AfterStart(func(s *Service) error {
			fmt.Println("AfterStart main: ", time.Now())
			return nil
		}),
	)

	service.AddService(&serviceChild)

	go func() {
		// wait for start
		wg.Wait()

		fmt.Println("Time : ", time.Now())
		//time.Sleep(time.Duration(4) * time.Second)
		//fmt.Println("Wait completed")
		// shutdown the service
		cancel()
	}()

	service.Run()
}

func TestService2(t *testing.T) {

	serviceChild := NewService(
		Name("children"),
		BeforeStart(func(s *Service) error {
			log.Println("BeforeStart children: ", time.Now())
			return nil
		}),
		AfterStart(func(s *Service) error {
			log.Println("AfterStart children: ", time.Now())
			return nil
		}),
	)

	serviceChildChild := NewService(
		Name("children2"),
		BeforeStart(func(s *Service) error {
			log.Println("BeforeStart children2: ", time.Now())
			return nil
		}),
		AfterStart(func(s *Service) error {
			log.Println("AfterStart children2: ", time.Now())
			return nil
		}),
	)

	service := NewService(
		Name("main"),
		AfterStart(func(s *Service) error {
			log.Println("AfterStart main: ", time.Now())
			return nil
		}),
	)


	serviceChild.AddService(&serviceChildChild)
	service.AddService(&serviceChild)

	service.Run()


}
