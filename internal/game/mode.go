package game

type Mode struct{
	id 			Identity
	prompt 		string
	properties 	map[Identity]string
	script 		string
}


func NewMode(id Identity, prompt string, properties map[Identity]string, script string) Mode {
	return Mode{
		id: id,
		prompt: prompt,
		properties: properties,
		script: script,
	}
}
