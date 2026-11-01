package game

type Mode struct{
	id 			Identity
	prompt 		string
	states 		map[Identity]string
	script 		CompiledScript
}

func NewMode(id Identity, prompt string, state map[Identity]string, script CompiledScript) Mode {
	return Mode{
		id: id,
		prompt: prompt,
		states: state,
		script: script,
	}
}

func (o Mode) state(id Identity) (string, bool) {
	result, ok := o.states[id]
	return result, ok
}

func (o Mode) hasState(id Identity) bool {
	_, ok := o.state(id)
	return ok
}
