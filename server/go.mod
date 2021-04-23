module github.com/bitcapybara/cuckoo/server

go 1.16

require (
	github.com/bitcapybara/cuckoo/core v0.0.1
	github.com/bitcapybara/raft v0.0.1
	github.com/emirpasic/gods v1.12.0 // indirect
)

replace (
	github.com/bitcapybara/cuckoo/core => ../core
	github.com/bitcapybara/raft => ../../raft
)
