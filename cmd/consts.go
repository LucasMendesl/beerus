package cmd

const (
	shortDescription = `
 /$$$$$$$
| $$__  $$
| $$  \ $$  /$$$$$$   /$$$$$$   /$$$$$$  /$$   /$$  /$$$$$$$
| $$$$$$$  /$$__  $$ /$$__  $$ /$$__  $$| $$  | $$ /$$_____/
| $$__  $$| $$$$$$$$| $$$$$$$$| $$  \__/| $$  | $$|  $$$$$$
| $$  \ $$| $$_____/| $$_____/| $$      | $$  | $$ \____  $$
| $$$$$$$/|  $$$$$$$|  $$$$$$$| $$      |  $$$$$$/ /$$$$$$$/
|_______/  \_______/ \_______/|__/       \______/ |_______/

Clean up your docker workspace, removing useless resources
`

	longDescription = `
Beerus is a command-line tool that automatically cleans up unused Docker resources,
such as containers and images, based on customizable labels.

It helps maintain a lean Docker environment by identifying and removing orphaned
or unnecessary resources without manual intervention.
`

	helpTemplate = `{{.Short}}

Description:
  {{.Long}}

{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}
`
)
