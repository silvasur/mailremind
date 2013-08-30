{{define "title"}}Login{{end}}

{{define "content"}}
	{{if .Error}}
		<div class="error">{{.Error}}</div>
	{{end}}
	
	{{if .Success}}
		<div class="success">{{.Success}}</div>
	{{else}}
		<form action="/login" method="post" accept-charset="UTF-8">
			<p><strong>E-Mail</strong> <input type="text" name="Mail" /></p>
			<p><strong>Password</strong> <input type="password" name="Password" /></p>
			<p><input type="submit" /></p>
			<p><a href="/forgotpw">Forgot your password?</a></p>
		</form>
	{{end}}
{{end}}