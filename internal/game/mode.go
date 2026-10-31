package game

type Mode struct{
	id 			Identity
	prompt 		string
	state 		map[Identity]string
	script 		CompiledScript
}

func NewMode(id Identity, prompt string, state map[Identity]string, script CompiledScript) Mode {
	return Mode{
		id: id,
		prompt: prompt,
		state: state,
		script: script,
	}
}
