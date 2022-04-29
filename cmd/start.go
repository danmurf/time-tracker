/*
Copyright © 2022 Dan Murfitt <dan@murfitt.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"database/sql"
	"fmt"
	"github.com/danmurf/time-tracker/internal/tasks"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start working on a task",
	Long: `Record that you have started working on a specific task, for example:

tt start task1`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.PrintErrln("command usage is `tt start <task-name>`")
			return
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			cmd.PrintErrln(fmt.Errorf("finding user home directory: %w", err))
			return
		}

		dbPath := fmt.Sprintf("%s/%s", homeDir, ".time-tracker")
		if err = os.MkdirAll(dbPath, os.ModePerm); err != nil {
			cmd.PrintErrln(fmt.Errorf("creating time tracker directory [%s]: %w", dbPath, err))
			return
		}

		dbFilePath := fmt.Sprintf("%s/%s", dbPath, "time-tracker.db")
		db, err := sql.Open("sqlite3", dbFilePath)
		if err != nil {
			cmd.PrintErrln(fmt.Errorf("creating database: %w", err))
			return
		}

		eventStore, err := tasks.NewSQLEventStore(cmd.Context(), db)
		if err != nil {
			cmd.PrintErrln(fmt.Errorf("creating event store: %w", err))
			return
		}

		starter := tasks.NewStarter(eventStore)
		taskName := args[0]
		if err = starter.Start(cmd.Context(), taskName); err != nil {
			cmd.PrintErrln(fmt.Errorf("creating task starter: %w", err))
			return
		}

		cmd.Printf("⏱  %s started. Run `tt finish %s` when you have finished work.\n", taskName, taskName)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
