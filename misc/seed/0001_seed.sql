-- delete all user_groups
delete from user_groups;

-- delete all groups
UPDATE users
SET "active_group_id" = null;
delete from item;
delete from location;
delete from manufacturer;
delete from groups;

-- delete all users
delete from users;
ALTER SEQUENCE users_id_seq RESTART WITH 1;
INSERT INTO users
("id","name","email","avatar","active_group_id")
values
(1,'Ugur Emirmustafaoglu','uguremirmustafa@gmail.com','https://lh3.googleusercontent.com/a/ACg8ocIw41m2yoEqfU1VsI8zaU6hn8xYA0Lm5wQAvf2rNhJtdOVm8i9u=s96-c',null),
(2,'Lionel Messi','leomessi@gmail.com','https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS8LhSodMPo5XTyGAnXTaOI_lfXJBeci-2MUQ&s',null);

INSERT INTO groups
("id","name","description","group_owner_id")
values
(1,'Ugur Emirmustafaoglu Family',null,1);

INSERT INTO user_groups
("user_id","group_id")
values
(1,1),
(2,1);

UPDATE users
SET "active_group_id" = 1
WHERE id in (1,2);


-- delete all items
delete from item;
ALTER SEQUENCE item_id_seq RESTART WITH 1;

-- insert item types
delete from item_type;
ALTER SEQUENCE item_type_id_seq RESTART WITH 1;
insert into item_type 
("name", "description", "icon_class") 
values
('clothing, shoes and accessories','You can store clothing items, shoes and accessories under this category.','ion:shirt-outline'),
('technological device','Laptops, sound systems, cables and mobile phones etc.','ph:devices-light'),
('furniture','Sofa, bed, chairs and carpets. Any furniture can be stored under this category.','solar:sofa-2-linear'),
('household appliances','Fridge, oven, microwave and your favourite toaster goes here.','solar:fridge-outline'),
('cooking utensils','Spoons, pots and forks. Dont forget the jars.','hugeicons:kitchen-utensils'),
('sports equipment','Bike, yoga mat and weights.','game-icons:weight-lifting-up'),
('stationery and books','Pencils, A4 papers, books and notebooks goes here.','lucide:notebook-pen');


-- insert manufacturers
ALTER SEQUENCE manufacturer_id_seq RESTART WITH 1;
insert into manufacturer 
("name", "logo_url", "description", "group_id") 
values
('levis','https://st-levis.mncdn.com/Content/img/levis_logo-500px.png','Levi Strauss & Co. | Jeans', 1),
('mavi','https://upload.wikimedia.org/wikipedia/commons/2/28/Logo_of_Mavi.png','Mavi, denim ürünleri ile tanınan tekstil firması.', 1);


-- insert location

ALTER SEQUENCE location_id_seq RESTART WITH 1;
insert into location 
("name", "description", "group_id","image_url") 
values
('Bedroom','Master bedroom for parents.', 1, 'uploads/seed/bedroom1.jpg'),
('Office','Cozy working place for my daily work.', 1, 'uploads/seed/office.jpg'),
('Living Room','Living room where I love watching TV.', 1,'uploads/seed/living_room.jpg');

-- insert item
insert into item
("name","description","user_id","group_id", "item_type_id","manufacturer_id")
values
('beyaz dar paca kot','en sevdigim beyaz pantolonum',1,1,1,1);
