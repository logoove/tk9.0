cd _examples/
for f in *.go
do
	go run $f &
	sleep 1
done
cd ..
