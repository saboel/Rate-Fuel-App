package clients

type getClientQuote struct {
	date string 
	gallons_req int 
	delivery_date string
	delivery_address string
	price_per_gall float64 
	total float64
}
