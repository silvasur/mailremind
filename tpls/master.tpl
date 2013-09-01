<html>
<head>
	<title>{{template "title" .Data}} – mailremind</title>
	<link rel="stylesheet" type="text/css" href="/static/style.css" />
</head>
<body>
	<div id="main">
		<div id="nav">
			<a href="/" class="apptitle">mailremind</a>
			<ul>
				{{if .Mail}}
					<li><a href="/jobedit">new job</a></li>
					<li><a href="/jobs">list jobs</a></li>
					<li><a href="/logout">logout</a></li>
				{{else}}
					<li><a href="/register">register</a></li>
					<li><a href="/login">login</a></li>
				{{end}}
			</ul>
		</div>
		
		<div id="content">
			<h1>{{template "title" .Data}}</h1>
			{{template "content" .Data}}
		</div>
		
		<div id="footer">
			© 2013 Bla bla....<br />
			Foo!
		</div>
	</div>
</body>
</html>