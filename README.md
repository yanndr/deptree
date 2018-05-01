#Deptree 

This library produces a fully resolve dependency tree for Perl distributions. 
A command line client is available on the cmd\deptree folder.

## Installing

As a private repository go get won't work :(
 ```
 got get bitbucket.org/yanndr/deptree 
 ```

If you have access to the repo, you can clone the repository in the folder $GOPATH/src/bitbucket.org/yanndr/

```
git clone git clone https://yanndr@bitbucket.org/yanndr/deptree.git $GOPATH/src/bitbucket.org/yanndr/
```


## Using

### Docker


```
docker run --rm -it  deptree -name  DateTime 
```

```
docker run --rm -it -v $PWD/cmd/deptree/data:/mydata  deptree -name  DateTime -path ./mydata
```