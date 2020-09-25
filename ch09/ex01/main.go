package bank

var deposits = make(chan int)
var balances = make(chan int)
var withdrawRequests = make(chan withdrawRequest)

type withdrawRequest struct {
	Amount int
	Result chan bool
}

func Deposit(amount int) {
	deposits <- amount
}

func Balance() int {
	return <-balances
}

func Withdraw(amount int) bool {
	req := withdrawRequest{Amount: amount, Result: make(chan bool)}
	withdrawRequests <- req
	return <- req.Result
}

func teller() {
	var balance int
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case req := <-withdrawRequests:
			if balance >= req.Amount {
				balance -= req.Amount
				req.Result <- true
			} else {
				req.Result <- false
			}
		}
	}
}

func init() {
	go teller()
}
