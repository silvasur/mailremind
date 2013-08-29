{{define "title"}}{{.Title}}{{end}}

{{define "content"}}
	{{if .Class}}
		<div class="{{.Class}}">{{.Msg}}</div>
	{{else}}
		<div>{{.Msg}}</div>
	{{end}}
{{end}}