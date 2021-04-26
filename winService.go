package main

import (
	"fmt"
	"github.com/shirou/gopsutil/winservices"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"time"
)

// 定义Service对象
type Service struct {
	Name   string
	srv    *winservices.Service
	Config mgr.Config
	Status winservices.ServiceStatus
}

// 得到Service信息
func (s *Service) getServiceDetail() error {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	err := s.srv.GetServiceDetail()
	if err != nil {
		// 服务打开失败
		return err
	}
	s.Name = s.srv.Name
	s.Status = s.srv.Status
	s.Config = s.srv.Config
	return nil
}

// 启动一个服务
func (s *Service) StartService() error {
	return StartService(s.Name)
}

// 停止一个服务
func (s *Service) StopService() error {
	return StopService(s.Name)
}

// 服务是否已经停止
func (s *Service) IsStop() bool {
	return s.Status.State == svc.Stopped
}

// 服务是否正在运行
func (s *Service) IsRunning() bool {
	return s.Status.State == svc.Running
}

// 新建一个Service对象
func NewWinService(serviceName string) (*Service, error) {
	winService, err := winservices.NewService(serviceName)
	if err != nil {
		return nil, err
	}
	result := &Service{
		Name: serviceName,
		srv:  winService,
	}
	err = result.getServiceDetail()
	return result, err
}

// 得到Service信息
func GetServiceInfo(serviceName string) (*winservices.Service, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	newservice, _ := winservices.NewService(serviceName)
	err := newservice.GetServiceDetail()
	if err != nil {
		fmt.Println(serviceName, "服务打开失败!")
		return nil, err
	}
	return newservice, nil
}

// 启动一个服务
func StartService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

// 停止一个服务
func StopService(name string) error {
	return ControlService(name, svc.Stop, svc.Stopped)
}

// 改一个服务状态，
func ControlService(name string, c svc.Cmd, to svc.State) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}
