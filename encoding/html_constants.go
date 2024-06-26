package encoding

const resourceHtml = `{{define "resource"}}
<table>
	{{range $dataKey, $dataValue := .Values}}
	<tr>
		<td>{{$dataKey}}</td><td>{{FormatValue $dataValue}}</td>
	</tr>
	{{end}}
</table>
{{if .Embedded }}
	<h3>Embedded</h3>
	<table>
		{{range $embeddedName, $resources := .Embedded}}
		<tr>
			<td>{{$embeddedName}}</td>
			<td>
			{{range $resource := GetEmbeddedList $resources}}
				<h2>{{$resource.Schema}}</h2>
				{{template "resource" $resource}}
			{{end}}
			</td>
		</tr>
		{{end}}
	</table>
{{end}}
{{if .Links }}
	<h3>Links</h3>
	<table>
		{{range $linkName, $link := .Links}}
		<tr>
			<td>{{$linkName}}</td>
			{{if or (ne $link.Verb "GET") $link.Parameters}}
				<td>
					<form action={{$link.Href}} {{if eq $link.Verb "GET"}} method="GET" {{else}} method="POST" {{end}}>
						{{if and (ne $link.Verb "GET") (ne $link.Verb "POST")}}
							<input type="hidden" name="_method" value="{{$link.Verb}}"></input>
						{{end}}

						{{range $parameter := $link.Parameters}}
							{{ if $parameter.ListOfValues }}
								<select name="{{$parameter.Name}}" placeholder="{{$parameter.Name}}" value="{{$parameter.DefaultValue}}">
									{{ range $value := SeparateListOfValues $parameter.ListOfValues }}
										<option value="$value" {{ if eq $value $parameter.DefaultValue }} selected="selected" {{ end }}>
											{{ $value }}
										</option>
									{{ end }}
								</select>
							{{ else }}
								<input name="{{$parameter.Name}}" placeholder="{{$parameter.Name}}"	value="{{$parameter.DefaultValue}}"></input>
							{{ end }}
							<br>
						{{end}}

						<input type="submit" class="btn" value="{{$link.Verb}}"></input>	
					</form>
			   </td>
			{{else}}
			  <td>
				<a id="{{$linkName}}" href="{{$link.Href}}">{{$link.Href}}</a>
				{{ range $templatedParameter := GetTemplatedParameters $link }}
					<br>
                    <input id="{{$linkName}}_{{$templatedParameter}}" placeholder="{{$templatedParameter}}" oninput="OnUpdateTemplatedUrl({{$linkName}}, {{$link.Href}}, {{ GetTemplatedParameters $link }})"></input>
                {{ end }}
			  </td>
			{{end}}
		</tr>
		{{end}}
	</table>
{{end}}
{{end}}<!DOCTYPE html>
<html lang="en">
<head>
    <title>{{.Schema}}</title>
	<style>
		body{
		    margin: 0;
    		padding: 0;
    		color: #333;
    		background-color: #eee;
    		font: 1em/1.2 'Helvetica Neue', Helvetica, Arial, Geneva, sans-serif;
		}

		h1,h2,h3 {
    		margin: 0 0 .5em;
    		font-weight: 500;
    		line-height: 1.1;
		}

		h1 { font-size: 2.25em; }
		h2 { font-size: 1.375em; background: #555; color: white; padding: 0.25em}
		h3 { font-size: 1.375em; background: lightgrey; padding: 0.25em }
		
		p {
			margin: 0 0 1.5em;
			line-height: 1.5;
		}

		table {
			background-color: transparent;
			border-spacing: 0;
			border-collapse: collapse;
			border-top: 1px solid #ddd;
			width: 100%
		}
		
		th, td {
			padding: .5em 1em;
			vertical-align: center;
			text-align: left;
			border-bottom: 1px solid #ddd;
		}
		
		td:last-child {
			width:100%;
		}
		
		
		.btn {
			color: #fff !important;
			background-color: GRAY;
			border-color: #222;
			display: inline-block;
			padding: .5em 1em;
			font-weight: 400;
			line-height: 1.2;
			text-align: center;
			white-space: nowrap;
			vertical-align: middle;
			cursor: pointer;
			border: 1px solid transparent;
			border-radius: .2em;
			text-decoration: none;
		}
		
		.btn:hover {
			color: #fff !important;
			background-color: #555;
		}
		
		.btn:focus {
			color: #fff !important;
			background-color: #555;
		}
		
		.btn:active {
			color: #fff !important;
			background-color: #555;
		}
		
		.container {
			max-width: 70em;
			margin: 0 auto;
			background-color: #fff;
		}
		
		.header {
			color: #fff;
			background: #555;
			p
		}
		
		.subheader {
			color: #fff;
			background: #555;
			p
		}
		
		.header-heading { margin: 0; }
		
		.content { padding: 1em 1.25em; }
		
		.embedded { padding-left: 1.5em }
		
		@media (min-width: 42em) {
			.header { padding: 1.5em 3em; }
			.subheader { padding: .2em 3em; }
			.content { padding: 2em 3em; }
		}
	</style>
    <script>
		function OnUpdateTemplatedUrl (linkId, originalHref, parameters) {
			var href = originalHref;
			for (x in parameters) {
				parameterElement = document.getElementById(linkId + "_" + parameters[x]);
                value = parameterElement.value;
				if (value === '') {
					continue; 
				}
				href = href.replace(parameters[x], parameterElement.value);
			}
			linkElement = document.getElementById(linkId) 
			linkElement.setAttribute('href', href)
			linkElement.innerHTML = href
		}
	</script>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1 class="header-heading">{{.Schema}}</h1>
		</div>
		<div class="content">
		{{template "resource" .}}
		</div>
    </div>
</body>
</HTML>`
