package entities

type Row struct {
	Value     string
	Timestamp string
}

type Statement struct {
	Table     string
	Timestamp string
	ID        int
}

type CapturedData struct {
	ID   int
	Rows []Row
}
