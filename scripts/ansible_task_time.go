// Package ansible_task_time sorts through a single verbose (-v) ansible-playbook log
// to show the execution time of each task.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	var ansibleLogTaskLines []string
	ansibleLogFile := "ansible.log"
	ansibleLogOpen, _ := os.Open(ansibleLogFile)
	ansibleLogScanner := bufio.NewScanner(ansibleLogOpen)
	defer ansibleLogOpen.Close()

	for ansibleLogScanner.Scan() {

		if strings.Contains(ansibleLogScanner.Text(), "TASK") {
			ansibleLogTaskLines = append(ansibleLogTaskLines, ansibleLogScanner.Text())
		}

	}

	var taskTimeStart time.Time
	var taskTimeEnd time.Time
	var taskTimeTotal time.Duration
	var taskName string
	for index, line := range ansibleLogTaskLines {
		// ansible.log time format: 2020-04-21 08:02:21
		// time.RFC3339 format: 2020-04-21T08:02:21"
		taskTime, taskTimeErr := time.Parse(time.RFC3339Nano, line[0:10]+"T"+line[11:19]+"."+line[20:23]+"Z")

		if taskTimeErr != nil {
			fmt.Println(taskTimeErr)
		}

		if index == 0 {
			taskTimeStart = taskTime
			taskName = line[46:len(line)]
			continue
		} else {
			taskTimeEnd = taskTime
		}

		taskTimeTotal = taskTimeEnd.Sub(taskTimeStart)
		fmt.Println(taskTimeTotal.Seconds(), taskName)
		taskTimeStart = taskTime
		taskName = line[46:len(line)]
	}

}
