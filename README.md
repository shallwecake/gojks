### 使用gojks命令来控制jenkins构建应用，使用gojks help 来查看使用命令

##### 1. 添加jenkins配置
gojks add jks http://localhost:8500 admin:admin
##### 2. 添加rancher配置
gojks add rcr http://localhost:443 admin:admin
##### 3. 添加飞书机器人配置
gojks add whk https://open.feishu.cn/open-apis/bot/v2/hook/test
##### 4. 构建应用时会模糊查询构建名称，使用,隔开要构建应用的序号
gojks pub app
##### 5. 构建多个应用，需要输入全称，使用,隔开
gojks pubs app1,app2

