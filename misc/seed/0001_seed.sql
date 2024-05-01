-- delete all items
delete from item;
ALTER SEQUENCE item_id_seq RESTART WITH 1;

-- insert item types
delete from item_type;
ALTER SEQUENCE item_type_id_seq RESTART WITH 1;
insert into item_type ("name") values('trousers');
insert into item_type ("name") values('technological device');

-- insert manufacturers
delete from manufacturer;
ALTER SEQUENCE manufacturer_id_seq RESTART WITH 1;
insert into manufacturer 
("name", "logo_url", "description") 
values
('levis','https://st-levis.mncdn.com/Content/img/levis_logo-500px.png','Levi Strauss & Co. | Jeans'),
('mavi','https://upload.wikimedia.org/wikipedia/commons/2/28/Logo_of_Mavi.png','Mavi, denim ürünleri ile tanınan tekstil firması.');


-- insert location
delete from location;
ALTER SEQUENCE location_id_seq RESTART WITH 1;
insert into location 
("name", "description") 
values
('depo','daire ici depo'),
('feyzanin gardrop','yatak odasi, feyzanin gardrop');

-- insert item
insert into item
("name","description","user_id","item_type_id","manufacturer_id")
values
('beyaz dar paca kot','en sevdigim beyaz pantolonum',1,1,1);
