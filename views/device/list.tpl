<html>
{{range $k, $v := .devices}}
{{$v.Id}}
{{$v.Name}}
{{$v.Ip}}
{{$v.Dtype}}
{{$v.Status}}
{{$v.Duser}}
<br/>
{{end}}
</html>