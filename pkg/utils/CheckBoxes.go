package utils

import (
	"dockerAddHost/types"
	"github.com/AlecAivazis/survey/v2"
)

func Checkboxes(label string, opts []types.ContainerData) []string {
	res := []string{}

	var names []string
	for _, opt := range opts {
		names = append(names, opt.ContainerName)
	}
	prompt := &survey.MultiSelect{
		Message: label,
		Options: names, //여기서 opts의 contanerName만 가져올 수 있는 방법,
	}
	survey.AskOne(prompt, &res)

	return res
}
