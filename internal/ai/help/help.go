package help

func GetHelp() string {
    return `Usage: ai [options] 
Options:
    --system <prompt>   Set system prompt
    --model <name>      Specify the model to use (default: llama2)
    --image <path>      Path to image file for image-based prompts
    --list             List available models
    --help             Show this help message

After running the command, enter your prompt when prompted.`
}

func GetErrorNoPrompt() string {
    return "Error: No prompt provided. Please enter a prompt when prompted."
}

