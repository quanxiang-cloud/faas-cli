# Run
Provide local running function tools.

# code package
Depends on the faas-lowcode project, and builds the basic code through this project.
```bash=
git clone https://github.com/quanxiang-cloud/faas-lowcode.git
cd faas-lowcode
tar -zcf lowcode-go116.tar.gz ./pkg/local_proxy \
./pkg/util.go \
./go.mod
```

local plugins
```bash=
tar -zcf plugin-quanxiang-lowcode-client.tar.gz ./plugin-quanxiang-lowcode-client
```

# hand drawn environment
```bash=
mkdir -p ~/.quanxiang/faas/go116
mkdir -p ~/.quanxiang/faas/go116/pkg/faas-lowcode
tar -zxvf lowcode-go116.tar.gz -C ~/.quanxiang/faas/go116/pkg/faas-lowcode

mkdir -p ~/.quanxiang/faas/go116/plugins
tar -zxvf plugin-quanxiang-lowcode-client.tar.gz -C ~/.quanxiang/faas/go116/plugins
```