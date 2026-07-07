sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
