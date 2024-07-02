CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Bảng lưu trữ thông tin về danh mục sách
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Bảng lưu trữ thông tin về các cuốn sách
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    author VARCHAR(100),
    description TEXT,
    category_id UUID REFERENCES categories(id),
    image_url TEXT,
    rarity BOOLEAN DEFAULT FALSE, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'auth_type') THEN
        CREATE TYPE auth_type AS ENUM ('credentials', 'github', 'google');
    END IF;
END$$;

-- Bảng lưu trữ thông tin về người dùng
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    balance NUMERIC(10, 2) DEFAULT 1000000.00,
    login_enabled BOOLEAN DEFAULT TRUE, 
    depot_limit INTEGER DEFAULT 10,
    auth_method auth_type NOT NULL DEFAULT 'credentials',
    role_id UUID REFERENCES roles(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng lưu trữ thông tin về sách của người dùng
CREATE TABLE user_books (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    current_price NUMERIC(10, 2) NOT NULL,
    quantity_stock INTEGER NOT NULL CHECK (quantity_stock >= 0),
    quantity_sell INTEGER NOT NULL DEFAULT 0 CHECK (quantity_sell >= 0),
    visible BOOLEAN DEFAULT TRUE,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, book_id)  -- Đảm bảo mỗi người dùng chỉ sở hữu một bản của mỗi sách
);

-- Bảng lưu trữ thông tin về các yêu cầu bán
CREATE TABLE sell_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    userbook_id UUID REFERENCES user_books(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'Pending', -- Trạng thái yêu cầu (Pending, Approved, Rejected)
    request_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng lưu trữ thông tin về thông báo
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    receiver_id UUID REFERENCES users(id) ON DELETE SET NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read BOOLEAN DEFAULT FALSE
);

-- Bảng lưu trữ thông tin về các đơn đặt hàng
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    userbook_id UUID REFERENCES user_books(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    total_price NUMERIC(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng liên kết giữa sách và danh mục sách
CREATE TABLE book_category (
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, category_id)
);

-- Bảng lưu trữ thông tin về giỏ hàng của người dùng
CREATE TABLE shopping_cart (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    userbook_id UUID REFERENCES user_books(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Bảng lưu trữ thông tin đánh giá
CREATE TABLE evaluate (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    userbook_id UUID REFERENCES user_books(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tokens (
    user_id UUID PRIMARY KEY,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    access_token_expiry TIMESTAMP NOT NULL,
    refresh_token_expiry TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT (now() AT TIME ZONE 'Asia/Ho_Chi_Minh'),
    updated_at TIMESTAMP DEFAULT (now() AT TIME ZONE 'Asia/Ho_Chi_Minh'),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Tạo danh mục sách
INSERT INTO categories (id, name) VALUES
  ('38c107fa-d0ec-4e1f-8db0-c012fb7f46f7', 'Trinh thám - Kinh dị'),
  ('0e0da8e1-95e4-4d9b-81a1-8590b9b37c27', 'Hài hước'),
  ('1f5873a0-15f6-4a6e-b8a5-27bb03e9245e', 'Lãng mạn'),
  ('356fb3a6-5183-4560-8d83-81f2fda33d82', 'Khoa học viễn tưởng');

-- Tạo dữ liệu cho bảng books với giá tiền theo tiền tệ Việt Nam (VND)
INSERT INTO books (id, title, author, description, category_id, image_url)
VALUES
  ('4e80ef43-3e94-46c8-8486-6e541f27a57d', 'Trở lại hang thôn', 'Sái Tuấn', 'Tôi nghe thấy linh hồn trong thân xác của mình đang hỏi : "Tôi còn sống không ?"', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7', 'kinhdi_2.jpg'),
  ('5f97ebbb-1538-4c7f-8d24-4b1db50dd865', 'Trò chơi tử thần', 'Liên Phạm', 'Liệu bạn có dám chơi đùa với tử thần', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7', 'kinhdi_1.jpg'),
  ('6c9839cf-f5f1-4f05-8b29-1a6a35b16688', 'Người gác đêm', 'Trung Nguyễn', 'Người gác đêm sẽ bảo vệ bạn khỏi những thứ nguy hiểm', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7', 'kinhdi_3.jpg'),
  ('7b2fe3b5-8f2a-4487-8622-1fd72734d6f0', 'Vui vẻ không quạu nha', 'Trung Nguyễn', 'Truyện hài hước giúp bạn xả xì stress', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27', 'haihuoc_1.jpg'),
  ('8d6c5108-2f08-4b43-a798-7283d6a1e8c2', 'Cuộc sống rất giống cuộc đời', 'Liên Ngô', 'Truyện hài hước giúp bạn xả xì stress', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27', 'haihuoc_2.jpg'),
  ('9e00935b-5157-4a77-8331-6b26b5d7931c', 'Cái nồi gì thế', 'Khánh huyền', 'Truyện hài hước giúp bạn xả xì stress', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27', 'haihuoc_3.jpg'),
  ('a8edc88f-77c0-4c09-bf21-1b9345dd02bc', 'Ba điều bí ẩn', 'Trà my', 'Ba điều bí ẩn nhất sẽ giúp bạn mở mang tâm trí', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7', 'kinhdi_4.jpg'),
  ('b2f7e127-b236-4861-b451-3ad7a5ec7b68', 'Sự im lặng của bầy cừu', 'Phạm Lâm', 'Sự im lặng luôn đem tới cái gì đó cực kì chết chóc', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7', 'kinhdi_5.jpg'),
  ('c3368a4e-1993-468f-9c6e-4e409ee1a396', 'Thăng cấp làm vợ', 'Trung Nguyễn', 'Mối tình của cặp đôi khiến cho bạn không thể nhịn cười !!!', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27', 'haihuoc_4.jpg'),
  ('d476b026-2594-47f6-8199-fa19a1b4982c', 'Đôrêmon - Máy hút chữ', 'Sasaki', 'Truyện tuổi thơ vui vẻ hài hước', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27', 'haihuoc_5.jpg');


-- Tạo dữ liệu cho bảng book_category
INSERT INTO book_category (book_id, category_id)
VALUES
  ('4e80ef43-3e94-46c8-8486-6e541f27a57d', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7'), -- Cuốn sách 1 thuộc thể loại trinh thám - kinh dị
  ('5f97ebbb-1538-4c7f-8d24-4b1db50dd865', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7'), -- Cuốn sách 2 thuộc thể loại trinh thám - kinh dị
  ('6c9839cf-f5f1-4f05-8b29-1a6a35b16688', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7'), -- Cuốn sách 3 thuộc thể loại trinh thám - kinh dị
  ('7b2fe3b5-8f2a-4487-8622-1fd72734d6f0', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27'), -- Cuốn sách 4 thuộc thể loại hài hước
  ('8d6c5108-2f08-4b43-a798-7283d6a1e8c2', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27'), -- Cuốn sách 5 thuộc thể loại hài hước
  ('9e00935b-5157-4a77-8331-6b26b5d7931c', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27'), -- Cuốn sách 6 thuộc thể loại hài hước
  ('a8edc88f-77c0-4c09-bf21-1b9345dd02bc', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7'), -- Cuốn sách 7 thuộc thể loại trinh thám - kinh dị
  ('b2f7e127-b236-4861-b451-3ad7a5ec7b68', '38c107fa-d0ec-4e1f-8db0-c012fb7f46f7'), -- Cuốn sách 8 thuộc thể loại trinh thám - kinh dị
  ('c3368a4e-1993-468f-9c6e-4e409ee1a396', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27'), -- Cuốn sách 9 thuộc thể loại hài hước
  ('d476b026-2594-47f6-8199-fa19a1b4982c', '0e0da8e1-95e4-4d9b-81a1-8590b9b37c27'); -- Cuốn sách 10 thuộc thể loại hài hước

-- Insert data into the roles table
INSERT INTO roles (id, name) VALUES
  ('38c107fa-d0ec-4e1f-8db0-c012fb7f46f7', 'admin'),
  ('0e0da8e1-95e4-4d9b-81a1-8590b9b37c27', 'user'),
  ('6c9839cf-f5f1-4f05-8b29-1a6a35b16688', 'employee'),
  ('7b2fe3b5-8f2a-4487-8622-1fd72734d6f0', 'manager');

-- Thêm dữ liệu cho bảng users
INSERT INTO users (id,username, email, password, balance,role_id)
VALUES
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78','trungnguyen1', 'trungnguyen1@gmail.com', '1232003', 1000000.00,'38c107fa-d0ec-4e1f-8db0-c012fb7f46f7'),
  (uuid_generate_v4(),'trungnguyen2', 'trungnguyen2@gmail.com', '1232003', 1000000.00,'0e0da8e1-95e4-4d9b-81a1-8590b9b37c27'),
  (uuid_generate_v4(),'trungnguyen_3', 'trungnguyen_3@example.com', '1232003', 1000000.00,'6c9839cf-f5f1-4f05-8b29-1a6a35b16688'),
  (uuid_generate_v4(),'trungnguyen_4', 'trungnguyen_4@example.com', '1232003', 1000000.00,'7b2fe3b5-8f2a-4487-8622-1fd72734d6f0');

INSERT INTO user_books (user_id, book_id, current_price, quantity_stock, quantity_sell)
VALUES
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', '4e80ef43-3e94-46c8-8486-6e541f27a57d', 68900, 20, 10),
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', '5f97ebbb-1538-4c7f-8d24-4b1db50dd865', 45977, 15, 10),
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', '6c9839cf-f5f1-4f05-8b29-1a6a35b16688', 57477, 25, 10),
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', '7b2fe3b5-8f2a-4487-8622-1fd72734d6f0', 80477, 10, 5),
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', '8d6c5108-2f08-4b43-a798-7283d6a1e8c2', 34477, 30, 10),
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', '9e00935b-5157-4a77-8331-6b26b5d7931c', 91977, 18, 10),
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', 'a8edc88f-77c0-4c09-bf21-1b9345dd02bc', 10337, 22, 10),
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', 'b2f7e127-b236-4861-b451-3ad7a5ec7b68', 50670, 27, 10),
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', 'c3368a4e-1993-468f-9c6e-4e409ee1a396', 45000, 40, 10),
  ('D123ADBF-DBBD-42FC-86F8-E28C0D3C6B78', 'd476b026-2594-47f6-8199-fa19a1b4982c', 68500, 35, 10);


-- DO $$ DECLARE
--     r RECORD;
-- BEGIN
--     FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
--         EXECUTE 'DROP TABLE IF EXISTS ' || quote_ident(r.tablename) || ' CASCADE';
--     END LOOP;
-- END $$;
