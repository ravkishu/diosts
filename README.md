# securitytxt-check
Takes a list of domains as the input, checks if they have a security.txt, outputs the results.

# build
```
git clone https://github.com/hakluke/securitytxt-check
cd securitytxt-check
go build *.go
```

# usage
```
cat domains.txt | ./main
```