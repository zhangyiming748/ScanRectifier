sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list.d/debian.sources
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
apt-get update && apt-get install -y \
    libopencv-dev \
    build-essential \
    cmake \
    pkg-config
pkg-config --modversion opencv4
apt-get install -y gcc g++ make cmake pkg-config build-essential
apt-get install -y libopencv-dev
export CGO_ENABLED=1
go env -w CGO_ENABLED=1
