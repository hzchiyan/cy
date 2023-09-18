# cy
~~~
go install github.com/hzchiyan/cy@main
~~~

# 生成ssl
~~~
cy deploy website-ssl --host=github.hzchiyan.com --dir=/data/wwwroot/hzchiyan
~~~

# 创建go站点
~~~
cy deploy  website-go --host=api.hzchiyan.com --port=8084 --dir=/data/wwwroot/api

cy deploy website-ssl --host=api.hzchiyan.com  --dir=/data/wwwroot/api
~~~



