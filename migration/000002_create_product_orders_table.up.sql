CREATE TABLE product_orders (
    order_id TEXT NOT NULL,
    product_id INTEGER NOT NULL,
    product_name VARCHAR(100) NOT NULL,
    image VARCHAR(300) NOT NULL,
    price INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    CONSTRAINT product_orders_pkey PRIMARY KEY (order_id, product_id)
);

ALTER TABLE product_orders ADD CONSTRAINT product_orders_id_fkey FOREIGN KEY(order_id) REFERENCES orders(order_id) ON DELETE RESTRICT ON UPDATE CASCADE;