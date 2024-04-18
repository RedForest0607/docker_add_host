package runner

import (
	"context"
	"dockerAddHost/pkg/utils"
	"dockerAddHost/types"
	"fmt"
	types2 "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"os"
	"strconv"
	"strings"
)

func Run() {
	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	suffix := utils.StringPrompt("적용할 컨테이너의 suffix를 입력해주세요 (공백 시 전체 검색)")

	var filteredContainerList []types.ContainerData

	for _, containerList := range containers {
		for _, name := range containerList.Names {
			if strings.HasSuffix(name, suffix) {
				var data types.ContainerData
				data.ContainerName = containerList.Names[0]
				data.ContainerId = containerList.ID
				filteredContainerList = append(filteredContainerList, data)
			}
		}
	}

	if len(filteredContainerList) == 0 {
		fmt.Println("컨테이너가 없습니다.")
		utils.WriteLogToFile("Container Not Found")
		utils.WriteLogToFile("Suffix : " + suffix)
		return
	}

	answers := utils.Checkboxes("적용할 컨테이너를 선택해주세요. 스페이스바로 선택, 오른쪽 마우스로 전체선택 합니다.",
		filteredContainerList,
	)
	filteredContainerList = filterContainerByName(filteredContainerList, answers)
	hostList := inputHostInformation()

	fmt.Println(filteredContainerList)

	err = addHostToYml(filteredContainerList, hostList)
	if err != nil {
		return
	}

	addHostToContainer(apiClient, filteredContainerList, hostList)
}

func filterContainerByName(containers []types.ContainerData, names []string) []types.ContainerData {

	containerMap := make(map[string]types.ContainerData)
	for _, containerInfo := range containers {
		containerMap[containerInfo.ContainerName] = containerInfo
	}

	var filteredContainers []types.ContainerData
	for _, name := range names {
		if containerInfo, ok := containerMap[name]; ok {
			filteredContainers = append(filteredContainers, containerInfo)
		}
	}

	return filteredContainers
}

func inputHostInformation() []types.HostsType {
	var hostDatum []types.HostsType
	for {
		var hostData types.HostsType

		hostData.HostName = utils.StringPrompt("host name을 입력해주세요")
		if hostData.HostName == "" {
			break
		}
		hostData.HostIp = utils.StringPrompt("host ip를 입력해주세요")
		for !utils.IsValidIPAddress(hostData.HostIp) {
			hostData.HostIp = utils.StringPrompt("올바른 ip 형식으로 입력해주세요")
		}
		hostDatum = append(hostDatum, hostData)
	}
	return hostDatum
}

func addHostToContainer(apiClient *client.Client, targetContainerList []types.ContainerData, hostDatum []types.HostsType) {
	for _, targetContainer := range targetContainerList {

		var commandInLine strings.Builder
		for _, hostData := range hostDatum {
			commandInLine.Reset()
			commandInLine.WriteString(hostData.HostIp)
			commandInLine.WriteString(" ")
			commandInLine.WriteString(hostData.HostName)
			command := []string{"bash", "-c", "echo " + commandInLine.String() + ">> /etc/hosts"}
			utils.WriteLogToFile(fmt.Sprintf("bash -c 'echo %s >> /etc/hosts'", commandInLine.String()))
			execConfig := types2.ExecConfig{
				AttachStdout: true,
				AttachStderr: true,
				Cmd:          command,
			}

			execResp, err := apiClient.ContainerExecCreate(context.Background(), targetContainer.ContainerId, execConfig)
			if err != nil {
				panic(err)
			}

			execID := execResp.ID
			execAttachConfig := types2.ExecStartCheck{}
			resp, err := apiClient.ContainerExecAttach(context.Background(), execID, execAttachConfig)
			if err != nil {
				panic(err)
			}
			utils.WriteLogToFile("Output:")
			result, _ := io.Copy(os.Stdout, resp.Reader)
			utils.WriteLogToFile(strconv.FormatInt(result, 10))

			resp.Close()
		}
	}
}

func addHostToYml(targetContainerList []types.ContainerData, hostDatum []types.HostsType) error {

	for _, targetContainer := range targetContainerList {
		var ymlContent string

		path := os.Getenv("DOCKER_FILES")

		// 기존 yml 파일 내용 읽어오기
		content, err := os.ReadFile(path + "/extra_hosts" + targetContainer.ContainerName + ".yml")
		utils.WriteLogToFile(fmt.Sprintf(path + "/extra_hosts/" + targetContainer.ContainerName + ".yml"))
		if err != nil {
			return err
		}
		ymlContent = string(content)

		// 각 컨테이너 정보를 yml 파일에 추가
		for _, hostData := range hostDatum {
			if hostData.HostIp != "" && hostData.HostName != "" {
				// hostType 추가
				newLine := fmt.Sprintf("      - \"%s:%s\"\n", hostData.HostName, hostData.HostIp)
				ymlContent += newLine
			}
		}

		// 변경된 내용을 yml 파일에 쓰기
		err = os.WriteFile(path+"/extra_hosts/"+targetContainer.ContainerName+".yml", []byte(ymlContent), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
