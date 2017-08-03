<html>
{{range $k, $v := .accounts}}
{{$v.Id}}
{{$v.Name}}
{{$v.Ip}}
{{$v.Status}}
{{$v.Priv}}
{{$v.Sync}}
{{$v.Mapping}}
<br/>
{{end}}
</html>