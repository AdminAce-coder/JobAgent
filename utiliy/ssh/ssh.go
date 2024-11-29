package ssh

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SshConfig struct {
	User     string
	Password string
	Remote   string
	Port     string
}

func SshConnect() {
	// SSH 配置
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("nissan@123"), // 使用密码认证
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略主机密钥检查（生产环境不建议）
	}

	// 连接到远程服务器
	client, err := ssh.Dial("tcp", "1.92.75.225:22", config)
	if err != nil {
		log.Fatalf("Failed to dial: %s", err)
	}
	defer client.Close()

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %s", err)
	}
	defer session.Close()

	// 请求伪终端
	err = session.RequestPty("xterm", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          1,     // 启用回显
		ssh.TTY_OP_ISPEED: 14400, // 输入速度
		ssh.TTY_OP_OSPEED: 14400, // 输出速度
	})
	if err != nil {
		log.Fatalf("Failed to request pty: %s", err)
	}

	// 连接标准输入输出
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// 获取当前终端状态
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to set terminal to raw mode: %s", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState) // 在退出时恢复终端状态

	// 捕获 Ctrl+C 等信号以安全退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		term.Restore(int(os.Stdin.Fd()), oldState) // 恢复终端状态
		os.Exit(0)
	}()

	// 启动交互式 Shell
	err = session.Shell()
	if err != nil {
		log.Fatalf("Failed to start shell: %s", err)
	}

	// 等待会话结束
	err = session.Wait()
	if err != nil {
		log.Printf("Session ended with error: %s", err)
	}
}
