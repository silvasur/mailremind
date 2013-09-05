{{define "title"}}Jobs{{end}}

{{define "content"}}
	{{if .Error}}
		<div class="error">{{.Error}}</div>
	{{end}}
	
	{{if not .Fatal}}
		{{if .Success}}
			<div class="success">{{.Success}}</div>
		{{end}}

		<form action="/jobs" method="post">
			<table class="fullwidth">
				<thead>
					<tr>
						<th>&nbsp;</th>
						<th>Subject</th>
						<th>Excerpt</th>
						<th>Next</th>
					</tr>
				</thead>
				<tbody>
					{{range .Jobs}}<tr>
						<td><input type="checkbox" name="Jobs" value="{{.ID}}" /></td>
						<td><a href="/jobedit/{{.ID}}">{{.Subject}}</a></td>
						<td>{{.Excerpt}}</td>
						<td>{{.Next}}</td>
					</tr>
					{{else}}<tr>
						<td colspan="4" class="emptytab">No jobs found</td>
					</tr>{{end}}
				</tbody>
			</table>
			<p>
				Delete selected:
				<select name="Delconfirm" size="0">
					<option value="no" selected="selected">No</option>
					<option value="yes">Yes</option>
				</select>
				<input type="submit" />
			</p>
		</form>
	{{end}}
{{end}}