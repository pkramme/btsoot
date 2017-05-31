package main

func CreateMasterProcessList() map[int]Process {
  pmap := make(map[int]Process)

  pmap[0] = Process {
    Channel: make(chan int),
    Level: 0,
    Description: "Update thread",
  }
  pmap[1] = Process {
    Channel: make(chan int),
    Level: 0,
    Description: "Master process for all scanning operations",
  }
  pmap[2] = Process {
    Channel: make(chan int),
    Level: 0,
    Description: "Webserver",
  }
  return pmap
}

func ListProcessList() {
  return
}

func (p Process) AddProcessToList() {
  return
}

func (p Process) Kill() {
  return
}

type Process struct {
  Channel chan int
  Level int
  Description string
}
