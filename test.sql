SELECT 
    c.name
    , i.customer_id
    , i.date
    , p.name as product
    , ii.quantity
    , ii.quantity * p.price as total 
FROM invoice i 
JOIN customer c ON i.customer_id = c.id 
JOIN invoice_item ii ON i.id = ii.invoice_id JOIN product p ON ii.product_id = p.id;

