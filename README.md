# slic

## straight scp
$ time scp -r data/ seb@Nacho2:~/out
1GB-1                                                                                                                                                                     100% 1024MB 286.8MB/s   00:03    
1GB-4                                                                                                                                                                     100% 1024MB 114.0MB/s   00:08    
1GB-5                                                                                                                                                                     100% 1024MB 136.7MB/s   00:07    
1GB-7                                                                                                                                                                     100% 1024MB 220.7MB/s   00:04    
1GB-0                                                                                                                                                                     100% 1024MB 296.4MB/s   00:03    
1GB-6                                                                                                                                                                     100% 1024MB 108.3MB/s   00:09    
1GB-3                                                                                                                                                                     100% 1024MB 113.7MB/s   00:09    
1GB-2                                                                                                                                                                     100% 1024MB 236.8MB/s   00:04    
1GB-9                                                                                                                                                                     100% 1024MB 109.5MB/s   00:09    
1GB-8                                                                                                                                                                     100% 1024MB  82.3MB/s   00:12    

real	1m12.985s
user	0m42.438s
sys	0m17.559s

## slic
$time go run slic.go -src data -host seb@Nacho2 -dest ~/out
sending data -> /home/seb/out
executing: ssh seb@Nacho2 "mkdir /home/seb/out"
mkdir: cannot create directory ‘/home/seb/out’: File exists

executing: ssh seb@Nacho2 "mkdir /home/seb/out/a"
mkdir: cannot create directory ‘/home/seb/out/a’: File exists

executing: ssh seb@Nacho2 "mkdir /home/seb/out/a/z"
mkdir: cannot create directory ‘/home/seb/out/a/z’: File exists

executing: ssh seb@Nacho2 "mkdir /home/seb/out/b"
mkdir: cannot create directory ‘/home/seb/out/b’: File exists

copying: {1GB-0 data/1GB-0 /1GB-0}
copying: {1GB-6 data/1GB-6 /1GB-6}
copying: {1GB-7 data/1GB-7 /1GB-7}
copying: {1GB-8 data/1GB-8 /1GB-8}
executing: scp data/1GB-0 seb@Nacho2:/home/seb/out/1GB-6
executing: scp data/1GB-6 seb@Nacho2:/home/seb/out/1GB-9
executing: scp data/1GB-7 seb@Nacho2:/home/seb/out/1GB-9
copying: {1GB-9 data/1GB-9 /1GB-9}
executing: scp data/1GB-8 seb@Nacho2:/home/seb/out/1GB-9
executing: scp data/1GB-9 seb@Nacho2:/home/seb/out/a/1GB-1
copying: {1GB-1 data/a/1GB-1 /a/1GB-1}
copying: {1GB-4 data/a/z/1GB-4 /a/z/1GB-4}
copying: {1GB-5 data/a/z/1GB-5 /a/z/1GB-5}
copying: {1GB-2 data/b/1GB-2 /b/1GB-2}
copying: {1GB-3 data/b/1GB-3 /b/1GB-3}
executing: scp data/a/1GB-1 seb@Nacho2:/home/seb/out/b/1GB-2
executing: scp data/b/1GB-3 seb@Nacho2:/home/seb/out/b/1GB-3
executing: scp data/a/z/1GB-5 seb@Nacho2:/home/seb/out/b/1GB-3
executing: scp data/b/1GB-2 seb@Nacho2:/home/seb/out/b/1GB-3
executing: scp data/a/z/1GB-4 seb@Nacho2:/home/seb/out/b/1GB-3

real	0m44.745s
user	1m32.844s
sys	0m36.513s

## comparison
* disks write speed 1.8 GBps
* slic 10 parallel jobs -> 10GB / 44secs = 227.272727 MBps
* scp -r -> 10GB / 72secs = 138.888889 MBps