package logic

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/AdminAce-coder/jobAgent/internal/svc"
	"github.com/AdminAce-coder/jobAgent/pb/jobAgent"
	"github.com/zeromicro/go-zero/core/logx"
)

type DoJobLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDoJobLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DoJobLogic {
	return &DoJobLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DoJobLogic) DoJob(in *jobAgent.Request) (*jobAgent.Response, error) {
	// 修改执行命令的方式
	fmt.Printf("执行的命令是%s", in.Command)
	cmd := exec.Command("bash", "-c", "free -m") // 使用 bash -c 来执行命令
	memoryUsage, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return &jobAgent.Response{Result: string(memoryUsage)}, nil
}
