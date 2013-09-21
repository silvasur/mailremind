{{define "title"}}User Settings{{end}}

{{define "content"}}
	{{if .Error}}<div class="error">{{.Error}}</div>{{end}}
	
	{{if not .Fatal}}
		{{if .Success}}<div class="success">{{.Success}}</div>{{end}}
		
		<form action="/settings?M=setpasswd" method="post" accept-charset="UTF-8" class="fancy">
			<h2>Set Password</h2>
			
			<p><strong>Password:</strong> <input type="password" name="Password" /></p>
			<p><strong>Repeat Password:</strong> <input type="password" name="RepeatPassword" /></p>
			<p><input type="submit" /></p>
		</form>
		
		<form action="/settings?M=settimezone" method="post" accept-charset="UTF-8" class="fancy">
			<h2>Set Timezone</h2>
			
			<p>
				<strong>Timezone:</strong>
				<select size="0" name="Timezone">
					{{range $tz, $active := .Timezones}}<option value="{{$tz}}"{{if $active}} selected="selected"{{end}}>{{$tz}}</option>{{end}}
				</select>
			</p>
			<p><input type="submit" /></p>
		</form>
		
		<form action="/delete-acc" method="get" class="fancy">
			<h2>Delete Account</h2>
			<p><input type="submit" value="Delete Account" /></p>
		</form>
	{{end}}
{{end}}