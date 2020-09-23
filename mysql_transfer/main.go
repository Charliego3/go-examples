package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/whimthen/kits/logger"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	sourceDB      *sql.DB
	targetDB      *sql.DB
	transferFiles []string
	sourceLink    string
	sourceTables  []string
	targetLink    string
	targetTables  []string
	truncate      bool
)

func main() {
	root := &cobra.Command{
		Use:   "transfer",
		Short: "Transfer source data to target from file or db.",
		Run:   transfer,
	}

	root.AddCommand(fileCommand(), dbCommand())

	err := root.Execute()
	if err != nil {
		logger.Fatal(err)
	}
}

func fileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "file",
		Aliases: []string{"f"},
		Short:   "Transfer file to target mysql",
		Run:     fileTransfer,
	}

	cmd.Flags().StringSliceVarP(&transferFiles, "file", "f", nil, "csv or sql file")
	bindFlags(cmd)
	return cmd
}

func dbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Transfer db to target mysql",
		Run:   dbTransfer,
	}

	cmd.Flags().StringVarP(&sourceLink, "sourceLink", "s", "", "The source database connection address")
	cmd.Flags().StringSliceVarP(&sourceTables, "sourceTable", "o", nil, "The source table names")
	bindFlags(cmd)
	return cmd
}

func bindFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&targetLink, "targetLink", "t", "", "The target database connection address")
	cmd.Flags().StringSliceVarP(&targetTables, "targetTable", "a", nil, "The target table name, which needs to match the source table name")
	cmd.Flags().BoolVarP(&truncate, "truncate", "e", false, "Do you want to do data transfer after truncate table?")
}

func fileTransfer(cmd *cobra.Command, _ []string) {
	if transferFiles == nil || len(transferFiles) == 0 {
		logger.Error("Please carry the filepath parameter to import the data in the file\n")
		_ = cmd.Help()
		return
	}

	if targetTables == nil || len(targetTables) == 0 {
		logger.Warn("No target table is entered, no migration will be done\n")
		_ = cmd.Help()
		return
	}

	checkLen(transferFiles, true)

	targetDB = initDb(targetLink, true)

	for _, file := range transferFiles {
		ext := filepath.Ext(file)
		lext := strings.ToLower(ext)
		if lext == "csv" || lext == "tsv" || lext == "psv" {
			bs, err := ioutil.ReadFile(file)
			if err != nil {
				logger.Error("Read file(%s) error: %+v", file, err)
				return
			}

			reader := csv.NewReader(bufio.NewReader(strings.NewReader(string(bs))))

			for {
				records, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					logger.Error("Read record error: %+v", err)
					return
				}

				print(records)
			}
		}
	}
}

func dbTransfer(*cobra.Command, []string) {
	flen := len(transferFiles)
	tlen := len(targetTables)
	if flen != tlen {
		var ignoreType string
		var ignore []string
		if flen > tlen {
			ignoreType = "File"
			ignore = transferFiles[tlen:]
		} else {
			ignoreType = "Table"
			ignore = targetTables[flen:]
		}
		if ignore != nil && len(ignore) > 0 {
			logger.Warn("%s [%s] will be ignore to transfer", ignoreType, strings.Join(ignore, ", "))
		}
	}
}

func checkLen(source []string, isFile bool) {
	flen := len(source)
	tlen := len(targetTables)
	if flen != tlen {
		var ignore []string
		var ignoreType string
		if flen > tlen {
			if isFile {
				ignoreType = "File"
			} else {
				ignoreType = "Source table"
			}
			ignore = source[tlen:]
		} else {
			if isFile {
				ignoreType = "Table"
			} else {
				ignoreType = "Target table"
			}
			ignore = targetTables[flen:]
		}
		if ignore != nil && len(ignore) > 0 {
			logger.Warn("%s [%s] will be ignore to transfer", ignoreType, strings.Join(ignore, ", "))
		}
	}
}

func transfer(cmd *cobra.Command, _ []string) {
	_ = cmd.Help()
}

func getFields(db *sql.DB) {

}

func initDb(ds string, isTarget bool) *sql.DB {
	db, err := sql.Open("mysql", ds)
	checkErrFunc(err, func() string {
		msg := "Connect database(`%s`) error"
		if isTarget {
			msg = fmt.Sprintf(msg, targetLink)
		} else {
			msg = fmt.Sprintf(msg, sourceLink)
		}
		return msg
	})
	err = db.Ping()
	checkErrFunc(err, func() string {
		msg := "Database(`%s`) ping error"
		if isTarget {
			msg = fmt.Sprintf(msg, targetLink)
		} else {
			msg = fmt.Sprintf(msg, sourceLink)
		}
		return msg
	})
	return db
}

func checkErrFunc(err error, f func() string) {
	checkErr(err, f())
}

func checkErr(err error, msg string) {
	if err != nil {
		logger.Error("%s -> %+v", msg, err)
		os.Exit(1)
	}
}
