package main

type Conuter interface {
	Add()
	Del()
}

type Conut struct {
}

func (Conut) Add(Addr string)  {
	for i := 0; i < len(ServerLists.ServerList); i++ {
		if ServerLists.ServerList[i].Addr == Addr{
			ServerLists.ServerList[i].Count = ServerLists.ServerList[i].Count + 1
			return
		}
	}
}

func (Conut) Del(Addr string)  {
	for i := 0; i < len(ServerLists.ServerList); i++ {
		if ServerLists.ServerList[i].Addr == Addr{
			ServerLists.ServerList[i].Count = ServerLists.ServerList[i].Count - 1
			return
		}
	}
}