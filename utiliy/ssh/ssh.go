package ssh

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
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
	client, err := ssh.Dial("tcp", "remote-server:22", config)
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

	// 执行命令
	output, err := session.CombinedOutput("uptime")
	if err != nil {
		log.Fatalf("Failed to execute command: %s", err)
	}

	// 打印命令输出
	fmt.Println(string(output))
}
