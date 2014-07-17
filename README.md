# KellyBackend

[![Build Status](https://secure.travis-ci.org/missdeer/KellyBackend.png)](https://travis-ci.org/missdeer/KellyBackend)

Yiili community which is a part of [Kelly project](https://github.com/missdeer/kelly), it is the server side, AKA backend, of the project. The source code is based on WeTalk project, thanks those guys for there great work.

### Usage

```
go get -u github.com/missdeer/KellyBackend
cd $GOPATH/src/github.com/missdeer/KellyBackend
```

I suggest you [update all Dependencies](#dependencies)

Copy `conf/global/app.ini` to `conf/app.ini` and edit it. All configure has comment in it.

The files in `conf/` can overwrite `conf/global/` in runtime.


**Run KellyBackend**

```
bee run watchall
```

### Dependencies

Contrib

* Beego [https://github.com/astaxie/beego](https://github.com/astaxie/beego) 
* Social-Auth [https://github.com/beego/social-auth](https://github.com/beego/social-auth)
* Compress [https://github.com/beego/compress](https://github.com/beego/compress)
* i18n [https://github.com/beego/i18n](https://github.com/beego/i18n)
* Mysql [https://github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)
* goconfig [https://github.com/Unknwon/goconfig](https://github.com/Unknwon/goconfig)
* fsnotify [https://github.com/howeyc/fsnotify](https://github.com/howeyc/fsnotify)
* resize [https://github.com/nfnt/resize](https://github.com/nfnt/resize)
* blackfriday [https://github.com/slene/blackfriday](https://github.com/slene/blackfriday)

Update all Dependencies

```
go get -u github.com/beego/social-auth
go get -u github.com/beego/compress
go get -u github.com/beego/i18n
go get -u github.com/go-sql-driver/mysql
go get -u github.com/Unknwon/goconfig
go get -u github.com/howeyc/fsnotify
go get -u github.com/nfnt/resize
go get -u github.com/slene/blackfriday
```

### Static Files

KellyBackend use `Google Closure Compile` and `Yui Compressor` compress js and css files.

So you could need Java Runtime. Or close this feature in code by yourself.

### Contact

Maintain by [missdeer](http://minidump.info/)

## License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).
