package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var (
	url    string
	table  string
	limit  int64
	logger *log.Logger
	cmd    = &cobra.Command{
		Use: "deleter",
		Run: deletePublicNoticeHandler,
	}
)

func main() {
	cmd.Flags().StringVar(&url, "url", "", `--url "user:password@tcp(ip:port)/dbname"`)
	cmd.Flags().StringVar(&table, "table", "publicnoticecontainer", `--table "publicnoticecontainer"`)
	cmd.Flags().Int64Var(&limit, "limit", 1000, "--limit 1000")
	_ = cmd.Flags().Parse(os.Args)

	if url == "" {
		_ = cmd.Help()
		return
	}

	file, err := os.Create(fmt.Sprintf("delete_%s.log", table))
	if err != nil {
		log.Println("创建日志文件失败", err.Error())
		return
	}
	logger = log.New(file, "", log.LstdFlags)

	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func deletePublicNoticeHandler(*cobra.Command, []string) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		logger.Println("链接数据库失败", err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		logger.Println("数据库链接失败", err.Error())
		return
	}

	logger.Printf("开始删除 [%s] 表数据, 每次删除 [%d] 条\n", table, limit)

	for {
		result, err := db.Exec(fmt.Sprintf("DELETE FROM %s LIMIT ?", table), limit)
		if err != nil {
			logger.Println("删除出错", err.Error())
			return
		}

		affected, err := result.RowsAffected()
		if err != nil {
			logger.Println("获取执行行数出错", err.Error())
		} else if affected > 0 {
			logger.Printf("成功删除 [%d] 条数据, 休眠100ms\n", affected)
			time.Sleep(time.Millisecond * 100)
		}

		if affected < limit {
			break
		}
	}

	logger.Println("所有数据删除完成.....")
}
