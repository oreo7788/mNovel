// dbtest 仅测试数据库连接与迁移，不启动 HTTP 服务
package main

import (
	"fmt"
	"log"
	"os"

	"yidaiku-server/internal/config"
	"yidaiku-server/internal/repository"
)

func main() {
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	dbCfg := &cfg.Database
	fmt.Printf("正在连接: %s:%d/%s ...\n", dbCfg.Host, dbCfg.Port, dbCfg.Database)

	if err := repository.InitDB(dbCfg); err != nil {
		log.Fatalf("连接或迁移失败: %v\n请检查：\n  1. MySQL 是否已启动\n  2. config.yaml 中 database 的 host/port/username/password 是否正确\n  3. 数据库 %q 是否已创建（若未创建可执行: CREATE DATABASE %s;）", err, dbCfg.Database, dbCfg.Database)
	}
	defer repository.CloseDB()

	fmt.Println("数据库连接成功，迁移已完成。")
}
