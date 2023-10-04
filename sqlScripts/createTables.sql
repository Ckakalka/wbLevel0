CREATE TABLE public.Order
(
	Uid varchar PRIMARY KEY,
	TrackNumber varchar,
	Entry varchar,
	Locale varchar,
	InternalSignature varchar,
	CustomerId varchar,
	DeliveryService varchar,
	ShardKey int,
	SmId int,
	DateCreated timestamp with time zone,
	OofShard int
);
CREATE TABLE public.Item
(
	ChrtId int PRIMARY KEY,
	TrackNumber varchar,
	Price int,
	Rid varchar,
	Name varchar,
	Sale int,
	Size int,
	TotalPrice int,
	NmId int,
	Brand varchar,
	Status int
);
CREATE TABLE public.OrderItem
(
	OrderUid varchar,
	ItemChrtId int,
	FOREIGN KEY(OrderUid) REFERENCES public.Order(Uid),
	FOREIGN KEY(ItemChrtId) REFERENCES public.Item(ChrtId)
);
CREATE TABLE public.Delivery
(
	Name varchar,
	Phone varchar PRIMARY KEY,
	OrderUid varchar,
	Zip int,
	City varchar,
	Address varchar,
	Region varchar,
	Email varchar,
	FOREIGN KEY(OrderUid) REFERENCES public.Order(Uid)
);
CREATE TABLE public.Payment
(
	Transaction varchar primary key,
	OrderUid varchar unique,
	RequestId varchar,
	Currency varchar,
	Provider varchar,
	Amount int,
	PaymentDt int,
	Bank varchar,
	DeliveryCost int,
	GoodsTotal int,
	CustomFee int,
	FOREIGN KEY(OrderUid) REFERENCES public.Order(Uid)
);





















