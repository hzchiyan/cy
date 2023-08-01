# cy
~~~
go install github.com/hzchiyan/cy@main
~~~

#生成ssl
~~~
cy deploy website-ssl --host=github.hzchiyan.com --dir=/data/wwwroot/hzchiyan
~~~

#创建go站点
~~~
cy deploy website-go --host=github.hzchiyan.com --dir=/data/wwwroot/hzchiyan --port=8080
~~~