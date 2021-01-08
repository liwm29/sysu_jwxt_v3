# SYSU_JWXT
![](https://img.shields.io/badge/sysu_jwxt-v3.0.1-519dd9.svg) ![](https://img.shields.io/badge/language-Golang-blue.svg) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)   
:rocket: version 3 of [sysu_jwxt_v2](https://github.com/liwm29/sysu_jwxt_v2) 
## TODO
从v2到现在,学到了更多的技术,因此打算升级v2,要做的事情:
- 前后端分离,后端仅作为api服务器
  - 可能涉及跨域的问题
- 前端界面重写,打算基于gva二次开发一下
  - 要从头开始写,还是太难了,水平有限只能写写组件
- 单例模式修改为支持多用户登陆
  - 加入cookie/session,支持多客户端
- 后端重构一下代码,整合一下基于sysu.edu.cn/jwxt登陆和基于portal.sysu.edu.cn的webvpn登陆
  - 此前由于jwxt对外网开发,不再需要通过portal跳转webvpn登陆,便丢弃了portal
- 加入mock用户,用于测试,不必真的登陆jwxt
- 教师照片的加载问题,应该改为在可视区域时自动加载,而不是hover()时

## Comment
当然,最重要的是前后端分离和重构后端代码,现在的代码都是当初便学便写的时候的遗留代码,大部分又臭又长

## ChangeLog
2021/01/08 初始化任务目标, 计划考试后开始工作