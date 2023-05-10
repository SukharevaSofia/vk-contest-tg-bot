package main

func main() {
	db := createDb("vkcontest")
	tgManageRequests(db)

}

//TODO: different dbs for different users
//TODO: layers of commands
//TODO: commands should be hidden from users
