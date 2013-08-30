{{define "title"}}Reset Password{{end}}

{{define "content"}}
	{{if .Success}}
		<div class="success">{{.Success}}</div>
	{{else}}
		{{if .Error}}
			<div class="error">{{.Error}}</div>
		{{end}}
		
		<form action="/forgotpw" method="post" accept-charset="UTF-8">
			<p><strong>E-Mail:</strong> <input type="text" name="Mail" /></p>
			<p><input type="submit" /></p>
		</form>
	{{end}}
{{end}}