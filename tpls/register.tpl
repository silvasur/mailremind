{{define "title"}}Register{{end}}

{{define "content"}}
	{{if .Success}}
		<div class="success">{{.Success}}</div>
	{{else}}
		{{if .Error}}<div class="error">{{.Error}}</div>{{end}}
		<form action="/register" method="post" accept-charset="UTF-8">
			<p><strong>E-Mail:</strong> <input type="text" name="Mail" /></p>
			<p><strong>Password:</strong> <input type="password" name="Password" /></p>
			<p><strong>Retype Password:</strong> <input type="password" name="RetypePassword" /></p>
			<p>
				<strong>Timezone:</strong>
				<select size="0" name="Timezone">
					{{range .Timezones}}<option value="{{.}}">{{.}}</option>{{end}}
				</select>
			</p>
			<p><input type="submit" /></p>
		</form>
	{{end}}
{{end}}