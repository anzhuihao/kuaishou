package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

var t int
var err error
var counts int
var cmd *exec.Cmd

func waitUntilLanched() {
	for count:=0;count<10;count++{
		cmd = exec.Command("cmd", "/K", "adb shell dumpsys window | findstr mCurrentFocus")
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run() // 无需等待执行结果，如果使用Run()方法需要等待写入stdout而一直卡住
		if err != nil {
			log.Fatalln("执行命令失败:", err)
		}
		s, _ := out.ReadString('\n')
		if strings.Contains(s, "vphone.launcher.Launcher") {
			log.Println("已经入nox主页面")
			return
		}
		time.Sleep(time.Second * 5)
	}
	log.Fatalln("无法进入nox主页面")
}

func startNox() {
	// 开启夜神模拟器 todo 支持多开
	cmd = exec.Command("cmd", "/K", "Nox.exe")
	err = cmd.Start() // 无需等待执行结果，如果使用Run()方法需要等待写入stdout而一直卡住
	if err != nil {
		log.Fatalln("执行命令失败", err)
	}

	// 等待直到进入主页面
	waitUntilLanched()
}

func quitNox() {
	cmd = exec.Command("cmd", "/K", "Nox.exe -quit")
	err = cmd.Start() // 无需等待执行结果，如果使用Run()方法需要等待写入stdout而一直卡住
	if err != nil {
		log.Fatalln("执行关闭命令失败，请手动关闭", err)
	}
	log.Fatalln("已关闭nox")
}

func openKuaishou() {
	cmd = exec.Command("cmd", "/K", "adb shell am start -n com.kuaishou.nebula/com.yxcorp.gifshow.HomeActivity")
	cmd.Run()
	time.Sleep(time.Second * 20)
}

func checkWindow() {
	// 检查窗口状态并纠正错误窗口
	cmd = exec.Command("cmd", "/K", "adb shell dumpsys window | findstr mCurrentFocus")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	s, _ := out.ReadString('\n')
	if !strings.Contains(s, "kuaishou") { // 表示未进入快手主页面
		fmt.Println("意外退出，将重新进入快手页面")
		openKuaishou()
	}
}

func run() {
	openKuaishou()
	for i:=1;i<counts+1;i++{
		checkWindow()
		cmd = exec.Command("cmd", "/K", "adb shell input swipe 500 1500 500 1000")
		cmd.Run()
		log.Printf("第%d次刷视频\n", i)
		time.Sleep(time.Second * time.Duration(t))
	}
}

func init() {
	flag.IntVar(&counts, "c", 100, "指定次数")
	flag.IntVar(&t, "t", 8, "指定刷视频间隔秒")
	flag.Parse()
}

func main() {
	startNox()
	run()
	quitNox()
}
