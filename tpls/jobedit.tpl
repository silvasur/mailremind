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
			{{range .Schedules}}
				<p>
					<strong>Start:</strong>
						<input type="text" name="Start" value="{{.Start}}" /><br />
					<input type="checkbox" name="RepetitionEnabled"{{if .RepetitionEnabled}} checked="checked"{{end}} />
					<strong>Repetition:</strong>
						<input type="text" name="Count" value="{{.Count}}" />
						<select name="Unit" size="0">
							<option value="Minute"{{if .UnitIsMinute}} selected="selected"{{end}}>Minute(s)</option>
							<option value="Hour"{{if .UnitIsHour}} selected="selected"{{end}}>Hour(s)</option>
							<option value="Day"{{if .UnitIsDay}} selected="selected"{{end}}>Day(s)</option>
							<option value="Week"{{if .UnitIsWeek}} selected="selected"{{end}}>Week(s)</option>
							<option value="Month"{{if .UnitIsMonth}} selected="selected"{{end}}>Month(s)</option>
							<option value="Year"{{if .UnitIsYear}} selected="selected"{{end}}>Year(s)</option>
						</select>
						<br />
					<input type="checkbox" name="EndEnabled"{{if .EndEnabled}} checked="checked"{{end}} />
					<strong>End:</strong>
						<input type="text" name="End" value="{{.End}}" />
				</p>
			{{end}}
			
			<h2>Save</h2>
			<input type="submit" />
		</form>
	{{end}}
{{end}}