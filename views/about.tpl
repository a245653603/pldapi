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
        <a class="navbar-brand" href="/pldapi">pldapi {{.Version}}</a>
      </div>
	
	<div id="navbar" class="navbar-collapse collapse">
		<ul class="nav navbar-nav">
		    <li><a href="/pldapi">Home</a></li>
            <li class="active"><a href="#">About</a></li>
		</ul>
	</div>
	
	</div>
  </nav>

  <div class="container">
  <div class="page-header">
    <h1>变更记录</h1>
  </div>

  <h2>pldapi v0.0.3(2017.07.11)</h2>
  <h3>新增功能</h3>
  <p>程序界面做出一定的调整</p>
  <p>加入了 批量改密 功能</p>
  <h3>修正</h3>
  <p>解决了密码内包含特殊字符的问题</p>

  <h2>pldapi v0.0.2(2017.07.04)</h2>
  <h3>新增功能</h3>
  <p>加入了对数据库的支持，使用 mysql memory 数据库保存临时数据，每四小时刷新一次</p>
  <p>增加了 device 库</p>
  <p>增加了 user 库</p>
  <p>增加了 group 库</p>
  <p>加入了 设备列表 功能</p>
  <p>加入了 账号列表 功能</p>
  <p>加入了 添加账号 功能</p>
  <p>加入了 mapping类各操作 功能</p>
  <p>加入了 设备、用户、组 信息载入数据库功能</p>

  <h2>pldapi v0.0.1(2017.06.14)</h2>
  <h3>新增功能</h3>
  <p>新建项目</p>
  <p>添加自动识别验证码的能力</p>
  <p>添加了 session 的保持功能，每 20 分钟刷新一次 session</p>
  <p>加入了 添加设备 功能</p>
  <p>加入了 设备改名  功能</p>
  <p></p>


  </div>
  <footer class="footer">
    
      <p class="text-muted">有问题请联系：{{.Email}}</p>
	
  </footer>
  <script src="/static/js/bootstrap.min.js"></script>
</body>
</html>