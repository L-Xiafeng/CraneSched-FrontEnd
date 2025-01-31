package ccontrol

import (
	"CraneFrontEnd/generated/protos"
	"CraneFrontEnd/internal/util"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	stub protos.CraneCtldClient
)

func mapToKVString(data map[string]uint64) string {
	var kvStrings []string

	for key, value := range data {
		kvStrings = append(kvStrings, fmt.Sprintf("%s:%d", key, value))
	}
	return strings.Join(kvStrings, ", ")
}

func ShowNodes(nodeName string, queryAll bool) {
	var req *protos.QueryCranedInfoRequest

	req = &protos.QueryCranedInfoRequest{CranedName: nodeName}
	reply, err := stub.QueryCranedInfo(context.Background(), req)
	if err != nil {
		panic("QueryNodeInfo failed: " + err.Error())
	}

	var B2MBRatio uint64 = 1024 * 1024

	if queryAll {
		if len(reply.CranedInfoList) == 0 {
			fmt.Printf("No node is avalable.\n")
		} else {
			for _, nodeInfo := range reply.CranedInfoList {
				fmt.Printf("NodeName=%v State=%v CPU=%.2f AllocCPU=%.2f FreeCPU=%.2f\n"+
					"\tRealMemory=%dM AllocMem=%dM FreeMem=%dM\n"+
					"\tGres={%s} AllocGres={%s} FreeGres={%s}\n"+
					"\tPatition=%s RunningJob=%d\n\n",
					nodeInfo.Hostname, nodeInfo.State.String()[6:], nodeInfo.Cpu,
					math.Abs(nodeInfo.AllocCpu),
					math.Abs(nodeInfo.FreeCpu),
					nodeInfo.RealMem/B2MBRatio, nodeInfo.AllocMem/B2MBRatio, nodeInfo.FreeMem/B2MBRatio,
					mapToKVString(nodeInfo.Device), mapToKVString(nodeInfo.AllocDevice), mapToKVString(nodeInfo.AvailDevice),
					strings.Join(nodeInfo.PartitionNames, ","), nodeInfo.RunningTaskNum)
			}
		}
	} else {
		if len(reply.CranedInfoList) == 0 {
			fmt.Printf("Node %s not found.\n", nodeName)
		} else {
			for _, nodeInfo := range reply.CranedInfoList {
				fmt.Printf("NodeName=%v State=%v CPU=%.2f AllocCPU=%.2f FreeCPU=%.2f\n"+
					"\tRealMemory=%dM AllocMem=%dM FreeMem=%dM\n"+
					"\tGres={%s} AllocGres={%s} FreeGres={%s}\n"+
					"\tPatition=%s RunningJob=%d\n\n",
					nodeInfo.Hostname, nodeInfo.State.String()[6:], nodeInfo.Cpu, nodeInfo.AllocCpu, nodeInfo.FreeCpu,
					nodeInfo.RealMem/B2MBRatio, nodeInfo.AllocMem/B2MBRatio, nodeInfo.FreeMem/B2MBRatio,
					mapToKVString(nodeInfo.Device), mapToKVString(nodeInfo.AllocDevice), mapToKVString(nodeInfo.AvailDevice),
					strings.Join(nodeInfo.PartitionNames, ","), nodeInfo.RunningTaskNum)
			}
		}
	}
}

func ShowPartitions(partitionName string, queryAll bool) {
	var req *protos.QueryPartitionInfoRequest

	req = &protos.QueryPartitionInfoRequest{PartitionName: partitionName}
	reply, err := stub.QueryPartitionInfo(context.Background(), req)
	if err != nil {
		panic("QueryPartitionInfo failed: " + err.Error())
	}

	var B2MBRatio uint64 = 1024 * 1024

	if queryAll {
		if len(reply.PartitionInfo) == 0 {
			fmt.Printf("No node is avalable.\n")
		} else {
			for _, partitionInfo := range reply.PartitionInfo {
				fmt.Printf("PartitionName=%v State=%v\n"+
					"\tTotalNodes=%d AliveNodes=%d\n"+
					"\tTotalCPU=%.2f AvailCPU=%.2f AllocCPU=%.2f\n"+
					"\tTotalMem=%dM AvailMem=%dM AllocMem=%dM\n"+
					"\tTotalGres={%s} AvailGres={%s} AllocGres={%s}"+
					"\tHostList=%v\n\n",
					partitionInfo.Name, partitionInfo.State.String()[10:],
					partitionInfo.TotalNodes, partitionInfo.AliveNodes,
					partitionInfo.TotalCpu, partitionInfo.AvailCpu, partitionInfo.AllocCpu,
					partitionInfo.TotalMem/B2MBRatio, partitionInfo.AvailMem/B2MBRatio, partitionInfo.AllocMem/B2MBRatio,
					mapToKVString(partitionInfo.Device), mapToKVString(partitionInfo.AvailDevice), mapToKVString(partitionInfo.AllocDevice),
					partitionInfo.Hostlist)
			}
		}
	} else {
		if len(reply.PartitionInfo) == 0 {
			fmt.Printf("Partition %s not found.\n", partitionName)
		} else {
			for _, partitionInfo := range reply.PartitionInfo {
				fmt.Printf("PartitionName=%v State=%v\n"+
					"\tTotalNodes=%d AliveNodes=%d\n"+
					"\tTotalCPU=%.2f AvailCPU=%.2f AllocCPU=%.2f\n"+
					"\tTotalMem=%dM AvailMem=%dM AllocMem=%dM\n"+
					"\tTotalGres={%s} AvailGres={%s} AllocGres={%s}"+
					"\tHostList=%v\n\n",
					partitionInfo.Name, partitionInfo.State.String()[10:],
					partitionInfo.TotalNodes, partitionInfo.AliveNodes,
					partitionInfo.TotalCpu, partitionInfo.AvailCpu, partitionInfo.AllocCpu,
					partitionInfo.TotalMem/B2MBRatio, partitionInfo.AvailMem/B2MBRatio, partitionInfo.AllocMem/B2MBRatio,
					mapToKVString(partitionInfo.Device), mapToKVString(partitionInfo.AvailDevice), mapToKVString(partitionInfo.AllocDevice),
					partitionInfo.Hostlist)
			}
		}
	}
}

func ShowTasks(taskId uint32, queryAll bool) {
	var req *protos.QueryTasksInfoRequest
	var taskIdList []uint32
	taskIdList = append(taskIdList, taskId)
	if queryAll {
		req = &protos.QueryTasksInfoRequest{}
	} else {
		req = &protos.QueryTasksInfoRequest{FilterTaskIds: taskIdList}
	}

	reply, err := stub.QueryTasksInfo(context.Background(), req)
	if err != nil {
		panic("QueryTasksInfo failed: " + err.Error())
	}

	if !reply.GetOk() {
		log.Errorf("Failed to retrive the information of job %d", taskId)
		os.Exit(1)
	}

	if len(reply.TaskInfoList) == 0 {
		if queryAll {
			fmt.Printf("No job is running.\n")
		} else {
			fmt.Printf("Job %d is not running.\n", taskId)
		}

	} else {
		for _, taskInfo := range reply.TaskInfoList {
			timeSubmit := taskInfo.SubmitTime.AsTime()
			timeSubmitStr := "unknown"
			if !timeSubmit.Before(time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)) {
				timeSubmitStr = timeSubmit.Local().String()
			}
			timeStart := taskInfo.StartTime.AsTime()
			timeStartStr := "unknown"
			runTime := "unknown"
			if !timeStart.Before(time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)) {
				timeStartStr = timeStart.Local().String()
				runTime = time.Now().Sub(timeStart).String()
			}
			timeEnd := taskInfo.EndTime.AsTime()
			timeEndStr := "unknown"
			if timeEnd.After(timeStart) {
				timeEndStr = timeEnd.Local().String()
				runTime = timeEnd.Sub(timeStart).String()
			}

			fmt.Printf("JobId=%v JobName=%v\n\tUserId=%d GroupId=%d Account=%v\n\tJobState=%v RunTime=%v "+
				"TimeLimit=%s SubmitTime=%v\n\tStartTime=%v EndTime=%v Partition=%v NodeList=%v "+
				"NumNodes=%d\n\tCmdLine=%v Workdir=%v\n",
				taskInfo.TaskId, taskInfo.Name, taskInfo.Uid, taskInfo.Gid,
				taskInfo.Account, taskInfo.Status.String(), runTime, util.SecondTimeFormat(taskInfo.TimeLimit.Seconds),
				timeSubmitStr, timeStartStr, timeEndStr, taskInfo.Partition,
				taskInfo.CranedList, taskInfo.NodeNum, taskInfo.CmdLine, taskInfo.Cwd)
		}
	}
}

func ChangeTaskTimeLimit(taskId uint32, timeLimit string) {
	re := regexp.MustCompile(`((.*)-)?(.*):(.*):(.*)`)
	result := re.FindAllStringSubmatch(timeLimit, -1)

	if result == nil || len(result) != 1 {
		log.Fatalf("Time format error")
	}
	var dd uint64
	if result[0][2] != "" {
		dd, _ = strconv.ParseUint(result[0][2], 10, 32)
	}
	hh, err := strconv.ParseUint(result[0][3], 10, 32)
	if err != nil {
		log.Fatalf("The hour time format error")
	}
	mm, err := strconv.ParseUint(result[0][4], 10, 32)
	if err != nil {
		log.Fatalf("The minute time format error")
	}
	ss, err := strconv.ParseUint(result[0][5], 10, 32)
	if err != nil {
		log.Fatalf("The second time format error")
	}

	seconds := int64(60*60*24*dd + 60*60*hh + 60*mm + ss)

	var req *protos.ModifyTaskRequest

	req = &protos.ModifyTaskRequest{
		Uid:       uint32(os.Getuid()),
		TaskId:    taskId,
		Attribute: protos.ModifyTaskRequest_TimeLimit,
		Value: &protos.ModifyTaskRequest_TimeLimitSeconds{
			TimeLimitSeconds: seconds,
		},
	}
	reply, err := stub.ModifyTask(context.Background(), req)
	if err != nil {
		panic("ModifyTask failed: " + err.Error())
	}

	if reply.Ok {
		fmt.Println("Change time limit success")
	} else {
		fmt.Printf("Chang time limit failed: %s\n", reply.GetReason())
	}
}
