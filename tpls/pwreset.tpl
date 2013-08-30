{{define "title"}}Reset Password{{end}}

{{define "content"}}
	{{if .Success}}
		<div class="success">{{.Success}}</div>
	{{else}}
		{{if .Error}}
			<div class="error">{{.Error}}</div>
		{{end}}
		
		<form action="/pwreset?Code={{.Code}}&amp;U={{.UID}}" method="post" accept-charset="UTF-8">
			<p><strong>Password:</strong> <input type="password" name="Password" /></p>
			<p><strong>Retype Password:</strong> <input type="password" name="PasswordAgain" /></p>
			<p><input type="submit" /></p>
		</form>
	{{end}}
{{end}}