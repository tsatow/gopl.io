package bank

var deposits = make(chan int)
var balances = make(chan int)
type withdraw struct {
	Amount int
	Result chan<- bool
}
var withdraws = make(chan withdraw)

func Deposit(amount int) {
	deposits <- amount
}

func Balance() int {
	return <- balances
}

func WithDraw(amount int) bool {
	result := make(chan bool)
	withdraws <- withdraw{amount, result}
	return <- result
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case withdraw := <-withdraws:
			if balance >= withdraw.Amount {
				balance -= withdraw.Amount
				withdraw.Result <- true
			} else {
				withdraw.Result <- false
			}
		}
	}
}

func init() {
	go teller()
}