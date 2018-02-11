# 一个简单的权限设计demo(入口权限)
### 思路
主要有三个模型，分别是user、role、permission。
* User 是用户模型，存储用户相关信息。
* role 是 角色模型，存储所有的角色信息。
* Permission 是具体的权限，这里是做的入口权限的限定，通过http请求的方式和请求路径确定是否有访问权限。
User和role、role 和permission 之间都是 多对多的关联关系。

在整个请求响应中，应该有一个服务能提供当前用户的信息以及该用户所能够访问的资源。



### 管理界面
`python3 app.py`
使用 flask-admin 做的简陋crud界面，但是还没有关联关系

