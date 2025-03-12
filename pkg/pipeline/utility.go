package pipeline

import (
	"github.com/StoneG24/slape/cmd/slape/prompt"
	"github.com/fatih/color"
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/gpu"
)

func processPrompt(mode string) (string, int64) {
	var promptChoice string
	var maxtokens int64

	switch mode {
	case "simple":
		promptChoice = prompt.SimplePrompt
		maxtokens = 1024
	case "cot":
		promptChoice = prompt.CoTPrompt
		maxtokens = 4096
	case "tot":
		promptChoice = prompt.ToTPrompt
		maxtokens = 32768
	case "got":
		promptChoice = prompt.GoTPrompt
		maxtokens = 32768
	case "moe":
		promptChoice = prompt.MoEPrompt
		maxtokens = 32768
	case "thinkinghats":
		promptChoice = prompt.SixThinkingHats
		maxtokens = 32768
	case "goe":
		promptChoice = prompt.GoEPrompt
		maxtokens = 32768
	default:
		promptChoice = prompt.SimplePrompt
		maxtokens = 100
	}

	return promptChoice, maxtokens
}

func IsGPU() bool {
	gpuInfo, err := ghw.GPU()
	// if there is an error continue without using a GPU
	if err != nil {
		color.Red("%s", err)
		color.Yellow("Continuing without GPU...")
	}

	// for debugging
	color.Yellow("%s", gpuInfo.GraphicsCards)

	if len(gpuInfo.GraphicsCards) == 0 {
		color.Yellow("No GPUs to use, switching to cpu only")
		return false
	}

	// This guy nil derefernce panics when the gpu isn't actually a graphics card
	// fmt.Println(gpuInfo.GraphicsCards[0].Node.Memory)

	return false
}

// CheckMemoryUsage is used to check the availble memory of a machine.
func GetAmountofMemory() (int64, error) {
	memory, err := ghw.Memory()
	if err != nil {
		return 0, err
	}
	return memory.TotalUsableBytes, nil
}

// If thread count is low we probably should run only a small amount of models
func GetNumThreads() (uint32, error) {
	cpu, err := ghw.CPU()
	if err != nil {
		return 0, err
	}
	return cpu.TotalThreads, nil
}

// Gather the gpus on the system a select their vendors.
// Normally one would only have two but this interface assumes
// that the internal graphics module on the CPU is a graphics card.
func GatherGPUs() ([]*gpu.GraphicsCard, error) {
	gpuInfo, err := ghw.GPU()
	if err != nil {
		return nil, err
	}

	return gpuInfo.GraphicsCards, nil
}
