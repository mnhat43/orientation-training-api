-- ------------------------------------------ user_roles ----------------------------------------------
INSERT INTO user_roles (id, name, description, created_at, updated_at, deleted_at) VALUES
(1, 'Admin', 'Administrator role', '2024-11-12 11:00:00', '2024-11-12 11:00:00', NULL),
(2, 'Manager', 'Manager role', '2024-11-12 11:00:00', '2024-11-12 11:00:00', NULL),
(3, 'User', 'Regular user role', '2024-11-12 11:00:00', '2024-11-12 11:00:00', NULL),
(4, 'General Manager', 'General Manager role', '2024-11-12 11:00:00', '2024-11-12 11:00:00', NULL);

-- ------------------------------------------ users ----------------------------------------------
INSERT INTO users (id, email, password, role_id, last_login_time, created_at, updated_at, deleted_at) VALUES
(1, 'dmn@gmail.com', 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3', 2, '2024-11-10 10:00:00', '2024-01-01 10:00:00', '2024-01-01 10:00:00', NULL),
(2, 'user2@example.com', 'password123', 2, '2024-11-12 11:00:00', '2024-02-01 10:00:00', '2024-02-01 10:00:00', NULL),
(3, 'user3@example.com', 'password123', 3, '2024-11-15 12:00:00', '2024-03-01 10:00:00', '2024-03-01 10:00:00', NULL),
(4, 'user4@example.com', 'password123', 4, '2024-11-17 14:00:00', '2024-04-01 10:00:00', '2024-04-01 10:00:00', NULL);

---------------------------------------------user_profiles----------------------------------------
INSERT INTO user_profiles (id, user_id, avatar, first_name, last_name, birthday, phone_number, personal_email, company_joined_date, introduce, gender, created_at, updated_at, deleted_at) VALUES
(1, 1, 'https://placekitten.com/911/302', 'Nhat', 'Dao', '2003-04-03', '157.156.1000', 'antonio.rivera@example.com', '2020-01-01', 'Introduction text for Antonio.', 1,'2024-11-12 11:00:00', '2024-11-12 11:00:00', NULL),
(2, 2, 'https://placekitten.com/911/302', 'Virginia', 'Boone', '1940-09-09', '+1-544-389-3962x8888', 'virginia.boone@example.com', '2021-01-01', 'Introduction text for Virginia.', 0,'2024-11-12 11:00:00', '2024-11-12 11:00:00', NULL),
(3, 3, 'https://placekitten.com/911/302', 'Jessica', 'Howard', '1973-06-06', '864-807-0007x869', 'jessica.howard@example.com', '2022-01-01', 'Introduction text for Jessica.', 1,'2024-11-12 11:00:00', '2024-11-12 11:00:00', NULL),
(4, 4, 'https://placekitten.com/911/302', 'Melanie', 'Heath', '1965-03-09', '+1-540-750-4656', 'melanie.heath@example.com', '2023-01-01', 'Introduction text for Melanie.', 0,'2024-11-12 11:00:00', '2024-11-12 11:00:00', NULL);

-- ------------------------------------------- courses ---------------------------------------------
INSERT INTO courses (id, title, description, thumbnail, created_by, created_at, updated_at, deleted_at) VALUES
(1,'Introduction to PostgreSQL','Learn the basics of PostgreSQL, from setup to query writing.','https://via.placeholder.com/150',1,'2025-03-09 11:46:12.543623','2025-03-09 11:46:12.543623','2025-03-09 11:58:50.904155'),
(2,'Advanced SQL Queries','Master advanced SQL queries including joins, subqueries, and performance optimization.','https://via.placeholder.com/150',2,'2025-03-09 11:46:12.543623','2025-03-09 11:46:12.543623','2025-03-10 11:49:28.070764'),
(3,'Web Development with Go','A full-stack web development course using Go, focusing on web servers and APIs.','https://via.placeholder.com/150',1,'2025-03-09 11:46:12.543623','2025-03-09 11:46:12.543623','2025-03-10 11:49:29.728733'),
(4,'Data Science with Python','Learn the fundamentals of Data Science using Python, including data analysis and machine learning.','https://via.placeholder.com/150',3,'2025-03-09 11:46:12.543623','2025-03-09 11:46:12.543623','2025-03-10 11:49:31.019082'),
(5,'Machine Learning Basics','An introduction to machine learning algorithms and how to implement them using Python.','https://via.placeholder.com/150',4,'2025-03-09 11:46:12.543623','2025-03-09 11:46:12.543623','2025-03-10 11:49:32.58576'),
(11,'123','123',NULL,1,'2025-03-18 14:57:10.87656','2025-03-18 14:57:10.87656','2025-03-18 15:39:20.393913'),
(6,'test','123',NULL,1,'2025-03-11 10:58:01.10642','2025-03-11 10:58:01.10642','2025-03-18 15:39:31.978324'),
(7,'123','213','1_1742033167349',1,'2025-03-15 10:06:10.420183','2025-03-15 10:06:10.420183','2025-03-10 11:49:32.585'),
(8,'Khóa học PostgreSQL','123','1_1742232731640',1,'2025-03-17 17:32:13.77903','2025-03-17 17:32:13.77903','2025-03-10 11:49:32.585'),
(9,'Khóa học PostgreSQL','12312312','1_1742232837446',1,'2025-03-17 17:33:59.578737','2025-03-17 17:33:59.578737','2025-03-10 11:49:32.585'),
(10,'Khóa học PostgreSQL','123','1_1742233058031',1,'2025-03-17 17:37:40.060172','2025-03-17 17:37:40.060172','2025-03-10 11:49:32.585'),
(12,'123','123','1_1742316409109.png',1,'2025-03-18 16:46:49.710585','2025-03-18 16:46:49.710585','2025-03-10 11:49:32.585'),
(13,'test','test','',1,'2025-03-18 17:02:31.577293','2025-03-18 17:02:31.577293','2025-03-23 14:27:00.786403'),
(14,'123','123','1_1742741550006.png',1,'2025-03-23 14:52:30.540187','2025-03-23 14:52:30.540187',NULL),
(15,'Lập Trình JavaScript Cơ Bản','JavaScript Cơ Bản','1_1742964470395.png',1,'2025-03-26 04:47:50.629297','2025-03-26 04:47:50.629297',NULL);

---------------------------------------------user_courses----------------------------------------
INSERT INTO user_courses (id, user_id, course_id, created_at, updated_at, deleted_at)
VALUES
    (1, 2, 14, '2025-03-18 14:57:10.87656', '2025-03-18 14:57:10.87656', NULL), 
    (2, 2, 15, '2025-03-18 14:57:10.87656', '2025-03-18 14:57:10.87656', NULL);

-- ------------------------------------------- modules ---------------------------------------------
INSERT INTO modules (id, title, course_id, position, created_at, updated_at, deleted_at) VALUES 
(1,'Tuần 1: Giới thiệu về lập trình',1, 1,'2025-03-09 11:46:12.579403','2025-03-09 11:46:12.579403',NULL),
(2,'Tuần 2: Cấu trúc dữ liệu cơ bản',1, 2,'2025-03-09 11:46:12.579403','2025-03-09 11:46:12.579403',NULL),
(3,'Tuần 3: Thuật toán và tối ưu hóa',1, 3,'2025-03-09 11:46:12.579403','2025-03-09 11:46:12.579403',NULL),
(4,'1',6, 1,'2025-03-11 10:58:07.564676','2025-03-11 10:58:07.564676',NULL),
(5,'Tuần 1',14, 1,'2025-03-23 14:52:38.528537','2025-03-23 14:52:38.528537',NULL),
(6,'Tuần 2',14, 2,'2025-03-25 14:50:53.158651','2025-03-25 14:50:53.158651',NULL),
(7,'1. Giới thiệu',15, 1,'2025-03-26 04:48:13.640161','2025-03-26 04:48:13.640161',NULL),
(8,'2. Biến, comments, built-in',15, 2, '2025-03-26 15:42:48.346678','2025-03-26 15:42:48.346678',NULL),
(9,'3. Toán tử, kiểu dữ liệu',15, 3, '2025-03-26 15:47:30.752282','2025-03-26 15:47:30.752282',NULL),
(10,'4. Làm việc với hàm',15, 4,'2025-03-26 15:50:27.736618','2025-03-26 15:50:27.736618',NULL);

-- ------------------------------------------- module_items ---------------------------------------------
INSERT INTO module_items (id, title, item_type, resource, module_id, position, created_at, updated_at, deleted_at) VALUES
(1,'Video 1: Giới thiệu về lập trình','video','https://example.com/video1',1,1,'2025-03-09 11:46:12.601886','2025-03-09 11:46:12.601886',NULL),
(2,'Tài liệu 1: Tổng quan về lập trình','file','https://example.com/file1',1,2,'2025-03-09 11:46:12.601886','2025-03-09 11:46:12.601886',NULL),
(3,'Video 2: Cấu trúc dữ liệu cơ bản','video','https://example.com/video2',2,1,'2025-03-09 11:46:12.604838','2025-03-09 11:46:12.604838',NULL),
(4,'Tài liệu 2: Cấu trúc dữ liệu mảng và danh sách liên kết','file','https://example.com/file2',2,2,'2025-03-09 11:46:12.604838','2025-03-09 11:46:12.604838',NULL),
(5,'Video 3: Giới thiệu thuật toán và tối ưu hóa','video','https://example.com/video3',3,1,'2025-03-09 11:46:12.607187','2025-03-09 11:46:12.607187',NULL),
(6,'Tài liệu 3: Các kỹ thuật tối ưu hóa thuật toán','file','https://example.com/file3',3,2,'2025-03-09 11:46:12.607187','2025-03-09 11:46:12.607187',NULL),
(7,'Lời khuyên trước khóa học Node Express','video','z2f7RHgvddc',4,1,'2025-03-11 10:58:40.186152','2025-03-11 10:58:40.186152',NULL),
(8,'2','video','_52DDwyU_Pc',4,2,'2025-03-11 10:59:03.594678','2025-03-11 10:59:03.594678',NULL),
(11,'05_プロブレムインタビュー(1).pdf','file','4_1741695322598',4,3,'2025-03-11 12:15:39.710608','2025-03-11 12:15:39.710608','2025-03-11 12:56:57.257749'),
(10,'05_プロブレムインタビュー(1).pdf','file','4_1741695330746',4,4,'2025-03-11 12:15:34.848054','2025-03-11 12:15:34.848054','2025-03-11 12:56:58.877552'),
(9,'05_プロブレムインタビュー(1).pdf','file','4_1741691952664',4,5,'2025-03-11 11:19:16.134139','2025-03-11 11:19:16.134139','2025-03-11 12:57:00.392776'),
(12,'ai_quiz_1.pdf','file','4_1741697840778',4,6,'2025-03-11 12:57:22.864004','2025-03-11 12:57:22.864004',NULL),
(15,'Dao Minh Nhat PGNV-ĐATN.xlsx','file','5_1742750722048',5,1,'2025-03-23 17:25:22.800362','2025-03-23 17:25:22.800362','2025-03-23 17:31:22.734826'),
(14,'DaoMinhNhat_GR1_v2.pdf','file','55_1742747389040',5,2,'2025-03-23 16:29:49.507412','2025-03-23 16:29:49.507412',NULL),
(16,'123','video','hYbldLetD8M',5,3,'2025-03-23 17:59:35.294316','2025-03-23 17:59:35.294316',NULL),
(17,'video 2','video','bBEo3sdzf8w',6,1,'2025-03-25 14:51:11.150723','2025-03-25 14:51:11.150723',NULL),
(18,'week8.pdf','file','6_1742914289537',6,2,'2025-03-25 14:51:30.159097','2025-03-25 14:51:30.159097',NULL),
(21,'Bộ luật lao động .pdf','file','7_1742964968576',7,1,'2025-03-26 04:56:10.290709','2025-03-26 04:56:10.290709',NULL),
(20,'Lời khuyên trước khóa học','video','-jV06pqjUUc',7,2,'2025-03-26 04:49:29.032798','2025-03-26 04:49:29.032798','2025-03-26 04:56:29.300828'),
(19,'Javascript có thể làm được gì? ','video','0SJE9dYdpps',7,3,'2025-03-26 04:48:52.347421','2025-03-26 04:48:52.347421','2025-03-26 04:56:31.118803'),
(22,' Lời khuyên trước khóa học','video','-jV06pqjUUc',7,4,'2025-03-26 15:39:32.417685','2025-03-26 15:39:32.417685',NULL),
(23,'Cài đặt môi trường, công cụ phù hợp để học JavaScript','video','efI98nT8Ffo',7,5,'2025-03-26 15:40:15.376286','2025-03-26 15:40:15.376286','2025-03-26 16:34:42.44514'),
(24,'Cách sử dụng JS trong file HTML','video','W0vEUmyvthQ',8,1,'2025-03-26 15:43:06.296531','2025-03-26 15:43:06.296531',NULL),
(25,'Khai báo biến','video','CLbx37dqYEI',8,2,'2025-03-26 15:44:02.596127','2025-03-26 15:44:02.596127',NULL),
(26,'Sử dụng Comments trong JavaScript','video','xRpXBEq6TOY',8,3,'2025-03-26 15:44:22.171347','2025-03-26 15:44:22.171347',NULL),
(27,'Một số hàm built-in trong JavaScript','video','rSV33HGotgE',8,4,'2025-03-26 15:44:56.091325','2025-03-26 15:44:56.091325',NULL),
(28,'  VN  Bỏ qua điều hướng Tìm kiếm     Tạo  6  Hình ảnh đại diện Làm quen với toán tử trong JavaScript','video','SZb-N7TfPlw',9,1,'2025-03-26 15:48:18.395063','2025-03-26 15:48:18.395063','2025-03-26 15:48:26.783096'),
(29,'Làm quen với toán tử trong JavaScript','video','SZb-N7TfPlw',9,2,'2025-03-26 15:48:48.80639','2025-03-26 15:48:48.80639',NULL),
(30,'Toán tử số học trong JavaScript','video','m_h7-dgKnMU',9,3,'2025-03-26 15:49:24.911648','2025-03-26 15:49:24.911648',NULL),
(31,'Toán tử ++ -- với tiền tố & hậu tố (Prefix & Postfix) ','video','aM-DUx6Qnc8',9,4,'2025-03-26 15:49:51.333958','2025-03-26 15:49:51.333958',NULL),
(32,' Toán tử gán trong JavaScript','video','ncRmjazgsE8',9,5,'2025-03-26 15:50:10.479047','2025-03-26 15:50:10.479047',NULL),
(33,'Hàm trong JavaScript','video','4g9ENVc2KLA',10,1,'2025-03-26 15:53:17.861411','2025-03-26 15:53:17.861411',NULL),
(34,'Tham số trong hàm','video','jE6UPl17Nvo',10,2,'2025-03-26 15:53:47.596092','2025-03-26 15:53:47.596092',NULL)
