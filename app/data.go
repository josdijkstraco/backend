package app

type Product struct {
	ID           int `json:"_id"`
	UserID       int
	Name         string  `json:"name"`
	Image        string  `json:"image"`
	Description  string  `json:"description"`
	Brand        string  `json:"brand"`
	Price        float32 `json:"price"`
	CountInStock int     `json:"countInStock"`
	NumReviews   int     `json:"numReviews"`
	Rating       float32 `json:"rating"`
	Category     string  `json:"category"`
	CreatedAt    string  `json:"createdAt"`
}

type Listing struct {
	ID                  int     `json:"id"`
	RealtorID           int     `json:"realtorId"`
	Address             string  `json:"address"`
	City                string  `json:"city"`
	PostalCode          string  `json:"postalcode"`
	Latitude            float32 `json:"lat"`
	Longitude           float32 `json:"lon"`
	FirstName           string  `json:"first_name"`
	LastName            string  `json:"last_name"`
	Commission          float32 `json:"commission"`
	Image               string  `json:"image"`
	Price               string  `json:"price"`
	NumberBedrooms      int     `json:"number_bedrooms"`
	NumberFullbacts     int     `json:"number_full_baths"`
	NumberThreeQtrBaths int     `json:"number_three_qtr_baths"`
	NumberHalfBaths     int     `json:"number_half_baths"`
	FinishedSquareFt    int     `json:"finished_sq_ft"`
	Owner               string  `json:"owner"`
	LegalDescription    string  `json:"legal_description"`
	Strap               string  `json:"strap"`
	ParcelNumber        string  `json:"parcel_number"`
}

var data = []Product{
	{ID: 1, Name: "Airpods Wireless Bluetooth Headphones", Image: "/images/airpods.jpg", Description: "Bluetooth technology lets you connect it with compatible devices wirelessly High-quality AAC audio offers immersive listening experience Built-in microphone allows you to take calls while working",
		Brand: "Apple", Price: 89.99, CountInStock: 23, NumReviews: 2,
	},
	{ID: 2, Name: "iPhone 11 Pro 256GB Memory", Image: "/images/phone.jpg", Description: "Introducing the iPhone 11 Pro. A transformative triple-camera system that adds tons of capability without complexity. An unprecedented leap in battery life",
		Brand: "Apple", Price: 589.99, CountInStock: 12, NumReviews: 5,
	},
}

/*
{
    '_id': '1',
    'name': 'Airpods Wireless Bluetooth Headphones',
    'image': '/images/airpods.jpg',
    'description':
      'Bluetooth technology lets you connect it with compatible devices wirelessly High-quality AAC audio offers immersive listening experience Built-in microphone allows you to take calls while working',
    'brand': 'Apple',
    'category': 'Electronics',
    'price': 89.99,
    'countInStock': 10,
    'rating': 4.5,
    'numReviews': 12,
  },
  {
    '_id': '2',
    'name': 'iPhone 11 Pro 256GB Memory',
    'image': '/images/phone.jpg',
    'description':
      'Introducing the iPhone 11 Pro. A transformative triple-camera system that adds tons of capability without complexity. An unprecedented leap in battery life',
    'brand': 'Apple',
    'category': 'Electronics',
    'price': 599.99,
    'countInStock': 7,
    'rating': 4.0,
    'numReviews': 8,
  },
  {
    '_id': '3',
    'name': 'Cannon EOS 80D DSLR Camera',
    'image': '/images/camera.jpg',
    'description':
      'Characterized by versatile imaging specs, the Canon EOS 80D further clarifies itself using a pair of robust focusing systems and an intuitive design',
    'brand': 'Cannon',
    'category': 'Electronics',
    'price': 929.99,
    'countInStock': 5,
    'rating': 3,
    'numReviews': 12,
  },
  {
    '_id': '4',
    'name': 'Sony Playstation 4 Pro White Version',
    'image': '/images/playstation.jpg',
    'description':
      'The ultimate home entertainment center starts with PlayStation. Whether you are into gaming, HD movies, television, music',
    'brand': 'Sony',
    'category': 'Electronics',
    'price': 399.99,
    'countInStock': 11,
    'rating': 5,
    'numReviews': 12,
  },
  {
    '_id': '5',
    'name': 'Logitech G-Series Gaming Mouse',
    'image': '/images/mouse.jpg',
    'description':
      'Get a better handle on your games with this Logitech LIGHTSYNC gaming mouse. The six programmable buttons allow customization for a smooth playing experience',
    'brand': 'Logitech',
    'category': 'Electronics',
    'price': 49.99,
    'countInStock': 7,
    'rating': 3.5,
    'numReviews': 10,
  },
  {
    '_id': '6',
    'name': 'Amazon Echo Dot 3rd Generation',
    'image': '/images/alexa.jpg',
    'description':
      'Meet Echo Dot - Our most popular smart speaker with a fabric design. It is our most compact smart speaker that fits perfectly into small space',
    'brand': 'Amazon',
    'category': 'Electronics',
    'price': 29.99,
    'countInStock': 0,
    'rating': 4,
    'numReviews': 12,
  }
`
*/
