package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"mxshop/user_srv/global"
	"mxshop/user_srv/handler"
	"mxshop/user_srv/initialize"
	"mxshop/user_srv/proto"
	"mxshop/user_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	zap.S().Info(global.ServerConfig)

	flag.Parse()
	zap.S().Info("IP:", *IP)
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("Port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("10.60.83.138:%d", *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = []string{"GUET", "zhangjiayuan", "user", "srv"}
	registration.Check = check
	registration.Address = "10.60.83.138"
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收中止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
	fmt.Println("123123")
}
