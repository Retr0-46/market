CREATE TYPE product_type AS ENUM ('sport', 'electronics', 'books', 'clothes');

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    p_type product_type NOT NULL,
    seller_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),
    stock_quantity BIGINT NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0)
);

INSERT INTO products (p_type, seller_id, name, price, stock_quantity) VALUES
-- sport (10)
('sport', 1, 'Беговые кроссовки RunFlex 5',          899000,  42),
('sport', 1, 'Скакалка SpeedRope Pro',                59000, 310),
('sport', 1, 'Гантели разборные 2x20 кг',           1249000,  17),
('sport', 1, 'Коврик для йоги TPE 6 мм',             219000,  85),
('sport', 1, 'Велотренажёр HomeFit HT-300',         3299000,   6),
('sport', 1, 'Бутылка для воды 750 мл',               79000, 240),
('sport', 1, 'Фитнес-резинки, набор 5 шт.',          149000, 130),
('sport', 1, 'Штанга 20 кг + гриф',                 1599000,   9),
('sport', 1, 'Пульсометр Polar H9',                  699000,  28),
('sport', 1, 'Спортивная сумка 40 л',                329000,  64),

-- electronics (10)
('electronics', 1, 'Смартфон Pixel 8 128GB',        5999000,  23),
('electronics', 1, 'Наушники Sony WH-1000XM5',      3499000,  41),
('electronics', 1, 'Ноутбук MacBook Air M3',       11999000,  12),
('electronics', 1, 'Мышь Logitech MX Master 3S',     899000,  97),
('electronics', 1, 'Клавиатура Keychron K2',        1199000,  35),
('electronics', 1, 'Монитор Dell 27" 4K',           4499000,  18),
('electronics', 1, 'Внешний SSD 1TB Samsung T7',    1099000,  56),
('electronics', 1, 'Умная колонка Яндекс Станция',   899000,  74),
('electronics', 1, 'Роутер TP-Link AX55',            649000,  44),
('electronics', 1, 'Powerbank 20000 mAh',            399000, 112),

-- books (10)
('books', 1, 'Чистый код. Роберт Мартин',            189000, 120),
('books', 1, 'Грокаем алгоритмы',                    169000,  98),
('books', 1, 'Паттерны проектирования (GoF)',        219000,  55),
('books', 1, 'Совершенный код. Стив Макконнелл',     299000,  37),
('books', 1, 'Дюна. Фрэнк Герберт',                   99000, 210),
('books', 1, '1984. Джордж Оруэлл',                   79000, 305),
('books', 1, 'Мастер и Маргарита',                    89000, 180),
('books', 1, 'Преступление и наказание',             109000, 142),
('books', 1, 'Атлант расправил плечи (3 тома)',      249000,  61),
('books', 1, 'Sapiens. Краткая история человечества',159000,  88),

-- clothes (10)
('clothes', 1, 'Футболка хлопок, чёрная, M',         129000, 220),
('clothes', 1, 'Джинсы Slim Fit, 32/32',             499000,  68),
('clothes', 1, 'Худи oversize, серое',               349000,  95),
('clothes', 1, 'Куртка зимняя пуховая',             1299000,  14),
('clothes', 1, 'Кроссовки casual белые, 42',         599000,  47),
('clothes', 1, 'Носки, набор 10 пар',                 69000, 380),
('clothes', 1, 'Рубашка классическая белая, L',      279000,  52),
('clothes', 1, 'Пальто шерстяное, 48',              1499000,   8),
('clothes', 1, 'Кепка бейсболка хлопок',              89000, 160),
('clothes', 1, 'Шарф кашемировый',                   329000,  33);