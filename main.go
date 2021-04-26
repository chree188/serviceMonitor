// +build windows
// +build 386
package main

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

type myApp struct {
	title string
	msg   *walk.TextEdit
	mw    *walk.MainWindow
}

type myService struct {
	text        string
	serviceName string
	labelState  *walk.Label
	btnStart    *walk.PushButton
	btnStop     *walk.PushButton
}

var app myApp
var service1, service2, service3 myService

func init() {
	app.title = "XXXX管理系统-运行监控"
	service1 = myService{
		text:        "MySql（MySql数据库）",
		serviceName: "MySql",
	}
	service2 = myService{
		text:        "Web服务(nginx)",
		serviceName: "nginx",
	}
	service3 = myService{
		text:        "XXXX系统",
		serviceName: "xxxx-server",
	}
}

func main() {
	_ = getWindows()

	walk.App().SetProductName(app.title)
	walk.App().SetOrganizationName("dwt")

	_ = service1.labelState.SetText("未安装")
	_ = service2.labelState.SetText("未安装")
	_ = service3.labelState.SetText("未安装")
	service1.btnStart.Clicked().Attach(func() {
		startService(service1)
	})
	service1.btnStop.Clicked().Attach(func() {
		stopService(service1)
	})
	service2.btnStart.Clicked().Attach(func() {
		startService(service2)
	})
	service2.btnStop.Clicked().Attach(func() {
		stopService(service2)
	})
	service3.btnStart.Clicked().Attach(func() {
		startService(service3)
	})
	service3.btnStop.Clicked().Attach(func() {
		stopService(service3)
	})

	go flushServiceStat(service1)
	go flushServiceStat(service2)
	go flushServiceStat(service3)

	app.mw.Show()
	app.mw.Run()
}

func setServiceState(service myService, msg string, btnStartStatus bool, btnStopStatus bool) {
	_ = service.labelState.SetText(msg)
	service.btnStart.SetEnabled(btnStartStatus)
	service.btnStop.SetEnabled(btnStopStatus)
}

// 刷新服务状态的协程程序
func flushServiceStat(service myService) {
	for {
		winService, err := NewWinService(service.serviceName)
		if winService == nil || err != nil {
			if err == windows.ERROR_SERVICE_DOES_NOT_EXIST {
				setServiceState(service, "未安装", false, false)
			} else {
				setServiceState(service, "服务打开失败", false, false)
			}
		} else {
			if winService.IsStop() {
				setServiceState(service, "已经停止", true, false)
			} else if winService.IsRunning() {
				setServiceState(service, "正在运行", false, true)
			}
		}
		time.Sleep(time.Second)
	}
}

// 启动服务
func startService(service myService) {
	s, err := NewWinService(service.serviceName)
	if s == nil || err != nil {
		return
	}
	showMsg(service.serviceName + " 服务开始启动......")
	err = s.StartService()
	if err != nil {
		showMsg(service.serviceName + " 服务启动失败！")
	} else {
		showMsg(service.serviceName + " 服务启动成功。")
	}
}

// 停止服务
func stopService(service myService) {
	s, err := NewWinService(service.serviceName)
	if s == nil || err != nil {
		return
	}
	showMsg(service.serviceName + " 服务开始停止......")
	err = s.StopService()
	if err != nil {
		showMsg(service.serviceName + " 服务停止失败！")
	} else {
		showMsg(service.serviceName + " 服务停止成功。")
	}
}

func showMsg(msg string) {
	app.msg.AppendText(time.Now().Format("2006-01-02 15:04:05 "))
	app.msg.AppendText(msg)
	app.msg.AppendText("\r\n")
}

// 初始始化窗体
func getWindows() error {
	icon, _ := walk.NewIconFromResourceId(3)
	err := MainWindow{
		Visible:  false,
		AssignTo: &app.mw,
		Title:    app.title,
		Size:     Size{500, 360},
		Font:     Font{Family: "微软雅黑", PointSize: 9},
		Icon:     icon,
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				Title:  "基础服务状态",
				Layout: Grid{Columns: 3},
				Children: []Widget{
					Label{Text: service1.text, MinSize: Size{220, 30}, TextColor: walk.RGB(255, 255, 0)},
					Label{AssignTo: &service1.labelState, Text: "正在运行", MinSize: Size{80, 30}},
					Composite{
						Layout:  HBox{},
						MaxSize: Size{132, 30},
						Children: []Widget{
							PushButton{
								AssignTo: &service1.btnStop,
								MaxSize:  Size{60, 30},
								Text:     "停止",
							},
							PushButton{
								AssignTo: &service1.btnStart,
								MaxSize:  Size{60, 30},
								Text:     "启动",
							},
						},
					},
					Label{Text: service2.text},
					Label{AssignTo: &service2.labelState, Text: "正在运行"},
					Composite{
						Layout:  HBox{},
						MaxSize: Size{132, 30},
						Children: []Widget{
							PushButton{
								AssignTo: &service2.btnStop,
								MaxSize:  Size{60, 30},
								Text:     "停止",
							},
							PushButton{
								AssignTo: &service2.btnStart,
								MaxSize:  Size{60, 30},
								Text:     "启动",
							},
						},
					},
				},
			},
			GroupBox{
				Title:  "业务服务状态",
				Layout: Grid{Columns: 3},
				Children: []Widget{
					Label{Text: service3.text, MinSize: Size{220, 30}},
					Label{AssignTo: &service3.labelState, Text: "正在运行", MinSize: Size{80, 30}},
					Composite{
						Layout:  HBox{},
						MaxSize: Size{132, 30},
						Children: []Widget{
							PushButton{
								AssignTo: &service3.btnStop,
								MaxSize:  Size{60, 30},
								Text:     "停止",
							},
							PushButton{
								AssignTo: &service3.btnStart,
								MaxSize:  Size{60, 30},
								Text:     "启动",
							},
						},
					},
				},
			},
			TextEdit{AssignTo: &app.msg, VScroll: true, ReadOnly: true},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						MinSize: Size{160, 30},
						Text:    "打开windows服务管理程序",
						OnClicked: func() {
							c := exec.Command("cmd", "/C", "SERVICES.MSC")
							c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 不显示命令窗口
							if err := c.Start(); err != nil {
								showMsg(fmt.Sprintf("打开windows服务管理程序失败, 错误信息: %s", err))
							}
						},
					},
					HSpacer{},
					PushButton{
						MinSize: Size{121, 30},
						Text:    "关闭本监控窗口",
						OnClicked: func() {
							walk.App().Exit(0)
						},
					},
				},
			},
		},
		OnSizeChanged: func() {
			_ = app.mw.SetSize(walk.Size(Size{500, 360}))
		},
	}.Create()
	winLong := win.GetWindowLong(app.mw.Handle(), win.GWL_STYLE)
	// 不能调整窗口大小，禁用最大化按钮
	win.SetWindowLong(app.mw.Handle(), win.GWL_STYLE, winLong & ^win.WS_SIZEBOX & ^win.WS_MAXIMIZEBOX & ^win.WS_SIZEBOX)
	// 设置窗体生成在屏幕的正中间，并处理高分屏的情况
	// 窗体横坐标 = ( 屏幕宽度 - 窗体宽度 ) / 2
	// 窗体纵坐标 = ( 屏幕高度 - 窗体高度 ) / 2
	_ = app.mw.SetX((int(win.GetSystemMetrics(0)) - app.mw.Width()) / 2 / app.mw.DPI() * 96)
	_ = app.mw.SetY((int(win.GetSystemMetrics(1)) - app.mw.Height()) / 2 / app.mw.DPI() * 96)
	return err
}
