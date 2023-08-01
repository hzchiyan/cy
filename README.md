# cy
~~~
go install github.com/hzchiyan/cy@main
~~~

#生成ssl
~~~
cy deploy ssl --domain=github.hzchiyan.com --path=/data/wwwroot/hzchiyan
~~~
# ssl nginx 
~~~
location /.well-known {root /data/wwwroot/hzchiyan;}
~~~
