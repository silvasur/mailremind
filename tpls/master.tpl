<html>
<head>
	<title>{{template "title" .Data}} – mailremind</title>
	<link rel="stylesheet" type="text/css" href="/static/style.css" />
	<script type="text/javascript" src="/static/jquery-1.10.2.min.js"></script>
	<script type="text/javascript" src="/static/mailremind.js"></script>
</head>
<body>
	<div id="main">
		<div id="nav">
			<a href="/" class="apptitle">mailremind</a>
			<ul>
				{{if .Mail}}
					<li><a href="/jobedit">new job</a></li>
					<li><a href="/jobs">list jobs</a></li>
					<li><a href="/settings">settings</a></li>
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
{{/* START EDITING FOOTER HERE */}}
			<strong>PLEASE EDIT THIS FOOTER TO CONTAIN SOME CONTACT DATA. YOU CAN DO SO BY EDITING THE tpls/master.tpl TEMPLATE.</strong>
{{/* STOP EDITING FOOTER HERE */}}
		</div>
	</div>
</body>
</html>