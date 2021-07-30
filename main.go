package main

func main() {
	r := InitRouter()
	_ = r.Run(":8000")
}
