{{define "title"}}Delete Account{{end}}

{{define "content"}}
	{{if .OK}}
		<a href="/delete-acc/yes">Click here if you really want to delete your account</a>
	{{else}}
		<div class="error">
			You need to be logged in to do that
		</div>
	{{end}}
{{end}}