package postgres

import (
	"context"

	"github.com/Ckakalka/wbLevel0/db"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Manager struct {
	db *sqlx.DB
}

func NewManager() (db.Manager, error) {
	var manager Manager
	var err error
	manager.db, err = sqlx.Open("pgx",
		"user=valeriy password=12345 host=localhost port=5432 dbname=wblevel0")
	if err != nil {
		return nil, err
	}
	err = manager.db.Ping()
	if err != nil {
		return nil, err
	}
	return manager, nil
}

func (m Manager) Close() error {
	return m.db.Close()
}

func (m Manager) InsertOrder(order db.Order) error {
	tx, err := m.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	orderQuery := `INSERT INTO public.Order
			(Uid, TrackNumber, Entry, Locale,
			 InternalSignature, CustomerId, DeliveryService, Shardkey,
			 SmId, DateCreated, Oofshard)
	VALUES (:uid, :tracknumber, :entry, :locale,
			:internalsignature, :customerid, :deliveryservice, :shardkey,
			:smid, :datecreated, :oofshard)`

	if _, err := tx.NamedExec(orderQuery, order); err != nil {
		return err
	}
	paymentQuery := `INSERT INTO public.Payment
					(Transaction, RequestId, Currency, Provider,
					Amount, PaymentDt, Bank, DeliveryCost,
					GoodsTotal, CustomFee, OrderUid)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	p := &order.Payment
	paymentArgs := []any{p.Transaction, p.RequestId, p.Currency, p.Provider,
		p.Amount, p.PaymentDt, p.Bank, p.DeliveryCost,
		p.GoodsTotal, p.CustomFee, order.Uid}
	if _, err := tx.Exec(paymentQuery, paymentArgs...); err != nil {
		return err
	}
	deliveryQuery := `INSERT INTO public.Delivery
					(Name, Phone, Zip, City,
					Address, Region, Email, OrderUid)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	d := &order.Delivery
	deliveryArgs := []any{d.Name, d.Phone, d.Zip, d.City,
		d.Address, d.Region, d.Region, order.Uid}
	if _, err := tx.Exec(deliveryQuery, deliveryArgs...); err != nil {
		return err
	}
	orderItemQuery := `INSERT INTO public.OrderItem
						(OrderUid, ItemChrtId)
				VALUES ($1, $2)`
	for _, item := range order.Items {
		if _, err := tx.Exec(orderItemQuery, order.Uid, item.ChrtId); err != nil {
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (m Manager) GetOrder(uid string) (db.Order, error) {
	orderQuery := `select * from public.Order where uid=$1`
	var order db.Order
	if err := m.db.Get(&order, orderQuery, uid); err != nil {
		return order, err
	}

	deliveryQuery := `select name, phone, zip, city, address, region, email
					from public.delivery where orderuid=$1`
	if err := m.db.Get(&order.Delivery, deliveryQuery, uid); err != nil {
		return order, err
	}

	paymentQuery := `select transaction, requestid, currency, provider, amount,
							paymentdt, bank, deliverycost, goodstotal, customfee
					from public.payment where orderuid=$1`
	if err := m.db.Get(&order.Payment, paymentQuery, uid); err != nil {
		return order, err
	}

	orderItemQuery := `select itemchrtid as chrtid from public.orderitem where orderuid=$1`
	var items db.Items
	if err := m.db.Select(&items, orderItemQuery, uid); err != nil {
		return order, err
	}

	itemQuery := `select * from public.item where chrtid=$1`
	for i := range items {
		if err := m.db.Get(&items[i], itemQuery, items[i].ChrtId); err != nil {
			return order, err
		}
	}
	order.Items = items
	return order, nil
}

func (m Manager) GetAllOrders() ([]db.Order, error) {
	getOrdersUidQuery := `select uid from public.order`
	var uids []string
	if err := m.db.Select(&uids, getOrdersUidQuery); err != nil {
		return nil, err
	}
	orders := make([]db.Order, 0, len(uids))
	for _, uid := range uids {
		order, err := m.GetOrder(uid)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
