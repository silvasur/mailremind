{{define "title"}}{{if .JobID}}Edit Job{{else}}Create Job{{end}}{{end}}

{{define "content"}}
	{{if .Error}}
		<div class="error">{{.Error}}</div>
	{{end}}
	
	{{if not .Fatal}}
		{{if .Success}}
			<div class="success">{{.Success}}</div>
		{{end}}
		
		<form action="/jobedit{{if .JobID}}/{{.JobID}}{{end}}" method="post" accept-charset="UTF-8">
			<h2>Mail</h2>
			<p>
				<strong>Subject:</strong><br />
				<input type="text" name="Subject" value="{{.Subject}}" />
			</p>
			<p>
				<strong>Content:</strong><br />
				<textarea name="Content" cols="80" rows="20">{{.Content}}</textarea>
			</p>
			
			<h2>Schedule</h2>
			<table class="fullwidth schedules">
				<thead>
					<tr>
						<th>Start<br /><span class="hint">(Format: YYYY-MM-DD HH:MM:SS)</span></th>
						<th colspan="2">Repetition</th>
						<th colspan="2">End<br /><span class="hint">(Format: YYYY-MM-DD HH:MM:SS)</span></th>
					</tr>
				</thead>
				<tbody>
					{{range .Schedules}}<tr>
						<td><input type="text" name="Start" value="{{.Start}}" /></td>
						<td>
							<select name="RepetitionEnabled" size="0" class="enabler">
								<option value="no"{{if not .RepetitionEnabled}} selected="selected"{{end}}>Off</option>
								<option value="yes"{{if .RepetitionEnabled}} selected="selected"{{end}}>On</option>
							</select>
						</td>
						<td>
							<input type="text" name="Count" value="{{.Count}}" class="quant" />
							<select name="Unit" size="0">
								<option value="Minute"{{if .UnitIsMinute}} selected="selected"{{end}}>Minute(s)</option>
								<option value="Hour"{{if .UnitIsHour}} selected="selected"{{end}}>Hour(s)</option>
								<option value="Day"{{if .UnitIsDay}} selected="selected"{{end}}>Day(s)</option>
								<option value="Week"{{if .UnitIsWeek}} selected="selected"{{end}}>Week(s)</option>
								<option value="Month"{{if .UnitIsMonth}} selected="selected"{{end}}>Month(s)</option>
								<option value="Year"{{if .UnitIsYear}} selected="selected"{{end}}>Year(s)</option>
							</select>
						</td>
						<td>
							<select name="EndEnabled" size="0" class="enabler">
								<option value="no"{{if not .EndEnabled}} selected="selected"{{end}}>Off</option>
								<option value="yes"{{if .EndEnabled}} selected="selected"{{end}}>On</option>
							</select>
						</td>
						<td><input type="text" name="End" value="{{.End}}" /></td>
					</tr>{{end}}
				</tbody>
			</table>
			
			<h2>Save</h2>
			<input type="submit" />
		</form>
	{{end}}
{{end}}