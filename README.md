#Deptree 

This library produces a fully resolved dependency tree for Perl distributions. 
A command line client is available on the cmd\deptree folder.

## Installing

1. Clone the repository in the folder $GOPATH/src/bitbucket.org/yanndr/ `git clone git clone https://yanndr@bitbucket.org/yanndr/deptree.git $GOPATH/src/bitbucket.org/yanndr/` Note: as this is a private repository, `go get` won't work.
2. Build/install the application :`make install` or `go build ./cmd/deptree`

## Usage

This program needs access to a CPAN folder to run successfully; one is provided here ./cmd/deptree/data. 
You can define the path of the CPAN folder with the flag -path. By default the path is set to ./data. If you run the program directely from the ./cmd/deptree folder, you won't have to define the -path flag.

```
deptree -path path/to/cpan -name DateTime -name Specio
OPTIONS:
  -name value
        Distribition name to resolve; this flag is mandatory. You need to define it once but you can also define it multiple times.
  -path string
        The path to the CPAN folder. (default "./data")
```

Example:
```
deptree -name DateTime
```



## Docker
You can aslo use the program with Docker. 
First you'll have to make the docker image:
```
make docker
```

Once the docker image is created, you can run the program with the command:
```
docker run --rm -it  deptree -name  DateTime 
```
The docker image embeds a data folder on the root folder, so you don't have to specify the -path flag.

If you want to use a different path, use the following syntax:
```
docker run --rm -it -v $PWD/cmd/deptree/data:/mydata  deptree -name  DateTime -path ./mydata
```