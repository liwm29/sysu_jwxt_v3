# SYSU_JWXT
![](https://img.shields.io/badge/sysu_jwxt-v3.0.1-519dd9.svg) ![](https://img.shields.io/badge/language-Golang-blue.svg) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)   
:rocket: version 3 of [sysu_jwxt_v2](https://github.com/liwm29/sysu_jwxt_v2) 
## TODO
从v2到现在,学到了更多的技术,因此打算升级v2,要做的事情:
- 前后端分离,后端仅作为api服务器
  - 可能涉及跨域的问题
- 前端界面重写,打算基于vue-el-admin二次开发照猫画虎一下
- 单例模式修改为支持多用户登陆
  - 加入cookie/session,支持多客户端
- 后端重构一下代码,整合一下基于sysu.edu.cn/jwxt登陆和基于portal.sysu.edu.cn的webvpn登陆
  - 此前由于jwxt对外网开发,不再需要通过portal跳转webvpn登陆,便丢弃了portal
- 加入mock用户,用于测试,不必真的登陆jwxt
- 教师照片的加载问题,应该改为在可视区域时自动加载,而不是hover()时

## 项目目录结构
- proxyServer作为一个代理服务器,单纯接收来自前端的请求,调用jwxtClient的方法获得数据,然后返回给前端
```sh
.
├── backend
│   ├── jwxtClient
│   └── proxyServer
└── frontend
     └── dist
```

## 待决
1. 自己封装的httpClient,到底是应该嵌入http.Client还是组合?
   1. 目前是用组合,否则的话无法定义同名函数
2. Do()的返回值是用httpResp还是直接[]byte?
3. 某些地方应该使用接口吗? 比如req.Do(interface?),否则得传值c.HttpClient
4. 测试在test文件夹,应该用标准的go test,但是似乎go test无法得到用户输入,猜测是go test不用stdin,stdout,有时间再看看,目前直接改成main文件运行

## ChangeLog
- 2021/01/08 初始化任务目标, 计划考试后开始工作
- 2021/01/09 决定前后端分离的模式为:分开开发,合并部署,见 DevLog#1 ,添加了部署代码
- 2020/01/10 确定项目目录结构
- 2020/01/14 增加了request组件,为对http.request/response/client的简单封装
- 2020/01/15 增加了cookie管理,迁移了登陆实现,已测试

## DevLog
1. 前后端分离,肯定要分离开发,至于是否分离部署,看个人需要
   1. 如果分离部署,这是在说前端代码`npm run build`后,将`/dist`目录直接扔进nginx或tomcat,后端作为api服务器单独运行在另一个端口
      1. 由于端口不同,涉及CORS跨域资源共享问题,对xhr请求的发出没影响,主要是响应必须带有`Access-Control-Allow-Origin`,否则被浏览器拦截;dom的请求似乎直接禁止了,防止冒牌网站直接套壳iframe;具体如何,没试过
   2. 如果一起部署,也就是虽然后端服务器是作为api服务器,但是当请求`'/'`时,便返回`html`,其余的路由都是api
      1. 这在go中很容易实现,但其实不算太优雅,毕竟api服务器多了几条ServeFile代码,动态路由的html(指访问`/`而不是`/index.html`)和其他静态文件都由api服务器响应
