<!DOCTYPE html>

<html>
<head>
  <title>pldapi help</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <link href="/static/css/bootstrap.min.css" rel="stylesheet">

  <style>
	html {
	  position: relative;
	  min-height: 100%;
	}
	
    body {padding-top: 60px;margin-bottom: 50px;}
	
	.footer {
	  position: absolute;
	  bottom: 0;
	  width: 100%;
	  height: 50px;
	  background-color: #000000;
	  text-align: center;
	}
	
	.footer .text-muted {
    margin: 10px 0;
}
  </style>
 </head>

<body>
  <nav class="navbar navbar-inverse navbar-fixed-top">
    <div class="container">
	  <div class="navbar-header">
        <a class="navbar-brand" href="#">pldapi {{.Version}}</a>
      </div>
	
	<div id="navbar" class="navbar-collapse collapse">
		<ul class="nav navbar-nav">
		    <li class="active"><a href="#">Home</a></li>
            <li><a href="/pldapi/about">About</a></li>
		</ul>
	</div>
	
	</div>
  </nav>

  <div class="container">

    <div class="page-header">
	  <h1>pldapi {{.Version}}</h1>
	</div>
	
	<h2>一、列表类</h2>
	<h3>1、设备列表</h3>
	<p><span class="label label-primary">get</span> http://{{.Ip}}:21080/v1/device/list</p>
	<h3>2、账号列表</h3>
    <p><span class="label label-primary">get</span> http://{{.Ip}}:21080/v1/account/list</p>


    <h2>二、设备类</h2>
    <h3>1、添加设备 <span class="label label-danger"> hot</span></h3>
	<p>该 API 可添加指定的设备，设备类型为 linux，并会自动的为其增加 root,tomcat,log 用户。</p>
    <p><span class="label label-primary">post</span> curl -X POST -d'{"device":"WGQ00001","ip":"222.111.222.111"}' http://{{.Ip}}:21080/v1/device</p>
	<h3>2、设备改名 <span class="label label-danger"> hot</span></h3>
	<p>该 API 用来给指定的 IP 地址改名， 使用 device 指定新的名称。</p>
	<p><span class="label label-primary">put</span> curl -X PUT -d'{"device":"TGT","ip":"222.111.222.111"}' http://{{.Ip}}:21080/v1/device</p>


    <h2>三、账号类</h2>
    <h3>1、添加账号 <span class="label label-danger"> hot</span></h3>
	<p>该 API 可添加某 IP 地址下的不存在的账号（root,tomcat,log）</p>
    <p><span class="label label-primary">post</span> curl -X POST -d'192.168.1.1' http://{{.Ip}}:21080/v1/account/</p>
	<h3>2、账号改密</h3>
	<p>该 API 对指定的账号 id 进行改密（tomcat,log）,如果指定的账号 id 不是 tomcat或log 则报错</p>
	<p><span class="label label-primary">put</span> curl -X PUT http://{{.Ip}}:21080/v1/account/changepassword/3550</p>
	<h3>3、批量改密</h3>
	<p>该 API 对指定的账号（tomcat,log）进行改密</p>
	<p><span class="label label-primary">put</span> curl -X PUT http://{{.Ip}}:21080/v1/account/changeallpassword/tomcat</p>

    <h2>四、mapping 类</h2>
    <h3>1、查看账号下的 mapping</h3>
    <p><span class="label label-primary">get</span> http://{{.Ip}}:21080/v1/mapping/list/3550</p>

    <h3>2、添加指定账号ID下的所有用户的 mapping</h3>
    <p><span class="label label-primary">put</span> curl -X PUT -d'3550' http://{{.Ip}}:21080/v1/mapping/update</p>

    <h3>3、查询 IP 拥有的 mapping 账号</h3>
    <p><span class="label label-primary">get</span> http://{{.Ip}}:21080/v1/mapping/ip/10.40.40.60</p>

    <h3>4、设置 IP 下所有账号的 mapping <span class="label label-danger"> hot</span></h3>
    <p><span class="label label-primary">put</span> curl -X PUT http://{{.Ip}}:21080/v1/mapping/ip/10.40.40.60</p>

    <h3>5、刷新所有的 IP 的 mapping</h3>
    <p><span class="label label-primary">put</span> curl -X PUT http://{{.Ip}}:21080/v1/mapping/update/all</p>


    <h2>五、其它</h2>
    <h3>1、载入所有的设备</h3>
    <p><span class="label label-primary">get</span> http://{{.Ip}}:21080/v1/device/load</p>

    <h3>2、载入所有的账号</h3>
    <p><span class="label label-primary">get</span> http://{{.Ip}}:21080/v1/account/load</p>

    <h3>3、载入所有的组和设备ID映射</h3>
    <p><span class="label label-primary">get</span> http://{{.Ip}}:21080/v1/group/load</p>

    <h3>4、载入设备、账号、组</h3>
    <p><span class="label label-primary">get</span> http://{{.Ip}}:21080/v1/loadall</p>

  </div>

  <footer class="footer">
    
      <p class="text-muted">有问题请联系：{{.Email}}</p>
	
  </footer>
  <script src="/static/js/bootstrap.min.js"></script>
</body>
</html>
